package main

import (
	"encoding/json"
	"os"
	"sort"
	"sync"

	"github.com/nuclio/nuclio-sdk-go"
	"github.com/pkg/errors"
	"github.com/v3io/v3io-go-http"
	"github.com/v3io/v3io-tsdb/pkg/config"
	"github.com/v3io/v3io-tsdb/pkg/tsdb"
	"github.com/v3io/v3io-tsdb/pkg/utils"
)

// Example event:
//
// {
//		"metric": "cpu",
//		"labels": {
//			"dc": "7",
//			"hostname": "mybesthost"
//		},
//		"samples": [
//			{
//				"t": "1532595945142",
//				"v": {
//					"N": 95.2
//				}
//			},
//			{
//				"t": "1532595948517",
//				"v": {
//					"n": 86.8
//				}
//			}
//		]
// }

type value struct {
	N float64 `json:"n,omitempty"`
}

type sample struct {
	Time  string `json:"t,omitempty"`
	Value value  `json:"v,omitempty"`
}

type request struct {
	Metric  string            `json:"metric"`
	Labels  map[string]string `json:"labels,omitempty"`
	Samples []sample          `json:"samples"`
}

type userData struct {
	tsdbAppender tsdb.Appender
}

var adapter *tsdb.V3ioAdapter
var adapterLock sync.Mutex

func Ingest(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	var request request

	// parse body
	if err := json.Unmarshal(event.GetBody(), &request); err != nil {
		return "", nuclio.WrapErrBadRequest(err)
	}

	// convert the map[string]string -> []Labels
	labels := getLabelsFromRequest(request.Metric, request.Labels)

	// get user data from context, as initialized by InitContext
	userData := context.UserData.(*userData)

	// iterate over request samples
	for _, sample := range request.Samples {

		// if time is not specified assume "now"
		if sample.Time == "" {
			sample.Time = "now"
		}

		// convert time string to time int, string can be: now, now-2h, int (unix milisec time), or RFC3339 date string
		sampleTime, err := utils.Str2unixTime(sample.Time)
		if err != nil {
			return "", errors.Wrap(err, "Failed to parse time: "+sample.Time)
		}

		// append sample to metric
		_, err = userData.tsdbAppender.Add(labels, sampleTime, sample.Value.N)
		if err != nil {
			return "", errors.Wrap(err, "Failed to add sample")
		}
	}

	return "", nil
}

// InitContext runs only once when the function runtime starts
func InitContext(context *nuclio.Context) error {
	var err error
	var userData userData

	// get configuration from env
	tsdbAppenderPath := os.Getenv("INGEST_V3IO_TSDB_PATH")
	if tsdbAppenderPath == "" {
		return errors.New("INGEST_V3IO_TSDB_PATH must be set")
	}

	context.Logger.InfoWith("Initializing", "tsdbAppenderPath", tsdbAppenderPath)

	// create TSDB appender
	userData.tsdbAppender, err = createTSDBAppender(context, tsdbAppenderPath)
	if err != nil {
		return err
	}

	// set user data into the context
	context.UserData = &userData

	return nil
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

func createTSDBAppender(context *nuclio.Context, path string) (tsdb.Appender, error) {
	context.Logger.InfoWith("Creating TSDB appender", "path", path)

	adapterLock.Lock()
	defer adapterLock.Unlock()

	if adapter == nil {
		var err error

		v3ioConfig, err := config.GetOrLoadFromStruct(&config.V3ioConfig{
			TablePath: path,
		})

		if err != nil {
			return nil, err
		}

		// create adapter once for all contexts
		adapter, err = tsdb.NewV3ioAdapter(v3ioConfig, context.DataBinding["db0"].(*v3io.Container), context.Logger)
		if err != nil {
			return nil, err
		}
	}

	tsdbAppender, err := adapter.Appender()
	if err != nil {
		return nil, err
	}

	return tsdbAppender, nil
}
