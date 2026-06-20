package cli

import (
	"github.com/spf13/cobra"
)

func newInfoCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show blog statistics",
		Long: `info prints aggregate statistics about the blog.

Examples:
  lw info
  lw info -f json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			info, err := app.client.Stats(cmd.Context())
			if err != nil {
				return codeError(1, err)
			}
			return app.render(info)
		},
	}
}
