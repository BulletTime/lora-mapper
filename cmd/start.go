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
	"fmt"

	"github.com/apex/log"
	"github.com/bullettime/lora-mapper/daemon"
	"github.com/bullettime/lora-mapper/database/influxdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start web server",
	Long: `lora-mapper start will run the web server to view to coverage mappings`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")

		dbOptions := influxdb.InfluxOptions{
			Server:    viper.GetString("influxdb.server.url"),
			Username:  viper.GetString("influxdb.server.username"),
			Password:  viper.GetString("influxdb.server.password"),
			Database:  viper.GetString("influxdb.database"),
			Precision: viper.GetString("influxdb.precision"),
		}
		log.WithFields(log.Fields{
			"Server":    dbOptions.Server,
			"Username":  dbOptions.Username,
			"Database":  dbOptions.Database,
			"Precision": dbOptions.Precision,
		}).Debug("DB Options")

		server := daemon.Daemon{
			Address:          viper.GetString("web.address"),
			TLS:              viper.GetBool("web.tls"),
			CertFileLocation: viper.GetString("web.certfile"),
			KeyFileLocation:  viper.GetString("web.keyfile"),
			DBOptions:        dbOptions,
		}

		if err := server.Run(); err != nil {
			log.WithError(err).Error("stopped by error")
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
