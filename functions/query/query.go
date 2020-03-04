package main

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/nuclio/nuclio-sdk-go"
	"github.com/pkg/errors"
	v3ioerrors "github.com/v3io/v3io-go/pkg/errors"
	"github.com/v3io/v3io-tsdb/pkg/config"
	"github.com/v3io/v3io-tsdb/pkg/formatter"
	"github.com/v3io/v3io-tsdb/pkg/pquerier"
	"github.com/v3io/v3io-tsdb/pkg/tsdb"
	"github.com/v3io/v3io-tsdb/pkg/utils"
)

/* Example request:
{
	"metric": "cpu",
	"step": "1m",
	"start_time": "1532095945142",
	"end_time": "1642995948517"
}
*/
type request struct {
	Metric           string   `json:"metric"`
	Aggregators      []string `json:"aggregators"`
	FilterExpression string   `json:"filter_expression"`
	Step             string   `json:"step"`
	StartTime        string   `json:"start_time"`
	EndTime          string   `json:"end_time"`
	Last             string   `json:"last"`
}

var adapter *tsdb.V3ioAdapter
var adapterLock sync.Mutex

func Query(context *nuclio.Context, event nuclio.Event) (interface{}, error) {

	request, err := validateRequest(event.GetBody())
	if err != nil {
		return nil, nuclio.WrapErrBadRequest(err)
	}

	context.Logger.DebugWith("Got query request", "request", request)

	// convert string times (unix or RFC3339 or relative like now-2h) to unix milisec times
	from, to, step, err := utils.GetTimeFromRange(request.StartTime, request.EndTime, request.Last, request.Step)
	if err != nil {
		return nil, nuclio.WrapErrBadRequest(errors.Wrap(err, "Error parsing query time range"))
	}

	// Create TSDB Querier
	querier, err := adapter.QuerierV2()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize querier")
	}

	params := &pquerier.SelectParams{
		Name:      request.Metric,
		Functions: strings.Join(request.Aggregators, ","),
		Step:      step,
		Filter:    request.FilterExpression,
		From:      from,
		To:        to,
	}

	// Select query to get back a series set iterator
	seriesSet, err := querier.Select(params)
	if err != nil {
		if e, hasErrorCode := err.(v3ioerrors.ErrorWithStatusCode); hasErrorCode && e.StatusCode() >= 400 && e.StatusCode() < 500 {
			return nil, nuclio.WrapErrBadRequest(err, context)
		}
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
	v3ioAdapterPath := os.Getenv("QUERY_V3IO_TSDB_PATH")
	if v3ioAdapterPath == "" {
		return errors.New("QUERY_V3IO_TSDB_PATH must be set")
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

		v3ioConfig, err := config.GetOrLoadFromStruct(&config.V3ioConfig{TablePath: path})
		if err != nil {
			return errors.Wrap(err, "Failed to load v3io config")
		}

		v3ioUrl := os.Getenv("QUERY_V3IO_URL")
		accessKey := os.Getenv("QUERY_V3IO_ACCESS_KEY")
		username := os.Getenv("QUERY_V3IO_USERNAME")
		password := os.Getenv("QUERY_V3IO_PASSWORD")
		containerName := os.Getenv("QUERY_V3IO_CONTAINER")
		numWorkers, err := toNumber(os.Getenv("QUERY_V3IO_NUM_WORKERS"), 8)
		if err != nil {
			return errors.Wrap(err, "Failed to get number of workers")
		}

		if containerName == "" {
			containerName = "bigdata"
		}

		container, err := tsdb.NewContainer(v3ioUrl, numWorkers, accessKey, username, password, containerName, context.Logger)
		if err != nil {
			return errors.Wrap(err, "Failed to create container")
		}

		// create adapter once for all contexts
		adapter, err = tsdb.NewV3ioAdapter(v3ioConfig, container, context.Logger)
		if err != nil {
			return errors.Wrap(err, "Failed to v3io adapter")
		}
	}

	// adapter already exists, use it
	return nil
}

func toNumber(input string, defaultValue int) (int, error) {
	if input == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(input)
}

func validateRequest(eventBody []byte) (*request, error) {

	requestMap := make(map[string]interface{})
	if err := json.Unmarshal(eventBody, &requestMap); err != nil {
		return nil, errors.Wrap(err, "Request body must be a valid JSON object")
	}

	request := request{}
	var err error
	request.Metric, err = validateOptionalString(requestMap, "metric")
	if err != nil {
		return nil, err
	}
	request.FilterExpression, err = validateOptionalString(requestMap, "filter_expression")
	if err != nil {
		return nil, err
	}
	request.Step, err = validateOptionalString(requestMap, "step")
	if err != nil {
		return nil, err
	}
	request.StartTime, err = validateOptionalString(requestMap, "start_time")
	if err != nil {
		return nil, err
	}
	request.EndTime, err = validateOptionalString(requestMap, "end_time")
	if err != nil {
		return nil, err
	}
	request.Last, err = validateOptionalString(requestMap, "last")
	if err != nil {
		return nil, err
	}
	if request.Last != "" && (request.StartTime != "" || request.EndTime != "") {
		return nil, errors.New("'last' field must not be used in conjunction with 'start_time' or 'end_time'")
	}
	if request.Metric == "" && request.FilterExpression == "" {
		return nil, errors.New("Request must contain either a 'metric' field or 'filter_expression' field")
	}
	aggregators, ok := requestMap["aggregators"]
	delete(requestMap, "aggregators")
	if ok {
		request.Aggregators, ok = aggregators.([]string)
		if !ok {
			return nil, errors.New("'aggregators' field must be a string array")
		}
	}
	var unsupportedFields []string
	for key, _ := range requestMap {
		unsupportedFields = append(unsupportedFields, key)
	}
	if len(unsupportedFields) > 0 {
		return nil, errors.Errorf("Request must not contain unsupported fields: %s", strings.Join(unsupportedFields, ", "))
	}
	return &request, nil
}

func validateOptionalString(requestMap map[string]interface{}, fieldName string) (string, error) {
	value, ok := requestMap[fieldName]
	delete(requestMap, fieldName)
	if !ok {
		return "", nil
	}
	str, ok := value.(string)
	if !ok {
		return "", errors.Errorf("'%s' field must be a string", fieldName)
	}
	return str, nil
}
