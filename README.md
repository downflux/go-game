# game
Implementation of DownFlux, a collaborative RTS.

## Bazel

```bash
$ bazel --version
bazel 3.2.0

$ bazel test -c opt \
    $(bazel query //... except //client/...) \
    --runs_per_test=100
```
