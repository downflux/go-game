# game
Implementation of DownFlux, a collaborative RTS.

![Server CI](https://github.com/downflux/game/workflows/Server%20CI/badge.svg?branch=main)

## Setup

### Install Git LFS

See [official documentation](https://git-lfs.github.com/).

Install the `git-lfs` package

```bash
sudo apt install git-lfs
```

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

1. Install LFS for the repo

   ```bash
   git lfs install
   ```

### Installing Bazel

Install Bazel via
[official docs](https://docs.bazel.build/versions/master/install-ubuntu.html#install-on-ubuntu).

Current verified Bazel version is `3.7.2`.

```bash
bazel test -c opt \
  --features race ... \
  --nocache_test_results \
  --runs_per_test=10
```

### CPU Profiler

```bash
sudo apt install graphviz gv
bazel run -c opt \
  //server/grpc:main -- \
  --cpuprofile=${F}
go tool pprof -http=localhost:8888 ${F}
```
