package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PortainerAPI struct {
	ApiKey             string
	ApiBaseUrl         string
	EndpointId         int
	SwarmId            string
	InsecureSkipVerify bool
}

type PortainerStack struct {
	Id                       float64
	Name                     string   `json:"name"`
	RepositoryUrl            string   `json:"repositoryURL"`
	RepositoryReferenceName  string   `json:"repositoryReferenceName"`
	RepositoryUsername       string   `json:"RepositoryUsername"`
	RepositoryPassword       string   `json:"repositoryPassword"`
	RepositoryAuthentication bool     `json:"repositoryAuthentication"`
	ComposeFile              string   `json:"composeFile"`
	AdditionalFiles          []string `json:"additionalFiles"`
	FromAppTemplate          bool     `json:"fromAppTemplate"`
	EndpointId               float64
	SwarmId                  string                   `json:"swarmID"`
	TlsSkipVerify            bool                     `json:"tlsskipVerify"`
	Env                      []PortainerStackEnv      `json:"Env"`
	AutoUpdate               PortainerStackAutoUpdate `json:"AutoUpdate"`
}

type PortainerStackAutoUpdate struct {
	ForcePullImage bool   `json:"forcePullImage"`
	ForceUpdate    bool   `json:"forceUpdate"`
	Interval       string `json:"interval"`
	JobID          string `json:"jobID"`
	Webhook        string `json:"webhook"`
}

type PortainerStackEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (pApi *PortainerAPI) get(apiPath string) (data []interface{}, err error) {
	targetUrl := fmt.Sprintf("%s%s", pApi.ApiBaseUrl, apiPath)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: pApi.InsecureSkipVerify}
	client := &http.Client{}

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-Key", pApi.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var j interface{}
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		return nil, err
	}

	data, ok := j.([]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot parse the result")
	}
	return data, nil
}

func (pApi *PortainerAPI) post(apiPath string, jsonData []byte) (err error) {
	targetUrl := fmt.Sprintf("%s%s", pApi.ApiBaseUrl, apiPath)
	fmt.Printf("Target URL: %s\n", targetUrl)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: pApi.InsecureSkipVerify}
	client := &http.Client{}

	fmt.Println(string(jsonData))
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Add("X-API-Key", pApi.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error: %s", resp.Status)
	}
	fmt.Println(resp.Status)
	return nil
}

func (pApi *PortainerAPI) put(apiPath string, jsonData []byte) {
	targetUrl := fmt.Sprintf("%s%s", pApi.ApiBaseUrl, apiPath)
	fmt.Printf("Target URL: %s\n", targetUrl)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: pApi.InsecureSkipVerify}
	client := &http.Client{}

	req, err := http.NewRequest("PUT", targetUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Add("X-API-Key", pApi.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	return
}

func (pApi *PortainerAPI) delete(apiPath string) {
	targetUrl := fmt.Sprintf("%s%s", pApi.ApiBaseUrl, apiPath)
	// fmt.Printf("Target URL: %s\n", targetUrl)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: pApi.InsecureSkipVerify}
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", targetUrl, nil)
	if err != nil {
		return
	}
	req.Header.Add("X-API-Key", pApi.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error getting response: %s", err.Error())
		return
	}
	fmt.Println(string(body))
	return
}

func (pApi *PortainerAPI) GetStacks() (resp []interface{}, err error) {
	return pApi.get("/stacks")
}

func (pApi *PortainerAPI) IsStack(stackId int) (result bool) {
	targetUrl := fmt.Sprintf("/stack/%d", stackId)
	_, err := pApi.get(targetUrl)
	if err != nil {
		return false
	}
	return true
}

func (pApi *PortainerAPI) GetStackByName(stackName string) (foundStack PortainerStack, result bool) {
	resultStack := PortainerStack{}
	resp, err := pApi.GetStacks()

	if err != nil {
		return PortainerStack{}, false
	}
	for _, item := range resp {
		stack := item.(map[string]interface{})
		//fmt.Println(stack)
		if stack["Name"] == stackName {
			resultStack.Name = stack["Name"].(string)
			if stack["Id"] != nil {
				resultStack.Id = stack["Id"].(float64)
			}
			resultStack.EndpointId = stack["EndpointId"].(float64)

			if stack["GitConfig"] != nil {
				gitConfig := stack["GitConfig"].(map[string]interface{})
				resultStack.RepositoryUrl = gitConfig["URL"].(string)
				resultStack.RepositoryReferenceName = gitConfig["ReferenceName"].(string)
				// resultStack.RepositoryAuthentication = gitConfig[""]
				resultStack.ComposeFile = gitConfig["ConfigFilePath"].(string)
				resultStack.TlsSkipVerify = gitConfig["TLSSkipVerify"].(bool)
			}
			return resultStack, true
		}
	}
	return PortainerStack{}, false
}

func (pApi *PortainerAPI) CreateStack(newStack PortainerStack) (err error) {

	_, isStack := pApi.GetStackByName(newStack.Name)
	if isStack {
		return fmt.Errorf("stack %s already exists, skipping its creation", newStack.Name)
	}

	targetUrl := fmt.Sprintf("/stacks/create/swarm/repository?endpointId=%d", int(newStack.EndpointId))
	jsonData, err := json.Marshal(newStack)

	return pApi.post(targetUrl, jsonData)
}

func (pApi *PortainerAPI) DeleteStack(existingStack PortainerStack) (err error) {
	_, isStack := pApi.GetStackByName(existingStack.Name)
	if !isStack {
		return fmt.Errorf("stack %s doesn't exist", existingStack.Name)
	}

	targetUrl := fmt.Sprintf("/stacks/%d?endpointId=%d", int(existingStack.Id), int(existingStack.EndpointId))
	pApi.delete(targetUrl)

	return nil

}

/*
func (pApi *PortainerAPI) UpdateStack(existingStack PortainerStack) (err error) {
	_, isStack := pApi.GetStackByName(existingStack.Name)
	if !isStack {
		return fmt.Errorf("stack %s doesn't exist, skipping update", existingStack.Name)
	}

	targetUrl := fmt.Sprintf("/stacks/%d", int(existingStack.Id))

	jsonData, err := json.Marshal(existingStack)
	pApi.post(targetUrl, jsonData)

	return nil

}
*/

