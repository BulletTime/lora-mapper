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
	"github.com/paulmach/go.geojson"
	"fmt"
	"github.com/pkg/errors"
)

type gjson struct {
	callbackName string
	featureCollection *geojson.FeatureCollection
}

type GeoJSON interface {
	GetGeoJSON() (string, error)
	GetGeoJSONFromSF(string) (string, error)
	GetGeoJSONFromGateway(string) (string, error)
	GetGeoJSONFromGatewayAndSF(string, string) (string, error)
}

func NewGeoJSON(callbackName string) GeoJSON {
	return &gjson{
		callbackName: callbackName,
		featureCollection: geojson.NewFeatureCollection(),
	}
}

func (g *gjson) wrapCallbackFunction() (string, error) {
	json, err := g.featureCollection.MarshalJSON()
	if err != nil {
		return "", errors.Wrap(err, "marshalling json from featurecollection")
	}

	return fmt.Sprintf("%s(%s);", g.callbackName, json), nil
}

func (g *gjson) GetGeoJSON() (string, error) {
	return "", nil
}

func (g *gjson) GetGeoJSONFromSF(sf string) (string, error) {
	return "", nil
}

func (g *gjson) GetGeoJSONFromGateway(gw string) (string, error) {
	return "", nil
}

func (g *gjson) GetGeoJSONFromGatewayAndSF(gw string, sf string) (string, error) {
	return "", nil
}
