# PPROF

- Type profiles of node are available at `http://<NODE_IP>:<PPROF_PORT>/debug/pprof/`   
- You can get profile data to local by command:  
`curl http://<NODE_IP>:<PPROF_PORT>/debug/pprof/<PROFILE_NAME> > pprof.out`  
Ex: `go tool pprof ./gev http://<NODE_IP>:<PPROF_PORT>/debug/pprof/heap > pprof.out`   
- After that, you can see the visualization of `pprof.out` by this command `go tool pprof -web ./gev ./pprof.out`. It will show you the result after analyzing `./gev` on webbrowser.
- You can do more with `pprof.out` by `go tool pprof ./gev ./pprof.out`

Here is some profiles you can use directly:
## CPU Profile  
`go tool pprof http://<NODE_IP>:<PPROF_PORT>/debug/pprof/profile`  
The CPU profiler runs for 30 seconds by default. It uses sampling to determine which functions spend most of the CPU time. The Go runtime stops the execution every 10 milliseconds and records the current call stack of all running goroutines.

When pprof enters the interactive mode, type `top`, the command will show a list of functions that appeared most in the collected samples. In our case these are all runtime and standard library functions.

There is a much better way to look at the high-level performance overview - web command, it generates an SVG graph of hot spots and opens it in a web browser.   

## Heap Profile
`go tool pprof http://<NODE_IP>:<PPROF_PORT>/debug/pprof/heap`
By default it shows the amount of memory currently in-use.
But we are more interested in the number of allocated objects. Call pprof with -alloc_objects option:  
`go tool pprof -alloc_objects http://<NODE_IP>:<PPROF_PORT>/debug/pprof/heap`

## Goroutine Profile
Goroutine profile dumps the goroutine call stack and the number of running goroutines:  
`go tool pprof http://<NODE_IP>:<PPROF_PORT>/debug/pprof/goroutine`  

## Block Profile   
Blocking profile shows function calls that led to blocking on synchronization primitives like mutexes and channels.   
`go tool pprof http://<NODE_IP>:<PPROF_PORT>/debug/pprof/block`  

...