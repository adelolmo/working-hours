package cmd

import (
	"github.com/adelolmo/working-hours/timelog"
	"os"
	"path"
	"strings"
	"testing"
)

func TestAccountReportWithOneSickDay(t *testing.T) {
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
	f := createTimelog(tlog)
	defer os.Remove(f.Name())
	defer os.Remove(path.Dir(f.Name()))

	tl := timelog.New(f.Name())
	actual := accountReport(tl)

	expected := []string{"Total work done: 21:30",
		"Total working days: 3",
		"Balance: -00:00"}

	assertReport(t, expected, actual)
}

func TestAccountReportWithTwoSickDays(t *testing.T) {

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
	f := createTimelog(tlog)
	defer os.Remove(f.Name())
	defer os.Remove(path.Dir(f.Name()))

	tl := timelog.New(f.Name())
	actual := accountReport(tl)

	expected := []string{"Total work done: 23:00",
		"Total working days: 4",
		"Balance: -00:00"}

	assertReport(t, expected, actual)
}

func createTimelog(tlog string) *os.File {
	tmpDir, err := os.MkdirTemp("", "working-hours")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path.Join(tmpDir, "timelog.txt"))
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(tlog)
	if err != nil {
		f.Close()
		panic(err)
	}

	return f
}

func assertReport(t *testing.T, expected []string, report string) {
	r := strings.Split(report, "\n")

	fail := false
	for i := range expected {
		if r[i] != expected[i] {
			fail = true
		}
	}
	if fail {
		message := "\nExpected ->\t" + strings.Join(expected, "\t")
		message += "\nActual ->\t" + strings.Join(r, "\t")
		t.Fatal(message)
	}
}
