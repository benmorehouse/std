package edit

import (
	"strings"

	"github.com/benmorehouse/std/repo"
	"github.com/benmorehouse/std/utils"
	"github.com/spf13/cobra"
)

// Command will provide key editing capabilities
var Command = &cobra.Command{
	Use:   "open ",
	Short: "Open the current list",
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

	if bucketName == "backlog" {
		return nil
	}

	if err := user.RunLifeCycle(db, bucketName, user, false); err != nil {
		return err
	}

	return nil
}
