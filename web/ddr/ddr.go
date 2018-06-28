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

package ddr

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/apex/log"
	"github.com/bullettime/lora-mapper/model"
	"github.com/bullettime/lora-mapper/parser/csv"
	"github.com/bullettime/lora-mapper/web/utils"
	"github.com/spf13/viper"
)

type response struct {
	Datarate string `json:"datarate"`
}

type Handler struct {
	ddr model.DDR
}

func NewHandler(db model.Database) *Handler {
	metricName := viper.GetString("metric.name")

	if metricName == "" {
		metricName = csv.LocationData
	}

	radius := viper.GetFloat64("ddr.radius")

	if radius <= 0 {
		radius = 100.0
	}

	return &Handler{
		ddr: model.NewDDR(db, metricName, radius),
	}
}

func (h *Handler) Handle() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			h.handleGet().ServeHTTP(res, req)
		default:
			http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
}

func (h *Handler) handleGet() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var head string

		head, req.URL.Path = utils.ShiftPath(req.URL.Path)

		switch head {
		case "q":
			h.handleDDR().ServeHTTP(res, req)
		default:
			http.NotFound(res, req)
		}
	})
}

func (h *Handler) handleDDR() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		lat, err := strconv.ParseFloat(req.FormValue("lat"), 64)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"lat": req.FormValue("lat"),
			}).Error("handleDDR")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		lon, err := strconv.ParseFloat(req.FormValue("lon"), 64)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"lon": req.FormValue("lon"),
			}).Error("handleDDR")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		location := model.LatLon{Latitude: lat, Longitude: lon}

		sf, err := h.ddr.GetSF(location)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"lat": lat,
				"lon": lon,
			}).Error("handleDDR")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(response{sf})
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"response": response{sf},
			}).Error("handleDDR")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(js)
	})
}
