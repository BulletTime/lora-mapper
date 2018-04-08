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
	"github.com/bullettime/lora-mapper/database"
	"github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
	"strconv"
)

type gjson struct {
	db                database.Database
	callbackName      string
	measurementName   string
	featureCollection *geojson.FeatureCollection
}

type GeoJSON interface {
	GetGeoJSON() (string, error)
	GetGeoJSONFromSF(string) (string, error)
	GetGeoJSONFromGateway(string) (string, error)
	GetGeoJSONFromGatewayAndSF(string, string) (string, error)
}

func NewGeoJSON(db database.Database, measurementName string, callbackName string) GeoJSON {
	return &gjson{
		db:                db,
		callbackName:      callbackName,
		measurementName:   measurementName,
		featureCollection: geojson.NewFeatureCollection(),
	}
}

func (g *gjson) addPoints(metrics []Metric) error {
	for _, metric := range metrics {
		if !metric.HasTag("latitude") || !metric.HasTag("longitude") || !metric.HasTag("rssi") {
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

		feature := geojson.NewPointFeature([]float64{lat, lon})
		feature.SetProperty("rssi", rssi)

		g.featureCollection.AddFeature(feature)
	}

	return nil
}

func (g *gjson) wrapCallbackFunction() (string, error) {
	json, err := g.featureCollection.MarshalJSON()
	if err != nil {
		return "", errors.Wrap(err, "marshalling json from featurecollection")
	}

	return fmt.Sprintf("%s(%s);", g.callbackName, json), nil
}

func (g *gjson) GetGeoJSON() (string, error) {
	metrics, err := g.db.QueryMeasurement(g.measurementName)
	if err != nil {
		return "", err
	}

	if err := g.addPoints(metrics); err != nil {
		return "", err
	}

	return g.wrapCallbackFunction()
}

func (g *gjson) GetGeoJSONFromSF(sf string) (string, error) {
	filter := fmt.Sprintf("data_rate = %s", sf)

	metrics, err := g.db.QueryMeasurementWithFilter(g.measurementName, filter)
	if err != nil {
		return "", err
	}

	if err := g.addPoints(metrics); err != nil {
		return "", err
	}

	return g.wrapCallbackFunction()
}

func (g *gjson) GetGeoJSONFromGateway(gw string) (string, error) {
	filter := fmt.Sprintf("gateway_id = %s", gw)

	metrics, err := g.db.QueryMeasurementWithFilter(g.measurementName, filter)
	if err != nil {
		return "", err
	}

	if err := g.addPoints(metrics); err != nil {
		return "", err
	}

	return g.wrapCallbackFunction()
}

func (g *gjson) GetGeoJSONFromGatewayAndSF(gw string, sf string) (string, error) {
	filter := fmt.Sprintf("gateway_id = %s and data_rate = %s ", gw, sf)

	metrics, err := g.db.QueryMeasurementWithFilter(g.measurementName, filter)
	if err != nil {
		return "", err
	}

	if err := g.addPoints(metrics); err != nil {
		return "", err
	}

	return g.wrapCallbackFunction()
}
