package password

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"

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
		if viper.GetBool(getCLIKey) {
			return process(repo.PasswordConnector(), utils.DefaultInteractor(), getCLIKey)
		}
		if viper.GetBool(generateHashKey) {
			return process(repo.PasswordConnector(), utils.DefaultInteractor(), generateHashKey)
		}
		if viper.GetBool(listCLIKey) {
			return process(repo.PasswordConnector(), utils.DefaultInteractor(), listCLIKey)
		}

		return cmd.Help()
	},
}

const (
	putCLIKey       = "put"
	removeCLIKey    = "remove"
	getCLIKey       = "get"
	generateHashKey = "generate_password"
	listCLIKey      = "list"
)

var characterRunes = []rune("abcdefghipqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!?#$()-_=+")

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
	case getCLIKey:
		return getPassword(db, user)
	case generateHashKey:
		return generatePasswordHash()
	case listCLIKey:
		utils.DisplayPasswordList(db)
		return nil
	default:
		return fmt.Errorf("user_action_not_recognized")
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

func generatePasswordHash() error {
	rand.Seed(time.Now().UTC().UnixNano())
	randString := randomString(32)
	hash := sha1.New()
	hash.Write([]byte(randString))
	bs := hash.Sum(nil)

	fmt.Printf("%x\n", bs)
	return nil
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

func init() {
	Command.PersistentFlags().Bool(putCLIKey, false, "I want to put a password in")
	Command.PersistentFlags().Bool(removeCLIKey, false, "I want to delete a password")
	Command.PersistentFlags().Bool(generateHashKey, false, "I want to make a random hash")
	Command.PersistentFlags().Bool(getCLIKey, false, "I want to get a password")
	Command.PersistentFlags().Bool(listCLIKey, false, "I want to list my saved passwords")

	viper.BindPFlag(putCLIKey, Command.PersistentFlags().Lookup(putCLIKey))
	viper.BindPFlag(getCLIKey, Command.PersistentFlags().Lookup(getCLIKey))
	viper.BindPFlag(generateHashKey, Command.PersistentFlags().Lookup(generateHashKey))
	viper.BindPFlag(removeCLIKey, Command.PersistentFlags().Lookup(removeCLIKey))
	viper.BindPFlag(listCLIKey, Command.PersistentFlags().Lookup(listCLIKey))
}
