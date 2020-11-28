# game
Implementation of DownFlux, a collaborative RTS.

## Bazel

```bash
$ bazel --version
bazel 3.2.0

$ bazel test $(bazel query //... except //client/...)
```
