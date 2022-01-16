package main

import (
	"fmt"

	"mpourismaiel.dev/guts/store"
)

func main() {
	fmt.Println("Starting project...")
	_ = store.New("guts")
}
