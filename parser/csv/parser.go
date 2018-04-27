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
	"bytes"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/bullettime/lora-mapper/model"
	"github.com/bullettime/lora-mapper/parser"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	LocationData = "coverage"
	CSVFields    = 4
	CSVHeader    = "lat;lon;pwr;sf"
	DefaultSize  = 7
)

const (
	SF7 = 1 << iota
	SF8
	SF9
	SF10
	SF11
	SF12
)

var (
	ErrHeader = errors.New("invalid header")
)

type csvParser struct {
	MetricName  string
	DefaultTags map[string]string
}

func New() parser.Parser {
	metricName := viper.GetString("metric.name")

	if metricName == "" {
		metricName = LocationData
	}

	p := csvParser{
		MetricName: metricName,
	}

	return &p
}

func truncate(some float64) float64 {
	return float64(int(some*10000)) / 10000
}

func (p *csvParser) getMetricsFromRecord(record []string) ([]model.Metric, error) {
	var metrics []model.Metric

	if strings.Join(record, ";") == CSVHeader {
		return metrics, nil
	}

	tags := make(map[string]string, len(p.DefaultTags))
	for k, v := range p.DefaultTags {
		tags[k] = v
	}

	lat, err := strconv.ParseFloat(record[0], 64)
	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}

	tags["latitude"] = strconv.FormatFloat(truncate(lat/10000000), 'f', -1, 64)
	tags["longitude"] = strconv.FormatFloat(truncate(lon/10000000), 'f', -1, 64)
	tags["power"] = record[2]

	fields := map[string]interface{}{
		"size": DefaultSize,
		"rssi": 0,
		"snr":  0.0,
	}

	SF, err := strconv.ParseInt(record[3], 10, 8)
	if err != nil {
		return nil, err
	}

	//if int8(SF) == 0 {
	//	metric, err := model.NewMetric(p.MetricName, tags, fields, time.Time{})
	//	if err != nil {
	//		return nil, errors.Wrap(err, "[CSVParser] error creating metric")
	//	}
	//
	//	metrics = append(metrics, metric)
	//}

	if (int8(SF) & SF7) != 0 {
		tags7 := make(map[string]string, len(tags))
		for k, v := range tags {
			tags7[k] = v
		}

		m7, err := model.NewMetric(p.MetricName, tags7, fields, time.Time{})
		if err != nil {
			return nil, errors.Wrap(err, "[CSVParser] error creating metric")
		}

		m7.AddTag("data_rate", "SF7BW125")

		metrics = append(metrics, m7)
	}

	if (int8(SF) & SF8) != 0 {
		tags8 := make(map[string]string, len(tags))
		for k, v := range tags {
			tags8[k] = v
		}

		m8, err := model.NewMetric(p.MetricName, tags8, fields, time.Time{})
		if err != nil {
			return nil, errors.Wrap(err, "[CSVParser] error creating metric")
		}

		m8.AddTag("data_rate", "SF8BW125")

		metrics = append(metrics, m8)
	}

	if (int8(SF) & SF9) != 0 {
		tags9 := make(map[string]string, len(tags))
		for k, v := range tags {
			tags9[k] = v
		}

		m9, err := model.NewMetric(p.MetricName, tags9, fields, time.Time{})
		if err != nil {
			return nil, errors.Wrap(err, "[CSVParser] error creating metric")
		}

		m9.AddTag("data_rate", "SF9BW125")

		metrics = append(metrics, m9)
	}

	if (int8(SF) & SF10) != 0 {
		tags10 := make(map[string]string, len(tags))
		for k, v := range tags {
			tags10[k] = v
		}

		m10, err := model.NewMetric(p.MetricName, tags10, fields, time.Time{})
		if err != nil {
			return nil, errors.Wrap(err, "[CSVParser] error creating metric")
		}

		m10.AddTag("data_rate", "SF10BW125")

		metrics = append(metrics, m10)
	}

	if (int8(SF) & SF11) != 0 {
		tags11 := make(map[string]string, len(tags))
		for k, v := range tags {
			tags11[k] = v
		}

		m11, err := model.NewMetric(p.MetricName, tags11, fields, time.Time{})
		if err != nil {
			return nil, errors.Wrap(err, "[CSVParser] error creating metric")
		}

		m11.AddTag("data_rate", "SF11BW125")

		metrics = append(metrics, m11)
	}

	if (int8(SF) & SF12) != 0 {
		tags12 := make(map[string]string, len(tags))
		for k, v := range tags {
			tags12[k] = v
		}

		m12, err := model.NewMetric(p.MetricName, tags12, fields, time.Time{})
		if err != nil {
			return nil, errors.Wrap(err, "[CSVParser] error creating metric")
		}

		m12.AddTag("data_rate", "SF12BW125")

		metrics = append(metrics, m12)
	}

	return metrics, nil
}

func (p *csvParser) Parse(buf []byte) ([]model.Metric, error) {
	var metrics []model.Metric

	r := csv.NewReader(bytes.NewReader(buf))
	r.Comma = ';'
	r.FieldsPerRecord = CSVFields

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.WithError(err).Warn("[CSVParser] read error")
			continue
		}

		m, err := p.getMetricsFromRecord(record)
		if err != nil {
			log.WithError(err).Warn("[CSVParser] parse metrics error")
			continue
		}

		for _, metric := range m {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func (p *csvParser) SetDefaultTags(tags map[string]string) {
	p.DefaultTags = tags
}
