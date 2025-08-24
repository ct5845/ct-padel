package main

import (
	"ct-padel-s/src"
	_ "ct-padel-s/src/infrastructure/logging"
)

//go:generate npm run build-css
//go:generate go run ./cmd/copyassets .
//go:generate go run ./cmd/migrate .

func main() {
	src.App()
}
