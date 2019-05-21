package utils

import (
	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Post(endPoint, contentType string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorln("Error while trying to send request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	return resp, err
}
