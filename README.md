# Working Hours

CLI tool for tracking how long you work.

It provides reports with insights about your working days.

## Build
```
make
```

## Install
```
sudo make install
```

## Aliases
```shell script
alias morning='wh start morning'
alias back='wh start back'
alias lunch='wh stop lunch'
alias afk='wh stop afk'
alias bye='wh stop bye'
alias whr='wh report'
```

## How does it work
For the usage examples below please make sure to set the above aliases.

* When you start working in the morning type `morning`.
* If you want to have a brake, type `afk`.
* Once you're back from your brake, type `back`.
* At the end of the working day type `bye`.

A regular day could look like this:
1. `morning`
2. `lunch`
3. `back`
4. `afk`
5. `back`
6. `bye`

The generated `timelog.txt` for the above example looks like this:
```
2021-02-1108:17:11+0100 0 morning
2021-02-1112:08:50+0100 1 lunch
2021-02-1112:40:32+0100 0 back
2021-02-1113:38:47+0100 1 afk
2021-02-1114:28:29+0100 0 back
2021-02-1117:45:53+0100 1 bye
```

## Reporting
You can have insights about your working time.

### Daily
Use the command `whr day` to get an overview of the day:
```
Total work done today: 05:27
Finish work at 18:24
```

### Weekly
Use the command `whr week` to get an overview of the week:
```
Total work done this week: 23:24
Total working days: 3
Balance: -00:36
```

### Monthly
Use the command `whr month` to get an overview of the month:
```
Total work done this month: 103:49
Total working days: 13
Balance: -00:11
```

### Yearly
Use the command `whr year` to get an overview of the year:
```
Total work done this year: 257:00
Total working days: 32
Balance: 01:00
```

### Account
Use the command `whr account` to get an overview of your account:
```
Total work done: 777:28
Total working days: 98
Balance: -01:30
```
`Balance` indicates if you have overtime (positive value), or you miss working time (negative value).
