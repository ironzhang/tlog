package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/zaplog"
)

func LoadConfig(file string) (conf zaplog.Config, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return conf, err
	}
	if err = json.Unmarshal(data, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func main() {
	file := "../configs/std.json"
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
	//tlog.Panic("panic")
	//tlog.Fatal("fatal")

	log := tlog.GetLogger("access")
	log.Debug("access debug")
	log.Info("access info")
	//	log.Warn("warn")
	//	log.Error("error")
}
