package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client is a GitHub API client.
type Client interface {
	CreateLabel(ctx context.Context, owner, repo string, label Label) error
	GetLabels(ctx context.Context, owner, repo string) ([]Label, error)
	UpdateLabel(ctx context.Context, owner, repo string, label Label) error
	DeleteLabel(ctx context.Context, owner, repo, name string) error
}

type client struct {
	githubClient *github.Client
	token        string
}

// NewClient returns Client.
func NewClient(token string) Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &client{
		githubClient: github.NewClient(tc),
	}
}

func (c *client) CreateLabel(ctx context.Context, owner, repo string, label Label) error {
	l := &github.Label{
		Name:        &label.Name,
		Description: &label.Description,
		Color:       &label.Color,
	}
	_, _, err := c.githubClient.Issues.CreateLabel(ctx, owner, repo, l)
	return err
}

func (c *client) GetLabels(ctx context.Context, owner, repo string) ([]Label, error) {
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

func (c *client) UpdateLabel(ctx context.Context, owner, repo string, label Label) error {
	l := &github.Label{
		Name:        &label.Name,
		Description: &label.Description,
		Color:       &label.Color,
	}
	_, _, err := c.githubClient.Issues.EditLabel(ctx, owner, repo, label.Name, l)
	return err
}

func (c *client) DeleteLabel(ctx context.Context, owner, repo, name string) error {
	_, err := c.githubClient.Issues.DeleteLabel(ctx, owner, repo, name)
	return err
}
