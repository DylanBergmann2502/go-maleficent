// cmd/server/commands/root.go
package commands

import "github.com/spf13/cobra"

// NewRootCommand creates the server CLI command tree.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "server",
		Short: "Go Maleficent - API Server and Management CLI",
		Long:  "A comprehensive Go API server with database migrations, background workers, and more.",
	}

	root.AddCommand(NewAPICommand())
	root.AddCommand(NewMigrateCommand())

	return root
}
