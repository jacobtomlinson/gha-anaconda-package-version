# Anaconda Package Version

[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Anaconda%20Package%20Version-blue.svg?colorA=24292e&colorB=0366d6&style=flat&longCache=true&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAA4AAAAOCAYAAAAfSC3RAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAM6wAADOsB5dZE0gAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAERSURBVCiRhZG/SsMxFEZPfsVJ61jbxaF0cRQRcRJ9hlYn30IHN/+9iquDCOIsblIrOjqKgy5aKoJQj4O3EEtbPwhJbr6Te28CmdSKeqzeqr0YbfVIrTBKakvtOl5dtTkK+v4HfA9PEyBFCY9AGVgCBLaBp1jPAyfAJ/AAdIEG0dNAiyP7+K1qIfMdonZic6+WJoBJvQlvuwDqcXadUuqPA1NKAlexbRTAIMvMOCjTbMwl1LtI/6KWJ5Q6rT6Ht1MA58AX8Apcqqt5r2qhrgAXQC3CZ6i1+KMd9TRu3MvA3aH/fFPnBodb6oe6HM8+lYHrGdRXW8M9bMZtPXUji69lmf5Cmamq7quNLFZXD9Rq7v0Bpc1o/tp0fisAAAAASUVORK5CYII=)](https://github.com/jacobtomlinson/gha-anaconda-package-version)
[![Actions Status](https://github.com/jacobtomlinson/gha-anaconda-package-version/workflows/Build/badge.svg)](https://github.com/jacobtomlinson/gha-anaconda-package-version/actions)
[![Actions Status](https://github.com/jacobtomlinson/gha-anaconda-package-version/workflows/Integration%20Test/badge.svg)](https://github.com/jacobtomlinson/gha-anaconda-package-version/actions)

A GitHub Action to get the latest version of a package from [Anaconda](https://anaconda.org).

## Usage

### Example workflow

```yaml
name: My Workflow
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: Run action
      uses: jacobtomlinson/gha-anaconda-package-version@master
      with:
        org: anaconda
        package: python
```

### Inputs

| Input                                             | Description                                        |
|------------------------------------------------------|-----------------------------------------------|
| `org`  | The Anaconda user or organization    |
| `package` | The name of the Python package    |

### Outputs

| Output                                             | Description                                        |
|------------------------------------------------------|-----------------------------------------------|
| `version`  | The version of the package    |

## Examples

### Using outputs

Here is an example of getting the version of the latest Anaconda Python distribution and printing it out in the next step.

```yaml
name: My Workflow
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - name: Get latest Anaconda Python version
      id: anaconda
      uses: jacobtomlinson/gha-anaconda-package-version@master
      with:
        org: anaconda
        package: python

    - name: Check outputs
        run: |
          echo "The latest version of Python is ${{ steps.anaconda.outputs.version }}."
```