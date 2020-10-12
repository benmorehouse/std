package password

import (
	"fmt"

	"github.com/benmorehouse/std/repo"
	"github.com/benmorehouse/std/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command is the exported command for creating a list
var Command = &cobra.Command{
	Use:     "password",
	Short:   "Manage personal passwords",
	Example: "./std password",
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool(putCLIKey) && viper.GetBool(removeCLIKey) {
			return fmt.Errorf("please select either remove or put")
		}

		if viper.GetBool(putCLIKey) {
			return process(repo.PasswordConnector(), utils.DefaultInteractor(), putCLIKey)
		}

		if viper.GetBool(removeCLIKey) {
			return process(repo.PasswordConnector(), utils.DefaultInteractor(), removeCLIKey)
		}

		return process(repo.PasswordConnector(), utils.DefaultInteractor(), getDefaultKey)
	},
}

const (
	putCLIKey     = "put"
	removeCLIKey  = "remove"
	getDefaultKey = "get"
)

func process(connector repo.Connector, user utils.Interactor, userAction string) (err error) {
	db, err := connector.Connect()
	if err != nil {
		return fmt.Errorf("unable_to_connect: %s", err.Error())
	}

	switch userAction {
	case putCLIKey:
		return putPassword(db, user)
	case removeCLIKey:
		return removePassword(db, user)
	default:
		return getPassword(db, user)
	}
}

func removePassword(db repo.Repo, user utils.Interactor) error {
	fmt.Print("Name of password:")
	key := user.Input()
	for db.Get(key) == "" {
		fmt.Println("Password doesn't exist")
		fmt.Print("Name of password:")
		key = user.Input()
	}

	if err := db.Remove(key); err != nil {
		return fmt.Errorf("unable_to_put_vault_secret: %s", err.Error())
	}

	return nil
}

func putPassword(db repo.Repo, user utils.Interactor) error {
	fmt.Print("Name of password:")
	key := user.Input()
	for key == "" {
		fmt.Println("You must give the password a default key")
		key = user.Input()
	}

	fmt.Print("Value of password:")
	value := user.Input()
	for value == "" {
		fmt.Println("You must give the password a default value")
		value = user.Input()
	}

	if err := db.Put(key, value); err != nil {
		return fmt.Errorf("unable_to_put_vault_secret: %s", err.Error())
	}

	return nil
}

func getPassword(db repo.Repo, user utils.Interactor) error {
	fmt.Print("Name of password:")
	key := user.Input()
	password := db.Get(key)
	for password == "" {
		fmt.Println("Password not found. Please try again.")
		key = user.Input()
		password = db.Get(key)
	}
	fmt.Println(password)
	return nil
}

func init() {
	Command.PersistentFlags().Bool(putCLIKey, false, "I want to put a password in")
	Command.PersistentFlags().Bool(removeCLIKey, false, "I want to delete a password")

	viper.BindPFlag(putCLIKey, Command.PersistentFlags().Lookup(putCLIKey))
	viper.BindPFlag(removeCLIKey, Command.PersistentFlags().Lookup(removeCLIKey))
}
