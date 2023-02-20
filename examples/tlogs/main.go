package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"git.xiaojukeji.com/pearls/tlog"
	"git.xiaojukeji.com/pearls/tlog/zaplog"
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
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		return
	}
	logger, err := zaplog.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new logger: %v\n", err)
		return
	}
	defer logger.Close()

	tlog.SetLogger(logger)

	//run()
	RunTestCase()
}

func RunTestCase() {
	n := 0
	a := strings.Repeat("a", 100)
	b := strings.Repeat("b", 100)

	for i := 0; i < 10; i++ {
		tlog.Debugf("%d: %s", n, a)
		n++

		tlog.Debugf("%d: %s", n, b)
		n++

		time.Sleep(20 * time.Second)
	}

	time.Sleep(time.Hour)
}

func run() {
	tlog.Debug("debug")
	tlog.Info("info")
	tlog.Warn("warn")
	tlog.Error("error")

	log := tlog.Named("access")
	log.Debug("access debug")
	log.Info("access info")
	log.Warn("access warn")
	log.Error("access error")

	//tlog.Panic("panic")
	//tlog.Fatal("fatal")
}
