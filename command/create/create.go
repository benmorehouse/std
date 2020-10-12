package create

import (
	"fmt"
	"log"
	"strings"

	"github.com/benmorehouse/std/repo"
	"github.com/benmorehouse/std/utils"
	"github.com/spf13/cobra"
)

// Command is the exported command for creating a list
var Command = &cobra.Command{
	Use:     "create",
	Short:   "create a list",
	Example: "./std create work",
	RunE: func(cmd *cobra.Command, args []string) error {
		return process(repo.ListConnector(), utils.DefaultInteractor(), args)
	},
}

func process(connector repo.Connector, user utils.Interactor, args []string) error {
	db, err := connector.Connect()
	if err != nil {
		log.Println("Error connecting")
		return err
	}

	var newBucketName string
	if len(args) == 0 {
		utils.DisplayBucketList(db)
		newBucketName = user.Input()
	} else {
		newBucketName = strings.Join(args, " ")
	}

	for db.Get(newBucketName) != "" {
		fmt.Printf("%s already exists.\n", newBucketName)
		utils.DisplayBucketList(db)
		newBucketName = user.Input()
	}

	if err := user.RunLifeCycle(db, newBucketName, user, true); err != nil {
		return err
	}
	return nil
}
