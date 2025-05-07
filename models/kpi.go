package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type KPI struct {
	Name           string   `json:"name"`
	DisplayName    string   `json:"displayName"`
	Description    string   `json:"description"`
	Formula        string   `json:"formula"`
	Unit           string   `json:"unit"`
	Type           string   `json:"type"`
	Object         []string `json:"object"`
	PrometheusType string   `json:"prometheus_type"`
	NFType         string   `json:"nf_type"`
	Increment      bool     `json:"increment"`
	Decrement      bool     `json:"decrement"`
}

func LoadKPIsFromFile(filePath string) ([]KPI, error) {
	file, err := os.Open("/home/amantya/Music/amantya_metrics/models/kpi.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var kpis []KPI
	if err := json.Unmarshal(byteValue, &kpis); err != nil {
		return nil, err
	}

	return kpis, nil
}
