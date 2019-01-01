package main

import (
	"github.com/nuclio/handler/format"
	"os"
	"strconv"
	"sync"

	"github.com/nuclio/nuclio-sdk-go"
	"github.com/pkg/errors"
	"github.com/v3io/v3io-tsdb/pkg/config"
	"github.com/v3io/v3io-tsdb/pkg/tsdb"
)

type UserData struct {
	TsdbAppender tsdb.Appender
	ingester     format.Ingester
}

var adapter *tsdb.V3ioAdapter
var adapterLock sync.Mutex

func Ingest(context *nuclio.Context, event nuclio.Event) (interface{}, error) {

	// get user data from context, as initialized by InitContext
	userData := context.UserData.(*UserData)

	if err := userData.ingester.Ingest(userData.TsdbAppender, event); err != nil {
		return nil, nuclio.WrapErrBadRequest(err)
	}
	return nil, nil
}

// InitContext runs only once when the function runtime starts
func InitContext(context *nuclio.Context) error {
	var err error
	var userData UserData

	// get input format
	formatName := os.Getenv("INPUT_FORMAT")
	userData.ingester = format.IngesterForName(formatName)

	// get configuration from env
	tsdbAppenderPath := os.Getenv("INGEST_V3IO_TSDB_PATH")
	if tsdbAppenderPath == "" {
		return errors.New("INGEST_V3IO_TSDB_PATH must be set")
	}

	context.Logger.InfoWith("Initializing", "tsdbAppenderPath", tsdbAppenderPath)

	// create TSDB appender
	userData.TsdbAppender, err = createTSDBAppender(context, tsdbAppenderPath)
	if err != nil {
		return err
	}

	// set user data into the context
	context.UserData = &userData

	return nil
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
		v3ioUrl := os.Getenv("V3IO_URL")
		numWorkersStr := os.Getenv("V3IO_NUM_WORKERS")
		var numWorkers int
		if len(numWorkersStr) > 0 {
			numWorkers, err = strconv.Atoi(numWorkersStr)
			if err != nil {
				return nil, err
			}
		} else {
			numWorkers = 8
		}
		username := os.Getenv("V3IO_USERNAME")
		password := os.Getenv("V3IO_PASSWORD")
		containerName := os.Getenv("V3IO_CONTAINER")
		if containerName == "" {
			containerName = "bigdata"
		}
		container, err := tsdb.NewContainer(v3ioUrl, numWorkers, username, password, containerName, context.Logger)
		if err != nil {
			return nil, err
		}
		// create adapter once for all contexts
		adapter, err = tsdb.NewV3ioAdapter(v3ioConfig, container, context.Logger)
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
