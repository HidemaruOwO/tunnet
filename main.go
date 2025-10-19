package main

import (
	"context"
	"log"

	"github.com/charmbracelet/fang"

	"github.com/HidemaruOwO/bridge/cmd"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	root := cmd.NewRootCmd()
	if err := fang.Execute(
		context.Background(),
		root,
		fang.WithVersion(version),
		fang.WithCommit(commit),
	); err != nil {
		log.Fatal(err)
	}
}
