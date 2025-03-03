package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Layer-Edge/light-node/config"
	"github.com/Layer-Edge/light-node/node"
)

func Worker(ctx context.Context, wg *sync.WaitGroup, id int, verifier *node.Verifier) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d is shutting down\n", id)
			return
		default:
			fmt.Printf("Worker %d is running...\n", id)
			err := verifier.CollectSampleAndVerify()
			if err != nil {
				fmt.Printf("Worker %d encountered error: %v\n", id, err)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	// Load configuration
	cfg := config.GetConfig()

	// Create verifier with configuration
	verifier := node.NewVerifier(&cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM)

	wg.Add(1)
	go Worker(ctx, &wg, 1, verifier)

	<-signalChan
	fmt.Println("\nReceived interrupt signal. Shutting down gracefully...")

	cancel()

	wg.Wait()
	fmt.Println("Worker has shut down. Exiting..")
}