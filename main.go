package main

import (
	"github.com/datarootsio/cheek/cmd"
)

//go:generate npm run build
func main() {
	cmd.Execute()
}
