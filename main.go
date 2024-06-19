package main

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "github.com/gen2brain/beeep"
    mathRand "math/rand"
    "os"
    "runtime"
    "sync"
    "sync/atomic"
    "time"
)

var (
    username       string    // Username set from command-line argument
    minNonceLength = 1
    maxNonceLength int       // Adjusted dynamically based on username length
    numWorkers     = 1       // Adjust based on your CPU
    batchSize      = 1000
    lowestHash     []byte    // Stores the current lowest numerical hash
    lowestInput    string    // Stores the input that generated the lowest hash
    maxZeros       int64     // Stores the maximum number of leading zeros found so far
    iteration      int64     // Tracks the total number of iterations
    startTime      time.Time // Stores the start time of the search
    mu             sync.Mutex // Protects lowestHash and lowestInput
)

func countLeadingHexZeros(hash []byte) int {
    zeros := 0
    for _, b := range hash {
        if b == 0 {
            zeros += 2
        } else {
            if b&0xF0 == 0 {
                zeros++
            }
            break
        }
    }
    return zeros
}

func worker(ctx context.Context, id int) {
    nonceBuf := make([]byte, maxNonceLength)
    mathRand.Seed(time.Now().UnixNano() + int64(id))

    for {
        select {
        case <-ctx.Done():
            return
        default:
            for i := 0; i < batchSize; i++ {
                nonceLength := mathRand.Intn(maxNonceLength-minNonceLength+1) + minNonceLength
                _, err := rand.Read(nonceBuf[:nonceLength])
                if err != nil {
                    fmt.Println("Error generating nonce:", err)
                    continue
                }
                nonce := hex.EncodeToString(nonceBuf[:nonceLength])

                // Ensure the combined length of username/nonce is <= 64
                if (len(username) + 1 + len(nonce)) > 64 {
                    nonce = nonce[:64-len(username)-1]
                }

                input := username + "/" + nonce
                hash := sha256.Sum256([]byte(input))

                leadingZeros := countLeadingHexZeros(hash[:])
                compareHash(leadingZeros, hash[:], input)
            }
        }
    }
}

func formatHashWithSpaces(hash string) string {
    var formattedHash string
    for i, c := range hash {
        if i > 0 && i%8 == 0 {
            formattedHash += " "
        }
        formattedHash += string(c)
    }
    return formattedHash
}

func compareHash(leadingZeros int, hash []byte, input string) {
    mu.Lock()
    defer mu.Unlock()

    if lowestHash == nil || hex.EncodeToString(hash) < hex.EncodeToString(lowestHash) {
        lowestHash = make([]byte, len(hash))
        copy(lowestHash, hash)
        lowestInput = input

        formattedHash := formatHashWithSpaces(hex.EncodeToString(hash))
        elapsedTime := time.Since(startTime).Seconds()
        iterationsPerSecond := float64(atomic.LoadInt64(&iteration)) / elapsedTime
        fmt.Printf("\n\nIteration: %d\nNew lowest hash: %s\nInput: %s\nLeading zeros: %d\nIterations per second: %.2f\n\n",
            atomic.LoadInt64(&iteration), formattedHash, input, leadingZeros, iterationsPerSecond)

        // Send notification using beeep
        go func() {
            err := beeep.Alert("New Lowest Hash Found", fmt.Sprintf("Hash: %s\nInput: %s\nLeading zeros: %d",
                formattedHash, input, leadingZeros), "")
            if err != nil {
                fmt.Println("Error sending notification:", err)
            }
        }()
    }
    atomic.AddInt64(&iteration, 1)
}

func printIterationsPerSecond(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            currentIterations := atomic.LoadInt64(&iteration)
            elapsedTime := time.Since(startTime).Seconds()
            iterationsPerSecond := float64(currentIterations) / elapsedTime
            fmt.Printf("\rIterations: %d | Iterations per second: %.2f", currentIterations, iterationsPerSecond)
        }
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: program <username>")
        os.Exit(1)
    }
    username = os.Args[1]
    maxNonceLength = 64 - len(username) - 1

    runtime.GOMAXPROCS(runtime.NumCPU())
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    startTime = time.Now()

    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(ctx, id)
        }(i)
    }

    go printIterationsPerSecond(ctx)

    wg.Wait()
    cancel() // Ensure the print goroutine exits
}
