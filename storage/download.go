package storage

import (
	"errors"
	"io"
	"net/http"
)

func downloadFile(fileUrl string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fileUrl, nil)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		content, err := io.ReadAll(req.Response.Body)
		res.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(content))
	}

	return res, err
}
