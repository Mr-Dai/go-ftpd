package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/Mr-Dai/go-ftpd/log"
	"gopkg.in/urfave/cli.v1"
)

const CategoryPerm = "Permission Management"

func chownAction(c *cli.Context) {
	_, perm := prepareAuth(c)

	ownerAndGroup := c.Args().Get(0)
	if strings.TrimSpace(ownerAndGroup) == "" {
		log.Fatalf("Owner cannot be empty.")
	}
	path := c.Args().Get(1)
	if strings.TrimSpace(path) == "" {
		log.Fatalf("Path cannot be empty.")
	}

	parts := strings.Split(ownerAndGroup, ":")
	owner := parts[0]
	group := ""
	if len(parts) > 1 {
		group = parts[1]
	}

	if err := perm.ChOwner(path, owner); err != nil {
		log.Fatalf("Failed to change owner of `%s` to `%s`: %v", path, owner, err)
	}

	if group != "" {
		if err := perm.ChGroup(path, group); err != nil {
			log.Fatalf("Failed to change owner group of `%s` to `%s`: %v", path, group, err)
		}
	}
}

func chmodAction(c *cli.Context) {
	_, perm := prepareAuth(c)

	modeExpr := c.Args().Get(0)
	if strings.TrimSpace(modeExpr) == "" {
		log.Fatalf("Mode cannot be empty.")
	}

	path := c.Args().Get(1)
	if strings.TrimSpace(path) == "" {
		log.Fatalf("Path cannot be empty.")
	}

	mode, err := strconv.ParseUint(modeExpr, 8, 32)
	if err != nil {
		log.Fatalf("`%s` is not a valid file mode.", modeExpr)
	}

	if err := perm.ChMode(path, os.FileMode(uint32(mode))); err != nil {
		log.Fatalf("Failed to set mode of `%s` to `%s`: %v", path, modeExpr, err)
	}
}

func addPermCommand(app *cli.App) {
	// chown
	chown := cli.Command{}
	chown.Name = "chown"
	chown.Usage = "Changes file owner and group"
	chown.Category = CategoryPerm
	chown.ArgsUsage = "owner:[group] file"
	chown.Action = chownAction

	// chmod
	chmod := cli.Command{}
	chmod.Name = "chmod"
	chmod.Usage = "Changes file mode bits"
	chmod.Category = CategoryPerm
	chmod.ArgsUsage = "mode file"
	chmod.Action = chmodAction

	app.Commands = append(app.Commands, chown, chmod)
}
