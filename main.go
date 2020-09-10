package main

import (
	"fmt"
	"github.com/adelolmo/working-hours/timelog"
	"log"
	"math"
	"os"
	"path"
	"time"
)

var timeLogFile string

func main() {
	if len(os.Args) < 2 {
		fmt.Print(`Usage: working-hours <command>

List of commands:
  start [message]  Adds an entry in the timelog for starting work. An additional message is optional.
  stop [message]   Adds an entry in the timelog for stopping work. An additional message is optional.
  report <type>    Shows a report for the selected type.
                   List of types:
                     * day
                     * week
                     * month
                     * year
`)
		os.Exit(0)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	_ = os.Mkdir(path.Join(configDir, "working-hours"), os.ModePerm)
	timeLogFile = path.Join(configDir, "working-hours", "timelog.txt")
	if _, err := os.Stat(timeLogFile); os.IsNotExist(err) {
		_, err = os.Create(timeLogFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	tl := timelog.New(timeLogFile)

	option := os.Args[1]
	switch option {
	case "start":
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

	case "stop":
		fmt.Printf("Now is: %v\n", time.Now().Format("15:04:05"))

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

	default:
		fmt.Println("Not a valid option")
		os.Exit(1)
	}
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
