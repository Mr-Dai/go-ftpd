package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Mr-Dai/go-ftpd/log"
	"github.com/goftp/ftpd/web"
	"github.com/goftp/leveldb-auth"
	"github.com/goftp/leveldb-perm"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/urfave/cli.v1"
)

const CategoryUser = "User/Group Management"

func prepareAuth(c *cli.Context) (auth *ldbauth.LDBAuth, perm *ldbperm.LDBPerm) {
	authdb := c.GlobalString("authdb")
	db, err := leveldb.OpenFile(authdb, nil)
	if err != nil {
		log.Fatalf("Failed to setup auth DB on `%s`: %v", authdb, err)
	}

	auth = &ldbauth.LDBAuth{db}
	perm = ldbperm.NewLDBPerm(db, "root", "root", os.ModePerm)
	return
}

func listUserAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	users := make([]web.User, 0)
	if err := auth.UserList(&users); err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}
	groups := make([]string, 0)
	if err := auth.GroupList(&groups); err != nil {
		log.Fatalf("Failed to list groups: %v", err)
	}

	userGroups := make(map[string][]string)
	for _, user := range users {
		userGroups[user.Name] = make([]string, 0)
	}
	groupUsers := make([]string, 0)
	for _, group := range groups {
		groupUsers = groupUsers[:0]
		if err := auth.GroupUser(group, &groupUsers); err != nil {
			log.Fatalf("Failed to list users of group `%s`: %v", group, err)
		}
		for _, user := range groupUsers {
			userGroups[user] = append(userGroups[user], group)
		}
	}

	fmt.Println("USER\t\tGROUPS")
	for user, groups := range userGroups {
		fmt.Printf("%s\t\t%s\n", user, strings.Join(groups, ", "))
	}
}

func addUserAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	user := c.Args().Get(0)
	if strings.TrimSpace(user) == "" {
		log.Fatalf("Username cannot be empty.")
	}

	pass := c.Args().Get(1)

	if err := auth.AddUser(user, pass); err != nil {
		log.Fatalf("Failed to add new user `%s`: %v", user, err)
	}
	log.Infof("Successfully created new user `%s`.", user)
}

func delUserAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	user := c.Args().Get(0)
	if strings.TrimSpace(user) == "" {
		log.Fatalf("Username cannot be empty.")
	}

	if err := auth.DelUser(user); err != nil {
		log.Fatalf("Failed to delete user `%s`: %v", user, err)
	}
	log.Infof("Successfully deleted user `%s`.", user)
}

func passwdAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	user := c.Args().Get(0)
	if strings.TrimSpace(user) == "" {
		log.Fatalf("Username cannot be empty.")
	}
	pass := c.Args().Get(1)

	if err := auth.ChgPass(user, pass); err != nil {
		log.Fatalf("Failed to change password of user `%s`: %v", user, err)
	}
	log.Infof("Successfully changed password of user `%s`.", user)
}

func addUserCommand(app *cli.App) {
	// user
	user := cli.Command{}
	user.Name = "user"
	user.Usage = "Manages users of the FTP server"
	user.Category = CategoryUser

	// user list
	listUser := cli.Command{}
	listUser.Name = "list"
	listUser.Usage = "Lists users of the FTP server"
	listUser.Action = listUserAction

	// user add
	addUser := cli.Command{}
	addUser.Name = "add"
	addUser.Usage = "Adds user to the FTP server"
	addUser.ArgsUsage = "<username> [passwd]"
	addUser.Action = addUserAction

	// user del
	delUser := cli.Command{}
	delUser.Name = "del"
	delUser.Usage = "Deletes user from the FTP server"
	delUser.ArgsUsage = "<username>"
	delUser.Action = delUserAction

	// user passwd
	passwd := cli.Command{}
	passwd.Name = "passwd"
	passwd.Usage = "Updates user's password"
	passwd.ArgsUsage = "<username> [<passwd>]"
	passwd.Action = passwdAction

	user.Subcommands = []cli.Command{listUser, addUser, delUser, passwd}
	app.Commands = append(app.Commands, user)
}
