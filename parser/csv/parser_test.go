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

package csv

import (
	"encoding/csv"
	"testing"
)

const (
	csvData = `lat;lon;pwr;sf
508609281;46818589;1;0
508629196;46837878;1;63
508632782;46846425;1;32
508632408;46856960;1;24
`
	csvData2 = `lat;lon;pwr;sf;start
508609281;46818589;1;0;2018-04-01T12:06:00Z
508629196;46837878;1;63;2018-04-01T12:06:00Z
508632782;46846425;1;32;2018-04-01T12:06:00Z
508632408;46856960;1;24;2018-04-01T12:06:00Z
`
)

func TestNew(t *testing.T) {
	p := New()

	if p.(*csvParser).MetricName != LocationData {
		t.Error("metric name should be set automatically")
	}
}

func TestCsvParser_SetDefaultTags(t *testing.T) {
	p := New()

	tags := map[string]string{
		"test": "a",
	}

	p.SetDefaultTags(tags)

	if v, ok := p.(*csvParser).DefaultTags["test"]; !ok {
		t.Error("default tags is missing key 'test'")
	} else {
		if v != "a" {
			t.Error("default tags has wrong value for key 'test'")
		}
	}
}

func TestCsvParser_Parse(t *testing.T) {
	i := 0
	p := New()

	metrics, err := p.Parse([]byte(csvData))
	if err != nil {
		t.Fatal(err)
	}

	if len(metrics) != 10 {
		t.Error("there should be 10 metrics")
	}

	if metrics[i].HasTag("data_rate") {
		t.Error("the first metric should not have set a data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF7BW125" {
		t.Error("expected different data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF8BW125" {
		t.Error("expected different data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF9BW125" {
		t.Error("expected different data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF10BW125" {
		t.Error("expected different data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF11BW125" {
		t.Error("expected different data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF12BW125" {
		t.Error("expected different data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF12BW125" {
		t.Error("expected different data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF10BW125" {
		t.Error("expected different data rate")
	}

	i++

	if metrics[i].Tags()["data_rate"] != "SF11BW125" {
		t.Error("expected different data rate")
	}
}

func TestCsvParser_Parse2(t *testing.T) {
	p := New()

	_, err := p.Parse([]byte(csvData2))
	if err == nil {
		t.Fatalf("expecting error: %s", csv.ErrFieldCount)
	}
}
