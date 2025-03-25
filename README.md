# go-particle-life-simulation
A simple project I did to replicate particle life with go using the ebitengine library
# Particle Life Engine in Golang
This was a project for school.

## The Concept
The idea was to make a particle simulator that uses attraction and repulsion between particles to create interesting interactions and behaviours between the particles.

## How I Did It
The Engine
After a bit of research I found the simple game engine called Ebitengine, a minimalistic game/simulation engine, that only has the most basic functions, and easy to use setup functions. I found out how to open a window, and draw to the screen, and was up and running with a static ball.

### Motion and basic development.
There is an update function in the engine where I can put all of my logic to be run on each frame. I created a struct for each circle that contained all the info about each particle, like velocity, position, radius, and more. I then got the velocity to work, with the velocity of X and Y being added to the particle at the end of every frame. I got the circles moving, and then some code to find the distance between every particle in the simulation and, and applied some force according to that difference (assuming it was under the attractionThreshold defined at the top of the file that tells us how far away particles can interact). The next thing to do was implement collisions. That took a lot of googling, and a lot of math, but I eventually got it.

### Colour Attractions
Now all that was left was to change the attraction and repulsion of particles based on their colour. Looking back, I did it terribly, and should have used a hash table for the colours, where the colour was they key, and the value would be a data object with colour and forces, but the way it did it was to have two separate tables for both the colour interactions and the actual colours for some reason.

### Other Parameters
At the top of the source file you will see adjustable parameters for gravity, speed multipliers, and other. These are constants defined at compile time that are factored into physics calculations. You can change these values and compile to see what changes and happens.
