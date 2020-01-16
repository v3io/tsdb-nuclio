package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type QueryTestSuite struct {
	suite.Suite
}

func (suite *QueryTestSuite) TestValidateRequest() {
	requestString := `{
		"metric": "cpu",
		"step": "1m",
		"start_time": "1532095945142",
		"end_time": "1642995948517"
    }`
	expected := &request{
		Metric:    "cpu",
		Step:      "1m",
		StartTime: "1532095945142",
		EndTime:   "1642995948517",
	}
	req, err := validateRequest([]byte(requestString))
	suite.NoError(err)
	suite.Equal(expected, req)
}

func (suite *QueryTestSuite) TestValidateRequestBadAggregators() {
	requestString := `{
		"metric": "cpu",
		"aggregators": "not a json array",
		"step": "1m",
		"start_time": "1532095945142",
		"end_time": "1642995948517"
    }`
	_, err := validateRequest([]byte(requestString))
	suite.Error(err)
	expectedErrorMessage := "'aggregators' field must be a string array"
	suite.Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestBadFieldName() {
	requestString := `{
    	"M3tric": "cpu",
    	"Aggregators": ["max", "stdvar"],
    	"Start": 1542111395000,
    	"End": "now"
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Error(err)
	expectedErrorMessage := "Request object is missing 'metric' field"
	suite.Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestBadFieldType() {
	requestString := `{
    	"metric": "cpu",
    	"start_time": 1542111395000
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Error(err)
	expectedErrorMessage := "'start_time' field must be a string"
	suite.Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestMinimal() {
	requestString := `{
    	"metric": "cpu"
	}`
	expected := &request{
		Metric: "cpu",
	}
	req, err := validateRequest([]byte(requestString))
	suite.NoError(err)
	suite.Equal(expected, req)
}

func (suite *QueryTestSuite) TestValidateRequestMissingMetric() {
	requestString := `{
    	"end_time": "1542111395000"
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Error(err)
	expectedErrorMessage := "Request object is missing 'metric' field"
	suite.Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestIntAggregators() {
	requestString := `{
    	"metric": "cpu",
		"aggregators": [1, 2, 3] 
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Error(err)
	expectedErrorMessage := "'aggregators' field must be a string array"
	suite.Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestUnsupportedField() {
	requestString := `{
    	"metric": "cpu",
		"3nd_t1me": "1542111395000"
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Error(err)
	expectedErrorMessage := "Request must not contain unsupported fields: 3nd_t1me"
	suite.Equal(expectedErrorMessage, err.Error())
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
