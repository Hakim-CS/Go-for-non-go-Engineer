# 🎯 Concurrency Exercises - Learn by Doing!

## 📝 Exercise 1: Fix the Original Code

**Problem:** The current code might not print the letter 'e'.

**Your Task:** Modify `main.go` to use a WaitGroup instead of `time.Sleep`.

**Starter Code:**
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func numbers(wg *sync.WaitGroup) {
    // TODO: Add defer wg.Done()
    for i := 1; i <= 5; i++ {
        time.Sleep(250 * time.Millisecond)
        fmt.Printf("%d ", i)
    }
}

func alphabets(wg *sync.WaitGroup) {
    // TODO: Add defer wg.Done()
    for i := 'a'; i <= 'e'; i++ {
        time.Sleep(400 * time.Millisecond)
        fmt.Printf("%c ", i)
    }
}

func main() {
    // TODO: Create WaitGroup
    // TODO: Add 2 to WaitGroup
    // TODO: Launch goroutines with &wg
    // TODO: Wait for completion
    fmt.Println("\nmain terminated")
}
```

<details>
<summary>💡 Solution</summary>

```go
func numbers(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 1; i <= 5; i++ {
        time.Sleep(250 * time.Millisecond)
        fmt.Printf("%d ", i)
    }
}

func alphabets(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 'a'; i <= 'e'; i++ {
        time.Sleep(400 * time.Millisecond)
        fmt.Printf("%c ", i)
    }
}

func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    go numbers(&wg)
    go alphabets(&wg)
    wg.Wait()
    fmt.Println("\nmain terminated")
}
```
</details>

---

## 📝 Exercise 2: Concurrent Sum Calculator

**Task:** Calculate the sum of numbers 1-1000 using 4 goroutines (each handling 250 numbers).

**Starter Code:**
```go
package main

import (
    "fmt"
    "sync"
)

func sumRange(start, end int, result chan int) {
    sum := 0
    for i := start; i <= end; i++ {
        sum += i
    }
    // TODO: Send sum to result channel
}

func main() {
    // TODO: Create result channel
    // TODO: Launch 4 goroutines:
    //   - sumRange(1, 250, result)
    //   - sumRange(251, 500, result)
    //   - sumRange(501, 750, result)
    //   - sumRange(751, 1000, result)
    
    // TODO: Receive 4 results and add them up
    // TODO: Print total (should be 500500)
}
```

<details>
<summary>💡 Solution</summary>

```go
func sumRange(start, end int, result chan int) {
    sum := 0
    for i := start; i <= end; i++ {
        sum += i
    }
    result <- sum
}

func main() {
    result := make(chan int, 4)
    
    go sumRange(1, 250, result)
    go sumRange(251, 500, result)
    go sumRange(501, 750, result)
    go sumRange(751, 1000, result)
    
    total := 0
    for i := 0; i < 4; i++ {
        total += <-result
    }
    
    fmt.Printf("Sum of 1-1000: %d\n", total)
}
```
</details>

---

## 📝 Exercise 3: Race Condition Bug Hunt

**Task:** Find and fix the race condition in this code.

**Buggy Code:**
```go
package main

import (
    "fmt"
    "sync"
)

var balance int = 1000

func withdraw(amount int, wg *sync.WaitGroup) {
    defer wg.Done()
    if balance >= amount {
        balance -= amount
        fmt.Printf("Withdrew %d, balance: %d\n", amount, balance)
    }
}

func main() {
    var wg sync.WaitGroup
    
    // 10 people try to withdraw $100 each
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go withdraw(100, &wg)
    }
    
    wg.Wait()
    fmt.Printf("Final balance: %d (should be 0)\n", balance)
}
```

**What's wrong?** Run it multiple times - balance might be negative!

<details>
<summary>💡 Solution</summary>

```go
var balance int = 1000
var mutex sync.Mutex  // Add mutex!

func withdraw(amount int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    mutex.Lock()  // Lock before accessing shared data
    if balance >= amount {
        balance -= amount
        fmt.Printf("Withdrew %d, balance: %d\n", amount, balance)
    } else {
        fmt.Printf("Insufficient funds for %d\n", amount)
    }
    mutex.Unlock()  // Unlock after we're done
}
```
</details>

---

## 📝 Exercise 4: Timeout Handler

**Task:** Create a function that times out if it takes longer than 2 seconds.

**Starter Code:**
```go
package main

import (
    "fmt"
    "time"
)

func slowOperation(result chan string) {
    time.Sleep(3 * time.Second)  // Simulates slow operation
    result <- "Operation complete!"
}

func main() {
    result := make(chan string, 1)
    go slowOperation(result)
    
    // TODO: Use select to wait for result OR timeout after 2 seconds
    // If timeout: print "Operation timed out!"
    // If success: print the result
}
```

<details>
<summary>💡 Solution</summary>

```go
func main() {
    result := make(chan string, 1)
    go slowOperation(result)
    
    select {
    case msg := <-result:
        fmt.Println("Success:", msg)
    case <-time.After(2 * time.Second):
        fmt.Println("Operation timed out!")
    }
}
```
</details>

---

## 📝 Exercise 5: Worker Pool

**Task:** Process 20 jobs using 3 workers.

**Starter Code:**
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    // TODO: Loop through jobs channel
    // TODO: Process each job (multiply by 2)
    // TODO: Sleep 500ms to simulate work
    // TODO: Send result to results channel
}

func main() {
    const numJobs = 20
    const numWorkers = 3
    
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    var wg sync.WaitGroup
    
    // TODO: Start 3 workers
    
    // TODO: Send 20 jobs (numbers 1-20)
    // TODO: Close jobs channel
    
    // TODO: Wait for all workers to finish
    // TODO: Close results channel
    
    // TODO: Print all results
}
```

<details>
<summary>💡 Solution</summary>

```go
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(500 * time.Millisecond)
        results <- job * 2
    }
}

func main() {
    const numJobs = 20
    const numWorkers = 3
    
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    var wg sync.WaitGroup
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }
    
    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Close results when done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    for result := range results {
        fmt.Println("Result:", result)
    }
}
```
</details>

---

## 📝 Exercise 6: Ping Pong

**Task:** Two goroutines passing a ball back and forth.

**Requirements:**
- Goroutine 1 receives from `ping`, prints "ping", sends to `pong`
- Goroutine 2 receives from `pong`, prints "pong", sends to `ping`
- Do this 5 times

**Starter Code:**
```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ping := make(chan int)
    pong := make(chan int)
    
    // TODO: Goroutine 1 - receives from ping, sends to pong
    
    // TODO: Goroutine 2 - receives from pong, sends to ping
    
    // TODO: Start the game by sending 1 to ping
    
    time.Sleep(2 * time.Second)
}
```

<details>
<summary>💡 Solution</summary>

```go
func main() {
    ping := make(chan int)
    pong := make(chan int)
    
    go func() {
        for i := 0; i < 5; i++ {
            ball := <-ping
            fmt.Println("Ping!", ball)
            time.Sleep(300 * time.Millisecond)
            pong <- ball + 1
        }
        close(pong)
    }()
    
    go func() {
        for ball := range pong {
            fmt.Println("Pong!", ball)
            time.Sleep(300 * time.Millisecond)
            ping <- ball + 1
        }
    }()
    
    ping <- 1
    time.Sleep(4 * time.Second)
}
```
</details>

---

## 📝 Exercise 7: URL Fetcher (Practical)

**Task:** Fetch multiple URLs concurrently and measure total time.

**Starter Code:**
```go
package main

import (
    "fmt"
    "time"
)

func fetchURL(url string, result chan string) {
    // Simulate HTTP request
    time.Sleep(1 * time.Second)
    result <- fmt.Sprintf("Fetched: %s", url)
}

func main() {
    urls := []string{
        "https://google.com",
        "https://github.com",
        "https://stackoverflow.com",
        "https://reddit.com",
    }
    
    start := time.Now()
    
    // TODO: Fetch all URLs concurrently
    // TODO: Collect and print all results
    
    elapsed := time.Since(start)
    fmt.Printf("\nTotal time: %.2f seconds\n", elapsed.Seconds())
    fmt.Printf("Sequential would take: %d seconds\n", len(urls))
}
```

<details>
<summary>💡 Solution</summary>

```go
func main() {
    urls := []string{
        "https://google.com",
        "https://github.com",
        "https://stackoverflow.com",
        "https://reddit.com",
    }
    
    start := time.Now()
    result := make(chan string, len(urls))
    
    // Fetch all concurrently
    for _, url := range urls {
        go fetchURL(url, result)
    }
    
    // Collect results
    for i := 0; i < len(urls); i++ {
        fmt.Println(<-result)
    }
    
    elapsed := time.Since(start)
    fmt.Printf("\nTotal time: %.2f seconds\n", elapsed.Seconds())
    fmt.Printf("Sequential would take: %d seconds\n", len(urls))
}
```
</details>

---

## 🎓 Bonus Challenge: Build a Simple Cache

**Task:** Create a thread-safe cache with concurrent reads/writes.

**Requirements:**
- `Set(key, value)` - store value
- `Get(key)` - retrieve value
- Multiple goroutines can access simultaneously
- No race conditions!

<details>
<summary>💡 Solution</summary>

```go
package main

import (
    "fmt"
    "sync"
)

type Cache struct {
    data map[string]string
    mu   sync.RWMutex  // RWMutex allows multiple readers
}

func NewCache() *Cache {
    return &Cache{
        data: make(map[string]string),
    }
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
    fmt.Printf("Set: %s = %s\n", key, value)
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()  // Read lock (multiple readers OK)
    defer c.mu.RUnlock()
    value, exists := c.data[key]
    return value, exists
}

func main() {
    cache := NewCache()
    var wg sync.WaitGroup
    
    // 5 writers
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", n)
            cache.Set(key, fmt.Sprintf("value%d", n))
        }(i)
    }
    
    wg.Wait()
    
    // 10 readers
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", n%5+1)
            if value, ok := cache.Get(key); ok {
                fmt.Printf("Read: %s = %s\n", key, value)
            }
        }(i)
    }
    
    wg.Wait()
}
```
</details>

---

## 🎯 Progress Checklist

- [ ] Exercise 1: Fixed original code with WaitGroup
- [ ] Exercise 2: Built concurrent sum calculator
- [ ] Exercise 3: Found and fixed race condition
- [ ] Exercise 4: Implemented timeout handler
- [ ] Exercise 5: Created worker pool
- [ ] Exercise 6: Built ping-pong game
- [ ] Exercise 7: Fetched URLs concurrently
- [ ] Bonus: Built thread-safe cache

## 🚀 Next Steps

1. Run `go run -race main.go` on your solutions to check for race conditions
2. Benchmark sequential vs concurrent versions
3. Try building your own projects using these patterns
4. Study the `context` package for cancellation
5. Learn about buffered vs unbuffered channels

**Keep practicing! Concurrency is a skill that improves with experience.** 🎓
