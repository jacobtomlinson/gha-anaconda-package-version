# How to make a release

This action is published to the [GitHub Actions Marketplace](https://github.com/marketplace/actions/get-anaconda-package-version).

To make a release you need push rights to the [jacobtomlinson/gha-anaconda-package-version GitHub repository](https://github.com/jacobtomlinson/gha-anaconda-package-version).

## Steps to make a release

1. Create a new _GitHub release_ by clicking on `Draft a new release` at
   https://github.com/jacobtomlinson/gha-anaconda-package-version/releases and
   creating a new tag for `master`. Remember to prefix the tag with `v`, e.g.
   `v1.0.5`

1. Check everything again, then click `Publish release`.
   This will create a tag, and trigger a CI action to update tags like `v1`.
