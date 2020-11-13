/*
Copyright © 2020 Andoni del Olmo <andoni.delolmo@gmail.com>

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
	"time"
)

var startCmd = &cobra.Command{
	Use:   "start [message]",
	Short: "Starts a working session",
	Long: `Creates an entry in the timelog indicating that a work session starts.
It can be the beginning of the working day or coming back from a break (e.g lunch).`,
	Example:
	`  You can and a message when starting the day:
    
    wk start 'good morning'

  Or to indicate that you're back from lunch or a break:
    
    wk start back`,
	Args:                  cobra.MaximumNArgs(1),
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		tl := timelog.New(timelogFilename())

		message, err := tl.LastMessage()
		if err != nil {
			log.Fatal(err)
		}

		if message.Type == timelog.StartWorking {
			fmt.Printf("Attention! A work session was already started at %s.\nYou have to end this session before starting a new one.\n",
				message.Timestamp.Format("15:04"))
			return
		}

		fmt.Printf("Now is: %v\n", time.Now().Format("15:04"))

		if message.Timestamp.Day() != time.Now().Day() {
			fmt.Printf("Finish work at %v\n", time.Now().Add(8*time.Hour).Format("15:04"))

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
}
