package main

import (
	"os"

	"github.com/cwaits6/apk-datasource/cmd/apk-index/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
