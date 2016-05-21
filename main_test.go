package main_test

import (
    "net/http"
    "testing"
    "bytes"
    "time"
    "fmt"
    "io/ioutil"
    
    "github.com/stretchr/testify/assert"
)

type Test struct {
    description         string
    payload             string
    url                 string
    expectedBody        string
}
type Tests []Test

func TestMain(t *testing.T) {
	tests := &Tests{
        {
            description:        "Invalid user login",
            payload:            `{"username":"adm1n","password":"admin"}`,
            url:                "http://localhost:5000/api/v1/login",
            expectedBody:       `{"status":"error","message":"Login failed","data":{"token":""}}`,
        },
        {
            description:        "Default admin login",
            payload:            `{"username":"admin","password":"admin"}`,
            url:                "http://localhost:5000/api/v1/login",
            expectedBody:       `{"status":"success","message":"","data":{"token":"[0-9A-Za-z._-]+"}}`,
        },
        {
            description:        "Auth",
            payload:            `{"token":"asf98jjasfj"}`,
            url:                "http://localhost:5000/api/v1/auth",
            expectedBody:       `{"status":"error","message":"Login failed","data":{"token":""}}`,
        },
	}
	for _,test  := range *tests {
        jsonResponse := test.postReq()
		assert.Regexp(
            t, 
            test.expectedBody,
            jsonResponse,
        )
	}
}

func (t *Test) postReq() string {
    req, err := http.NewRequest("POST", t.url, bytes.NewBuffer([]byte(t.payload)))
    if err != nil {
        return fmt.Sprintln("Request parse err: ", err)
    }
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{
        Timeout: time.Duration(3 * time.Second),
    }
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Sprintln("Response err: ", err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return string(body)
}