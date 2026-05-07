package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Kivayan/azctx/go/internal/cli"
	"github.com/Kivayan/azctx/go/internal/version"
)

func main() {
	root := cli.NewRootCommand(version.Version)
	if err := root.Execute(); err != nil {
		var exitErr cli.ExitError
		switch {
		case errors.As(err, &exitErr):
			if exitErr.Err != nil {
				fmt.Fprintln(os.Stderr, exitErr.Err)
			}
			os.Exit(exitErr.Code)
		case errors.Is(err, cli.ErrVersionRequested):
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}