package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"

	totuicli "github.com/giorgtarkha/totui/tui/cli"
)

func main() {
	app := &cli.App{
		Name: "totui",
		Commands: []*cli.Command{
			{
				Name:  "cli",
				Flags: []cli.Flag{},
				Action: func(c *cli.Context) error {

					p := tea.NewProgram(&totuicli.CLITUI{}, tea.WithAltScreen())

					if _, err := p.Run(); err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Failed to run totui: %s", err.Error())
		os.Exit(1)
	}
}
