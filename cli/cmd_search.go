package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

func newSearchCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search posts by title, summary, or tags",
		Long: `search filters posts from the RSS feed whose title, summary, or tags
contain the query string (case-insensitive).

Examples:
  lw search "attention"
  lw search "reinforcement learning" -n 5
  lw search diffusion -f json`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")
			posts, err := app.client.Search(cmd.Context(), query, app.limit)
			if err != nil {
				return err
			}
			if len(posts) == 0 {
				return noResults("no posts matched " + query)
			}
			return app.render(posts)
		},
	}
}
