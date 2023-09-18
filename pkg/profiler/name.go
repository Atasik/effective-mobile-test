package profiler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type NameProfiler struct {
	agifyResource       string
	genderizeResource   string
	nationalizeResource string
	client              http.Client
}

func NewNameProfiler(agifyResource, genderizeResource, nationalizeResource string, timeout time.Duration) *NameProfiler {
	return &NameProfiler{
		client: http.Client{
			Timeout: timeout,
		},
		agifyResource:       agifyResource,
		genderizeResource:   genderizeResource,
		nationalizeResource: nationalizeResource,
	}
}

type AgifyResponse struct {
	Age int `json:"age"`
}

type GenderizeResponse struct {
	Gender string `json:"gender"`
}

type NationalizeResponse struct {
	Country []Country `json:"country"`
}

type Country struct {
	CountryID string `json:"country_id"`
}

func (p *NameProfiler) AgifyPerson(name string) (int, error) {
	req, err := http.NewRequest("GET", p.agifyResource, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()
	resp, err := p.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	var agifyResponse AgifyResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&agifyResponse); err != nil {
		return 0, errors.New("error while decoding")
	}

	return agifyResponse.Age, nil
}

func (p *NameProfiler) GenderizePerson(name string) (string, error) {
	req, err := http.NewRequest("GET", p.genderizeResource, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create new request due to error: %v", err)
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	var genderizeResponse GenderizeResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&genderizeResponse); err != nil {
		return "", errors.New("error while decoding")
	}

	return genderizeResponse.Gender, nil
}

func (p *NameProfiler) NationalizePerson(name string) (string, error) {
	req, err := http.NewRequest("GET", p.nationalizeResource, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create new request due to error: %v", err)
	}

	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	var nationalizeResponse NationalizeResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&nationalizeResponse); err != nil {
		return "", errors.New("error while decoding")
	}

	if len(nationalizeResponse.Country) == 0 {
		return "", nil
	}

	return nationalizeResponse.Country[0].CountryID, nil
}
