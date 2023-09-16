package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rhysd/abspath"
	"encoding/json"
	"io/ioutil"

	"crypto/tls"
	"net/http"
)

const (
	tblApiVersion        		= "/api/now/"
)

// Client structure
type Client struct {
	httpClient := &http.Client{}

	snowInstance string
	snowUsername string
	snowPassword string

	debug bool
}

// NewClient - Creates a new client
func NewClient(snowUsername, snowPassword, snowInstance string, debug bool) *Client {

	var c *Client

	if snowInstance != "" && snowUsername != "" && snowPassword != "" {

		c = &Client{}

		if strings.Contains(snowInstance, "service-now.com") {
			c.snowInstance = snowInstance
		} else {
			c.snowInstance = "https://" + snowInstance + ".service-now.com"
		}
		c.snowUsername = snowUsername
		c.snowPassword = snowPassword

		tr  := &http.Transport{
			TLSClientConfig := &tls.Config(InsecureSkipVerify: true),
		}
		c.httpClient := http.Client{Transport: tr}

		c.debug = debug

	}

	return c

}

func (c *Client) API(method string, apiHeader map[string][]string, snowApiEndPoint string, params url.Values, apiBody io.Reader) (httpResp interface{}, httpCode int, err error) {

	if c.debug {
		fmt.Println("\nAPI(", method, ",", apiHeader, ",", snowApiEndPoint, ", ", params, ", ", apiBody, ")")
	}

	apiURL := c.snowInstance + tblApiVersion + snowApiEndPoint
	if params != nil {
		apiURL = apiURL + "?" + params.Encode()
	}
	if apiHeader == nil {
		apiHeader = url.Values{}
	}
	apiHeader["Accept"] = []string{"application/json"}
	apiHeader["Content-Type"] = []string{"application/json"}

	if method == "" {
		method = "GET"
	}
	apiMethod := strings.ToUpper(method)

	if c.debug {
		fmt.Println("======================================================================================")
		fmt.Println("API URL : ", apiURL)
		fmt.Println("API Method : ", apiMethod)
		fmt.Println("API Header : ", apiHeader)
		fmt.Println("API UserName/Password : [", c.snowUsername, "]-[", c.snowPassword, "]")
		fmt.Println("API BODY :", apiBody)
		fmt.Println("======================================================================================")
	}

        httpReq, httpReqErr := http.NewRequest(apiMethod, apiURL, apiBody)
	if httpReqErr != nil {
		fmt.Println("Failed to create http.Request object")
		return jsonResp, httpReq.StatusCode, errors.New(string(apiResponse))
	}
	if c.debug {
		fmt.Println("Request URL: " + apiURL)
	}	
	
	httpReq.SetBasicAuth(c.snowUsername, c.snowPassword)
	if c.debug {
		fmt.Println(apiHeader)
		fmt.Println(httpReq.URL)
	}	

	httpResp, httpRespErr := c.client.Do(httpReq)
	if httpRespErr != nil {
		fmt.Printf("Error(httpRespCode:%d) making %s request: %v\n", httpResp.StatusCode, apiMethod, httpRespErr)
	} else {

		jsonResp, err := utils.GetJson(httpResp)
		if err != nil {
			return jsonResp, httpResp.StatusCode, err
		}

		if httpResp.StatusCode == 200 || httpResp.StatusCode == 201 {
			return jsonResp, httpResp.StatusCode, nil
		}

	}

	return nil, httpResp.StatusCode, errors.New(string(apiResponse))

}

// snowTable - Gets the request based on method for accessiong SNOW Table API
func (c *Client) snowTable(method string, sTable string, qParams map[string]string, pParams io.Reader) ([]interface{}, error) {

	if c.debug {
		fmt.Println("\nsnowTable(", method, ",", sTable, ",", qParams, ", ", pParams, ")")
	}

	qparams := url.Values{}
	for key, val := range qParams {
		if c.debug {
			fmt.Println("qParams : ", key, " - ", val)
		}
		qparams.Set(key, val)
	}

	if c.debug {
		fmt.Println("\npParams : ", pParams, "\n")
	}

	if method == "" {
		method = "GET"
	}
	apiMethod := strings.ToUpper(method)
	if c.debug {
		fmt.Println("Method : ", apiMethod)
		fmt.Println("Table : table/", sTable)
	}

	apiResponse, _, err := c.API(apiMethod, nil, "table/"+sTable, qparams, pParams)
	if err != nil {
		fmt.Println("ERROR: ", fmt.Errorf("snowTable(table/%s) %v", sTable, err))
		return nil, fmt.Errorf("snowTable(table/%s) %v", sTable, err)
	}

	if apiResponse == nil {
		fmt.Println("ERROR: apiResonse is nil")
		return nil, nil
	}

	var sResponse []interface{}
	if apiMethod != "GET" {
		sResponse = append(sResponse, apiResponse.(map[string]interface{})["result"])
	} else {
		sResponse = apiResponse.(map[string]interface{})["result"].([]interface{})
	}

	if c.debug {
		fmt.Println("\nResponse : ", sResponse, "\n")
	}

	return sResponse, nil

}
