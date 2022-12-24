package cli

import (
	"gazur/pkg/common"
	"gazur/pkg/server"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Start() {
	app := &cli.App{
		Name:  "gazur",
		Usage: "Utility web server which queries azure api for some information",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "tenant-id",
				Required: true,
				EnvVars:  []string{"AZURE_TENANT_ID"},
			},
			&cli.StringFlag{
				Name:     "client-id",
				Required: true,
				EnvVars:  []string{"AZURE_CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:     "client-secret",
				Required: true,
				EnvVars:  []string{"AZURE_CLIENT_SECRET"},
			},
			&cli.PathFlag{
				Name:     "cfg-file",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			identity := common.NewIdentity(
				c.String("tenant-id"),
				c.String("client-id"),
				c.String("client-secret"),
			)

			path := c.Path("cfg-file")

			gazur := server.New(&identity, &path)
			return gazur.Run()
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
