package main

import (
	"bytes"
	"encoding/json"
	"github.com/v3io/v3io-tsdb/pkg/formatter"
	"os"
	"strings"
	"sync"

	"github.com/nuclio/nuclio-sdk-go"
	"github.com/pkg/errors"
	"github.com/v3io/v3io-go-http"
	"github.com/v3io/v3io-tsdb/pkg/config"
	"github.com/v3io/v3io-tsdb/pkg/tsdb"
	"github.com/v3io/v3io-tsdb/pkg/utils"
)

// Example request:
//
// {
//     "metric": "cpu",
//     "step": "1m",
//     "from": "1537724620",
//     "to": "1537724630"
// }

type request struct {
	Metric      string   `json:"metric"`
	Aggregators []string `json:"aggregators"`
	Step        string   `json:"step"`
	Filter      string   `json:"filter"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Last        string   `json:"last"`
}

var adapter *tsdb.V3ioAdapter
var adapterLock sync.Mutex

func Query(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	request := request{}

	// try to unmarshal the request. return bad request if failed
	if err := json.Unmarshal(event.GetBody(), &request); err != nil {
		return nil, nuclio.WrapErrBadRequest(err)
	}

	context.Logger.DebugWith("Got query request", "request", request)

	// convert string times (unix or RFC3339 or relative like now-2h) to unix milisec times
	from, to, step, err := utils.GetTimeFromRange(request.From, request.To, request.Last, request.Step)
	if err != nil {
		return nil, nuclio.WrapErrBadRequest(errors.Wrap(err, "Error parsing query time range"))
	}

	// Create TSDB Querier
	querier, err := adapter.Querier(nil, from, to)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize querier")
	}

	// Select query to get back a series set iterator
	seriesSet, err := querier.Select(request.Metric, strings.Join(request.Aggregators, ","), step, request.Filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute query select")
	}

	// convert SeriesSet to JSON (Grafana simpleJson format)
	jsonFormatter, err := formatter.NewFormatter("json", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start json formatter")
	}

	var buffer bytes.Buffer
	err = jsonFormatter.Write(&buffer, seriesSet)

	return buffer.String(), err
}

// InitContext runs only once when the function runtime starts
func InitContext(context *nuclio.Context) error {

	// get configuration from env
	v3ioAdapterPath := os.Getenv("INGEST_V3IO_TSDB_PATH")
	if v3ioAdapterPath == "" {
		return errors.New("INGEST_V3IO_TSDB_PATH must be set")
	}

	context.Logger.InfoWith("Initializing", "v3ioAdapterPath", v3ioAdapterPath)

	// create v3io adapter
	return createV3ioAdapter(context, v3ioAdapterPath)
}

func createV3ioAdapter(context *nuclio.Context, path string) error {
	context.Logger.InfoWith("Creating v3io adapter", "path", path)

	adapterLock.Lock()
	defer adapterLock.Unlock()

	if adapter == nil {
		var err error

		v3ioConfig := config.V3ioConfig{}
		config.InitDefaults(&v3ioConfig)
		v3ioConfig.Path = path

		// create adapter once for all contexts
		adapter, err = tsdb.NewV3ioAdapter(&v3ioConfig,
			context.DataBinding["db0"].(*v3io.Container),
			context.Logger)

		return err
	}

	// adapter already exists, use it
	return nil
}
