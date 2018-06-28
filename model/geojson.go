// The MIT License (MIT)
//
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

package model

import (
	"fmt"
	"strconv"

	"github.com/apex/log"
	"github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
)

const (
	InfluxSF    = `select max("rssi") as "rssi" from (select mean("rssi") as "rssi", mean("snr") as "snr" from %s where data_rate='%s' group by latitude, longitude, gateway_id) group by latitude, longitude`
	InfluxAllSF = `select distinct(data_rate) as "data_rate" from (select rssi, snr, data_rate from %s where rssi < 0 group by latitude, longitude) group by latitude, longitude`
)

type gjson struct {
	db                Database
	measurementName   string
	featureCollection *geojson.FeatureCollection
}

type GeoJSON interface {
	GetGeoJSONFromSF(string, string) (string, error)
	GetGeoJSONFromAllSF(string) (string, error)
}

func NewGeoJSON(db Database, measurementName string) GeoJSON {
	return &gjson{
		db:                db,
		measurementName:   measurementName,
		featureCollection: geojson.NewFeatureCollection(),
	}
}

func (g *gjson) getJSON(callback string) (string, error) {
	json, err := g.featureCollection.MarshalJSON()
	if err != nil {
		return "", errors.Wrap(err, "marshalling json from featurecollection")
	}

	if callback != "" {
		return fmt.Sprintf("%s(%s);", callback, json), nil
	} else {
		return fmt.Sprintf("%s", json), nil
	}
}

func (g *gjson) GetGeoJSONFromSF(sf string, callback string) (string, error) {
	g.featureCollection = geojson.NewFeatureCollection()

	command := fmt.Sprintf(InfluxSF, g.measurementName, sf)

	metrics, err := g.db.Query(command)
	if err != nil {
		return "", err
	}

	for _, series := range metrics {
		for _, metric := range series {
			if !metric.HasTag("latitude") || !metric.HasTag("longitude") || !metric.HasField("rssi") {
				log.WithField("metric", metric).Warn("invalid metric")
				continue
			}

			lat, err := strconv.ParseFloat(metric.Tags()["latitude"], 64)
			if err != nil {
				log.WithField("latitude", metric.Tags()["latitude"]).Warn("invalid latitude")
				continue
			}
			lon, err := strconv.ParseFloat(metric.Tags()["longitude"], 64)
			if err != nil {
				log.WithField("latitude", metric.Tags()["longitude"]).Warn("invalid longitude")
				continue
			}
			rssi := metric.Fields()["rssi"]

			feature := geojson.NewPointFeature([]float64{lat, lon})
			feature.SetProperty("rssi", rssi)

			g.featureCollection.AddFeature(feature)
		}
	}

	return g.getJSON(callback)
}

func (g *gjson) GetGeoJSONFromAllSF(callback string) (string, error) {
	g.featureCollection = geojson.NewFeatureCollection()

	command := fmt.Sprintf(InfluxAllSF, g.measurementName)

	metrics, err := g.db.Query(command)
	if err != nil {
		return "", err
	}

	for _, series := range metrics {
		var lat, lon float64

		sf := 12

		for _, metric := range series {
			var err error

			if !metric.HasTag("latitude") || !metric.HasTag("longitude") || !metric.HasField("data_rate") {
				log.WithField("metric", metric).Warn("invalid metric")
				continue
			}

			if lat == 0 {
				lat, err = strconv.ParseFloat(metric.Tags()["latitude"], 64)
				if err != nil {
					log.WithField("latitude", metric.Tags()["latitude"]).Warn("invalid latitude")
					continue
				}
			}

			if lon == 0 {
				lon, err = strconv.ParseFloat(metric.Tags()["longitude"], 64)
				if err != nil {
					log.WithField("latitude", metric.Tags()["longitude"]).Warn("invalid longitude")
					continue
				}
			}

			dr := metric.Fields()["data_rate"].(string)
			i, err := strconv.Atoi(dr[2 : len(dr)-5])
			if err != nil {
				log.WithField("data rate", metric.Fields()["data_rate"]).Warn("invalid data rate")
				continue
			}

			if i < sf {
				sf = i
			}
		}

		if lat > 0 && lon > 0 {
			feature := geojson.NewPointFeature([]float64{lat, lon})
			feature.SetProperty("sf", fmt.Sprintf("sf%v", sf))

			g.featureCollection.AddFeature(feature)
		}
	}

	return g.getJSON(callback)
}
