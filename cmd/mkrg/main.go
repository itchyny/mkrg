package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/itchyny/mkrg"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-client-go"
)

const cmdName = "mkrg"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", cmdName, err)
		os.Exit(1)
	}
}

func run() error {
	client, hostID, err := setupClientHostID()
	if err != nil {
		return err
	}
	app := mkrg.NewApp(client, hostID)
	if err := app.Run(); err != nil {
		return err
	}
	return nil
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
