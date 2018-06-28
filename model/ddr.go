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
	"bytes"
	"container/list"
	"fmt"
	"math"
	"strconv"

	"github.com/apex/log"
)

const (
	REarth                 = 6373.0
	LatitudeSecondInMeters = LatitudeMinuteInMeters / 60.0
	LatitudeMinuteInMeters = 1853.0
	LatitudeDegreeInMeters = LatitudeMinuteInMeters * 60.0
	floatToIntPrecision    = 10000
)

const (
	InfluxDDR = `select distinct(data_rate) as "data_rate" from (select rssi, snr, data_rate from %s where rssi < 0 and latitude=~/%s/ and longitude=~/%s/ group by latitude, longitude) group by latitude, longitude`
)

type ddr struct {
	db              Database
	measurementName string
	minRadius       float64
}

type LatLon struct {
	Latitude  float64
	Longitude float64
}

type DDR interface {
	GetSF(lon LatLon) (string, error)
}

func NewDDR(db Database, measurementName string, minRadius float64) DDR {
	return &ddr{
		db:              db,
		measurementName: measurementName,
		minRadius:       minRadius,
	}
}

func (d *ddr) GetSF(ll LatLon) (string, error) {
	var datarates = make(map[LatLon]int)

	top, bottom := d.getBoundsTopBottom(ll)
	latitude := getRegex(bottom, top)

	left, right := d.getBoundsLeftRight(ll)
	longitude := getRegex(left, right)

	command := fmt.Sprintf(InfluxDDR, d.measurementName, latitude, longitude)

	metrics, err := d.db.Query(command)
	if err != nil {
		return "", err
	}

	// TODO calculate distance to get weighted best datarate

	for _, series := range metrics {
		var lat, lon float64

		sf := 12

		for _, metric := range series {
			var err error

			if !metric.HasTag("latitude") || !metric.HasTag("longitude") || !metric.HasField("data_rate") {
				log.WithField("metric", metric).Warn("invalid metric")
				continue
			}

			if lat == 0 {
				lat, err = strconv.ParseFloat(metric.Tags()["latitude"], 64)
				if err != nil {
					log.WithField("latitude", metric.Tags()["latitude"]).Warn("invalid latitude")
					continue
				}
			}

			if lon == 0 {
				lon, err = strconv.ParseFloat(metric.Tags()["longitude"], 64)
				if err != nil {
					log.WithField("latitude", metric.Tags()["longitude"]).Warn("invalid longitude")
					continue
				}
			}

			dr := metric.Fields()["data_rate"].(string)
			i, err := strconv.Atoi(dr[2 : len(dr)-5])
			if err != nil {
				log.WithField("data rate", metric.Fields()["data_rate"]).Warn("invalid data rate")
				continue
			}

			if i < sf {
				sf = i
			}
		}

		if lat > 0 && lon > 0 {
			location := LatLon{lat, lon}
			datarates[location] = sf
		}
	}

	bestDatarate := 12

	for _, dr := range datarates {
		if dr < bestDatarate {
			bestDatarate = dr
		}
	}

	return fmt.Sprintf("SF%dBW125", bestDatarate), nil
}

func radians(degree float64) float64 {
	return degree * math.Pi / 180
}

func (d *ddr) getBoundsTopBottom(ll LatLon) (top float64, bottom float64) {
	difference := d.minRadius / LatitudeDegreeInMeters

	top = ll.Latitude + difference
	bottom = ll.Latitude - difference

	return
}

func (d *ddr) getBoundsLeftRight(ll LatLon) (left float64, right float64) {
	difference := d.minRadius / (LatitudeDegreeInMeters * math.Cos(radians(ll.Latitude)))

	left = ll.Longitude - difference
	right = ll.Longitude + difference

	return
}

func getRegex(start, end float64) string {
	var result bytes.Buffer

	left := leftBounds(int(start*floatToIntPrecision), int(end*floatToIntPrecision))
	lastLeft := left.Remove(left.Back()).(*Range)

	right := rightBounds(lastLeft.Start, int(end*floatToIntPrecision))
	firstRight := right.Remove(right.Front()).(*Range)

	merged := list.New()
	merged.PushBackList(left)

	if !lastLeft.Overlaps(firstRight) {
		merged.PushBack(lastLeft)
		merged.PushBack(firstRight)
	} else {
		merged.PushBack(Join(lastLeft, firstRight))
	}

	merged.PushBackList(right)

	for e := merged.Front(); e != nil; e = e.Next() {
		if result.Len() != 0 {
			result.WriteByte('|')
		}
		result.WriteString(e.Value.(*Range).Regex(floatToIntPrecision))
	}

	return result.String()
}

func leftBounds(start, end int) *list.List {
	left := list.New()

	for start < end {
		r := NewRangeFromStart(start)
		left.PushBack(r)
		start = r.End + 1
	}

	return left
}

func rightBounds(start, end int) *list.List {
	right := list.New()

	for start < end {
		r := NewRangeFromEnd(end)
		right.PushFront(r)
		end = r.Start - 1
	}

	return right
}

type Range struct {
	Start int
	End   int
}

func NewRangeFromStart(start int) *Range {
	startS := []byte(strconv.Itoa(start))

	for i := len(startS) - 1; i >= 0; i-- {
		if startS[i] == '0' {
			startS[i] = '9'
		} else {
			startS[i] = '9'
			break
		}
	}

	end, err := strconv.Atoi(string(startS))
	if err != nil {
		end = start
	}

	return &Range{
		Start: start,
		End:   end,
	}
}

func NewRangeFromEnd(end int) *Range {
	endS := []byte(strconv.Itoa(end))

	for i := len(endS) - 1; i >= 0; i-- {
		if endS[i] == '9' {
			endS[i] = '0'
		} else {
			endS[i] = '0'
			break
		}
	}

	start, err := strconv.Atoi(string(endS))
	if err != nil {
		start = end
	}

	return &Range{
		Start: start,
		End:   end,
	}
}

func Join(r1, r2 *Range) *Range {
	return &Range{
		Start: r1.Start,
		End:   r2.End,
	}
}

func (r Range) Overlaps(r2 *Range) bool {
	return r.End > r2.Start && r2.End > r.Start
}

func (r Range) Regex(precision int) string {
	startS := []byte(strconv.FormatFloat(float64(r.Start)/float64(precision), 'f', 4, 64))
	endS := []byte(strconv.FormatFloat(float64(r.End)/float64(precision), 'f', 4, 64))

	var result bytes.Buffer

	for pos := 0; pos < len(startS); pos++ {
		if startS[pos] == endS[pos] {
			result.WriteByte(startS[pos])
		} else {
			result.WriteByte('[')
			result.WriteByte(startS[pos])
			result.WriteByte('-')
			result.WriteByte(endS[pos])
			result.WriteByte(']')
		}
	}

	return result.String()
}
