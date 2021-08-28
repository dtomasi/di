# fakr

[![CodeFactor](https://www.codefactor.io/repository/github/dtomasi/fakr/badge)](https://www.codefactor.io/repository/github/dtomasi/di)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/dtomasi/fakr/main.svg)](https://results.pre-commit.ci/latest/github/dtomasi/fakr/main)
![Go Unit Tests](https://github.com/dtomasi/fakr/actions/workflows/build.yml/badge.svg)
![CodeQL](https://github.com/dtomasi/fakr/actions/workflows/codeql-analysis.yml/badge.svg)

A fake sink for logr interface

See: https://github.com/go-logr/logr

## Installation

    go get -u github.com/dtomasi/fakr

## Usage

```go
import (
    "github.com/dtomasi/fakr"
)

logger := fakr.New()
```
