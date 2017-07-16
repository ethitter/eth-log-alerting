package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/42wim/matterbridge/matterhook"
	"github.com/asaskevich/govalidator"
	"github.com/hpcloud/tail"
)

type config struct {
	DebugDest string      `json:"debug-dest"`
	Debug     bool        `json:"debug"`
	Logs      []logConfig `json:"logs"`
}

type logConfig struct {
	LogPath     string `json:"log_path"`
	WebhookURL  string `json:"webhook_url"`
	Username    string `json:"username"`
	Channel     string `json:"channel"`
	Color       string `json:"color"`
	IconURL     string `json:"icon_url"`
	SearchRegex string `json:"search"`
}

type attachment struct {
	Fallback string `json:"fallback,omitempty"`
	Pretext  string `json:"pretext,omitempty"`
	Text     string `json:"text"`
	Color    string `json:"color,omitempty"`
}

type attachments []attachment

var (
	configPath string
	logConfigs []logConfig

	logger    *log.Logger
	debugDest string
	debug     bool
)

func init() {
	flag.StringVar(&configPath, "config", "./config.json", "Path to configuration file")
	flag.Parse()

	cfgPathValid := validatePath(&configPath)
	if !cfgPathValid {
		usage()
	}

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		usage()
	}

	cfg := config{}
	if err = json.Unmarshal(configFile, &cfg); err != nil {
		usage()
	}

	debugDest = cfg.DebugDest
	debug = cfg.Debug
	logConfigs = cfg.Logs

	setUpLogger()
}

func main() {
	logger.Printf("Starting log monitoring with config %s", configPath)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for _, logCfg := range logConfigs {
		go tailLog(logCfg)
	}

	caughtSig := <-sig

	logger.Printf("Stopping, got signal %s", caughtSig)
}

func tailLog(logCfg logConfig) {
	if logPathValid := validatePath(&logCfg.LogPath); !logPathValid {
		if debug {
			logger.Println("Invalid path: ", fmt.Sprintf("%+v\n", logCfg))
		}

		return
	}

	if !govalidator.IsURL(logCfg.WebhookURL) || len(logCfg.Username) == 0 || len(logCfg.Channel) < 2 {
		if debug {
			logger.Println("Invalid webhook, username, channel: ", fmt.Sprintf("%+v\n", logCfg))
		}

		return
	}

	t, err := tail.TailFile(logCfg.LogPath, tail.Config{
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: true,
		ReOpen:    true,
		Logger:    logger,
	})
	if err != nil {
		logger.Println(err)
		t.Cleanup()
		return
	}

	mh := matterhook.New(logCfg.WebhookURL, matterhook.Config{DisableServer: true})

	parseLinesAndSend(t, mh, logCfg)

	t.Stop()
	t.Cleanup()
}

func parseLinesAndSend(t *tail.Tail, mh *matterhook.Client, logCfg logConfig) {
	for line := range t.Lines {
		if line.Err != nil {
			continue
		}

		if len(logCfg.SearchRegex) == 0 {
			go sendLine(line, mh, logCfg)
		} else if matched, _ := regexp.MatchString(logCfg.SearchRegex, line.Text); matched {
			go sendLine(line, mh, logCfg)
		}
	}
}

func sendLine(line *tail.Line, mh *matterhook.Client, logCfg logConfig) {
	atts := attachments{
		attachment{
			Fallback: fmt.Sprintf("New entry in %s", logCfg.LogPath),
			Pretext:  fmt.Sprintf("In `%s` at `%s`:", logCfg.LogPath, line.Time),
			Text:     fmt.Sprintf("    %s", line.Text),
			Color:    logCfg.Color,
		},
	}

	mh.Send(matterhook.OMessage{
		Channel:     logCfg.Channel,
		UserName:    logCfg.Username,
		Attachments: atts,
		IconURL:     logCfg.IconURL,
	})
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

func validatePath(path *string) bool {
	if len(*path) <= 1 {
		return false
	}

	var err error
	*path, err = filepath.Abs(*path)

	if err != nil {
		logger.Printf("Error: %s", err.Error())
		return false
	}

	if _, err = os.Stat(*path); os.IsNotExist(err) {
		return false
	}

	return true
}

func usage() {
	flag.Usage()
	os.Exit(3)
}
