package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func GetHttp(cmd string) ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", cmd, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func PostHttp(cmd string, jsonBody string) (*http.Response, error) {
	c := &http.Client{}
	b := bytes.NewBufferString(jsonBody)
	req, err := http.NewRequest("POST", cmd, b)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
