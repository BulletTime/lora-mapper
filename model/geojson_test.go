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

package model_test

import (
	"github.com/bullettime/lora-mapper/model"
	"github.com/pkg/errors"
	"strings"
	"testing"
	"time"
)

type mockingDB struct{}

func (db mockingDB) Connect() error {
	return nil
}

func (db mockingDB) Write(metric []model.Metric) error {
	return nil
}

func (db mockingDB) QueryMeasurement(measurement string) ([]model.Metric, error) {
	if measurement != measurementName {
		return nil, errors.New("measurement not found")
	}

	return metrics, nil
}

func (db mockingDB) QueryMeasurementWithFilter(measurement string, filter string) ([]model.Metric, error) {
	if measurement != measurementName {
		return nil, errors.New("measurement not found")
	}

	if filter == "data_rate = SF12BW125" {
		return metrics[:1], nil
	} else if filter == "data_rate = SF11BW125" {
		return metrics[1:], nil
	} else if filter == "gateway_id = eui-008000000000b88d" {
		return metrics[:1], nil
	} else if filter == "gateway_id = eui-e4a7a0ffffd4bbaa" {
		return metrics[1:], nil
	} else if filter == "gateway_id = eui-008000000000b88d and data_rate = SF12BW125" {
		return metrics[:1], nil
	} else if filter == "gateway_id = eui-e4a7a0ffffd4bbaa and data_rate = SF11BW125" {
		return metrics[1:], nil
	}

	return nil, errors.New("no measurements")
}

func (db mockingDB) QueryMeasurementWithMaxAge(measurement string, after string) ([]model.Metric, error) {
	return nil, nil
}

func (db mockingDB) QueryMeasurementWithMaxAgeAndFilter(measurement string, filter string, after string) ([]model.Metric, error) {
	return nil, nil
}

func (db mockingDB) HasMetric(metric model.Metric, after time.Time) bool {
	return false
}

func (db mockingDB) Close() error {
	return nil
}

var (
	db      model.Database
	gjson   model.GeoJSON
	metrics []model.Metric

	measurementName = "testing"
	callbackName    = "testCallback"
)

func setup() {
	db = mockingDB{}
	gjson = model.NewGeoJSON(db, measurementName, callbackName)
	metrics = make([]model.Metric, 0)

	tags := map[string]string{
		"data_rate":  "SF12BW125",
		"device_id":  "sodaq_one_test",
		"frequency":  "867.7",
		"gateway_id": "eui-008000000000b88d",
		"latitude":   "51.0019",
		"longitude":  "4.7135",
		"power":      "1",
	}
	fields := map[string]interface{}{
		"size": 7,
		"rssi": -117,
		"snr":  -1.8,
	}
	t, _ := time.Parse(time.RFC3339, "2018-03-30T13:26:03.978Z")
	metric, _ := model.NewMetric(measurementName, tags, fields, t)
	metrics = append(metrics, metric)

	tags = map[string]string{
		"data_rate":  "SF11BW125",
		"device_id":  "sodaq_one_test",
		"frequency":  "867.5",
		"gateway_id": "eui-e4a7a0ffffd4bbaa",
		"latitude":   "51.0019",
		"longitude":  "4.7135",
		"power":      "1",
	}
	fields = map[string]interface{}{
		"size": 7,
		"rssi": -118,
		"snr":  -3.2,
	}
	t, _ = time.Parse(time.RFC3339, "2018-03-30T13:26:07.146Z")
	metric, _ = model.NewMetric(measurementName, tags, fields, t)
	metrics = append(metrics, metric)
}

func TestGjson_GetGeoJSON(t *testing.T) {
	setup()

	json, err := gjson.GetGeoJSON()
	if err != nil {
		t.Fatal(err)
	}

	if !(strings.HasPrefix(json, callbackName+"(") && strings.HasSuffix(json, ");")) {
		t.Error("callback function isn't wrapped around the json code")
	}

	t.Log(json)
}

func TestGjson_GetGeoJSONFromSF(t *testing.T) {
	setup()

	json, err := gjson.GetGeoJSONFromSF("SF12BW125")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(json, `"rssi":-117`) {
		t.Error("probably wrong sf..")
	}

	if !(strings.HasPrefix(json, callbackName+"(") && strings.HasSuffix(json, ");")) {
		t.Error("callback function isn't wrapped around the json code")
	}

	t.Log(json)

	json, err = gjson.GetGeoJSONFromSF("SF11BW125")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(json, `"rssi":-118`) {
		t.Error("probably wrong sf..")
	}

	if !(strings.HasPrefix(json, callbackName+"(") && strings.HasSuffix(json, ");")) {
		t.Error("callback function isn't wrapped around the json code")
	}

	t.Log(json)

	json, err = gjson.GetGeoJSONFromSF("SF10BW125")
	if err == nil {
		t.Fatal("should give error")
	}

	t.Log(err)
}

func TestGjson_GetGeoJSONFromGateway(t *testing.T) {
	setup()

	json, err := gjson.GetGeoJSONFromGateway("eui-008000000000b88d")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(json, `"rssi":-117`) {
		t.Error("probably wrong sf..")
	}

	if !(strings.HasPrefix(json, callbackName+"(") && strings.HasSuffix(json, ");")) {
		t.Error("callback function isn't wrapped around the json code")
	}

	t.Log(json)

	json, err = gjson.GetGeoJSONFromGateway("eui-e4a7a0ffffd4bbaa")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(json, `"rssi":-118`) {
		t.Error("probably wrong sf..")
	}

	if !(strings.HasPrefix(json, callbackName+"(") && strings.HasSuffix(json, ");")) {
		t.Error("callback function isn't wrapped around the json code")
	}

	t.Log(json)

	json, err = gjson.GetGeoJSONFromGateway("eui-e4a7a0ffffd4dddd")
	if err == nil {
		t.Fatal("should give error")
	}

	t.Log(err)
}

func TestGjson_GetGeoJSONFromGatewayAndSF(t *testing.T) {
	setup()

	json, err := gjson.GetGeoJSONFromGatewayAndSF("eui-008000000000b88d", "SF12BW125")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(json, `"rssi":-117`) {
		t.Error("probably wrong sf..")
	}

	if !(strings.HasPrefix(json, callbackName+"(") && strings.HasSuffix(json, ");")) {
		t.Error("callback function isn't wrapped around the json code")
	}

	t.Log(json)

	json, err = gjson.GetGeoJSONFromGatewayAndSF("eui-e4a7a0ffffd4bbaa", "SF11BW125")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(json, `"rssi":-118`) {
		t.Error("probably wrong sf..")
	}

	if !(strings.HasPrefix(json, callbackName+"(") && strings.HasSuffix(json, ");")) {
		t.Error("callback function isn't wrapped around the json code")
	}

	t.Log(json)

	json, err = gjson.GetGeoJSONFromGatewayAndSF("eui-e4a7a0ffffd4bbaa", "SF12BW125")
	if err == nil {
		t.Fatal("should give error")
	}

	t.Log(err)

	json, err = gjson.GetGeoJSONFromGatewayAndSF("eui-008000000000b88d", "SF11BW125")
	if err == nil {
		t.Fatal("should give error")
	}

	t.Log(err)
}
