/*
Copyright © 2021 ALEKSEY MAKSIMKIN <maximkin@mail.ru>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"bsl2sonar/finder"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:        "bsl2sonar",
	Aliases:    []string{},
	SuggestFor: []string{},
	Short:      "bsl files finder to sonarscanner",
	Long: `bsl2sonar is a CLI application for Sonar-scanner that find files with .bsl extension.
This application is a tool to generate long string with paths to .bsl files and substitute to 
sonar-properties file`,
	Example: `bsl2sonar <srcdir> <parsephrases> [flags]
bsl2sonar "/src/cf" "рн_, рнт_общая" -f "src/sonar-project.properties" -a -u`,
	ValidArgs: []string{"src", "reg"},
	Args:      checkArgs,
	Version:   "0.0.1",
	Run:       bsl2sonar,
}

// Execute is the method to run root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	rootCmd.Flags().StringP("file", "f", "", "absolute path to file sonar-project.properties")
	rootCmd.Flags().BoolP("absolute", "a", false, "output absolute files path")
	rootCmd.Flags().BoolP("unicode", "u", false, "transform cyrillic charactes to unicode")
	rootCmd.Flags().BoolP("generate", "g", false, "generate sonar-project.properties, use only with -f flag")
	rootCmd.Flags().BoolP("logging", "l", false, "output log info to stdout")

}

// Check cmd arguments
func checkArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("requires only two arguments: srcdir [string] and parsephrases [string with comma separate]")
	}
	fileFlag, _ := cmd.Root().Flags().GetString("file")
	genFlag, _ := cmd.Root().Flags().GetBool("generate")
	checkResult, errText := isArgsValid(args, fileFlag, genFlag)
	if checkResult {
		return nil
	}
	return errors.New(errText)
}

func isArgsValid(args []string, fileFlag string, genFlag bool) (result bool, errText string) {

	fileInfo, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		errText := fmt.Sprintf("Path \"%s\" dosn't exist", args[0])
		return false, errText
	}

	if !fileInfo.IsDir() {
		errText := fmt.Sprintf("File \"%s\" is not directory", args[0])
		return false, errText
	}

	if len([]rune(args[1])) < 3 {
		errText := "must be at least 3 characters of parsephrases"
		return false, errText
	}

	if genFlag {
		if len(fileFlag) == 0 {
			errText := "Can't use flag -g without flag -f because need to know path to save template"
			return false, errText
		}
	} else {
		if len(fileFlag) != 0 {
			if file, err := os.Stat(fileFlag); os.IsNotExist(err) || file.IsDir() {
				errText := "file not found"
				return false, errText
			}
		}
	}

	return true, ""
}

func bsl2sonar(cmd *cobra.Command, args []string) {

	fndr := finder.NewFinder(args[0], args[1])
	fndr.Sfile, _ = cmd.Flags().GetString("file")
	fndr.Abspath, _ = cmd.Flags().GetBool("absolute")
	fndr.Unicode, _ = cmd.Flags().GetBool("unicode")
	fndr.Generate, _ = cmd.Flags().GetBool("generate")
	fndr.Logging, _ = cmd.Flags().GetBool("logging")

	if fndr.Logging {
		fndr.Logger.Printf(">>> Абсолютный путь к исходным файлам проекта: %s", args[0])
		fndr.Logger.Printf(">>> Абсолютный путь к исходным файлу sonar-project.properties: %s", fndr.Sfile)
	}

	fndr.DataToSonarQube()

}
