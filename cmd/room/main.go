package main

import (
	"log"

	"github.com/thockin/go-build-template/pkg/version"
)

func main() {
	log.Printf("version: %s\n", version.Version)
}
