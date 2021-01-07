# game
Implementation of DownFlux, a collaborative RTS.

![Server CI](https://github.com/downflux/game/workflows/Server%20CI/badge.svg?branch=main)

## Bazel

```bash
$ bazel --version
bazel 3.2.0

$ bazel test -c opt \
    $(bazel query //... except //client/...) \
    --nocache_test_results \
    --runs_per_test=100
```
