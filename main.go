package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var ColorLoop string = "colorloop"

type Group struct {
	ID     string
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Lights []string `json:"lights"`
}

type Light struct {
	ID   string
	Name string `json:"name"`
}

type Config struct {
	ApiKey string
	HubIP  string
}

type Hue struct {
	client *http.Client
	config Config
}

func NewHue(config Config) *Hue {
	return &Hue{
		client: &http.Client{Timeout: time.Second * 60},
		config: config,
	}
}

func (h *Hue) FetchLights() ([]Light, error) {
	resp, err := h.client.Get(fmt.Sprintf("http://%s/api/%s/lights", h.config.HubIP, h.config.ApiKey))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var body map[string]Light

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	var lights []Light
	for k, v := range body {
		v.ID = k
		lights = append(lights, v)
	}

	return lights, nil
}

func (h *Hue) FetchGroups() ([]Group, error) {
	resp, err := h.client.Get(fmt.Sprintf("http://%s/api/%s/groups", h.config.HubIP, h.config.ApiKey))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var body map[string]Group

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	var groups []Group
	for k, v := range body {
		v.ID = k
		groups = append(groups, v)
	}

	return groups, nil
}

func prettyPrint(data interface{}) {
	j, _ := json.MarshalIndent(data, "", "\t")

	fmt.Println(string(j))
}

func main() {
	apiKey := os.Getenv("HUE_API_KEY")
	hubIP := os.Getenv("HUE_HUB_IP")

	if len(apiKey) == 0 {
		fmt.Println("HUE_API_KEY not provided")
		return
	}

	if len(hubIP) == 0 {
		fmt.Println("HUE_HUB_IP not provided")
		return
	}

	config := Config{
		ApiKey: apiKey,
		HubIP:  hubIP,
	}

	hue := NewHue(config)

	_, err := hue.FetchLights()
	if err != nil {
		fmt.Printf("could not fetch lights: %v", err)
		return
	}

	groups, err := hue.FetchGroups()
	if err != nil {
		fmt.Printf("could not fetch groups: %v", err)
		return
	}

	// prettyPrint(lights)
	prettyPrint(groups)

	// fmt.Println(groups)
}
