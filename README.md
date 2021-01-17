# game
Implementation of DownFlux, a collaborative RTS.

![Server CI](https://github.com/downflux/game/workflows/Server%20CI/badge.svg?branch=main)

## Setup

### Cloning Repo

1. Set up `git`

   ```bash
   git config --global push.default current
   ```

1. Copy the SSH keys into `~/.ssh` and set permission to `600` for `id_rsa` and
   `644` for `id_rsa.pub`.

1. Add to SSH keychain

   ```bash
   ssh-add
   ```

1. Test SSH access to GitHub

   ```bash
   ssh -T git@github.com
   ```

1. Clone

   ```bash
   git clone git@github.com:downflux/game.git
   ```

## Bazel

```bash
$ bazel --version
bazel 3.2.0

$ bazel test -c opt \
    $(bazel query //... except //client/...) \
    --nocache_test_results \
    --runs_per_test=100
```
