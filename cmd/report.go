/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/adelolmo/working-hours/timelog"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report [day, week, month, year]",
	Short: "Shows a report for the selected type",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	ValidArgs: []string{"day", "week", "month", "year"},
	Run:  report(),
}

func report() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		tl := timelog.New(timelogFilename())

		switch args[0] {
		case "day":
			messages, err := tl.MessagesForDate(time.Now())
			if err != nil {
				log.Fatal(err)
			}
			workedTimeSoFar := workedTimeSoFar(messages)
			fmt.Printf("Total work done today: %v\n", fmtDuration(workedTimeSoFar))

		case "week":
			fmt.Printf("Total work done this week: ")
			fmt.Println("not implemented")

		case "month":
			fmt.Printf("Total work done this month: ")
			fmt.Println("not implemented")

		case "year":
			fmt.Printf("Total work done this year: ")
			fmt.Println("not implemented")

		default:
			fmt.Println("Not a valid report type.")
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//reportCmd.Flags().StringVar("type", []string{"day","week","month","year"},  "Help message for toggle")
}
