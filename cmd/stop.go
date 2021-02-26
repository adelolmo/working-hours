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
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [message]",
	Short: "Stops the current working session",
	Long: `Creates an entry in the timelog indicating that the active work session ends now.
It can be the end of the working day or having a break (e.g lunch).`,
	Example: `  You can and a message when ending the day:
    
    wh stop 'finish work'

  Or to indicate that you go for lunch or a break:
    
    wk stop lunch`,
	Args:                  cobra.MaximumNArgs(1),
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		tl := timelog.New(timelogFilename())

		message, err := tl.LastMessage()
		if err != nil {
			log.Fatal(err)
		}
		if message.Type == timelog.StopWorking {
			fmt.Println("Attention! There is no work session started.")
			fmt.Println("You have to start a new working session before.")
			return
		}

		fmt.Printf("Now is: %v\n", time.Now().Format("15:04"))

		messageContent := "afk"
		if len(os.Args) == 3 {
			messageContent = os.Args[2]
		}
		err = tl.Append(messageContent, timelog.StopWorking)
		if err != nil {
			log.Fatal(err)
		}

		messages, err := tl.MessagesForDate(time.Now())
		if err != nil {
			log.Fatal(err)
		}
		workedTimeSoFar := workedTimeSoFar(messages)
		fmt.Printf("Total work done: %v\n", fmtDuration(workedTimeSoFar))

		if workedTimeSoFar < 8*time.Hour {
			fmt.Printf("Time left at work: %v\n", fmtDuration(8*time.Hour-workedTimeSoFar))
		} else {
			fmt.Printf("Worked overtime: %v\n", fmtDuration(workedTimeSoFar-8*time.Hour))
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
