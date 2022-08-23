/*
Copyright 2018 Iguazio Systems Ltd.

Licensed under the Apache License, Version 2.0 (the "License") with
an addition restriction as set forth herein. You may not use this
file except in compliance with the License. You may obtain a copy of
the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing
permissions and limitations under the License.

In addition, you may not use the software for any purposes that are
illegal under applicable law, and the grant of the foregoing license
under the Apache 2.0 license is conditioned upon your compliance with
such restriction.
*/
package format

import (
	"github.com/nuclio/nuclio-sdk-go"
	"github.com/v3io/v3io-tsdb/pkg/tsdb"
	"github.com/v3io/v3io-tsdb/pkg/utils"
	"sort"
	"strings"
)

const tcollector string = "tcollector"

type Ingester interface {
	Ingest(tsdbAppender tsdb.Appender, event nuclio.Event) interface{}
}

func IngesterForName(formatName string) Ingester {
	if strings.ToLower(formatName) == tcollector {
		return tcollectorFormat{}
	} else {
		return defaultTsdb{}
	}
}

// convert map[string]string -> utils.Labels
func getLabelsFromRequest(metricName string, labelsFromRequest map[string]string) utils.Labels {

	// adding 1 for metric name
	labels := make(utils.Labels, 0, len(labelsFromRequest)+1)

	// add the metric name
	labels = append(labels, utils.Label{
		Name:  "__name__",
		Value: metricName,
	})

	for labelKey, labelValue := range labelsFromRequest {
		labels = append(labels, utils.Label{
			Name:  labelKey,
			Value: labelValue,
		})
	}

	sort.Sort(labels)

	return labels
}

func BadRequest(msg string) nuclio.Response {
	return nuclio.Response{
		StatusCode:  400,
		ContentType: "application/text",
		Body:        []byte(msg),
	}
}

func InternalError(msg string) nuclio.Response {
	return nuclio.Response{
		StatusCode:  500,
		ContentType: "application/text",
		Body:        []byte(msg),
	}
}
