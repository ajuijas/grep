/*
Copyright © 2025 NAME HERE ijas.ahmd.ap@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var searchString string

func readFiles(target string, lines chan string){
	fmt.Print("To be implimented")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grep",
	Short: "",
	Long: ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("please provide search string and file/directory name")
		}
		if _, err := os.Stat(args[1]); os.IsNotExist(err) {
			return fmt.Errorf("%v: open: No such file or directory", args[1])
		}
		if _, err := os.Stat(args[1]); os.IsPermission(err) {
			return fmt.Errorf("%v: open: Permission denied", args[1])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		searchString = args[0]
		lines := make(chan string)
		readFiles(args[1], lines)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.grep.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


