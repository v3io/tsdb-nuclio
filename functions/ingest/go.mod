module github.com/v3io/tsdb-nuclio/functions/ingest

go 1.14

require (
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/nuclio/errors v0.0.3 // indirect
	github.com/nuclio/nuclio-sdk-go v0.1.0
	github.com/nuclio/zap v0.0.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/v3io/frames v0.7.10 // indirect
	github.com/v3io/v3io-go v0.1.6 // indirect
	github.com/v3io/v3io-tsdb v0.10.3
	github.com/valyala/fasthttp v1.14.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200317114155-1f3552e48f24 // indirect
	google.golang.org/grpc v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace (
	github.com/v3io/frames => github.com/v3io/frames v0.7.10
	github.com/v3io/v3io-tsdb => github.com/v3io/v3io-tsdb v0.10.8
)
