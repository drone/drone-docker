package dockerhub

import (
	"fmt"
	"net/http"
)

// Client defines the client for retriving information from
// the Dockerhub API.
type Client struct {
	Username string
	Password string
}

// New creates a new Client with the given credentials.
func New(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
	}
}

// DeleteTag deletes a tag from Dockerhub.
func (c *Client) DeleteTag(image, tag string) error {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/%s/", image, tag)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if res.StatusCode > 299 {
		fmt.Errorf("Got status code: %d", res.StatusCode)
	}
	return err
}
