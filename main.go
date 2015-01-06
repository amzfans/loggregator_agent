package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/amzfans/loggregator_agent/listener"
	"github.com/amzfans/loggregator_agent/sender"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

var configFile = flag.String("config", "config/loggregator_agent.json", "Location of the loggregator agent config json file")

type Config struct {
	Index           uint64
	AppId           string
	DropplerAddress string
	SharedSecret    string
}

func main() {
	flag.Parse()

	// get the configuration.
	config, err := readConfigFromFile(*configFile)
	if err != nil {
		panic(err)
	}

	logListener := listener.NewLogListener()
	logListener.Start()
	defer logListener.Stop()

	logSender, err := sender.NewLogSender(config.AppId, config.DropplerAddress,
		config.SharedSecret, config.Index, logListener)
	if err != nil {
		panic(err)
	}
	logSender.Start()

	killChan := make(chan os.Signal)
	signal.Notify(killChan, os.Kill, os.Interrupt)
	for {
		select {
		case <-logSender.StopChan:
			log.Println("Shutdown the loggregator agent.")
			logListener.Stop()
			os.Exit(0)
		case <-killChan:
			log.Println("Shutdown the loggregator agent.")
			logListener.Stop()
			logSender.Stop()
			os.Exit(0)
		}
	}
}

func readConfigFromFile(cfgPath string) (config *Config, err error) {
	configBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can not read config file [%s]: %s", cfgPath, err))
	}
	config = &Config{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can not parse config file %s: %s", cfgPath, err))
	}
	return
}
