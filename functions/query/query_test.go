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
	suite.Require().NoError(err)
	suite.Require().Equal(expected, req)
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
	suite.Require().Error(err)
	expectedErrorMessage := "'aggregators' field must be an array"
	suite.Require().Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestBadFieldName() {
	requestString := `{
    	"m3tric": "cpu",
    	"filter_expression": "1==1"
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Require().Error(err)
	expectedErrorMessage := "Request must not contain unsupported fields: m3tric"
	suite.Require().Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestBadFieldType() {
	requestString := `{
    	"metric": "cpu",
    	"start_time": 1542111395000
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Require().Error(err)
	expectedErrorMessage := "'start_time' field must be a string"
	suite.Require().Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestMinimal() {
	requestString := `{
    	"metric": "cpu"
	}`
	expected := &request{
		Metric: "cpu",
	}
	req, err := validateRequest([]byte(requestString))
	suite.Require().NoError(err)
	suite.Require().Equal(expected, req)
}

func (suite *QueryTestSuite) TestValidateRequestWithoutMetric() {
	requestString := `{
		"filter_expression": "1==1",
    	"end_time": "1542111395000"
	}`
	expected := &request{
		FilterExpression: "1==1",
		EndTime:          "1542111395000",
	}
	req, err := validateRequest([]byte(requestString))
	suite.Require().NoError(err)
	suite.Require().Equal(expected, req)
}

func (suite *QueryTestSuite) TestValidateRequestIntAggregators() {
	requestString := `{
    	"metric": "cpu",
		"aggregators": [1, 2, 3] 
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Require().Error(err)
	expectedErrorMessage := "'aggregators' array must contain only strings"
	suite.Require().Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestStringAggregators() {
	requestString := `{
		"metric": "cpu",
		"aggregators": ["a", "b", "c"]
    }`
	expected := &request{
		Metric:      "cpu",
		Aggregators: []string{"a", "b", "c"},
	}
	req, err := validateRequest([]byte(requestString))
	suite.Require().NoError(err)
	suite.Require().Equal(expected, req)
}

func (suite *QueryTestSuite) TestValidateRequestUnsupportedField() {
	requestString := `{
    	"metric": "cpu",
		"3nd_t1me": "1542111395000"
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Require().Error(err)
	expectedErrorMessage := "Request must not contain unsupported fields: 3nd_t1me"
	suite.Require().Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateRequestLastAndEndTime() {
	requestString := `{
		"metric": "cpu",
    	"last": "123",
		"end_time": "1542111395000"
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Require().Error(err)
	expectedErrorMessage := "'last' field must not be used in conjunction with 'start_time' or 'end_time'"
	suite.Require().Equal(expectedErrorMessage, err.Error())
}

func (suite *QueryTestSuite) TestValidateEmptyRequest() {
	requestString := `{
	}`
	_, err := validateRequest([]byte(requestString))
	suite.Require().Error(err)
	expectedErrorMessage := "Request must contain either a 'metric' field or 'filter_expression' field"
	suite.Require().Equal(expectedErrorMessage, err.Error())
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
