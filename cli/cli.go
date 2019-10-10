package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbotower/utils"
	"github.com/turbonomic/turbotower/version"
	"github.com/urfave/cli"
	"os"
	"path"
	"time"
)

func Run() {
	app := &cli.App{}
	app.Name = path.Base(os.Args[0])
	app.Usage = "turbotower controls Turbonomic Operations Manager."
	app.Version = version.Version + " (" + version.GitCommit + ")"
	if buildTime, err := time.Parse(time.RFC1123Z, version.BuildTime); err == nil {
		app.Compiled = buildTime
	} else {
		app.Compiled = time.Time{}
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Value:  utils.GetLocalIP() + ":443",
			Usage:  "specify the endpoint of the Turbonomic server",
			EnvVar: "TURBO_SERVER",
		},
		cli.StringFlag{
			Name:   "influxdb",
			Value:  utils.GetLocalIP() + ":8086",
			Usage:  "specify the endpoint of the InfluxDB server",
			EnvVar: "INFLUXDB_SERVER",
		},
		cli.StringFlag{
			Name:   "database",
			Value:  "metron",
			Usage:  "specify the database name",
			EnvVar: "INFLUXDB_DATABASE",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "enable debug mode",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "log-level, l",
			Value:  "info",
			Usage:  "specify log level (debug, info, warn, error, fatal, panic)",
			EnvVar: "LOG_LEVEL",
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.SetLevel(level)
		// If a log level wasn't specified and we are running in debug mode,
		// enforce log-level=debug.
		if !c.IsSet("log-level") && !c.IsSet("l") &&
			(c.Bool("debug") || c.Bool("d")) {
			log.SetLevel(log.DebugLevel)
		}
		log.SetFormatter(
			&log.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: timeFormat,
			})
		return nil
	}

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
