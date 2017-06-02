package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/andersjanmyr/jobs/models"
	"github.com/pkg/errors"
)

type JobClient struct {
	baseUrl string
}

func (c *JobClient) request(method, path string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/%s", c.baseUrl, strings.TrimLeft(path, "/"))
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Request parsing failed %s", url))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Request failed %s", url))
	}

	return resp.Body, nil
}
func (c *JobClient) Index() ([]models.Job, error) {
	rc, err := c.request(http.MethodGet, "/")
	if err != nil {
		return nil, err
	}
	jobs, err := models.ParseJobs(rc)
	return jobs, err
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
	fmt.Printf("%#v\n", data)

}
