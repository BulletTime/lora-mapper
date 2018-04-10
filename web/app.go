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

package web

import (
	"net/http"

	"github.com/bullettime/lora-mapper/model"
	"github.com/bullettime/lora-mapper/web/adapter"
	"github.com/bullettime/lora-mapper/web/geojson"
	"github.com/bullettime/lora-mapper/web/index"
	"github.com/bullettime/lora-mapper/web/utils"
	"net"
	"time"
)

type App struct {
	IndexHandler   *index.Handler
	GeoJSONHandler *geojson.Handler
}

func (h *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string

	if req.URL.Path == "/" {
		adapter.Adapt(h.IndexHandler.Handle(), adapter.Log()).ServeHTTP(res, req)
	} else {
		head, req.URL.Path = utils.ShiftPath(req.URL.Path)

		switch head {
		case "geojson":
			adapter.Adapt(h.GeoJSONHandler.Handle(), adapter.Log()).ServeHTTP(res, req)
		default:
			http.NotFound(res, req)
		}
	}
}

func Start(listener net.Listener, db model.Database) {
	server := &http.Server{
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	app := &App{
		IndexHandler:   index.NewHandler(),
		GeoJSONHandler: geojson.NewHandler(db),
	}

	http.Handle("/", app)

	go server.Serve(listener)
}
