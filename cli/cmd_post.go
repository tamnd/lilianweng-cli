package cli

import (
	"github.com/spf13/cobra"
)

func newPostCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "post <slug>",
		Short: "Show details for a single post",
		Long: `post fetches a single blog post by its slug, partial path, or full URL.

Examples:
  lw post hallucination
  lw post 2024-07-07-hallucination
  lw post https://lilianweng.github.io/posts/2024-07-07-hallucination/`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			detail, err := app.client.PostBySlug(cmd.Context(), args[0])
			if err != nil {
				return codeError(3, err)
			}
			return app.render(detail)
		},
	}
}
