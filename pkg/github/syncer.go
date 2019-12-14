package github

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// LabelSyncer implements SyncLabels.
type LabelSyncer struct {
	client Client
}

func NewLabelSyncer(client Client) *LabelSyncer {
	return &LabelSyncer{
		client: client,
	}
}

// SyncLabels syncs the current GitHub labels with labels in the manifest.
func (s *LabelSyncer) SyncLabels(ctx context.Context, owner, repo string, labels []Label) error {
	labelMap := make(map[string]Label)
	for _, l := range labels {
		labelMap[l.Name] = l
	}

	currentLabels, err := s.client.GetLabels(ctx, owner, repo)
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
			return s.client.DeleteLabel(ctx, owner, repo, currentLabel.Name)
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
				return s.client.CreateLabel(ctx, owner, repo, l)
			}
			if currentLabel.Description != l.Description || currentLabel.Color != l.Color {
				return s.client.UpdateLabel(ctx, owner, repo, l)
			}
			return nil
		})
	}

	return eg.Wait()
}
