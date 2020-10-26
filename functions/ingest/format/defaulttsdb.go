package format

import (
	"encoding/json"
	"strings"

	"github.com/nuclio/nuclio-sdk-go"
	"github.com/pkg/errors"
	"github.com/v3io/v3io-tsdb/pkg/tsdb"
	"github.com/v3io/v3io-tsdb/pkg/utils"
)

/*
Example event:

[
	{
			"metric": "cpu",
			"labels": {
				"dc": "7",
				"hostname": "mybesthost"
			},
			"samples": [
				{
					"t": "1532595945142",
					"v": {
						"n": 95.2
					}
				},
				{
					"t": "1532595948517",
					"v": {
						"n": 86.8
					}
				}
			]
	}
]
*/

type value struct {
	N *float64 `json:"n"`
	S *string  `json:"s"`
}

type sample struct {
	Time  *string `json:"t"`
	Value *value  `json:"v"`
}

type request struct {
	Metric  *string           `json:"metric"`
	Labels  map[string]string `json:"labels"`
	Samples []sample          `json:"samples"`
}

//implements InputFormat
type defaultTsdb struct{}

func (Ingester defaultTsdb) Ingest(tsdbAppender tsdb.Appender, event nuclio.Event) interface{} {
	var requests []request

	body := event.GetBody()
	if body[0] == '[' {
		if err := json.Unmarshal(body, &requests); err != nil {
			return BadRequest(errors.Wrap(err, "Failed to deserialize requests").Error())
		}
	} else {
		var req request

		if err := json.Unmarshal(body, &req); err != nil {
			return BadRequest(errors.Wrap(err, "Failed to deserialize request").Error())
		}
		requests = []request{req}
	}

	for _, request := range requests {
		if request.Metric == nil {
			return BadRequest("Missing attribute: metric")
		}
		if *request.Metric == "" {
			return BadRequest("Attribute is empty: metric")
		}
		if request.Samples == nil { // if json contains an empty array, this will not be triggered
			return BadRequest("Missing attribute: samples")
		}

		// convert the map[string]string -> []Labels
		labels := getLabelsFromRequest(*request.Metric, request.Labels)

		var ref uint64
		// iterate over request samples
		for _, sample := range request.Samples {
			var value interface{}

			if sample.Time == nil {
				return BadRequest("Missing attribute in sample: t")
			}
			if sample.Value == nil {
				return BadRequest("Missing attribute in sample: v")
			}

			if sample.Value.N != nil {
				value = *sample.Value.N
			} else if sample.Value.S != nil {
				value = *sample.Value.S
			} else {
				return BadRequest("Missing attribute in sample value: 'n' or 's'")
			}

			// if time is not specified assume "now"
			var time = *sample.Time
			if time == "" {
				time = "now"
			}

			// convert time string to time int, string can be: now, now-2h, int (unix milisec time), or RFC3339 date string
			sampleTime, err := utils.Str2unixTime(time)
			if err != nil {
				return BadRequest(errors.Wrap(err, "Failed to parse time: "+time).Error())
			}

			// append sample to metric
			if ref == 0 {
				ref, err = tsdbAppender.Add(labels, sampleTime, value)
			} else {
				err = tsdbAppender.AddFast(ref, sampleTime, value)
				if err != nil && strings.Contains(err.Error(), "metric not found") {
					//retry with Add in case ref was evicted from cache
					ref, err = tsdbAppender.Add(labels, sampleTime, value)
				}
			}
			if err != nil {
				return BadRequest(errors.Wrap(err, "Failed to add sample").Error())
			}
		}
	}
	return nil
}
