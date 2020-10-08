package list

import (
	"github.com/benmorehouse/std/repo"
	"github.com/benmorehouse/std/utils"
	"github.com/spf13/cobra"

	"log"
)

// Command is the exported command for creating a list
var Command = &cobra.Command{
	Use:     "list",
	Short:   "show all lists entered",
	Example: "./std list",
	Run: func(cmd *cobra.Command, args []string) {
		process(repo.DefaultConnector(), args)
	},
}

func process(connector repo.Connector, args []string) (err error) {
	db, err := connector.Connect()
	if err != nil {
		log.Println("Error opening database at createlist command:", err)
		return
	}

	utils.DisplayBucketList(db)
	connector.Disconnect()
	return
}
