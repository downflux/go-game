# DownFlux Client

[g3n](https://github.com/g3n/engine) cannot be used in Bazel without extensive
work; as of now, we will have to accept running client code outside of Bazel,
and treat server / shared code as a normal Go repo.

g3n seems to require OpenGL 3.3 -- remember to use VMWare (vs. VirtualBox), as
VirtualBox does not support OpenGL 3+.

Run `glxgears` to test if the guest OS supports the correct OpenGL version.
