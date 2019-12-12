![logo](docs/assets/logo.png)

[![CI](https://github.com/micnncim/action-label-syncer/workflows/CI/badge.svg)](https://github.com/micnncim/action-label-syncer/actions)
[![Release](https://img.shields.io/github/v/release/micnncim/action-label-syncer.svg?logo=github)](https://github.com/micnncim/action-label-syncer/releases)
[![Marketplace](https://img.shields.io/badge/marketplace-label--syncer-blue?logo=github)](https://github.com/marketplace/actions/label-syncer)

GitHub Actions workflow to sync GitHub labels in the declarative way.  

By using this workflow, you can sync current labels with labels configured in a YAML manifest.

## Usage

### Create YAML manifest of GitHub labels

```yaml
- name: bug
  description: Something isn't working
  color: d73a4a
- name: documentation
  description: Improvements or additions to documentation
  color: 0075ca
- name: duplicate
  description: This issue or pull request already exists
  color: cfd3d7
```

![](./docs/assets/screenshot.png)

The default file path is `.github/labels.yml`, but you can specify any file path with `jobs.<job_id>.steps.with`.

To create manifest of the current labels easily, using [label-exporter](https://github.com/micnncim/label-exporter) is recommended.

### Create Workflow

An workflow example is here.

```yaml
name: Sync labels in the declarative way
on:
  push:
    branches:
      - master
    paths:
      - path/to/labels.yml
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@1.0.0
      - uses: micnncim/action-label-syncer@v0.4.0
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          manifest: labels.yml # optional; default path is `.github/labels.yml`
```

## See also

- [actions/labeler](https://github.com/actions/labeler)
- [lannonbr/issue-label-manager-action](https://github.com/lannonbr/issue-label-manager-action)
- [b4b4r07/github-labeler](https://github.com/b4b4r07/github-labeler)

## LICENSE

[MIT License](./LICENSE)

## Note

*Icon made by bqlqn from [www.flaticon.com](https://www.flaticon.com)*
