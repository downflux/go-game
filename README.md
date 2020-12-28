# game
Implementation of DownFlux, a collaborative RTS.

## Bazel

```bash
$ bazel --version
bazel 3.2.0

$ bazel test -c opt \
    $(bazel query //... except //client/...) \
    --nocache_test_results \
    --runs_per_test=100
```
