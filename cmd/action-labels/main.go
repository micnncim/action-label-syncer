package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micnncim/action-labels/pkg/github"
)

func main() {
	workspace := os.Getenv("GITHUB_WORKSPACE")
	// TODO: Allow to change filepath of manifest.
	path := filepath.Join(workspace, ".github", "labels.yml")
	labels, err := github.FromManifestToLabels(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load manifest: %v\n", err)
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(token)

	repository := os.Getenv("GITHUB_REPOSITORY")
	slugs := strings.Split(repository, "/")
	if len(slugs) != 2 {
		fmt.Fprintf(os.Stderr, "invalid repository: %v\n", repository)
		os.Exit(1)
	}
	owner, repo := slugs[0], slugs[1]

	ctx := context.Background()
	if err := client.SyncLabels(ctx, owner, repo, labels); err != nil {
		fmt.Fprintf(os.Stderr, "unable to sync labels: %v\n", err)
		os.Exit(1)
	}
}
