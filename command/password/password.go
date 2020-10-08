package password

import (
	"fmt"

	"github.com/benmorehouse/std/hashicorp"
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
		return process(c, utils.StdInteractor())
	},
}

func process(vaultClient hashicorp.VaultClient, user utils.Interactor) (err error) {
	if viper.GetBool("put") {
		return putPassword(vaultClient, user)
	}
	return getPassword(vaultClient, user)
}

func putPassword(vaultClient hashicorp.VaultClient, user utils.Interactor) error {
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

func getPassword(vaultClient hashicorp.VaultClient, user utils.Interactor) error {
	fmt.Print("Name of password:")
	key := user.Input()
	for key == "" {
		fmt.Println("You must give the password a default key")
		key = user.Input()
	}
	password, err := vaultClient.Get(key)
	if err != nil {
		return fmt.Errorf("vault_get_err: %s", err)
	}
	fmt.Println(password)
	return nil
}

func init() {
	Command.PersistentFlags().Bool("put", false, "I want to put a password in")

	viper.BindPFlag("put", Command.PersistentFlags().Lookup("put"))
}
