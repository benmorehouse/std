package delete

import (
	"fmt"
	"strings"

	"github.com/benmorehouse/std/repo"
	"github.com/benmorehouse/std/utils"
	"github.com/spf13/cobra"
)

// Command exports key delete functionality
var Command = &cobra.Command{
	Use:     "delete",
	Short:   "Delete the list from the database",
	Example: "./std delete work",
	RunE: func(cmd *cobra.Command, args []string) error {
		return process(repo.ListConnector(), utils.DefaultInteractor(), args)
	},
}

func process(connector repo.Connector, user utils.Interactor, args []string) error {
	db, err := connector.Connect()
	if err != nil {
		return err
	}

	var bucketName string
	if len(args) == 0 {
		utils.DisplayBucketList(db)
		bucketName = user.Input()
	} else {
		bucketName = strings.Join(args, " ")
	}

	for db.Get(bucketName) == "" {
		fmt.Printf("%s doesnt exist.\n", bucketName)
		utils.DisplayBucketList(db)
		bucketName = user.Input()
	}

	if err := db.Remove(bucketName); err != nil {
		return err
	}

	return nil
}
