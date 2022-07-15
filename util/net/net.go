package util

import (
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
