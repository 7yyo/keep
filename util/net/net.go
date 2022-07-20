package util

import (
	"bytes"
	"fmt"
	"github.com/wujiangweiphp/go-curl"
	"io/ioutil"
	"net/http"
)

const (
	get    string = "GET"
	post   string = "POST"
	delete string = "DELETE"
)

func GetHttp(c string) ([]byte, error) {
	cli := &http.Client{}
	req, err := http.NewRequest(get, c, nil)
	if err != nil {
		return nil, err
	}
	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func PostHttp(c string, jsonBody string) (*http.Response, error) {
	cli := &http.Client{}
	b := bytes.NewBufferString(jsonBody)
	req, err := http.NewRequest(post, c, b)
	if err != nil {
		return nil, err
	}
	return cli.Do(req)
}

func Curl(u string, postData map[string]interface{}) error {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	req := curl.NewRequest()

	resp, err := req.
		SetUrl(u).
		SetHeaders(headers).
		SetHeaders(headers).
		SetPostData(postData).
		Post()

	if err != nil {
		return err
	}

	if resp.Raw.StatusCode != 202 && resp.Raw.StatusCode != 200 {
		return fmt.Errorf(resp.Body)
	}
	return nil
}

func DeleteCurl(u string) error {
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if (resp.StatusCode != 202 && resp.StatusCode != 200) || err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}
