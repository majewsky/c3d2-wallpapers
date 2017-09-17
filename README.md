# C3D2 wallpapers

These are procedurally generated wallpapers based on the [C3D2](https://c3d2.de) [logo](https://wiki.c3d2.de/Logo).
Build the images with `make`, and look at the Makefile for details.

# Variants

* The `v1` wallpaper is based on a uniformly-distributed point cloud generated with Mitchell's best-candidate algorithm. I modified the algorithm to take in an image (of course, the C3D2 logo, although any other image will also work) and adjust the density of the point cloud based on the brightness of the input image at this location.

* More variants might appear in the future.

(TODO: add small screenshots here)
