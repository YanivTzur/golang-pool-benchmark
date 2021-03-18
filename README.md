# Object Pool Benchmarking Program

## Description
- An HTTP server with 3 dummy REST API's
- Each API's handler function:
    1. Expects to receive a buffer in the request.
    2. Reads the buffer.
    3. If an error occurred while reading, the handler returns status
       code 500.
- The difference in the API's is the method used to allocate a buffer
  to read in the content of the buffer in each of multiple,
  concurrent requests:
  - The first API `/basic-handler` allocates a new buffer each time.
  - The second API `/object-pool-handler` uses sync.Pool to optimize
    buffer allocation in terms of running time and allocated memory.
  - The third API `/bounded-pool-handler` uses a custom written pool to allocate buffers.

 ## Benchmark
 - Repository includes benchmark tests in `server_test.go`.
 - The benchmarks measure the running time, number of bytes allocated and number of allocations per operation.
 - Each benchmark does it for a different API.
 - The buffer sent in each request is a string consisting of 1024 "0"'s.
 - Each iteration of a benchmark runs `GOMAXPROCS` requests in parallel.
 - Current results on my machine:
 
 
 | Metric/API                              | /basic-handler      | /object-pool-handler   | /bounded-pool-handler |
 | --------------------------------------- | ------------------- | --------------------- | --------------------- |
 | Running time per operation (ns)         | 3,387               | 1,728                 | 1,892                 | 
 | Number of bytes allocated per operation | 11,858              | 4,950                 | 4,944                 |
 | Allocations per operation               | 14                  | 12                    | 12                    |
 

 ## Instructions
 - To run the benchmark there are two main options:
     1. Compile a binary from `main.go` and run it. This will run an HTTP server as described above that listens on port 8080. Then run your benchmarks on the server manually.
     2. Run the existing benchmarks in `server_test.go`:
         `go test -bench=. ./server`
