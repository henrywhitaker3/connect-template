package token

import (
	"fmt"
	"time"

	"github.com/henrywhitaker3/connect-template/internal/app"
	"github.com/spf13/cobra"
)

var expiry time.Duration

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Generate a new auth token",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := app.Jwt.New(expiry)
			if err != nil {
				return err
			}

			fmt.Println(token)

			return nil
		},
	}

	cmd.Flags().DurationVarP(&expiry, "expiry", "x", time.Hour*24*90, "How long the token is valid for")

	return cmd
}
