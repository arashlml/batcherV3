# Batcher V3 — Generic Concurrent Batcher in Go

BatcherV3 is a high-performance, generic, and thread-safe batch processor written in Go. It collects incoming items and flushes them in batches based on size or time. Flushed batches are delivered asynchronously through an output channel.

---

## Features

* Generic implementation using Go generics
* Size-based batch flushing
* Time-based batch flushing
* Fully concurrent and thread-safe
* Non-blocking Add operation
* Output via read-only channel
* Graceful shutdown with final flush
* Internal atomic counters for add and flush operations
* Designed for high-throughput workloads

---

## Installation

Using `go get`:

```bash
go get github.com/arashlml/batcherV3
```

Or manually import the package:

```go
import "your_project/batcherV3"
```

---

## API Overview

### Constructor

```go
b := NewBatcher[T](maxSize, interval,maxChanSize)
```

| Parameter   | Type          | Description                                   |
| ----------- | ------------- | --------------------------------------------- |
| maxSize     | int           | Maximum number of items per batch             |
| interval    | time.Duration | Maximum time before an automatic flush occurs |
| maxChanSize | int           | Maximum input and output channel size         |
---

### Add an Item

```go
err := b.Add(item)
```

* Non-blocking operation
* Returns an error if the internal buffer is full

---

### Read Output Batches

```go
for batch := range b.OutChan() {
    // process batch
}
```

---

### Close the Batcher

```go
b.Close()
```

* Flushes remaining items
* Shuts down internal goroutines safely
* Closes output channel

---

## Example Usage

```go
package main

import (
	"fmt"
	"time"

	"your_project/batcherV3"
)

func main() {
	b := batcherV3.NewBatcher[int](3, 2*time.Second,10)

	go func() {
		for batch := range b.OutChan() {
			fmt.Println("Received batch:", batch)
		}
	}()

	b.Add(1)
	b.Add(2)
	b.Add(3)

	time.Sleep(time.Second)

	b.Add(4)
	b.Add(5)

	time.Sleep(3 * time.Second)

	b.Close()
}
```

---

## Flush Behavior

| Condition              | Result               |
| ---------------------- | -------------------- |
| Batch reaches maxSize  | Immediate flush      |
| Timer reaches interval | Automatic flush      |
| Close is called        | Final flush and stop |

---

## Concurrency Model

* A single background goroutine manages batching
* Add is performed through a buffered channel
* Safe for concurrent producers
* Designed to avoid race conditions and blocking

---

## Internal Counters

The batcher tracks:

* addCounter: total number of added items
* flushCounter: total number of flush operations

These are used internally for monitoring and debugging.

---

## Use Cases

* Database bulk inserts
* Log aggregation
* Message queue batching
* Metrics collection
* Event streaming
* ETL pipelines
* High-throughput API backends

---

## Testing

The project includes unit tests that verify:

* Flush by size
* Flush by time
* Flush on Close
* Multiple consecutive batches
* Output data correctness

Run tests with:

```bash
go test ./... -v
```

---

## Package Structure

```text
batcherV3/
├── batcher.go
├── batcher_test.go
└── README.md
```

---

## Roadmap

* Context-aware Add operations
* Retry support on flush failures
* Dynamic batch resizing
* Metrics integration

---

## License

MIT License

---

## Author

Developed by Arash
