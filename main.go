package main

import (
	"os"

	"github.com/goftp/file-driver"
	"github.com/goftp/server"
	"github.com/Mr-Dai/go-ftpd/log"
	"gopkg.in/urfave/cli.v1"
)

const version = "0.1.0"

func serverAction(c *cli.Context) {
	// Setup auth DB
	auth, perm := prepareAuth(c)

	// Setup data directory
	datapath := c.String("datapath")
	_, err := os.Lstat(datapath)
	if os.IsNotExist(err) {
		os.MkdirAll(datapath, os.ModePerm)
	} else if err != nil {
		log.Fatalf("Failed to access data directory `%s`: %v", datapath, err)
	}
	factory := &filedriver.FileDriverFactory{
		RootPath: datapath,
		Perm:     perm,
	}

	opt := &server.ServerOpts{
		Name:    c.String("name"),
		Factory: factory,
		Port:    c.Int("port"),
		Auth:    auth,
		Logger:  log.FtpdLogger(),
	}

	// Start FTP server
	ftpServer := server.NewServer(opt)
	err = ftpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start FTP server: %v", err)
	}
}

func addRunCommand(app *cli.App) {
	runCommand := cli.Command{}
	runCommand.Name = "run"
	runCommand.Aliases = []string{"s"}
	runCommand.Usage = "Runs FTP server"
	runCommand.Flags = []cli.Flag {
		cli.StringFlag{Name: "name, n", Value: "my-ftp", Usage: "name of the FTP server"},
		cli.StringFlag{Name: "datapath, d", Value: "/tmp/go-ftpd/data",
			Usage: "data directory for the FTP server"},
		cli.IntFlag{Name: "port, p", Value: 21, Usage: "port to listen to"},
	}
	runCommand.Action = serverAction

	app.Commands = append(app.Commands, runCommand)
}

func prepareCLI() (app *cli.App) {
	app = cli.NewApp()
	app.HelpName = "go-ftpd"
	app.Version = version
	app.Usage = "A Simple FTP Server"
	app.Author = "Robert Peng"
	app.Email = "robert.peng@foxmail.com"

	// Setup common flags
	app.Flags = []cli.Flag {
		cli.StringFlag{Name: "authdb, a", Value: "/tmp/go-ftpd/auth.db",
			Usage: "path for auth DB of the FTP server"},
	}

	// Setup commands
	addRunCommand(app)
	addUserCommand(app)
	addGroupCommand(app)
	addPermCommand(app)
	return
}

func main() {
	app := prepareCLI()
	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
}
