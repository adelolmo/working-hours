package cmd

import (
	"fmt"
	"github.com/adelolmo/working-hours/timelog"
	"os"
	"path"
	"testing"
)

func TestAccountReportWithOneSickDay(t *testing.T) {

	tempDir := os.TempDir()
	tlog := `2021-04-1509:00:00+0200 0 morning
2021-04-1512:00:00+0200 1 afk
2021-04-1512:30:00+0200 0 back
2021-04-1515:00:00+0200 2 sick-one
2021-04-1609:00:00+0200 0 morning
2021-04-1612:00:00+0200 1 lunch
2021-04-1612:30:00+0200 0 back
2021-04-1617:30:00+0200 1 bye
2021-04-1909:00:00+0200 0 morning
2021-04-1912:00:00+0200 1 lunch
2021-04-1912:30:00+0200 0 back
2021-04-1917:30:00+0200 1 bye
`
	f, err := os.Create(path.Join(tempDir,"timelog-test.txt"))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(tlog)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	tl := timelog.New(f.Name())

	report := accountReport(tl)

	expected := "Total work done: 21:30\nTotal working days: 3\nBalance: -00:00\n"
	if report != expected {
		t.Fatal("argg" + report)
	}
}

func TestAccountReportWithTwoSickDays(t *testing.T) {

	tempDir := os.TempDir()
	tlog := `2021-04-1509:00:00+0200 0 morning
2021-04-1512:00:00+0200 1 afk
2021-04-1512:30:00+0200 0 back
2021-04-1515:00:00+0200 2 sick-one
2021-04-1609:00:00+0200 0 morning
2021-04-1612:00:00+0200 1 lunch
2021-04-1612:30:00+0200 0 back
2021-04-1617:30:00+0200 1 bye
2021-04-1909:00:00+0200 0 morning
2021-04-1912:00:00+0200 1 lunch
2021-04-1912:30:00+0200 0 back
2021-04-1917:30:00+0200 1 bye
2021-04-2009:00:00+0200 0 morning
2021-04-2010:30:00+0200 2 sick-two
`
	f, err := os.Create(path.Join(tempDir,"timelog-test.txt"))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(tlog)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	tl := timelog.New(f.Name())

	report := accountReport(tl)

	expected := "Total work done: 23:00\nTotal working days: 4\nBalance: -00:00\n"
	if report != expected {
		t.Fatal("argg" + report)
	}
}
