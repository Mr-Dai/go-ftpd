package main

import (
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/Mr-Dai/go-ftpd/log"
	"github.com/goftp/file-driver"
	"github.com/goftp/server"
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
		Name:         c.String("name"),
		Factory:      factory,
		Port:         c.Int("port"),
		Auth:         auth,
		Logger:       log.FtpdLogger(),
		PassivePorts: c.String("passive-ports"),
	}

	// Setup CPU profiling if needed
	if cpuprofile := c.String("cpuprofile"); cpuprofile != "" {
		f, err := os.OpenFile(cpuprofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatalf("Failed to open CPU profile file `%s`: %v", cpuprofile, err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("Failed not start CPU profiling: %v", err)
		}
		defer func() {
			pprof.StopCPUProfile()
			log.Infof("CPU profile was successfully written to `%s`.", cpuprofile)
		}()
	}

	// Start FTP server
	ftpServer := server.NewServer(opt)
	go func() {
		err := ftpServer.ListenAndServe()
		if err != nil && err != server.ErrServerClosed {
			log.Fatalf("Failed to start FTP server: %v", err)
		}
	}()

	// Wait for shutdown signal
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	Loop:
	for {
		select {
		case sig := <-sigs:
			log.Infof("Received `%s`. Warm shutdown...", sig)

			// Warm shutdown
			stoppedChan := make(chan struct{})
			go func() {
				ftpServer.Shutdown()
				close(stoppedChan)
			}()

			select {
			case <-sigs: // Cold shutdown
				log.Infof("Received `%s`. Cold shutdown!!!", sig)
			case <-stoppedChan: // WebSocket closed
			}
			break Loop
		}
	}

	// Write out memory profile
	if memprofile := c.String("memprofile"); memprofile != "" {
		f, err := os.OpenFile(memprofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatalf("Failed to open memory profile file `%s`: %v", memprofile, err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("Failed to write memory profile: %v", err)
		}
		log.Infof("Memory profile was successfully written to `%s`.", memprofile)
	}
}

func addRunCommand(app *cli.App) {
	runCommand := cli.Command{}
	runCommand.Name = "run"
	runCommand.Aliases = []string{"s"}
	runCommand.Usage = "Runs FTP server"
	runCommand.Flags = []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "my-ftp", Usage: "name of the FTP server"},
		cli.StringFlag{Name: "datapath, d", Value: "/tmp/go-ftpd/data",
			Usage: "data directory for the FTP server"},
		cli.IntFlag{Name: "port, p", Value: 21, Usage: "port to listen to"},
		cli.StringFlag{Name: "passive-ports", Value: "30000-31000", Usage: "range for passive ports"},
		cli.StringFlag{Name: "cpuprofile", Value: "", Usage: "write CPU profile to `file`"},
		cli.StringFlag{Name: "memprofile", Value: "", Usage: "write memory profile to `file`"},
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
	app.Flags = []cli.Flag{
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
