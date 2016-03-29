package main

import (
	"os"
	"github.com/codegangsta/cli"
	"nafue/utility"
)
// todo add delete temp function
func main() {
	// setup env as needed
	utility.Init()

	app := cli.NewApp()
	app.Name = "Nafue"
	app.Usage = "Anonymous, secure file transfers that self destruct after first use or 24 hours using client side encryption."
	app.Commands = []cli.Command{
		{
			Name:      "get",
			Usage:     "get [file]",
			ArgsUsage: "blah blah",
			Action: func(c *cli.Context) {
				url := c.Args().First()
				utility.GetFile(url)
			},
		},
		{
			Name:      "share",
			Usage:     "share [file]",
			Action: func(c *cli.Context) {

			},
		},

	};
	app.Action = func(c *cli.Context) {
		println("Please run with a sub-command. For more information try \"nafue help\"")
	}

	app.Run(os.Args)
}
