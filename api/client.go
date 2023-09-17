package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	tblApiVersion = "/api/now/"
)

// Client structure
type Client struct {
	client *http.Client

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

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.client = &http.Client{Transport: transport}
		c.debug = debug

	}

	return c

}

func (c *Client) API(method string, apiHeader map[string][]string, snowApiEndPoint string, snowApiParams url.Values, apiBody io.Reader) (apiresp interface{}, apirespcode int, apierr error) {

	if c.debug {
		fmt.Println("\nAPI(", method, ",", apiHeader, ",", snowApiEndPoint, ", ", snowApiParams, ", ", apiBody, ")")
	}

	apiURL := c.snowInstance + tblApiVersion + snowApiEndPoint
	if snowApiParams != nil {
		apiURL = apiURL + "?" + snowApiParams.Encode()
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
		log.Fatal("HTTP request creation error:", httpReqErr)
	}
	if c.debug {
		fmt.Println("Request URL: " + apiURL)
	}

	httpReq.SetBasicAuth(c.snowUsername, c.snowPassword)
	if c.debug {
		fmt.Println(apiHeader)
		fmt.Println(httpReq.URL)
	}

	statusCode := 0
	httpResp, httpRespErr := c.client.Do(httpReq)
	if httpRespErr != nil {
		log.Fatal("HTTP request error:", httpRespErr)
	} else {
		defer httpResp.Body.Close()

		// Access the HTTP status code
		statusCode := httpResp.StatusCode

		if statusCode == 200 || statusCode == 201 {
			return httpResp, statusCode, nil
		}
	}
	return nil, statusCode, httpRespErr
}

// api/client.go:114:32: httpResp.Body undefined (type interface{} has no field or method Body)
// api/client.go:116:25: httpResp.StatusCode undefined (type interface{} has no field or method StatusCode)
// api/client.go:119:15: httpResp.StatusCode undefined (type interface{} has no field or method StatusCode)
// api/client.go:120:30: httpResp.StatusCode undefined (type interface{} has no field or method StatusCode)
// api/client.go:125:23: httpResp.StatusCode undefined (type interface{} has no field or method StatusCode)

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

func GetJson(body []byte) (jsonSource interface{}) {
	if string(body) != "" && body != nil {
		err := json.Unmarshal(body, &jsonSource)
		if err != nil {
			log.Print("template executing error: ", err)
		}
	}
	return
}

func ReturnResponseBody(httpResponse *http.Response) (response []byte) {
	if httpResponse.ContentLength != 0 {
		contents, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			log.Fatal("%s", err)
		}
		return contents
	}
	return []byte("")
}
