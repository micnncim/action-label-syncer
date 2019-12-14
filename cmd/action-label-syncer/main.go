package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/micnncim/action-label-syncer/pkg/github"
)

func main() {
	manifest := os.Getenv("INPUT_MANIFEST")
	labels, err := github.FromManifestToLabels(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load manifest: %v\n", err)
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	cli := github.NewClient(token)

	repository := os.Getenv("GITHUB_REPOSITORY")
	slugs := strings.Split(repository, "/")
	if len(slugs) != 2 {
		fmt.Fprintf(os.Stderr, "invalid repository: %v\n", repository)
		os.Exit(1)
	}
	owner, repo := slugs[0], slugs[1]

	ctx := context.Background()
	s := github.NewLabelSyncer(cli)
	if err := s.SyncLabels(ctx, owner, repo, labels); err != nil {
		fmt.Fprintf(os.Stderr, "unable to sync labels: %v\n", err)
		os.Exit(1)
	}
}
