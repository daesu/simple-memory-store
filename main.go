package main

import (
	"simple-memory-store/cmd"
)

func main() {
	client := cmd.NewClient()
	client.Run()
}
