package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	SonarMeasuresApi = "%s/api/measures/search?projectKeys=%s&metricKeys=bugs,vulnerabilities,code_smells,duplicated_lines_density,coverage,ncloc"
)

var (
	SonarUsername string

	SonarPassword string

	SonarServer string
)

// SonarResult acquire project measures
func SonarResult(httpClient *http.Client, projectKey string) (map[string]string, error) {
	log.Printf("SonarResult.projectKey:%s", projectKey)
	//request sonar server
	reqMeasure, _ := http.NewRequest("GET", fmt.Sprintf(SonarMeasuresApi, SonarServer, projectKey), nil)
	reqMeasure.SetBasicAuth(SonarUsername, SonarPassword)

	sonarMeasuresResp, err := httpClient.Do(reqMeasure)
	if err != nil {
		log.Printf("reqMeasure.err:%s,%+v", err.Error(), sonarMeasuresResp)
		return nil, err
	}
	statusCode := sonarMeasuresResp.StatusCode
	if statusCode != http.StatusOK {
		log.Printf("sonarMeasuresResp.statusCode:%d", sonarMeasuresResp.StatusCode)
		return nil, err
	}

	measureByte, err := io.ReadAll(sonarMeasuresResp.Body)
	if err != nil {
		log.Println("Error reading measureByte body")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(sonarMeasuresResp.Body)

	var mb map[string][]map[string]any
	_ = json.Unmarshal(measureByte, &mb)
	if msgSli, ok := mb["errors"]; ok {
		alert := msgSli[0]
		log.Printf("SONAR_MEASURES_API.resp:%s", alert["msg"])
		return nil, err
	}
	mSli := mb["measures"]
	metricValueMap := make(map[string]string, 6)
	for _, val := range mSli {
		metricValueMap[val["metric"].(string)] = val["value"].(string)
	}
	log.Printf("SonarResult.metricValueMap:%s", metricValueMap)
	return metricValueMap, nil
}
