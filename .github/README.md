# introduction
According to [plaintext](https://github.com/TechEmpower/FrameworkBenchmarks/wiki/Project-Information-Framework-Tests-Overview#plaintext) to output http response。

> note: restart app after every benchmark.

## net
```bash
wrk -t6 -c256 -d10s http://127.0.0.1:3000
Running 10s test @ http://127.0.0.1:3000
  6 threads and 256 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.45ms  572.93us  10.50ms   76.39%
    Req/Sec    11.90k     1.82k   26.73k    77.45%
  714415 requests in 10.10s, 97.43MB read
  Socket errors: connect 5, read 94, write 0, timeout 0
Requests/sec:  70712.52
Transfer/sec:      9.64MB
```

## gnet
```bash
wrk -t6 -c256 -d10s http://127.0.0.1:3001
Running 10s test @ http://127.0.0.1:3001
  6 threads and 256 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.33ms  322.72us  11.63ms   88.50%
    Req/Sec    12.37k   736.67    14.00k    80.00%
  738705 requests in 10.00s, 101.45MB read
  Socket errors: connect 5, read 103, write 0, timeout 0
Requests/sec:  73851.67
Transfer/sec:     10.14MB
```

## gev
```bash
wrk -t6 -c256 -d10s http://127.0.0.1:3002
Running 10s test @ http://127.0.0.1:3002
  6 threads and 256 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.38ms  328.59us  14.46ms   89.43%
    Req/Sec    12.21k   627.96    13.30k    71.12%
  736337 requests in 10.10s, 100.42MB read
  Socket errors: connect 5, read 103, write 0, timeout 0
Requests/sec:  72888.42
Transfer/sec:      9.94MB
```
