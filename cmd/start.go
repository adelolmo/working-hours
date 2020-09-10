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
	"math"
	"os"
	"path"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [message]",
	Short: "Adds an entry in the timelog for starting work",
	Long: `Add to the timelog that work starts/resumes now.
It can be the beginning of the working day or coming back from a break (e.g lunch).`,
	Example:
`  You can and a message when starting the day:
    
    wk start 'good morning'

  Or to indicate that you're back from lunch or a break:
    
    wk start back`,
	Args: cobra.MaximumNArgs(1),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal(err)
		}

		_ = os.Mkdir(path.Join(configDir, "working-hours"), os.ModePerm)
		timeLogFile := path.Join(configDir, "working-hours", "timelog.txt")
		if _, err := os.Stat(timeLogFile); os.IsNotExist(err) {
			_, err = os.Create(timeLogFile)
			if err != nil {
				log.Fatal(err)
			}
		}

		tl := timelog.New(timeLogFile)
		fmt.Printf("Now is: %v\n", time.Now().Format("15:04:05"))

		message, err := tl.LastMessage()
		if err != nil {
			log.Fatal(err)
		}

		if message.Timestamp.Day() != time.Now().Day() {
			fmt.Printf("Finish work at %v\n", time.Now().Add(8*time.Hour).Format("15:04:05"))
			tl := timelog.New(timeLogFile)
			messageContent := "morning"
			if len(os.Args) == 3 {
				messageContent = os.Args[2]
			}
			err := tl.Append(messageContent, timelog.StartWorking)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		brake := time.Since(message.Timestamp)
		fmt.Printf("Brake took %v minutes\n", math.Round(brake.Minutes()))

		messages, err := tl.MessagesForDate(time.Now())
		if err != nil {
			log.Fatal(err)
		}
		workedTimeSoFar := workedTimeSoFar(messages)
		fmt.Printf("Total work done: %v\n", fmtDuration(workedTimeSoFar))
		timeToWork := time.Now().
			Add(8 * time.Hour).
			Add(-workedTimeSoFar)
		fmt.Printf("Finish work at %v\n", timeToWork.Format("15:04:05"))

		messageContent := "back"
		if len(os.Args) == 3 {
			messageContent = os.Args[2]
		}
		err = tl.Append(messageContent, timelog.StartWorking)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func workedTimeSoFar(messages []timelog.Message) time.Duration {
	workedHours := time.Duration(0)
	starBlock := time.Time{}
	for _, m := range messages {
		if m.Type == timelog.StartWorking {
			starBlock = m.Timestamp
		}
		if m.Type == timelog.StopWorking {
			var diff = m.Timestamp.Sub(starBlock).Minutes()
			workedHours = time.Duration(workedHours.Minutes()+diff) * time.Minute
		}
	}
	return workedHours
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
