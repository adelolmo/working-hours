package timelog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const dateFormat = "2006-01-0215:04:05-0700"

type MessageType int

const (
	StartWorking MessageType = iota
	StopWorking
	Unknown
)

func (t MessageType) String() string {
	return [...]string{"0", "1"}[t]
}

func ParseMessageType(t string) MessageType {
	switch t {
	case "0":
		return StartWorking
	case "1":
		return StopWorking
	default:
		return Unknown
	}
}

type Log struct {
	filename string
}

type Message struct {
	Timestamp time.Time
	Content   string
	Type      MessageType
}

func New(file string) *Log {
	return &Log{
		filename: file,
	}
}

func (l *Log) Append(message string, messageType MessageType) error {
	timeLog, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot append to timelog file. cause: %s", err.Error())
	}

	_, err = timeLog.WriteString(fmt.Sprintf("%v %d %s\n",
		time.Now().Format(dateFormat), messageType, message))
	if err != nil {
		return err //fmt.Errorf("cannot write filename %s. Error: %s", l.filename, err)
	}

	return nil
}

func (l *Log) LastMessage() (*Message, error) {
	f, err := os.Open(l.filename)
	if err != nil {
		return &Message{}, err
	}

	scanner := bufio.NewScanner(f)

	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	if len(line) == 0 {
		return &Message{}, nil
	}

	part := strings.Split(line, " ")
	date, err := time.Parse(dateFormat, part[0])
	if err != nil {
		return &Message{}, err
	}
	err = f.Close()
	if err != nil {
		return &Message{}, err
	}
	return &Message{
		Timestamp: date,
		Content:   part[1],
	}, nil
}

func (l *Log) MessagesForDate(date time.Time) ([]Message, error) {
	f, err := os.Open(l.filename)
	if err != nil {
		return []Message{}, err
	}

	scanner := bufio.NewScanner(f)

	var messages []Message
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) == 0 {
			continue
		}

		part := strings.Split(line, " ")
		d, err := time.Parse(dateFormat, part[0])
		if err != nil {
			panic(err)
		}
		if d.Day() == date.Day() {
			messages = append(messages, Message{
				Timestamp: d,
				Type:      ParseMessageType(part[1]),
				Content:   part[2],
			})
		}
	}
	return messages, nil
}