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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Usage: "host id",
		},
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help",
		},
	}
	app.HideHelp = true
	app.Action = func(ctx *cli.Context) error {
		if ctx.GlobalBool("help") {
			return cli.ShowAppHelp(ctx)
		}
		client, hostID, err := setupClientHostID(ctx)
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

func setupClientHostID(ctx *cli.Context) (*mackerel.Client, string, error) {
	var conf *config.Config
	var err error
	var apiKey, apiBase string

	apiKey = os.Getenv("MACKEREL_APIKEY")
	if apiKey == "" {
		confFile := config.DefaultConfig.Conffile
		conf, err = config.LoadConfig(confFile)
		if err != nil {
			return nil, "", errors.New("MACKEREL_APIKEY not set")
		}
		apiKey = conf.Apikey
		if apiKey == "" {
			return nil, "", errors.New("MACKEREL_APIKEY not set")
		}
		apiBase = conf.Apibase
	}
	if apiBase == "" {
		apiBase = config.DefaultConfig.Apibase
	}
	client, err := mackerel.NewClientWithOptions(apiKey, apiBase, false)
	if err != nil {
		return nil, "", err
	}

	hostID := ctx.GlobalString("host")
	if hostID == "" {
		hostID, err = loadHostID(conf.Root)
		if err != nil {
			return nil, "", errors.New("specify host id")
		}
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
