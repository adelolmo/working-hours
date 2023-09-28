wh 1 "November, 13 2020" wh "User Manual"
=========================================

# NAME
  wh - Tracking how long you work

# SYNOPSIS
  wh [help report start stop] [-h]

# DESCRIPTION
  *wh* is a tool for tracking how long you work. It provides reports with insights about your working days. A working day is calculated to have 8 working hours.

  The commands are as follows:
  help        Help about any command.
  report      Shows a report for the selected type.
  start       Starts a working session.
  stop        Stops the current working session.

  The options are as follows:
  -h, --help   help for wh

# EXAMPLES
  Check-in in the morning when you start working.

    wh start

    Now is: 07:45
    Finish work at 15:45

  Check-out when you want to have a brake, or you finish your working day.

    wh stop

    Now is: 16:40
    Total work done: 07:55
    Time left at work: 00:05

  Check how long you've worked today and how log is left.

    wh report day

    Total work done today: 03:30
    Finish work at 04:30

# AUTHOR
  Andoni del Olmo (andoni.delolmo@gmail.com) 