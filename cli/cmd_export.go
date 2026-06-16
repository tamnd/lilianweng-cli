package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func newExportCmd(app *App) *cobra.Command {
	var outFile string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export all posts as JSONL",
		Long: `export fetches the full post archive and writes one JSON record per line.

Examples:
  lw export > lilianweng.jsonl
  lw export --out lilianweng.jsonl
  lw export -f csv > lilianweng.csv`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			posts, err := app.client.Latest(cmd.Context(), 0)
			if err != nil {
				return codeError(1, err)
			}
			if outFile != "" {
				f, err := os.Create(outFile)
				if err != nil {
					return codeError(1, err)
				}
				defer f.Close()
				r := NewRendererTo(f, app)
				return r.Render(posts)
			}
			return app.render(posts)
		},
	}
	cmd.Flags().StringVar(&outFile, "out", "", "write output to FILE instead of stdout")
	return cmd
}
