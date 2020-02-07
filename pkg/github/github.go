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
	"io/ioutil"

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
	Name        string `yaml:"name"`
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
	return labels, err
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

func (c *Client) SyncLabels(ctx context.Context, owner, repo string, labels []Label) error {
	labelMap := make(map[string]Label)
	for _, l := range labels {
		labelMap[l.Name] = l
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
	for _, currentLabel := range currentLabels {
		currentLabel := currentLabel
		eg.Go(func() error {
			_, ok := labelMap[currentLabel.Name]
			if ok {
				return nil
			}
			return c.deleteLabel(ctx, owner, repo, currentLabel.Name)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	// Create and/or update labels.
	for _, l := range labels {
		l := l
		eg.Go(func() error {
			currentLabel, ok := currentLabelMap[l.Name]
			if !ok {
				return c.createLabel(ctx, owner, repo, l)
			}
			if currentLabel.Description != l.Description || currentLabel.Color != l.Color {
				return c.updateLabel(ctx, owner, repo, l)
			}
			return nil
		})
	}

	return eg.Wait()
}

func (c *Client) createLabel(ctx context.Context, owner, repo string, label Label) error {
	l := &github.Label{
		Name:        &label.Name,
		Description: &label.Description,
		Color:       &label.Color,
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

func (c *Client) updateLabel(ctx context.Context, owner, repo string, label Label) error {
	l := &github.Label{
		Name:        &label.Name,
		Description: &label.Description,
		Color:       &label.Color,
	}
	_, _, err := c.githubClient.Issues.EditLabel(ctx, owner, repo, label.Name, l)
	return err
}

func (c *Client) deleteLabel(ctx context.Context, owner, repo, name string) error {
	_, err := c.githubClient.Issues.DeleteLabel(ctx, owner, repo, name)
	return err
}
