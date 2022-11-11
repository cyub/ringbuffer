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


