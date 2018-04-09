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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/apex/log"
	"github.com/bullettime/lora-mapper/database/influxdb"
	"io/ioutil"
	"github.com/bullettime/lora-mapper/model"
	"github.com/bullettime/lora-mapper/parser/csv"
	"github.com/pkg/errors"
)

var (
	callback = "eqfeed_callback"
	output = "data_geo.json"
)

// geojsonCmd represents the geojson command
var geojsonCmd = &cobra.Command{
	Use:   "geojson",
	Short: "Create a geo jsonp file from the data",
	Long: `lora-mapper geojson creates a geo jsonp file from the data currently in the database.

This command takes one arguments:
	1. datarate [eg. SF7BW125]`,
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
			log.WithError(err).Fatal("can't connect to the influx database")
		}
		defer db.Close()

		err = writeGeoJSONFile(db, args[0])
		if err != nil {
			log.WithError(err).Fatal("can't write geojson file")
		}

		log.WithFields(log.Fields{
			"filename": output,
			"sf": args[0],
			"callback": callback,
		}).Info("geojson file written")
	},
}

func init() {
	RootCmd.AddCommand(geojsonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// geojsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// geojsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	geojsonCmd.Flags().StringVarP(&callback, "callback", "c", "eqfeed_callback", "name of the callback function")
	geojsonCmd.Flags().StringVarP(&output, "output", "o", "data_geo.json", "name of the output file")
}

func writeGeoJSONFile(db model.Database, sf string) error {
	g := model.NewGeoJSON(db, csv.LocationData, callback)

	data, err := g.GetGeoJSONFromSF(sf)
	if err != nil {
		return errors.Wrapf(err, "retrieving geojson data with sf: %s", sf)
	}

	return ioutil.WriteFile(output, []byte(data), 0644)
}
