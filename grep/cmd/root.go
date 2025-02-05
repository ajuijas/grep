/* trunk-ignore-all(gofmt) */
/*
Copyright © 2025 NAME HERE ijas.ahmd.ap@gmail.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/spf13/cobra"
)

var searchString string
const maxLineSize = 64 * 1024
var outputFile string

var getTextScanner = func (filepath string) (*bufio.Scanner) {
	var source io.Reader
	var scanner *bufio.Scanner
	if filepath != "" {
		source, _ = os.Open(filepath)
		scanner = bufio.NewScanner(source)
		buf := make([]byte, maxLineSize)
		scanner.Buffer(buf, maxLineSize)
	} else {  // Filename empty is redirected to stdin
		source = os.Stdin
		scanner = bufio.NewScanner(source)
	}
	return scanner
}


func readFiles(filePath string, wg *sync.WaitGroup){

	defer wg.Done()

	var lines []string
	// defer source.Close()

	re, err := regexp.Compile(searchString)
	if err!=nil {
		log.Fatal(err)
	}

	scanner := getTextScanner(filePath)

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line){
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if outputFile != "" {
		outFile, _ := os.Create(outputFile)
		for _, line := range lines {
			outFile.WriteString(line + "\n")
			}
	} else {
		for _, line := range lines {
			fmt.Println(filePath, line)
			}
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grep <string pattern> <file name or directory>",
	Short: "",
	Long: ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please provide search string and file/directory name")
		}
		if len(args) < 2 {
			return nil
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
		var wg sync.WaitGroup
		var err error
		if len(args) < 2 {

			wg.Add(1)
			go readFiles("", &wg)  // The empty file name

		} else {

			dir := args[1]

			err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error accessing path:", err)
				return err
			}
			if !info.IsDir() {
				wg.Add(1)
				readFiles(path, &wg)
			}
			return nil
		})}
	wg.Wait()

	if err != nil {
		fmt.Println("Error walking through the directory:", err)
	}
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
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "file path to save output")
}
