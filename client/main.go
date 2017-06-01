package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type JobClient struct {
	baseUrl string
}

func (c *JobClient) request(method, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseUrl, strings.TrimLeft(path, "/"))
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Request parsing failed %s", url))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Request failed %s", url))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Request body read failed %s", url))
	}

	return body, nil
}
func (c *JobClient) Index() (string, error) {
	data, err := c.request(http.MethodGet, "/")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *JobClient) Create()  {}
func (c *JobClient) Show()    {}
func (c *JobClient) Update()  {}
func (c *JobClient) Destroy() {}

func main() {
	client := JobClient{baseUrl: "http://localhost:5555/jobs"}
	data, err := client.Index()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

}
