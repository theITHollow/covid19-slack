package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
	"time"
	"encoding/json"
	"bytes"
	"errors"
    "github.com/aws/aws-lambda-go/lambda"
)
type slackRequestBody struct {
    Text string `json:"text"`
}

func getInfections(endpoint string) string {

    resp, err := http.Get(endpoint)
    if err != nil {
        log.Fatalln(err)
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln(err)
    }
    
    return string(body)
}

func postToSlack(webhookURL string, msg string) error {
	slackBody, _ := json.Marshal(slackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}

	return nil
}

func HandleRequest() (string, error) {
	//set slack webhook URL
	webhookUrl := "NEEDS_TO_BE_FILLED_IN"

    //Get today's date
    date := time.Now()
    today := date.Format("2006.01.02 15:04:05")

    //Get Global Infections
    globalEndpoint := "https://corona.lmao.ninja/all"
    worldInfections := getInfections(globalEndpoint)

    //Get US Infections
    usEndpoint := "https://corona.lmao.ninja/countries/usa"
    usInfections := getInfections(usEndpoint)

	//Post to Slack - World
	err := postToSlack(webhookUrl, worldInfections)
	//Post to Slack - USA
	err = postToSlack(webhookUrl, usInfections)
	if err != nil {
		log.Fatal(err)
	}

    fmt.Println("These are the current COVID-19 numbers as of", today)
    fmt.Println("World Infections -----" )
    fmt.Println(worldInfections)
    fmt.Println("US Infections -----")
    fmt.Println(usInfections)

	return "Infections Printed", nil
}

func main() {
    //start lambda function
    lambda.Start(HandleRequest)
}