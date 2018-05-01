package main

import (
	"os"

	"github.com/flaccid/snstxtr/sms"
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
)

var (
	VERSION = "v0.0.0-dev"
)

func beforeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	if c.NArg() < 1 {
		log.Fatal("message contents required")
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "snstxtr"
	app.Version = VERSION
	app.Usage = "send sms using sns"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "phone",
			Usage:  "phone number to send sms to",
			EnvVar: "PHONE",
		},
		cli.StringFlag{
			Name:   "aws-region",
			Usage:  "aws region for the sns service",
			EnvVar: "AWS_REGION",
		},
		cli.StringFlag{
			Name:   "aws-access-key-id",
			Usage:  "aws access key id",
			EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "aws-secret-access-key",
			Usage:  "aws secret access key",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
		cli.BoolFlag{
			Name:  "dry",
			Usage: "run in dry mode",
		},
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "run in debug mode",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	log.Info("send to ", c.String("phone"), " message: ", c.Args().Get(0))
	err := sms.Send(c.Args().Get(0), c.String("phone"))
	log.Info(err)

	return nil
}
