package format

import (
	"github.com/nuclio/nuclio-sdk-go"
	"github.com/v3io/v3io-tsdb/pkg/tsdb"
	"github.com/v3io/v3io-tsdb/pkg/utils"
	"sort"
)

var DEFAULT_FORMAT string = "DEFAULT"
var TCOLLECTOR_FORMAT string = "TCOLLECTOR"

type InputFormat interface  {
	Ingest(tsdbAppender tsdb.Appender, event nuclio.Event) (interface{}, error)
}

func InputFormatForName(formatName string) InputFormat {
	if formatName == TCOLLECTOR_FORMAT {
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

