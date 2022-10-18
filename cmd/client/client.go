package main

import (
	"ddia/src/client"
	"fmt"
	"log"
	"os"
)

func main() {
	err := runClient()
	if err != nil {
		log.Fatal(err)
	}
}

func runClient() error {
	logger := log.New(os.Stdout, "[client] ", 0)
	c := client.NewClient(logger, "localhost:6379")

	set, err := c.Set("hello2", "world")
	if err != nil {
		return err
	}
	logger.Printf("set: %q\n", string(set))

	for i := 0; i < 10; i++ {
		c.Set(fmt.Sprintf("hello %d", i), "world")
	}

	get, err := c.Get("hello")
	if err != nil {
		return err
	}
	logger.Printf("get: %q\n", string(get))

	return nil
}
