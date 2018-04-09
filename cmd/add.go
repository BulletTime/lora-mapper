// Copyright © 2018 Sven Agneessens <sven.agneessens@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/bullettime/lora-mapper/database/influxdb"
	"github.com/bullettime/lora-mapper/model"
	"github.com/bullettime/lora-mapper/parser/csv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	deviceID   string
	timeString string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add data from a file",
	Long: `lora-mapper add will process the data stored in a csv file format and produced
by the Sodaq-One logging device. It will only add the data that is missing in the database.
This is necessary if you want to plot a coverage map that includes the points that were
scanned and did not have reception.

This command takes one argument:
	- file name from the csv file [eg. data.csv]
It will parse the data and add the missing data to the influx database.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		influxOptions := influxdb.InfluxOptions{
			Server:    viper.GetString("influxdb.server.url"),
			Username:  viper.GetString("influxdb.server.username"),
			Password:  viper.GetString("influxdb.server.password"),
			Database:  viper.GetString("influxdb.database"),
			Precision: viper.GetString("influxdb.precision"),
		}
		log.WithFields(log.Fields{
			"Server":    influxOptions.Server,
			"Username":  influxOptions.Username,
			"Database":  influxOptions.Database,
			"Precision": influxOptions.Precision,
		}).Debug("InfluxDB Options")
		db := influxdb.New(influxOptions)

		err := db.Connect()
		if err != nil {
			panic(err)
		}

		addDataFromCSV(args[0], db)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringVar(&deviceID, "device-id", "", "adds the device id to data")
	addCmd.Flags().StringVar(&timeString, "time", "", "set the oldest time to compare data (in RFC3339 format)")
}

func addDataFromCSV(fileName string, db model.Database) {
	var metricsToAdd []model.Metric
	var t time.Time

	ctx := log.WithField("data-file", fileName)

	p := csv.New()

	if len(deviceID) > 0 {
		p.SetDefaultTags(map[string]string{
			"device_id": deviceID,
		})
	}

	if len(timeString) > 0 {
		var err error
		t, err = time.Parse(time.RFC3339, timeString)
		if err != nil {
			log.WithError(err).Fatal("parsing time")
		}
	}

	csvFile, err := os.Open(fileName)
	if err != nil {
		ctx.WithError(err).Fatal("opening csv file")
	}
	defer csvFile.Close()

	lineScanner := bufio.NewScanner(csvFile)

	for lineScanner.Scan() {
		metrics, _ := p.Parse(lineScanner.Bytes())

		for _, metric := range metrics {
			if !db.HasMetric(metric, t) {
				log.WithField("metric", metric).Debug("add metric")
				metricsToAdd = append(metricsToAdd, metric)
			}
		}
	}

	if len(metricsToAdd) > 0 {
		err := db.Write(metricsToAdd)
		if err != nil {
			log.WithError(err).Fatal("writing metrics")
		}
	}

	log.WithField("amount", len(metricsToAdd)).Info("metrics added")
}
