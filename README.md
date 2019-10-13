# Action Labels

[![CI](https://github.com/micnncim/action-labels/workflows/CI/badge.svg)](https://github.com/micnncim/action-labels/actions)

Action to sync GitHub labels in the declarative way.

## Usage

### Create `.github/labels.yml`

```
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

To create `.github/labels.yml` for the current status of labels, use [label-exporter](https://github.com/micnncim/label-exporter).

### Create Workflow

```
name: Sync labels in the declarative way
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@1.0.0
      - uses: micnncim/action-labels@latest
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GITHUB_REPOSITORY: ${{ github.repository }}
```

## See also

- [actions/labeler](https://github.com/actions/labeler)
- [lannonbr/issue-label-manager-action](https://github.com/lannonbr/issue-label-manager-action)

## LICENSE

[MIT License](./LICENSE)
