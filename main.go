package main

import (
	"os"

	"github.com/goftp/file-driver"
	"github.com/goftp/server"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/goftp/leveldb-perm"
	"github.com/goftp/leveldb-auth"
	"github.com/Mr-Dai/go-ftpd/log"
	"gopkg.in/urfave/cli.v1"
)

const version = "0.1.0"

func serverAction(c *cli.Context) {
	// Setup auth DB
	authdb := c.GlobalString("authdb")
	db, err := leveldb.OpenFile(authdb, nil)
	if err != nil {
		log.Fatalf("Failed to setup auth DB on `%s`: %v", authdb, err)
	}

	auth := &ldbauth.LDBAuth{db}
	perm := ldbperm.NewLDBPerm(db, "root", "root", os.ModePerm)

	// Setup data directory
	datapath := c.String("datapath")
	_, err = os.Lstat(datapath)
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

func addServerCommand(app *cli.App) {
	serverCommand := cli.Command{}
	serverCommand.Name = "server"
	serverCommand.Aliases = []string{"s"}
	serverCommand.Usage = "Run FTP server"
	serverCommand.Flags = []cli.Flag {
		cli.StringFlag{Name: "name, n", Value: "my-ftp", Usage: "name of the FTP server"},
		cli.StringFlag{Name: "datapath, d", Value: "/tmp/go-ftpd/data",
			Usage: "data directory for the FTP server"},
		cli.IntFlag{Name: "port, p", Value: 21, Usage: "port to listen to"},
	}
	serverCommand.Action = serverAction

	app.Commands = append(app.Commands, serverCommand)
}

func prepareCLI() (app *cli.App) {
	app = cli.NewApp()
	app.HelpName = "go-ftpd"
	app.Version = version
	app.Usage = "A Simple FTP Server"
	app.Author = "Robert Peng"
	app.Email = "robert.peng@foxmail.com"

	// Setup common flag
	app.Flags = []cli.Flag {
		cli.StringFlag{Name: "authdb, a", Value: "/tmp/go-ftpd/auth.db",
			Usage: "path for auth DB of the FTP server"},
	}

	// Setup `server` command
	addServerCommand(app)
	return
}

func main() {
	app := prepareCLI()
	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
}
