package wargaming

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/cufee/am-api/config"
)

// HTTP client
var clientHTTP = &http.Client{Timeout: 500 * time.Millisecond, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

// Mutex lock for rps counter
var waitGroup sync.WaitGroup
var limiterChan chan int = make(chan int, config.OutRPSlimit)

// GetLimited -
func GetLimited(url string, target interface{}) error {
	// Outgoing rate limiter
	start := time.Now()
	limiterChan <- 1
	defer func() {
		go func() {
			timer := time.Now().Sub(start)

			if timer < (time.Second * 1) {
				toSleep := (time.Second * 1) - timer
				time.Sleep(toSleep)
			}
			<-limiterChan
		}()
	}()

	var resData []byte
	res, err := clientHTTP.Get(url)

	if res == nil {
		// Change timeout to account for cold starts
		clientHTTP.Timeout = 2 * time.Second
		defer func() { clientHTTP.Timeout = 750 * time.Millisecond }()

		// Marshal a request
		proxyReq := struct {
			URL string `json:"url"`
		}{
			URL: url,
		}
		reqData, pErr := json.Marshal(proxyReq)
		if pErr != nil {
			return fmt.Errorf("no response recieved from WG API after proxy try, error: %v", pErr)
		}

		// Make request
		req, pErr := http.NewRequest("GET", config.WGProxyURL, bytes.NewBuffer(reqData))
		if pErr != nil {
			return fmt.Errorf("failed to make a proxy request, error: %v", pErr)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send request
		res, pErr = clientHTTP.Do(req)
		if res == nil {
			return fmt.Errorf("no response recieved from WG API after proxy try, error: %v", pErr)
		}
		resData, pErr = ioutil.ReadAll(res.Body)

		// Check for errors
		var proxyErr struct {
			Message string `json:"error"`
		}
		json.Unmarshal(resData, &proxyErr)
		if proxyErr.Message != "" {
			pErr = fmt.Errorf(proxyErr.Message)
		}

		// Set error to proxy error
		err = pErr
	} else {
		resData, err = ioutil.ReadAll(res.Body)
	}

	// Check error
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.Unmarshal(resData, target)
}
