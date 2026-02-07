package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================
// EXAMPLE 1: Basic Goroutines (Original)
// ============================================

func numbers() {
	for i := 1; i <= 5; i++ {
		time.Sleep(250 * time.Millisecond)
		fmt.Printf("%d ", i)
	}
}

func alphabets() {
	for i := 'a'; i <= 'e'; i++ {
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("%c ", i)
	}
}

func basicExample() {
	fmt.Println("=== Basic Goroutines (with sleep) ===")
	go numbers()
	go alphabets()
	time.Sleep(1700 * time.Millisecond)
	fmt.Println("\nmain terminated")
}

// ============================================
// EXAMPLE 2: Using WaitGroups (RECOMMENDED)
// ============================================

func numbersWithWG(wg *sync.WaitGroup) {
	defer wg.Done() // Tell WaitGroup we're done when function exits
	for i := 1; i <= 5; i++ {
		time.Sleep(250 * time.Millisecond)
		fmt.Printf("%d ", i)
	}
}

func alphabetsWithWG(wg *sync.WaitGroup) {
	defer wg.Done() // Tell WaitGroup we're done when function exits
	for i := 'a'; i <= 'e'; i++ {
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("%c ", i)
	}
}

func waitGroupExample() {
	fmt.Println("\n\n=== Using WaitGroups (Better!) ===")
	var wg sync.WaitGroup

	// Tell WaitGroup we're launching 2 goroutines
	wg.Add(2)

	go numbersWithWG(&wg)
	go alphabetsWithWG(&wg)

	// Wait for both goroutines to finish
	wg.Wait()
	fmt.Println("\nmain terminated (all goroutines finished)")
}

// ============================================
// EXAMPLE 3: Using Channels
// ============================================

func numbersWithChannel(done chan bool) {
	for i := 1; i <= 5; i++ {
		time.Sleep(250 * time.Millisecond)
		fmt.Printf("%d ", i)
	}
	done <- true // Send signal that we're done
}

func alphabetsWithChannel(done chan bool) {
	for i := 'a'; i <= 'e'; i++ {
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("%c ", i)
	}
	done <- true // Send signal that we're done
}

func channelExample() {
	fmt.Println("\n\n=== Using Channels ===")

	// Create channels for each goroutine
	done1 := make(chan bool)
	done2 := make(chan bool)

	go numbersWithChannel(done1)
	go alphabetsWithChannel(done2)

	// Wait for both to finish
	<-done1 // Block until numbers is done
	<-done2 // Block until alphabets is done

	fmt.Println("\nmain terminated (received all done signals)")
}

// ============================================
// EXAMPLE 4: Practical - Concurrent Downloads
// ============================================

func downloadFile(filename string, size int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("⬇️  Downloading %s (%dMB)...\n", filename, size)

	// Simulate download time based on file size
	time.Sleep(time.Duration(size*100) * time.Millisecond)

	fmt.Printf("✅ Completed %s\n", filename)
}

func downloadExample() {
	fmt.Println("\n\n=== Practical Example: Concurrent Downloads ===")

	files := []struct {
		name string
		size int
	}{
		{"video.mp4", 10},
		{"document.pdf", 5},
		{"image.jpg", 3},
		{"archive.zip", 8},
	}

	start := time.Now()
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go downloadFile(file.name, file.size, &wg)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("\n🎉 All downloads completed in %.2f seconds!\n", elapsed.Seconds())
	fmt.Printf("💡 Sequential would take: %.2f seconds\n", (10+5+3+8)*0.1)
}

// ============================================
// EXAMPLE 5: Worker Pool Pattern
// ============================================

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("🔧 Worker %d processing job %d\n", id, job)
		time.Sleep(500 * time.Millisecond) // Simulate work
		results <- job * 2                 // Send result
	}
}

func workerPoolExample() {
	fmt.Println("\n\n=== Worker Pool Pattern ===")

	const numWorkers = 3
	const numJobs = 9

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Start workers
	fmt.Printf("Starting %d workers...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	fmt.Printf("Sending %d jobs...\n", numJobs)
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Println("\nResults:")
	for result := range results {
		fmt.Printf("✓ Result: %d\n", result)
	}
}

// ============================================
// EXAMPLE 6: Race Condition (Problem & Solution)
// ============================================

var unsafeCounter int // Shared variable - UNSAFE!
var safeCounter int   // Shared variable - SAFE with mutex
var mutex sync.Mutex  // Mutex to protect safeCounter

func incrementUnsafe(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		unsafeCounter++ // Race condition! Multiple goroutines modifying
	}
}

func incrementSafe(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		mutex.Lock()   // Lock before accessing shared data
		safeCounter++  // Now it's safe!
		mutex.Unlock() // Unlock after we're done
	}
}

func raceConditionExample() {
	fmt.Println("\n\n=== Race Condition Demo ===")

	// Unsafe version
	unsafeCounter = 0
	var wg1 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg1.Add(1)
		go incrementUnsafe(&wg1)
	}
	wg1.Wait()
	fmt.Printf("❌ Unsafe counter (10 goroutines × 1000): %d (expected 10000)\n", unsafeCounter)

	// Safe version
	safeCounter = 0
	var wg2 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go incrementSafe(&wg2)
	}
	wg2.Wait()
	fmt.Printf("✅ Safe counter (with mutex): %d (correct!)\n", safeCounter)
}

// ============================================
// MAIN FUNCTION
// ============================================

func main() {
	fmt.Println("🚀 Go Concurrency Examples\n")
	fmt.Println("=" + string(make([]byte, 50)) + "=")

	// Run all examples
	basicExample()         // Original example with time.Sleep
	waitGroupExample()     // Better approach with WaitGroup
	channelExample()       // Using channels for synchronization
	downloadExample()      // Practical: concurrent downloads
	workerPoolExample()    // Advanced: worker pool pattern
	raceConditionExample() // Important: race conditions and how to fix

	fmt.Println("\n🎓 Key Lessons:")
	fmt.Println("1. Use 'go' keyword to launch goroutines")
	fmt.Println("2. Always wait for goroutines (WaitGroup or channels)")
	fmt.Println("3. Protect shared data with mutexes")
	fmt.Println("4. Channels are great for communication between goroutines")
	fmt.Println("5. Worker pools are efficient for processing many tasks")
}
