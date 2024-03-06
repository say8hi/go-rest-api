package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
)

func main() {
	stringPtr := flag.String("s", "", "String to hash")
	flag.Parse()

	if *stringPtr == "" {
		fmt.Println("Specify the string using flag -s")
		return
	}

  hash := sha256.Sum256([]byte(*stringPtr))
	hashString := hex.EncodeToString(hash[:])

	fmt.Println(string(hashString))
}
