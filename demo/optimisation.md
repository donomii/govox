I love optimisation articles.  The author gets to write terrible code, then pat himself on the back for making it better.  It's my turn now.

I recently devloped the obsession to take a simple voxel engine and slap it on an existing roguelike, like DCSS Crawl.  That sentence is basically completely made up of red flags, some of the worst being "a simple voxel engine" and also "slap it on", like code is suncreen and all I have to do is rub it in thoroughly and it will work.

So I grabbed a simple voxel engine from github called govox, and it worked reasonably well, although in hindsight, no.  It rendered a 50x50x50 cube of voxels fast enough for me to do some maths visualisations on it.

(insert visualisation screenshot here)

I wanted to have the standard crawl level display, which is a minimum of 10 squares radius around the player, plus some extra tiles "in the dark" to keep the player oriented in the level.

(insert crawl screenshot)

However some quick calculations made it clear that this wasn't going to work, like, at all.  A 10 tiles radius means a minimum of 21 tiles for each edge of the square.  This means I have a budget of 2x2x2 voxels for each tile.  Clearly, this game will have to be highly abstract.  I give it a try anyway, and it turns out to be kind of fun.  I resolve to keep this as a game mode.

(insert abstract rogue picture)

But I want a map that looks like a dungeon, with monsters and a guy with a sword.  I resolve, with absolutely no reason, that I will rewrite the voxel engine to use instancing, so I can draw billions of voxel per frame.  Instancing is where we tell the graphics card to draw duplicate cubes quickly.  We'll never know how much that would have improved the situation because I suck at graphics programming and failed to implement it.  Based on future revelations, it probably would not have improved things much.

At this point, the game is unplayable.  It frequently ignores keypresses, and runs at about 3fps.  I'm sad and frustrated.  I decide, again for very little reason, to make it multithreaded.  Luckily, not only does this work, it reveals where my problems are.

I break the program up into the obvious sub-components:

	voxel prep (blocks)  -> graphics prep (blocks to VBOs)  ->  draw screen (VBOs)

Each stage runs in its own thread.  While the renderer is drawing the screen, the next frame is being copied into (different) VBOs, and the other thread is rebuilding the voxel world.  At the same time, I switch the draw method.  Previously, it uploaded one cube, then drew it over and over again (confusingly, without using instancing).  This seems like a good idea, but it requires at least three calls to the graphics card per block (set colour, set position, draw).  This means for a 100x100x100 voxel world, there would be a million calls going out over the graphics bus.  This is too much.

To be clear, a modern graphics card can easily draw a million cubes, but it can't handle a million separate graphics calls.  To improve this, I copy the positions of all the cubes into one array, all the colours for those cubes into a separate array, then upload both arrays and tell the graphics card to run through the arrays, drawing cubes.  This takes us from 1 million calls per frame down to 3.  A significant improvement, although experienced readers will note that we would never actually draw a million voxels in one frame.  I was actually drawing around 11,000.  Still, going from 33,000 calls to 3 is still a big improvement.  Let's take a look at the timings.


Code                                      | Gl Draw | graphics prep | voxel prep
Initial implementation(!)                   | 600     |  0            | 0
Multithreaded drawarrays  slow laptop     | 21      | 23            | 302 
Multithreaded drawarrays  fast laptop     | 1       | 10            | 220

!  There is only one stage in the initial implementation, everything was mixed together in the one loop.

I split the 600ms loop into three loops, and while one of them takes 300ms to run, the other two take almost no time at all.  That was a surprise to me.  The graphics card is clearly more than capable, even on my old slow laptop.  The trouble is with the CPU, dealing with the large number of blocks.

Despite that, I made one additional improvement.  Instead of drawing cubes, I draw points.  so it is rendering a point cloud instead of voxels.  Provided the points are wide enough, this works nicely, and is almost indistinguishable from cubes.

So now, the challenge is to speed up the CPU voxel handling.

Unfortunately for this article, it turned out to be trivially simple.  In the voxel prep, I was looping over all 9 million voxels and setting them to inactive, each frame.  It turns out golang is not particularly fast at assigning to arrays, so it was taking around 200ms to clear the voxel array.  In fact, it was much faster to reallocate the array than it was to clear it, so that's what I did.

There are some fast routines in the C library for clearing memory, so later on I will investigate the best way to speed this up.  For now though, it is good enough.

Code                                      | Gl Draw | graphics prep | voxel prep
Initial implementation                    | 600     |  0            | 0
Multithreaded drawarrays  slow laptop     | 21      | 23            | 302 
Multithreaded drawarrays  fast laptop     | 10      | 10            | 220
MT + allocs  slow laptop                  | 17      | 14            | 17  
MT + allocs  fast laptop                  | 16      | 10            | 15

I have no idea why the Gl Draw thread is now the slowest part of the whole process.  Certainly something to look into.  The numbers for the graphics card are bouncing around all over the place, going from 0-22.
