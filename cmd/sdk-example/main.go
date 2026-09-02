package main

import (
	"context"
	"log"
	"os"

	"github.com/goark-projects/goark-observe-sdk-example/internal/example"
)

func main() {
	logger := log.New(os.Stderr, "observe-sdk-example: ", 0)
	if err := example.Run(context.Background(), os.Stdout); err != nil {
		logger.Fatal(err)
	}
}
