# 🎯 Go Concurrency Quick Reference

## 🔑 Core Concepts

### Goroutine
```go
go myFunction()  // Launches function concurrently
```

### Channel
```go
ch := make(chan int)     // Unbuffered channel
ch := make(chan int, 10) // Buffered channel (capacity 10)
ch <- 42                 // Send value to channel
value := <-ch            // Receive value from channel
close(ch)                // Close channel (no more sends)
```

### WaitGroup
```go
var wg sync.WaitGroup
wg.Add(1)     // Increment counter
wg.Done()     // Decrement counter
wg.Wait()     // Block until counter is 0
```

### Mutex
```go
var mu sync.Mutex
mu.Lock()     // Acquire lock
mu.Unlock()   // Release lock
```

## 📋 Common Patterns

### Pattern 1: Launch and Wait
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // Do work
}()
wg.Wait()
```

### Pattern 2: Channel Communication
```go
ch := make(chan string)
go func() {
    ch <- "result"
}()
result := <-ch
```

### Pattern 3: Select Statement
```go
select {
case msg := <-ch1:
    fmt.Println("Received:", msg)
case msg := <-ch2:
    fmt.Println("Received:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("Timeout")
}
```

### Pattern 4: Worker Pool
```go
jobs := make(chan int, 100)
results := make(chan int, 100)

// Start workers
for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}

// Send jobs
for j := 1; j <= 10; j++ {
    jobs <- j
}
close(jobs)
```

### Pattern 5: Fan-Out, Fan-In
```go
// Fan-out: Multiple goroutines reading from same channel
for i := 0; i < 10; i++ {
    go worker(jobs, results)
}

// Fan-in: Multiple goroutines writing to same channel
for i := 0; i < 10; i++ {
    go producer(output)
}
```

## ⚠️ Common Mistakes

### ❌ Forgetting to Wait
```go
// WRONG
func main() {
    go doSomething()
    // main exits, goroutine killed!
}
```

### ✅ Proper Wait
```go
// CORRECT
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        doSomething()
    }()
    wg.Wait()
}
```

### ❌ Race Condition
```go
// WRONG
var counter int
for i := 0; i < 10; i++ {
    go func() { counter++ }()
}
```

### ✅ Protected Access
```go
// CORRECT
var counter int
var mu sync.Mutex
for i := 0; i < 10; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

### ❌ Channel Deadlock
```go
// WRONG
ch := make(chan int)
ch <- 42  // Blocks forever! No receiver
```

### ✅ Send in Goroutine
```go
// CORRECT
ch := make(chan int)
go func() { ch <- 42 }()
value := <-ch
```

## 📊 When to Use What

| Use Case | Tool | Example |
|----------|------|---------|
| Wait for goroutines | WaitGroup | Parallel downloads |
| Pass data | Channel | Producer-consumer |
| Timeout | select + time.After | API with timeout |
| Protect shared data | Mutex | Counter, cache |
| Process tasks | Worker pool | Image processing |
| Cancel operation | Context | Long-running tasks |

## 🎯 Performance Tips

1. **Goroutines are cheap** - Don't worry about creating thousands
2. **Buffered channels** - Reduce blocking with `make(chan T, size)`
3. **Worker pools** - Limit concurrent operations
4. **Context** - Use for cancellation and timeouts
5. **Avoid sharing** - "Share memory by communicating"

## 🔍 Debugging Tools

```bash
# Check for race conditions
go run -race main.go

# Profile goroutines
import _ "net/http/pprof"
http.ListenAndServe("localhost:6060", nil)
# Visit: http://localhost:6060/debug/pprof/goroutine
```

## 📚 Remember

- **CSP Model**: "Don't communicate by sharing memory; share memory by communicating"
- **Goroutines are managed by Go runtime**, not OS
- **Channels are typed** - `chan int` only passes ints
- **Closed channels** - Reading returns zero value + false
- **nil channels** - Reading/writing blocks forever
