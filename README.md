# Heightmap terrain with an animated energy dome

**Disclaimer:** *This is nothing much more than a little fun project. I wanted to try some ideas and effects and decided to put it all in this project.
It turned out to be fun and nice looking. Don't ask me what that is for exactly ;)*

## How does it look? (Just to keep you interested :))

This is how:
![Energy dome](https://github.com/MauriceGit/RealHeightmapTerrain/blob/master/Screenshots/final_energy_sphere.png "a colored heightmap terrain with an energy dome")

## What is that?
This is a heightmap from a real-world-location (Böblingen - South Germany in this case), taken from [terrain.party](http://terrain.party/), that is
just colored to look interesting and to emphasize the height difference.

On top of that terrain is an animated energy dome (could be some kind of defense mechanism/range?). The energy flow will not repeat itself exactly and looks very natural and nice.
This idea is based on a [youtube video](https://www.youtube.com/watch?v=zLSPaE1qsvM&feature=youtu.be) (watched it once, did not look at any code, credits for the idea anyway!).

## What exactly is implemented here?

* Heightmap taken from a png, applied on the GPU
* Realtime terrain tessellation and Level-of-Detail (LOD) based on camera distance. Calculated in the tessellation shader.
* Compute shader for calculating the normals based on the heightmap for nicer shading (turns out 8bit color is not much, can't really use the normals without further post-processing).
* The energy dome, implemented using multiple simplex noise images. One for the energy, one for further animation.
* Image distortion where the energy is high. As the energy dome is rendered as a post-processing effect of the scene, there is a uv-distortion based on energy level.

## What now?

I hope, that some of you can take some ideas or implementation details and use it in your own projects. If so, please let me know, I am always interested :)

Most likely, I will use some parts of my own projects in later ones myself. Always good, to have some references.

## Some more Screenshots

A little further back:
![Energy dome](https://github.com/MauriceGit/RealHeightmapTerrain/blob/master/Screenshots/final_energy_sphere2.png "a colored heightmap terrain with an energy dome")

Visible LOD (level of detail) on the mesh:
![Energy dome](https://github.com/MauriceGit/RealHeightmapTerrain/blob/master/Screenshots/terrain_distance_based_lod.png "LOD")

How it looked at some point before:
![Energy dome](https://github.com/MauriceGit/RealHeightmapTerrain/blob/master/Screenshots/energy_sphere_multisampled.png "before color adjustments")

## How to run it:

To run this demo, you only need **Golang** installed on your system. And a graphics card, that supports **OpenGL 4.3**.

On linux, you can just run the *compile.sh* file in the main folder.

On Windows, you have to copy the command from the *compile.sh* and run it directly.
