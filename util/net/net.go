package util

import (
	"errors"
	"github.com/wujiangweiphp/go-curl"
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

func PostHttp(cmd string) error {
	c := &http.Client{}
	req, err := http.NewRequest("POST", cmd, nil)
	if err != nil {
		return err
	}
	if _, err = c.Do(req); err != nil {
		return err
	}
	return nil
}

func CurlHttp(url string, queries map[string]string, postData map[string]interface{}) (string, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	req := curl.NewRequest()
	resp, err := req.SetUrl(url).SetHeaders(headers).SetQueries(queries).SetPostData(postData).Get()
	if err != nil {
		return "", err
	}
	if resp.IsOk() {
		return resp.Body, nil
	} else {
		return "", errors.New("curl http failed")
	}
}
