# Label Syncer

[![CI](https://github.com/micnncim/action-label-syncer/workflows/CI/badge.svg)](https://github.com/micnncim/action-label-syncer/actions)
[![Release](https://img.shields.io/github/v/release/micnncim/action-label-syncer.svg?logo=github)](https://github.com/micnncim/action-label-syncer/releases)
[![Marketplace](https://img.shields.io/badge/marketplace-label--syncer-blue?logo=github)](https://github.com/marketplace/actions/label-syncer)

Action to sync GitHub labels in the declarative way.  
By using this action, you can sync current labels with configured labels in a YAML manifest.

## Usage

### Create YAML manifest of GitHub labels

```yaml
- color: d73a4a
  description: Something isn't working
  name: bug
- color: 0075ca
  description: Improvements or additions to documentation
  name: documentation
- color: cfd3d7
  description: This issue or pull request already exists
  name: duplicate
```

![](./docs/assets/screenshot.png)

The default file path is `.github/labels.yml`, but you can specify any file path with `jobs.<job_id>.steps.with`.  
To create manifest of GitHub labels for the current status of labels easily, we recommend using [label-exporter](https://github.com/micnncim/label-exporter).

### Create Workflow

```yaml
name: Sync labels in the declarative way
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@1.0.0
      - uses: micnncim/action-label-syncer@latest
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GITHUB_REPOSITORY: ${{ github.repository }}
        with:
          manifest: labels.yml # default: .github/labels.yml
```

## See also

- [actions/labeler](https://github.com/actions/labeler)
- [lannonbr/issue-label-manager-action](https://github.com/lannonbr/issue-label-manager-action)

## LICENSE

[MIT License](./LICENSE)
