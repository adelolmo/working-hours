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
	"os"
	"time"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report [day, week, month, year, account]",
	Short: "Shows a report for the selected type",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"day", "week", "month", "year", "account"},
	Run:                   report(),
}

func init() {
	rootCmd.AddCommand(reportCmd)
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
			fmt.Printf("Finish work at %v\n", time.Now().Add(8*time.Hour-workedTimeSoFar).Format("15:04"))

		case "week":
			start, end := weekRange(time.Now().ISOWeek())
			messages, err := tl.MessagesForDateRange(start, end)
			if err != nil {
				log.Fatal(err)
			}
			workedTimeSoFar := workedTimeSoFar(messages)
			fmt.Printf("Total work done this week: %v\n", fmtDuration(workedTimeSoFar))

			numberOfWorkingDays, numberOfWorkingHours := workedDaysAndHours(messages)
			fmt.Printf("Total working days: %d\n", numberOfWorkingDays)

			if workedTimeSoFar > numberOfWorkingHours {
				fmt.Printf("Balance: %v\n", fmtDuration(workedTimeSoFar-numberOfWorkingHours))
			} else {
				fmt.Printf("Balance: -%v\n", fmtDuration(numberOfWorkingHours-workedTimeSoFar))
			}

		case "month":
			now := time.Now()
			currentYear, currentMonth, _ := now.Date()
			currentLocation := now.Location()

			firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
			lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

			messages, err := tl.MessagesForDateRange(firstOfMonth, lastOfMonth)
			if err != nil {
				log.Fatal(err)
			}
			workedTimeSoFar := workedTimeSoFar(messages)
			fmt.Printf("Total work done this month: %v\n", fmtDuration(workedTimeSoFar))

			numberOfWorkingDays, numberOfWorkingHours := workedDaysAndHours(messages)
			fmt.Printf("Total working days: %d\n", numberOfWorkingDays)

			if workedTimeSoFar > numberOfWorkingHours {
				fmt.Printf("Balance: %v\n", fmtDuration(workedTimeSoFar-numberOfWorkingHours))
			} else {
				fmt.Printf("Balance: -%v\n", fmtDuration(numberOfWorkingHours-workedTimeSoFar))
			}

		case "year":
			now := time.Now()
			currentLocation := now.Location()

			firstOfJanuary := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, currentLocation)
			thirtyFirstOfDecember := time.Date(now.Year(), time.December, 31, 0, 0, 0, 0, currentLocation)

			messages, err := tl.MessagesForDateRange(firstOfJanuary, thirtyFirstOfDecember)
			if err != nil {
				log.Fatal(err)
			}
			workedTimeSoFar := workedTimeSoFar(messages)
			fmt.Printf("Total work done this year: %v\n", fmtDuration(workedTimeSoFar))

			numberOfWorkingDays, numberOfWorkingHours := workedDaysAndHours(messages)
			fmt.Printf("Total working days: %d\n", numberOfWorkingDays)

			if workedTimeSoFar > numberOfWorkingHours {
				fmt.Printf("Balance: %v\n", fmtDuration(workedTimeSoFar-numberOfWorkingHours))
			} else {
				fmt.Printf("Balance: -%v\n", fmtDuration(numberOfWorkingHours-workedTimeSoFar))
			}

		case "account":
			messages, err := tl.AllMessages()
			if err != nil {
				log.Fatal(err)
			}
			workedTimeSoFar := workedTimeSoFar(messages)
			fmt.Printf("Total work done: %v\n", fmtDuration(workedTimeSoFar))

			numberOfWorkingDays, maxWorkingDuration := workedDaysAndHours(messages)
			fmt.Printf("Total working days: %d\n", numberOfWorkingDays)

			sickDuration := sickHours(messages)

			totalWorkDuration := workedTimeSoFar + sickDuration
			if totalWorkDuration > maxWorkingDuration {
				fmt.Printf("Balance: %v\n", fmtDuration(totalWorkDuration-maxWorkingDuration))
			} else {
				fmt.Printf("Balance: -%v\n", fmtDuration(maxWorkingDuration-totalWorkDuration))
			}

		default:
			fmt.Println("Not a valid report type.")
			os.Exit(1)
		}
	}
}

// https://github.com/icza/gox/blob/master/timex/timex.go
func workedDaysAndHours(messages []timelog.Message) (int, time.Duration) {
	numberOfWorkingDays := 0
	dayOfMonth := 0
	for _, message := range messages {
		if dayOfMonth != message.Timestamp.Day() {
			dayOfMonth = message.Timestamp.Day()
			numberOfWorkingDays++
		}
	}
	numberOfWorkingHours := time.Duration(8*numberOfWorkingDays) * time.Hour
	return numberOfWorkingDays, numberOfWorkingHours
}

func sickHours(messages []timelog.Message) time.Duration {
	//numberOfSickMinutes := 0
	sickDuration := time.Duration(0)
	workDuration := time.Duration(0)
	dayOfMonth := 0
	//indexDayStart := 0

	for i, message := range messages {
		if dayOfMonth != message.Timestamp.Day() {
			dayOfMonth = message.Timestamp.Day()
			workDuration = time.Duration(0)
			//indexDayStart = i
		}
		if message.Type == timelog.StopWorking {
			workDuration += message.Timestamp.Sub(messages[i-1].Timestamp)
		}
		if message.Type == timelog.StopWorkingSick {
			workDuration += message.Timestamp.Sub(messages[i-1].Timestamp)
			sickDuration = 8*time.Hour - workDuration
		}
	}
	return sickDuration
}

func weekRange(year, week int) (start, end time.Time) {
	start = weekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}

func weekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}
