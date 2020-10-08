package main

import (
	"github.com/benmorehouse/std/command/create"
	"github.com/benmorehouse/std/command/delete"
	"github.com/benmorehouse/std/command/edit"
	"github.com/benmorehouse/std/command/list"
	"github.com/benmorehouse/std/command/password"
	"github.com/spf13/cobra"
)

var bucketName = []byte("Lists")

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "STD - Shit to do! Use this to manage all of your tasks, grocery lists, contacts or passwords all from the command line!",
}

func main() {
	rootCmd.AddCommand(create.Command)
	rootCmd.AddCommand(delete.Command)
	rootCmd.AddCommand(edit.Command)
	rootCmd.AddCommand(list.Command)
	rootCmd.AddCommand(password.Command)
	rootCmd.Execute()
}
