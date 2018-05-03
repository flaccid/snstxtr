package main

import (
	"os"

	"github.com/flaccid/snstxtr"
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

	// check if needed aws variables are available
	if len(c.String("aws-region")) < 1 || len(c.String("aws-access-key-id")) < 1 || len(c.String("aws-secret-access-key")) < 1 {
		log.Fatal("insufficient aws credentials")
	}

	if c.NArg() < 1 && !c.Bool("daemon") {
		log.Fatal("message contents required")
	}

	if len(c.String("phone")) < 1 && !c.Bool("daemon") {
		log.Fatal("phone number required")
	}

	// explicitly set the phone number in the env in case using an incoming webhook
	os.Setenv("PHONE", c.String("phone"))

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
			Name:   "phone,n",
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
			Name:  "daemon",
			Usage: "run in daemon mode (web service)",
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
	if c.Bool("daemon") {
		snstxtr.Serve()
	} else {
		log.Info("send to ", c.String("phone"), " message: ", c.Args().Get(0))
		err := snstxtr.Send(c.Args().Get(0), c.String("phone"))
		log.Info(err)
	}

	return nil
}
