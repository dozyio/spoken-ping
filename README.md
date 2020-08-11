# Go Spoken Ping

A TCP connect ping written in Go. Internet was down for the afternoon, so wrote this to let me know when it was working again. Requires OS X for the 'say' command

## Build
```sh
go build -o spoken-ping
```

## Run
```sh
./spoken-ping -p 80 -w 30 -s 5 1.1.1.1
```

## Usage

-p Port

-w Wait timeout

-s Stop after X successful pings


