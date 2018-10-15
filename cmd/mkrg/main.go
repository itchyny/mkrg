package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/itchyny/mkrg"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	cmdName     = "mkrg"
	description = "Mackerel graph viewer in terminal"
	version     = "0.0.0"
	author      = "itchyny"
)

func main() {
	if run(os.Args) != nil {
		os.Exit(1)
	}
}

func run(args []string) error {
	app := cli.NewApp()
	app.Name = cmdName
	app.HelpName = cmdName
	app.Usage = description
	app.Version = version
	app.Author = author
	app.Flags = []cli.Flag{}
	app.Action = func(c *cli.Context) error {
		client, hostID, err := setupClientHostID()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", cmdName, err)
			return err
		}
		err = mkrg.NewApp(client, hostID).Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", cmdName, err)
		}
		return err
	}
	return app.Run(args)
}

func setupClientHostID() (*mackerel.Client, string, error) {
	confFile := config.DefaultConfig.Conffile
	conf, err := config.LoadConfig(confFile)
	if err != nil {
		return nil, "", err
	}
	apiKey := conf.Apikey
	if key := os.Getenv("MACKEREL_APIKEY"); key != "" {
		apiKey = key
	}
	if apiKey == "" {
		return nil, "", errors.New("MACKEREL_APIKEY not set")
	}
	apiBase := conf.Apibase
	if apiBase == "" {
		apiBase = config.DefaultConfig.Apibase
	}
	client, err := mackerel.NewClientWithOptions(apiKey, apiBase, false)
	if err != nil {
		return nil, "", err
	}
	hostID, err := loadHostID(conf.Root)
	if err != nil {
		return nil, "", err
	}
	return client, hostID, nil
}

func loadHostID(root string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(root, "id"))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
