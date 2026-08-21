// cmd/server/main.go
package main

import (
	"log/slog"
	"os"

	"github.com/DylanBergmann2502/go-maleficent/cmd/server/commands"
)

func main() {
	if err := commands.NewRootCommand().Execute(); err != nil {
		slog.Error("Failed to execute command", "error", err)
		os.Exit(1)
	}
}
