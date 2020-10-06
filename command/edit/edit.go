package edit

import (
	"github.com/benmorehouse/std/repo"
	"github.com/benmorehouse/std/utils"
	"github.com/spf13/cobra"
)

// Command will provide key editing capabilities
var Command = &cobra.Command{
	Use:   "",
	Short: "Open the current list",
	RunE: func(cmd *cobra.Command, args []string) error {
		return process(repo.DefaultConnector(), utils.StdInteractor(), args)
	},
}

func process(connector repo.Connector, user utils.Interactor, args []string) error {
	db, err := connector.Connect()
	if err != nil {
		return err
	}

	var bucketName string
	if len(args) != 1 {
		utils.DisplayBucketList(db)
		bucketName = user.Input()
	} else {
		bucketName = args[0]
	}

	if bucketName == "backlog" {
		// then go ahead and run the backlog command
		return nil
	}

	if err := utils.RunLifeCycle(db, bucketName, user); err != nil {
		return err
	}

	return nil
}
