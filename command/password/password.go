package password

import (
	"fmt"

	"github.com/benmorehouse/std/hashicorp"
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
		c, err := hashicorp.DefaultVaultClient()
		if err != nil {
			return fmt.Errorf("hashicorp_password_client_fail: %s", err.Error())
		}
		return process(c, utils.DefaultInteractor(), viper.GetBool(putCLIKey))
	},
}

const putCLIKey = "put"

func process(vaultClient repo.Repo, user utils.Interactor, userWillPutPassword bool) (err error) {
	if userWillPutPassword {
		return putPassword(vaultClient, user)
	}
	return getPassword(vaultClient, user)
}

func putPassword(vaultClient repo.Repo, user utils.Interactor) error {
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

	if err := vaultClient.Put(key, value); err != nil {
		return fmt.Errorf("unable_to_put_vault_secret: %s", err.Error())
	}

	return nil
}

func getPassword(vaultClient repo.Repo, user utils.Interactor) error {
	fmt.Print("Name of password:")
	key := user.Input()
	password := vaultClient.Get(key)
	for password == "" {
		fmt.Println("Password not found. Please try again.")
		key = user.Input()
		password = vaultClient.Get(key)
	}
	fmt.Println(password)
	return nil
}

func init() {
	Command.PersistentFlags().Bool(putCLIKey, false, "I want to put a password in")

	viper.BindPFlag(putCLIKey, Command.PersistentFlags().Lookup(putCLIKey))
}
