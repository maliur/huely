package main

import (
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"net/http"
	"os"
)

type Light struct {
	ID   string
	Name string `json:"name"`
}

func main() {
	apiKey := os.Getenv("HUE_API_KEY")
	hubIP := os.Getenv("HUE_HUB_IP")

	if len(apiKey) == 0 {
		fmt.Println("HUE_API_KEY not provided")
		return
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/api/%s/lights", hubIP, apiKey))
	if err != nil {
		fmt.Printf("could not make request: %v", err)
		return
	}

	defer resp.Body.Close()

	var body map[string]Light

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		fmt.Printf("could decode json body: %v", err)
		return
	}

	var lights []Light
	for k, v := range body {
		v.ID = k
		lights = append(lights, v)
	}

	fmt.Println(lights)

	j, _ := json.MarshalIndent(lights, "", "\t")

	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Lights"), 0, 1, false).
		AddItem(tview.NewTextView().SetText(string(j)), 0, 1, false)
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
