package main

import (
	"cid.com/content-service/business"
	"cid.com/content-service/common/concurrency"
	"cid.com/content-service/common/secrets"
	"gopkg.in/yaml.v2"
	"os"
	"syscall"
)

func main() {
	cfg, vault := mustLoadConfiguration()
	ctx := concurrency.CancelOnSignal(syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	app, err := business.New(ctx, cfg, vault)
	if err != nil {
		panic(err)
	}
	if err = app.Run(ctx); err != nil {
		panic(err)
	}
}

func mustLoadConfiguration() (business.Configuration, secrets.Secrets) {
	conf := new(business.Configuration)
	vault := new(secrets.Secrets)
	if len(os.Args) != 5 {
		panic("missing required argument")
	}
	for i := 1; i < 5; i += 2 {
		file, err := os.Open(os.Args[i+1])
		if err != nil {
			panic(err)
		}
		dec := yaml.NewDecoder(file)
		arg := os.Args[i]
		switch arg {
		case "-cfg", "--configuration-file":
			if err = dec.Decode(conf); err != nil {
				panic(err)
			}
		case "-sec", "--secrets-file":
			if err = dec.Decode(vault); err != nil {
				panic(err)
			}
		}
		_ = file.Close()
	}
	if conf == nil {
		panic("missing required argument -cfg or --configuration-file")
	} else if vault == nil {
		panic("missing required argument -sec or --secrets-file")
	}
	return *conf, *vault
}
