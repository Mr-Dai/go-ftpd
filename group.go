package main

import (
	"fmt"
	"strings"

	"github.com/Mr-Dai/go-ftpd/log"
	"gopkg.in/urfave/cli.v1"
)

func listGroupAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	groups := make([]string, 0)
	if err := auth.GroupList(&groups); err != nil {
		log.Fatalf("Failed to list user groups: %v", err)
	}

	groupUsers := make(map[string][]string)
	for _, group := range groups {
		users := make([]string, 0)
		if err := auth.GroupUser(group, &users); err != nil {
			log.Fatalf("Failed to list users of group `%s`: %v", group, err)
		}

		groupUsers[group] = users
	}

	fmt.Println("GROUP\t\tUSERS")
	for group, users := range groupUsers {
		fmt.Printf("%s\t\t%s\n", group, strings.Join(users, ", "))
	}
}

func addGroupAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	group := c.Args().Get(0)
	if strings.TrimSpace(group) == "" {
		log.Fatalf("Group name cannot be empty.")
	}

	if err := auth.AddGroup(group); err != nil {
		log.Fatalf("Failed to create user group `%s`: %v", group, err)
	}
	log.Infof("Successfully created user group `%s`.", group)
}

func delGroupAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	group := c.Args().Get(0)
	if strings.TrimSpace(group) == "" {
		log.Fatalf("Group name cannot be empty.")
	}

	if err := auth.DelGroup(group); err != nil {
		log.Fatalf("Failed to delete user group `%s`: %v", group, err)
	}
	log.Infof("Successfully deleted user group `%s`.", group)
}

func addGroupUserAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	group := c.Args().Get(0)
	if strings.TrimSpace(group) == "" {
		log.Fatalf("Group name cannot be empty.")
	}
	user := c.Args().Get(1)
	if strings.TrimSpace(user) == "" {
		log.Fatalf("User name cannot be empty.")
	}

	if err := auth.AddUserGroup(user, group); err != nil {
		log.Fatalf("Failed to add user `%s` to group `%s`: %v", user, group, err)
	}
	log.Infof("Successfully added user `%s` to group `%s`.", user, group)
}

func delGroupUserAction(c *cli.Context) {
	auth, _ := prepareAuth(c)

	group := c.Args().Get(0)
	if strings.TrimSpace(group) == "" {
		log.Fatalf("Group name cannot be empty.")
	}
	user := c.Args().Get(1)
	if strings.TrimSpace(user) == "" {
		log.Fatalf("User name cannot be empty.")
	}

	if err := auth.DelUserGroup(user, group); err != nil {
		log.Fatalf("Failed to delete user `%s` from group `%s`: %v", user, group, err)
	}
	log.Infof("Successfully deleted user `%s` from group `%s`.", user, group)
}

func addGroupCommand(app *cli.App) {
	// group
	group := cli.Command{}
	group.Name = "group"
	group.Usage = "Manages user groups of the FTP server"
	group.Category = CategoryUser

	// group list
	listGroup := cli.Command{}
	listGroup.Name = "list"
	listGroup.Usage = "Lists user groups of the FTP server"
	listGroup.Action = listGroupAction

	// group add
	addGroup := cli.Command{}
	addGroup.Name = "add"
	addGroup.Usage = "Adds user group to the FTP server"
	addGroup.ArgsUsage = "group"
	addGroup.Action = addGroupAction

	// group del
	delGroup := cli.Command{}
	delGroup.Name = "del"
	delGroup.Usage = "Deletes user group from the FTP server"
	delGroup.ArgsUsage = "group"
	delGroup.Action = delGroupAction

	// group adduser
	addGroupUser := cli.Command{}
	addGroupUser.Name = "adduser"
	addGroupUser.Usage = "Adds user to user group"
	addGroupUser.ArgsUsage = "group user"
	addGroupUser.Action = addGroupUserAction

	// group deluser
	delGroupUser := cli.Command{}
	delGroupUser.Name = "deluser"
	delGroupUser.Usage = "Deletes user from user group"
	delGroupUser.ArgsUsage = "group user"
	delGroupUser.Action = delGroupUserAction

	group.Subcommands = []cli.Command{listGroup, addGroup, delGroup, addGroupUser, delGroupUser}
	app.Commands = append(app.Commands, group)
}
