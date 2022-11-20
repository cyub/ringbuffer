# ringbuffer

![GitHub](https://img.shields.io/github/license/cyub/ringbuffer)

Lock-free ring buffer in Go, support SPSC/SPMC/MPSC/MPMC implementations.

- SPSC - Single Producer/Single Consumer
- SPMC - Single Producer/Multi-Consumer
- MPSC - Multi-Producer/Single Consumer
- MPMC - Multi-Producer/Multi-Consumer

## Features

- Lock-free operations - they succeed or fail immediately without blocking or waiting.
- Thread-safe direct access to the internal ring buffer memory.
- Support SPSC/SPMC/MPSC/MPMC implementations. You can choose the best performing implementation based on your business scenario


## Benchmark

Machine information for benchmarks:

Apple M1 Pro 8 core

### MPMC ringbuffer vs channel

```bash
go test  -benchmem -run=^$ -bench="^BenchmarkRingMPMC|BenchmarkChanMPMC$" . github.com/cyub/ringbuffer
```

```
goos: darwin
goarch: arm64
pkg: github.com/cyub/ringbuffer
BenchmarkRingMPMC/100P100C-8         	 2634507	       423.5 ns/op	       7 B/op	       0 allocs/op
BenchmarkRingMPMC/4P4C_1CPU-8        	40298432	        27.80 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanMPMC/100P100C-8         	10558228	       104.3 ns/op	       7 B/op	       0 allocs/op
BenchmarkChanMPMC/4P4C_1CPU-8        	31277354	        35.78 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	github.com/cyub/ringbuffer	5.557s
```

### MPSC ringbuffer vs channel

```bash
go test  -benchmem -run=^$ -bench="^BenchmarkRingMPSC|BenchmarkChanMPSC$" . github.com/cyub/ringbuffer
```

```
goos: darwin
goarch: arm64
pkg: github.com/cyub/ringbuffer
BenchmarkRingMPSC/100P1C-8         	 3915708	       296.5 ns/op	       8 B/op	       1 allocs/op
BenchmarkRingMPSC/4P1C_1CPU-8      	40444438	        27.86 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanMPSC/100P1C-8         	11623476	       124.5 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanMPSC/4P1C_1CPU-8      	31618956	        35.88 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	github.com/cyub/ringbuffer	5.544s
```

### SPMC ringbuffer vs channel

```bash
go test  -benchmem -run=^$ -bench="^BenchmarkRingSPMC|BenchmarkChanSPMC$" . github.com/cyub/ringbuffer
```

```
goos: darwin
goarch: arm64
pkg: github.com/cyub/ringbuffer
BenchmarkRingSPMC/1P100C-8         	 4392514	       285.5 ns/op	       7 B/op	       0 allocs/op
BenchmarkRingSPMC/1P4C_1CPU-8      	54208296	        20.46 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanSPMC/1P100C-8         	 5216336	       394.5 ns/op	       7 B/op	       0 allocs/op
BenchmarkChanSPMC/1P4C_1CPU-8      	31566313	        35.75 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	github.com/cyub/ringbuffer	6.256s
```

### SPSC ringbuffer vs channel

```bash
go test  -benchmem -run=^$ -bench="^BenchmarkRingSPSC|BenchmarkChanSPSC$" . github.com/cyub/ringbuffer
```

```
goos: darwin
goarch: arm64
pkg: github.com/cyub/ringbuffer
BenchmarkRingSPSC/1P1C-8         	16629625	        68.52 ns/op	       7 B/op	       0 allocs/op
BenchmarkRingSPSC/1P1C_1CPU-8    	44770410	        24.89 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanSPSC/1P1C-8         	 8335009	       286.3 ns/op	       7 B/op	       0 allocs/op
BenchmarkChanSPSC/1P1C_1CPU-8    	31815259	        36.12 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	github.com/cyub/ringbuffer	6.204s
```