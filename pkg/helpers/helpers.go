package helpers

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func HttpGET(url string) ([]byte, *int, error) {

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	var client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return bodyText, &resp.StatusCode, nil
}

func DecodeResponseBody(data []byte, target interface{}) error {
	err := json.Unmarshal(data, &target)
	if err != nil {
		return err
	}

	return nil
}
