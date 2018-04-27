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
	"strings"

	"github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
)

type gjson struct {
	db                Database
	measurementName   string
	featureCollection *geojson.FeatureCollection
}

type GeoJSON interface {
	GetGeoJSON(string) (string, error)
	GetGeoJSONFromSF(string, string) (string, error)
	GetGeoJSONFromGateway(string, string) (string, error)
	GetGeoJSONFromGatewayAndSF(string, string, string) (string, error)
}

func NewGeoJSON(db Database, measurementName string) GeoJSON {
	return &gjson{
		db:                db,
		measurementName:   measurementName,
		featureCollection: geojson.NewFeatureCollection(),
	}
}

func (g *gjson) addPoints(metrics []Metric) error {
	for _, metric := range metrics {
		if !metric.HasTag("latitude") || !metric.HasTag("longitude") || !metric.HasField("rssi") {
			return errors.New("invalid metric")
		}

		lat, err := strconv.ParseFloat(metric.Tags()["latitude"], 64)
		if err != nil {
			return errors.New("invalid latitude")
		}
		lon, err := strconv.ParseFloat(metric.Tags()["longitude"], 64)
		if err != nil {
			return errors.New("invalid longitude")
		}
		rssi := metric.Fields()["rssi"]
		dataRate := metric.Tags()["data_rate"]
		sf := strings.ToLower(dataRate[:len(dataRate)-5])

		feature := geojson.NewPointFeature([]float64{lat, lon})
		feature.SetProperty("rssi", rssi)
		feature.SetProperty("sf", sf)

		g.featureCollection.AddFeature(feature)
	}

	return nil
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

func (g *gjson) GetGeoJSON(callback string) (string, error) {
	g.featureCollection = geojson.NewFeatureCollection()

	metrics, err := g.db.QueryMeasurementWithFilter(g.measurementName, "rssi != 0")
	if err != nil {
		return "", err
	}

	if err := g.addPoints(metrics); err != nil {
		return "", err
	}

	return g.getJSON(callback)
}

func (g *gjson) GetGeoJSONFromSF(sf string, callback string) (string, error) {
	g.featureCollection = geojson.NewFeatureCollection()

	filter := fmt.Sprintf("data_rate = '%s'", sf)

	metrics, err := g.db.QueryMeasurementWithFilter(g.measurementName, filter)
	if err != nil {
		return "", err
	}

	if err := g.addPoints(metrics); err != nil {
		return "", err
	}

	return g.getJSON(callback)
}

func (g *gjson) GetGeoJSONFromGateway(gateway string, callback string) (string, error) {
	g.featureCollection = geojson.NewFeatureCollection()

	filter := fmt.Sprintf("gateway_id = '%s'", gateway)

	metrics, err := g.db.QueryMeasurementWithFilter(g.measurementName, filter)
	if err != nil {
		return "", err
	}

	if err := g.addPoints(metrics); err != nil {
		return "", err
	}

	return g.getJSON(callback)
}

func (g *gjson) GetGeoJSONFromGatewayAndSF(gateway string, sf string, callback string) (string, error) {
	g.featureCollection = geojson.NewFeatureCollection()

	filter := fmt.Sprintf("gateway_id = '%s' and data_rate = '%s'", gateway, sf)

	metrics, err := g.db.QueryMeasurementWithFilter(g.measurementName, filter)
	if err != nil {
		return "", err
	}

	if err := g.addPoints(metrics); err != nil {
		return "", err
	}

	return g.getJSON(callback)
}
