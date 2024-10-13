package root

import (
	"github.com/henrywhitaker3/connect-template/cmd/migrate"
	"github.com/henrywhitaker3/connect-template/cmd/routes"
	"github.com/henrywhitaker3/connect-template/cmd/serve"
	"github.com/henrywhitaker3/connect-template/internal/app"
	"github.com/spf13/cobra"
)

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "api",
		Short:   "Golang template API",
		Version: app.Version,
	}

	cmd.AddCommand(serve.New(app))
	cmd.AddCommand(migrate.New(app))
	cmd.AddCommand(routes.New(app))

	cmd.PersistentFlags().StringP("config", "c", "connect-template.yaml", "The path to the api config file")

	return cmd
}
