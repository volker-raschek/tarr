package health

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var ErrNoAPIToken error = errors.New("No API token defined")
var ErrNoURL error = errors.New("No API token defined")

type ReadinessProbe struct {
	additionalQueryValues url.Values
	insecure              bool
	url                   string
}

func (rp *ReadinessProbe) QueryAdd(key, value string) *ReadinessProbe {
	if rp.additionalQueryValues.Has(key) {
		rp.additionalQueryValues.Add(key, value)
	} else {
		rp.additionalQueryValues.Set(key, value)
	}
	return rp
}

func (rp *ReadinessProbe) Insecure(insecure bool) *ReadinessProbe {
	rp.insecure = insecure
	return rp
}

func (rp *ReadinessProbe) Run(ctx context.Context) error {
	if len(rp.url) <= 0 {
		return ErrNoURL
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: rp.insecure,
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, rp.url, nil)
	if err != nil {
		return err
	}

	reValues := req.URL.Query()
	for key, values := range rp.additionalQueryValues {
		for _, value := range values {
			reValues.Add(key, value)
		}
	}
	req.URL.RawQuery = reValues.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Received unexpected HTTP status code %v", resp.StatusCode)
	}

	return nil
}

func NewReadinessProbe(url string) *ReadinessProbe {
	return &ReadinessProbe{
		additionalQueryValues: make(map[string][]string),
		url:                   url,
	}
}
