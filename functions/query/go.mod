module github.com/v3io/tsdb-nuclio/functions/query

go 1.14

require (
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/nuclio/errors v0.0.3 // indirect
	github.com/nuclio/nuclio-sdk-go v0.2.0
	github.com/nuclio/zap v0.0.4 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/v3io/v3io-go v0.2.5-0.20210113095419-6c806b8d5186
	github.com/v3io/v3io-tsdb v0.13.1
	golang.org/x/text v0.3.8 // indirect
	google.golang.org/genproto v0.0.0-20200204135345-fa8e72b47b90 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace (
	github.com/v3io/frames => github.com/v3io/frames v0.10.2
	github.com/v3io/v3io-go => github.com/v3io/v3io-go v0.2.3
)
