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
