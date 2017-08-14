//     Digota <http://digota.com> - eCommerce microservice
//     Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.
//
//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published
//     by the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.
//
//     You should have received a copy of the GNU Affero General Public License
//     along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

// prevents gofmt changes the imports order

// first init
import _ "github.com/digota/digota/config"

// standards imports
import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/server"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

const (
	//defaultLevel = logrus.WarnLevel
	defaultLevel = logrus.DebugLevel
	configPath   = "digota.yaml"
	version      = "0.1"
	name         = "digota"
	usage        = ""
	email        = "yaron@digota.com"
	Description  = "eCommerce microservice"
	port         = ":3051"
	copyright    = `
		Digota <http://digota.com> - eCommerce microservice
		Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.

		This program is free software: you can redistribute it and/or modify
		it under the terms of the GNU Affero General Public License as published
		by the Free Software Foundation, either version 3 of the License, or
		(at your option) any later version.

		This program is distributed in the hope that it will be useful,
		but WITHOUT ANY WARRANTY; without even the implied warranty of
		MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
		GNU Affero General Public License for more details.

		You should have received a copy of the GNU Affero General Public License
		along with this program.  If not, see <http://www.gnu.org/licenses/>.
`
)

var (
	app      = cli.NewApp()
	addr     = port
	confPath = configPath
)

func init() {
	logrus.SetLevel(defaultLevel)
	app.Version = version
	app.Name = name
	app.Usage = usage
	app.UsageText = usage
	app.Email = email
	app.Copyright = copyright
	app.Description = Description
}

func main() {
	// set flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "info",
			Usage: "Set log level to info",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Set log level to debug",
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Value:       configPath,
			Destination: &confPath,
		},
		cli.StringFlag{
			Name:        "addr, a",
			Usage:       "Address to bind",
			Value:       ":3051",
			Destination: &addr,
		},
	}
	// prepare things up
	app.Action = func(c *cli.Context) error {
		// set log to level
		if c.Bool("info") {
			logrus.SetLevel(logrus.InfoLevel)
		}
		// set log level to debug
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		logrus.Infof("Log Level: %s", logrus.GetLevel().String())
		// load config
		conf, err := config.LoadConfig(confPath)
		if err != nil {
			log.Fatalf("Could not load config => %s", err.Error())
		}
		// create new server and run on port , Run() will block
		server.New(addr, conf).Run()
		return nil
	}
	// run with os.args
	app.Run(os.Args)
}
