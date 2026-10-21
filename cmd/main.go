package main

import (
	"context"
	"fmt"
	"gahenna/internal/importer"
	"log"
	"os"

	"golang.org/x/sync/semaphore"
)

func processLine(line []string) {
	fmt.Println(line)
}

func readWithSemaphore() {
	sem := semaphore.NewWeighted(4)
	ctx := context.Background()


	sem.Acquire(ctx, 4)
}

func readWithoutSemaphore() {
}

func main() {
	args := os.Args[1:]
    file, err := os.Open(args[0])
    if err != nil {
        log.Fatalf("Failed to open CSV file: %v", err)
    }
    defer file.Close()

	lineChan := make(chan []string, 1)

	go importer.ImportCSV(file, lineChan)

	for r := range lineChan {
		processLine(r)
	}
}
