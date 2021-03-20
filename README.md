# C3D2 wallpapers

These are procedurally generated wallpapers based on the [C3D2](https://c3d2.de) [logo](https://wiki.c3d2.de/Logo).
Build the images with `make`, and look at the Makefile for details.

# Variants

More variants might appear in the future when I'm inspired.

![v1 wallpaper](https://static.bethselamin.de/c3d2-wallpapers/v1-1920x1080.png)

The `v1` wallpaper is based on a uniformly-distributed point cloud generated with Mitchell's best-candidate algorithm. I modified the algorithm to take in an image (of course, the C3D2 logo, although any other image will also work) and adjust the density of the point cloud based on the brightness of the input image at this location.

* [1920x1080 PNG](https://static.bethselamin.de/c3d2-wallpapers/v1-1920x1080.png)
* [3840x2160 PNG](https://static.bethselamin.de/c3d2-wallpapers/v1-3840x2160.png)

![v2 wallpaper](https://static.bethselamin.de/c3d2-wallpapers/v2-1920x1080.png)

The `v2` wallpaper is based on a grid of equilateral triangles, colorized with a smoothened white noise, in which a set of triangles forming the C3D2 logo is highlighted.

* [1920x1080 PNG](https://static.bethselamin.de/c3d2-wallpapers/v2-1920x1080.png)
* [3840x2160 PNG](https://static.bethselamin.de/c3d2-wallpapers/v2-3840x2160.png)
* [16:9 SVG](https://static.bethselamin.de/c3d2-wallpapers/v2.svg)

## TODO: Ideas for more wallpapers

- print the logo with really coarse offset printing
- make a complex Fourier transform of the logo path(s), truncate to highest-order powers only, draw several such approximations over the logo ([context](https://www.youtube.com/watch?v=r6sGWTCMz2k))

Will do these at some point. Maybe.
