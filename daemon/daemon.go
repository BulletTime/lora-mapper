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

package daemon

import (
	"net"
	"github.com/pkg/errors"
	"crypto/tls"
	"os"
	"os/signal"
	"syscall"
	"github.com/apex/log"
	"github.com/bullettime/lora-mapper/web"
)

type Daemon struct {
	Address string
	Listener net.Listener
	TLS bool
	CertFileLocation string
	KeyFileLocation string
}

func (d *Daemon) Run() error {
	var err error

	if len(d.Address) == 0 {
		return errors.New("listen address not found")
	}

	log.WithField("listening address", d.Address).Info("starting listener")

	if d.TLS {
		if len(d.CertFileLocation) == 0 || len(d.KeyFileLocation) == 0 {
			return errors.New("tls key pair not found")
		}

		cert, err := tls.LoadX509KeyPair(d.CertFileLocation, d.KeyFileLocation)
		if err != nil {
			return errors.Wrap(err, "error loading tls key pair")
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}

		d.Listener, err = tls.Listen("tcp", d.Address, config)
		if err != nil {
			return errors.Wrap(err, "error starting tls listener")
		}
	} else {
		d.Listener, err = net.Listen("tcp", d.Address)
		if err != nil {
			return errors.Wrap(err, "error starting net listener")
		}
	}

	defer d.Listener.Close()

	web.Start(d.Listener)

	waitForSignal()

	return nil
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.WithField("signal", s).Warn("exiting")
}
