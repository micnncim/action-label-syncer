package fake

import (
	"context"

	"github.com/micnncim/action-label-syncer/pkg/github"
)

type Client struct {
	FakeCreateLabel func(ctx context.Context, owner, repo string, label github.Label) error
	FakeGetLabels   func(ctx context.Context, owner, repo string) ([]github.Label, error)
	FakeUpdateLabel func(ctx context.Context, owner, repo string, label github.Label) error
	FakeDeleteLabel func(ctx context.Context, owner, repo, name string) error
}

var _ github.Client = (*Client)(nil)

func (c *Client) CreateLabel(ctx context.Context, owner, repo string, label github.Label) error {
	return c.FakeCreateLabel(ctx, owner, repo, label)
}

func (c *Client) GetLabels(ctx context.Context, owner, repo string) ([]github.Label, error) {
	return c.FakeGetLabels(ctx, owner, repo)
}

func (c *Client) UpdateLabel(ctx context.Context, owner, repo string, label github.Label) error {
	return c.FakeUpdateLabel(ctx, owner, repo, label)
}

func (c *Client) DeleteLabel(ctx context.Context, owner, repo, name string) error {
	return c.FakeDeleteLabel(ctx, owner, repo, name)
}
