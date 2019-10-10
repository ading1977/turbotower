package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbotower/version"
	"github.com/urfave/cli"
	"os"
	"path"
	"time"
)

func Run() {
	app := &cli.App{}
	app.Name = path.Base(os.Args[0])
	app.Description = "turbotower controls Turbonomic Operations Manager."
	app.Version = version.Version + " (" + version.GitCommit + ")"
	if buildtime, err := time.Parse(time.UnixDate, version.BuildTime); err == nil {
		app.Compiled = buildtime
	} else {
		app.Compiled = time.Time{}
	}
	app.Commands = commands
	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.WarnLevel)
		log.SetFormatter(
			&log.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: timeFormat,
			})
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
