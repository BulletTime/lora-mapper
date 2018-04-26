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

package index

import (
	"fmt"
	"net/http"
)

const indexHTML = `<!DOCTYPE HTML>
<html>
  <head>
    <meta charset="utf-8">
    <title>LoRa Coverage</title>
  </head>
  <body>
    <div id='root'>
	  <ul>
		<li><a href="/maps/coverage.html">All Spreading Factors</a></li>
	  	<li><a href="/maps/sf7bw125.html">SF7 BW125</a></li>
	  	<li><a href="/maps/sf8bw125.html">SF8 BW125</a></li>
	  	<li><a href="/maps/sf9bw125.html">SF9 BW125</a></li>
	  	<li><a href="/maps/sf10bw125.html">SF10 BW125</a></li>
	  	<li><a href="/maps/sf11bw125.html">SF11 BW125</a></li>
	  	<li><a href="/maps/sf12bw125.html">SF12 BW125</a></li>
	  </ul>
	</div>
  </body>
</html>
`

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
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

func(h *Handler) handleGet() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, indexHTML)
	})
}
