// Digota <http://digota.com> - eCommerce microservice
// Copyright (c) 2018 Yaron Sumel <yaron@digota.com>
//
// MIT License
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

// prevents gofmt changes the imports order

// first init
import _ "github.com/digota/digota/config"

// standards imports
import (
	"log"
	"os"

	"github.com/digota/digota/config"
	"github.com/digota/digota/server"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

const (
	//defaultLevel = logrus.WarnLevel
	defaultLevel = logrus.DebugLevel
	version      = "0.1"
	name         = "digota"
	usage        = ""
	email        = "yaron@digota.com"
	description  = "eCommerce microservice"
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
	app = cli.NewApp()
)

func init() {
	logrus.SetLevel(defaultLevel)
	app.Version = version
	app.Name = name
	app.Usage = usage
	app.UsageText = usage
	app.Email = email
	app.Copyright = copyright
	app.Description = description
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
		conf, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Could not load config => %s", err.Error())
		}
		// create new server and run on port , Run() will block
		server.New(conf).Run()
		return nil
	}
	// run with os.args
	app.Run(os.Args)
}
