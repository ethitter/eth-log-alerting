package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/asaskevich/govalidator"
	"github.com/hpcloud/tail"
)

var (
	logPath     string
	webhookURL  string
	username    string
	channel     string
	color       string
	iconURL     string
	searchRegex string

	logger    *log.Logger
	debugDest string
	debug     bool

	tailer *tail.Tail
)

func init() {
	flag.StringVar(&logPath, "log-path", "", "Log to monitor")
	flag.StringVar(&webhookURL, "webhook", "", "Webhook to forward log entries to")
	flag.StringVar(&username, "username", "logbot", "Username to post as")
	flag.StringVar(&channel, "channel", "", "Channel to post log entries to")
	flag.StringVar(&color, "color", "default", "Color for entry, either named or hex with `#`")
	flag.StringVar(&iconURL, "icon-url", "", "URL of icon to use for bot")
	flag.StringVar(&searchRegex, "search", "", "Search term or regex to match")
	flag.StringVar(&debugDest, "debug-dest", "os.Stdout", "Destination for debug and other messages, omit to log to Stdout")
	flag.BoolVar(&debug, "debug", false, "Include additional log data for debugging")
	flag.Parse()

	setUpLogger()

	validatePath(&logPath)

	if !govalidator.IsURL(webhookURL) || len(channel) < 2 {
		usage()
	}
}

func main() {
	logger.Printf("Monitoring %s", logPath)
	logger.Printf("Forwarding to channel \"%s\" as user \"%s\" at %s", channel, username, webhookURL)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go tailLog()

	caughtSig := <-sig
	tailer.Stop()
	tailer.Cleanup()
	logger.Printf("Stopping, got signal %s", caughtSig)
}

func tailLog() {
	t, err := tail.TailFile(logPath, tail.Config{Follow: true, MustExist: true, ReOpen: true})
	tailer = t

	if err != nil {
		logger.Println(err)
		usage()
	}

	for line := range t.Lines {
		if line.Err != nil {
			continue
		}

		if len(searchRegex) == 0 {
			go sendLine(line)
		} else {
			matched, _ := regexp.MatchString(searchRegex, line.Text)
			if matched {
				go sendLine(line)
			}
		}
	}
}

func sendLine(line *tail.Line) {
	logger.Println(line.Text)
}

func setUpLogger() {
	logOpts := log.Ldate | log.Ltime | log.LUTC | log.Lshortfile

	if debugDest == "os.Stdout" {
		logger = log.New(os.Stdout, "DEBUG: ", logOpts)
	} else {
		path, err := filepath.Abs(debugDest)
		if err != nil {
			logger.Fatal(err)
		}

		logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}

		logger = log.New(logFile, "", logOpts)
	}
}

func validatePath(path *string) {
	if len(*path) > 1 {
		var err error
		*path, err = filepath.Abs(*path)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(3)
		}

		if _, err = os.Stat(*path); os.IsNotExist(err) {
			usage()
		}
	} else {
		usage()
	}
}

func usage() {
	flag.Usage()
	os.Exit(3)
}
