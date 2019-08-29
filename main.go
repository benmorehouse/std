package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ // this is a global variable dont do this! Put it in main instead
	// or do a function that returns the cobra command 
	Use: "", // will run everytime you type nothing in 
	Short: "Task Manager",
	Example: "An example would be that you enter in test[Enter] and then the command that you want to test",
}

func main(){
	rootCmd.AddCommand(createList)
	rootCmd.AddCommand(writeList)
	rootCmd.AddCommand(deleteList)
	rootCmd.AddCommand(renameList)
	rootCmd.AddCommand(welcome)
	rootCmd.Execute()
}
