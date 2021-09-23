// Copyright 2020 micnncim
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

type Client struct {
	githubClient *github.Client
	token        string
}

type Label struct {
	// If "import" is present, all other fields are ignored.
	Import      string `yaml:"import"`
	Name        string `yaml:"name"`
	Alias       string `yaml:"alias"`
	Aliases   []string `yaml:"aliases"`
	Description string `yaml:"description"`
	Color       string `yaml:"color"`
}

func FromManifestToLabels(path string) ([]Label, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var labels []Label
	err = yaml.Unmarshal(buf, &labels)
	if err != nil {
		return nil, err
	}

	// Handle imports of labels from another file
	var flatLabels []Label
	for _, l := range labels {
		if l.Import == "" {
			if len(l.Description) > 100 {
				return nil, fmt.Errorf("Description of \"%s\" exceeds 100 characters", l.Name)
			}
			flatLabels = append(flatLabels, l)
		} else {
			importPath := filepath.Join(filepath.Dir(path), l.Import)
			importedLabels, err := FromManifestToLabels(importPath)
			if err != nil {
				return nil, err
			}
			flatLabels = append(flatLabels, importedLabels...)
		}
	}

	return flatLabels, err
}

func NewClient(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &Client{
		githubClient: github.NewClient(tc),
	}
}

func (c *Client) SyncLabels(ctx context.Context, owner, repo string, labels []Label, prune bool, dryRun bool) error {
	if dryRun {
		fmt.Printf("Dry run!  No actual changes will be made.\n")
	}

	labelMap := make(map[string]Label)
	aliasMap := make(map[string]Label)
	for _, l := range labels {
		labelMap[l.Name] = l
		if l.Alias != "" {
			aliasMap[l.Alias] = l
		}
		for _, alias := range l.Aliases {
			aliasMap[alias] = l
		}
	}

	currentLabels, err := c.getLabels(ctx, owner, repo)
	if err != nil {
		return err
	}
	currentLabelMap := make(map[string]Label)
	for _, l := range currentLabels {
		currentLabelMap[l.Name] = l
	}

	eg := errgroup.Group{}

	// Delete labels.
	if prune {
		for _, currentLabel := range currentLabels {
			currentLabel := currentLabel
			eg.Go(func() error {
				_, name_ok := labelMap[currentLabel.Name]
				_, alias_ok := aliasMap[currentLabel.Name]
				if (alias_ok && !name_ok) || name_ok {
					return nil
				}
				return c.deleteLabel(ctx, owner, repo, currentLabel.Name, dryRun)
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}
	}

	// Create and/or update labels.
	for _, l := range labels {
		l := l
		eg.Go(func() error {
			labelName := l.Name
			currentLabel, ok := currentLabelMap[l.Name]
			if !ok {
				currentLabel, ok = currentLabelMap[l.Alias]
				if !ok {
					return c.createLabel(ctx, owner, repo, l, dryRun)
				}
				labelName = l.Alias
			}
			if currentLabel.Description != l.Description || currentLabel.Color != l.Color || currentLabel.Name != l.Name {
				return c.updateLabel(ctx, owner, repo, labelName, l, dryRun)
			}
			//fmt.Printf("Not changed: \"%s\" on %s/%s\n", l.Name, owner, repo)
			return nil
		})
	}

	return eg.Wait()
}

func (c *Client) createLabel(ctx context.Context, owner, repo string, label Label, dryRun bool) error {
	l := &github.Label{
		Name:        &label.Name,
		Description: &label.Description,
		Color:       &label.Color,
	}
	fmt.Printf("Created: \"%s\" on %s/%s\n", label.Name, owner, repo)
	if dryRun {
		return nil
	}
	_, _, err := c.githubClient.Issues.CreateLabel(ctx, owner, repo, l)
	return err
}

func (c *Client) getLabels(ctx context.Context, owner, repo string) ([]Label, error) {
	opt := &github.ListOptions{
		PerPage: 50,
	}
	var labels []Label
	for {
		ls, resp, err := c.githubClient.Issues.ListLabels(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}
		for _, l := range ls {
			labels = append(labels, Label{
				Name:        l.GetName(),
				Description: l.GetDescription(),
				Color:       l.GetColor(),
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return labels, nil
}

func (c *Client) updateLabel(ctx context.Context, owner, repo, labelName string, label Label, dryRun bool) error {
	l := &github.Label{
		Name:        &label.Name,
		Description: &label.Description,
		Color:       &label.Color,
	}
	if labelName != label.Name {
		fmt.Printf("Renamed: \"%s\" => \"%s\" on %s/%s\n", labelName, label.Name, owner, repo)
	} else {
		fmt.Printf("Updated: \"%s\" on %s/%s\n", label.Name, owner, repo)
	}
	if dryRun {
		return nil
	}
	_, _, err := c.githubClient.Issues.EditLabel(ctx, owner, repo, labelName, l)
	return err
}

func (c *Client) deleteLabel(ctx context.Context, owner, repo, name string, dryRun bool) error {
	fmt.Printf("Deleted: \"%s\" on %s/%s\n", name, owner, repo)
	if dryRun {
		return nil
	}
	_, err := c.githubClient.Issues.DeleteLabel(ctx, owner, repo, name)
	return err
}
