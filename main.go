package main

import (
	"os"
	"github.com/urfave/cli"
	"nafue/utility"
	"log"
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
			Action: func(c *cli.Context) {
				url := c.Args().First()
				if url == "" {
					log.Println("You must enter a url")
					os.Exit(1)
				}
				// get file
				utility.GetFile(url)
			},
		},
		{
			Name:      "share",
			Usage:     "share [file]",
			Action: func(c *cli.Context) {
				file := c.Args().First()
				if file == "" {
					log.Println("You must enter a file")
					os.Exit(1)
				}
				// share file
				utility.PutFile(file)
			},
		},

	};
	app.Action = func(c *cli.Context) {
		println("Please run with a sub-command. For more information try \"nafue help\"")
	}

	app.Run(os.Args)
}

//salt:  [823823226, 793384892]
//salt b64:  MRqLei9KF7w=
//iv:  [-1726315756, -750854000, -575834770, -1005289789, -875997976, 1972462144, 348864951, -1912938221]
//iv b64:  mRqDFNM+4JDdrXVuxBR+w8vJVOh1kWJAFMtBt4364RM=