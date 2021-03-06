package cmd

import (
	"fmt"
	"github.com/adelolmo/working-hours/timelog"
	"log"
	"os"
	"path"
	"time"
)

func timelogFilename() string {
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
	return timeLogFile
}

func workedTimeSoFar(messages []timelog.Message) time.Duration {
	if len(messages) == 0 {
		return time.Duration(0)
	}

	if len(messages) == 1 { // TODO don't assume message is type start
		return time.Now().Sub(messages[0].Timestamp)
	}

	workedHours := time.Duration(0)
	starBlock := time.Time{}
	for _, m := range messages {
		switch m.Type {
		case timelog.StartWorking:
			starBlock = m.Timestamp

		case timelog.StopWorking, timelog.StopWorkingSick:
			var diff = m.Timestamp.Sub(starBlock).Minutes()
			workedHours = time.Duration(workedHours.Minutes()+diff) * time.Minute
		}
	}

	lastMessage := messages[len(messages)-1]
	if lastMessage.Type == timelog.StartWorking {
		return workedHours + time.Now().Sub(lastMessage.Timestamp)
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
