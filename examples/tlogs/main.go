package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/zaplog"
)

func LoadConfig(file string) (conf zaplog.Config, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return conf, err
	}

	unmarshal := json.Unmarshal
	if path.Ext(file) == ".yaml" {
		unmarshal = yaml.Unmarshal
	}

	if err = unmarshal(data, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func main() {
	file := "../configs/development.json"
	if len(os.Args) >= 2 {
		file = os.Args[1]
	}

	cfg, err := LoadConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v", err)
		return
	}
	logger, err := zaplog.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new logger: %v", err)
		return
	}
	defer logger.Close()

	tlog.SetFactory(logger)

	run()
}

func run() {
	tlog.Debug("debug")
	tlog.Info("info")
	tlog.Warn("warn")
	tlog.Error("error")

	log := tlog.GetLogger("access")
	log.Debug("access debug")
	log.Info("access info")
	log.Warn("access warn")
	log.Error("access error")

	//tlog.Panic("panic")
	//tlog.Fatal("fatal")
}
