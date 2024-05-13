package main

import (
	"cid.com/content-service/business"
	"cid.com/content-service/common/concurrency"
	"gopkg.in/yaml.v2"
	"os"
	"syscall"
)

func main() {
	cfg := mustLoadConfiguration()
	ctx := concurrency.CancelOnSignal(syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	app, err := business.New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	if err = app.Run(ctx); err != nil {
		panic(err)
	}
}

func mustLoadConfiguration() business.Configuration {
	if len(os.Args) != 2 {
		panic("missing required argument for configuration file")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()
	dec := yaml.NewDecoder(file)
	var conf business.Configuration
	if err = dec.Decode(&conf); err != nil {
		panic(err)
	}
	return conf
}
