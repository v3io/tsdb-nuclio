module github.com/v3io/tsdb-nuclio/functions/ingest

go 1.14

require (
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/nuclio/errors v0.0.3 // indirect
	github.com/nuclio/nuclio-sdk-go v0.1.0
	github.com/nuclio/zap v0.0.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/v3io/v3io-tsdb v0.10.12
	github.com/valyala/fasthttp v1.14.0 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200204135345-fa8e72b47b90 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace (
	github.com/v3io/frames => github.com/v3io/frames v0.7.36
	github.com/v3io/v3io-tsdb => github.com/v3io/v3io-tsdb v0.11.5
)
