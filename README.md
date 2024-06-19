# SHAllenge-Go

This program is a multi-threaded hash miner that searches for the lowest numerical hash with a given username and random nonce. It uses SHA-256 hashing and aims to find the input that generates the lowest hash value. The program also tracks and prints the number of iterations per second.

## Features

- Multi-threaded processing
- Dynamic nonce length
- Leading zero count in the hash
- Notification on finding a new lowest hash
- Performance metrics (iterations per second)

## Requirements

- Go 1.15 or higher
- `github.com/gen2brain/beeep` package for notifications

## Installation

1. Install the Go programming language from [golang.org](https://golang.org/).

2. Set up your Go workspace and get the required package:

   ```sh
   go get github.com/gen2brain/beeep
   ```

3. Clone the repository or copy the source code into a `.go` file.

## Usage

1. Compile the program:

   ```sh
   go build -o shallenge_go main.go
   ```

2. Run the program with your desired username:
   ```sh
   ./shallenge_go <username>
   ```

### Example

```sh
./shallenge_go bbbenji
```

This will start the hash cracking process with the username `bbbenji`.

## Configuration

- `numWorkers`: Adjust the number of workers based on your CPU capacity. Default is 1.
- `batchSize`: Number of hash calculations per batch for each worker. Default is 1000.
- `minNonceLength`: Minimum length of the nonce. Default is 1.
- `maxNonceLength`: Adjusted dynamically based on the username length.

## How It Works

1. **Initialization**: The program starts by checking the command-line arguments for the username. It calculates the maximum nonce length based on the username length to ensure the combined length does not exceed 64 characters.

2. **Workers**: The program spawns multiple worker goroutines to generate random nonces and calculate the SHA-256 hash of the concatenated username and nonce.

3. **Hash Comparison**: Each worker compares the newly generated hash with the current lowest hash. If the new hash is lower, it updates the lowest hash and prints the details.

4. **Notification**: Upon finding a new lowest hash, the program sends a desktop notification using the `beeep` package.

5. **Performance Tracking**: A separate goroutine tracks and prints the number of iterations per second.

## Notifications

The program uses the `beeep` package to send desktop notifications when a new lowest hash is found. Ensure you have a compatible notification system for `beeep` to work (e.g., macOS Notification Center, Windows Notification Center, or Linux with libnotify).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## Acknowledgments

- `github.com/gen2brain/beeep` for the notification package.
- The Go programming language community for their support and documentation.
