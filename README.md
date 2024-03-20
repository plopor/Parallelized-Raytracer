# Parallelized Raytracer

## Benchmarking graphs and analysis of results can be found in the rt/benchmark directory.

The input command for this project is as follows:

go run main.go numthreads steal skew

numthreads represents the number of threads
steal is a boolean representing whether or not we do work stealing
skew determines if the work on each thread is evenly distributed

If numthreads is 1, then the program will run the sequential implementation. This means steal and skew are irrelevant.

Steal set to true means each worker thread will attempt to rebalance the DEQueues of itself with another victim at random, weighted based on the 
current number of tasks it has in its own queue. e.g. fewer tasks left in its queue -> greater chance it will attempt to balance work with another thread

Skew set to true means that the first numthreads/2 threads will have 1/3 the number of tasks as the other half of threads.

Customization on the scene is defined through a json found in scene.json. Here you can specify the size, placement, and color of spheres in the scene 
as well as placement the color of a light. The specular field determines the shininess of the sphere for the illumination model described in the Analysis.
You can also specify how far back the camera is in the scene's Z field and its FOV along with the final image dimensions, Gaussian kernel parameters, and 
max number of reflections per casted ray.

The 3D scene that is rendered comes with a bounding box that matches the dimensions of the image, with the scene's depth field determining how far back
the box extends. The reason the Z field is initialized so far back at -5000 is because the image needs to be projected onto a plane which sits in front
of where the camera ("eye") is placed. The camera is always placed in the center of the x,y plane which is defined as width/2, height/2.

Please note that since the parallel implementation relies on a bounded DEQueue, the capacity of each queue I have set by default to be 48,000,000 which 
represents the max number of tasks that a single thread can have (by default the 6k x 8k image). If you increase the image resolution by more than a
a factor of numthreads, you may overflow the DEQueue. Similarly, since each queue is initialized as such, increasing the number of threads too much may
result in the heap running out of memory to allocate.

The scene.json file must be present in the same directory where you made the call to main.go. Thus, running the benchmark-rt.sh file for the cluster
under rt/benchmark will use the scene.json file in rt/benchmark. Similarly, running it locally by invoking the above command in the rt directory
will use the local scene.json file to generate the image (which outputs in the same directory).

To run an analysis experiment of this project on the cluster is as follows: From the rt/benchmark directory,

sbatch benchmark-rt.sh

This command will run the sequential implementation as well as each of the 4 combinations of steal/skew options for each of the {2, 4, 6, 8, 12} threads
5 times. Using the average of the 5 runs for each of the 4 option combinations, it will then generate a speedup graph over the number of threads.

Note that you may need to convert DOS line endings to UNIX line endings with a command similar to:

sed 's/\r$//' benchmark-rt.sh > benchmark-rt-unix.sh
