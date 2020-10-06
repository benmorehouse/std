package main

import (
	"github.com/benmorehouse/std/command/create"
	"github.com/benmorehouse/std/command/delete"
	"github.com/benmorehouse/std/command/edit"
	"github.com/benmorehouse/std/command/list"
	"github.com/spf13/cobra"
)

var bucketName = []byte("Lists")

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A CLI Task Manager",
}

func main() {
	rootCmd.AddCommand(create.Command)
	rootCmd.AddCommand(delete.Command)
	rootCmd.AddCommand(edit.Command)
	rootCmd.AddCommand(list.Command)
	rootCmd.Execute()
}
