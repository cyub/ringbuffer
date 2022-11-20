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
BenchmarkRingMPMC/100P100C-8         	 9562984	       126.6 ns/op	       7 B/op	       0 allocs/op
BenchmarkRingMPMC/4P4C_1CPU-8        	40235656	        27.85 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanMPMC/100P100C-8         	10317433	       108.6 ns/op	       7 B/op	       0 allocs/op
BenchmarkChanMPMC/4P4C_1CPU-8        	31602544	        35.97 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	github.com/cyub/ringbuffer	7.001s
```

### MPSC ringbuffer vs channel

```bash
go test  -benchmem -run=^$ -bench="^BenchmarkRingMPSC|BenchmarkChanMPSC$" . github.com/cyub/ringbuffer
```

```
goos: darwin
goarch: arm64
pkg: github.com/cyub/ringbuffer
BenchmarkRingMPSC/100P1C-8         	10631840	       111.9 ns/op	       7 B/op	       0 allocs/op
BenchmarkRingMPSC/4P1C_1CPU-8      	40517098	        27.78 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanMPSC/100P1C-8         	11523838	       103.5 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanMPSC/4P1C_1CPU-8      	31493272	        35.93 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	github.com/cyub/ringbuffer	5.060s
```

### SPMC ringbuffer vs channel

```bash
go test  -benchmem -run=^$ -bench="^BenchmarkRingSPMC|BenchmarkChanSPMC$" . github.com/cyub/ringbuffer
```

```
goos: darwin
goarch: arm64
pkg: github.com/cyub/ringbuffer
BenchmarkRingSPMC/1P100C-8         	11944942	        99.32 ns/op	       7 B/op	       0 allocs/op
BenchmarkRingSPMC/1P4C_1CPU-8      	53785629	        20.31 ns/op	       8 B/op	       1 allocs/op
BenchmarkChanSPMC/1P100C-8         	 5077400	       405.6 ns/op	       7 B/op	       0 allocs/op
BenchmarkChanSPMC/1P4C_1CPU-8      	31406179	        35.81 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	github.com/cyub/ringbuffer	5.993s
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