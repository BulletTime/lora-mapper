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

package geojson

import (
	"net/http"
	"github.com/bullettime/lora-mapper/web/utils"
	"fmt"
	"github.com/bullettime/lora-mapper/model"
	"github.com/bullettime/lora-mapper/parser/csv"
	"github.com/apex/log"
)

type Handler struct {
	geoJSON model.GeoJSON
}

func NewHandler(db model.Database) *Handler {
	// TODO add geoJSON.SetCallback and check for param callback in GET
	return &Handler{
		geoJSON: model.NewGeoJSON(db, csv.LocationData, "eqfeed_callback"),
	}
}

func (h *Handler) Handle() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var head string

		head, req.URL.Path = utils.ShiftPath(req.URL.Path)

		switch head {
		case "sf7":
			h.handleSF7().ServeHTTP(res, req)
		case "sf8":
			h.handleSF8().ServeHTTP(res, req)
		case "sf9":
			h.handleSF9().ServeHTTP(res, req)
		case "sf10":
			h.handleSF10().ServeHTTP(res, req)
		case "sf11":
			h.handleSF11().ServeHTTP(res, req)
		case "sf12":
			h.handleSF12().ServeHTTP(res, req)
		default:
			http.NotFound(res, req)
		}
	})
}

func (h *Handler) handleSF7() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		json, err := h.geoJSON.GetGeoJSONFromSF("SF7BW125")
		if err != nil {
			log.WithError(err).Error("handle sf7")
			http.NotFound(res, req)
		}

		h.writeJSON(json).ServeHTTP(res, req)
	})
}

func (h *Handler) handleSF8() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		json, err := h.geoJSON.GetGeoJSONFromSF("SF8BW125")
		if err != nil {
			log.WithError(err).Error("handle sf8")
			http.NotFound(res, req)
		}

		h.writeJSON(json).ServeHTTP(res, req)
	})
}

func (h *Handler) handleSF9() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		json, err := h.geoJSON.GetGeoJSONFromSF("SF9BW125")
		if err != nil {
			log.WithError(err).Error("handle sf9")
			http.NotFound(res, req)
		}

		h.writeJSON(json).ServeHTTP(res, req)
	})
}

func (h *Handler) handleSF10() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		json, err := h.geoJSON.GetGeoJSONFromSF("SF10BW125")
		if err != nil {
			log.WithError(err).Error("handle sf10")
			http.NotFound(res, req)
		}

		h.writeJSON(json).ServeHTTP(res, req)
	})
}

func (h *Handler) handleSF11() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		json, err := h.geoJSON.GetGeoJSONFromSF("SF11BW125")
		if err != nil {
			log.WithError(err).Error("handle sf11")
			http.NotFound(res, req)
		}

		h.writeJSON(json).ServeHTTP(res, req)
	})
}

func (h *Handler) handleSF12() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		json, err := h.geoJSON.GetGeoJSONFromSF("SF12BW125")
		if err != nil {
			log.WithError(err).Error("handle sf12")
			http.NotFound(res, req)
		}

		h.writeJSON(json).ServeHTTP(res, req)
	})
}

func (h *Handler) writeJSON(json string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		fmt.Fprint(res, json)
	})
}
