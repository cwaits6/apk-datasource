package main

import (
	"os"

	"github.com/cwaits6/apk-datasource/cmd/apk-datasource/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
