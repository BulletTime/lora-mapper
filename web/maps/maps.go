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

package maps

import (
	"net/http"

	"github.com/spf13/viper"
)

type Handler struct {
	BaseURL string
	Assets  string
}

func NewHandler(base string) *Handler {
	assets := viper.GetString("web.Assets")
	if assets == "" {
		assets = "."
	}
	return &Handler{
		BaseURL: base,
		Assets:  assets,
	}
}

func (h *Handler) Handle() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/":
			http.NotFound(res, req)
		default:
			req.URL.Path = req.RequestURI
			h.handleSubfiles().ServeHTTP(res, req)
		}
	})
}

func (h *Handler) handleSubfiles() http.Handler {
		return http.StripPrefix(h.BaseURL, http.FileServer(http.Dir(h.Assets)))
}
