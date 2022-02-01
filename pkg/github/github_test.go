package github_test

import (
	"path/filepath"
	"testing"

	"github.com/micnncim/action-label-syncer/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromManifestToLabels(t *testing.T) {
	labels, err := github.FromManifestToLabels(filepath.Join("testdata", "labels.yml"))
	require.NoError(t, err)
	assert.ElementsMatch(t, labels, []github.Label{
		{
			Name:        "bug",
			Description: "Something isn't working",
			Color:       "d73a4a",
		},
		{
			Name:        "documentation",
			Description: "Improvements or additions to documentation",
			Color:       "0075ca",
		},
		{
			Name:        "duplicate",
			Description: "This issue or pull request already exists",
			Color:       "cfd3d7",
		},
	})
}
