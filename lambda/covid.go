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

	var data map[string]interface{}
	err := json.Unmarshal([]byte(worldInfections), &data)
	if err != nil {
		log.Fatal(err)
	}
	//for debugging
	//fmt.Println("These are the current COVID-19 numbers as of", today)
	//fmt.Println("World Cases :", data["cases"])
	//fmt.Println("World Deaths :", data["deaths"])
	//fmt.Println("World Recovered :", data["recovered"])
	//fmt.Println("Updated: ", data["updated"])
	
    usEndpoint := "https://corona.lmao.ninja/countries/usa"
    usInfections := getInfections(usEndpoint)

	var usData map[string]interface{}
	userr := json.Unmarshal([]byte(usInfections), &usData)
	if userr != nil {
		log.Fatal(userr)
	}

	//For Debugging
	//fmt.Println("US cases :", usData["cases"])
	//fmt.Println("US cases today:", usData["todayCases"])
	//fmt.Println("US deaths :", usData["deaths"])
	//fmt.Println("US deaths today:", usData["todayDeaths"])
	//fmt.Println("US recovered :", usData["recovered"])


	//Post to Slack - World
	err = postToSlack(webhookUrl, "These are the current COVID-19 numbers as of " + today + "\n" +
	"World Cases : " + fmt.Sprintf("%v", data["cases"]) + "\n" + 
	"World Deaths : " + fmt.Sprintf("%v", data["deaths"]) + "\n" + 
	"World Recovered : " + fmt.Sprintf("%v", data["recovered"]) + "\n" +
	"US Cases : " + fmt.Sprintf("%v", usData["cases"]) + "\n" +
	"US Cases today : " + fmt.Sprintf("%v", usData["todayCases"]) + "\n" +
	"US Deaths : " + fmt.Sprintf("%v", usData["deaths"]) + "\n" +
	"US Deaths today : " + fmt.Sprintf("%v", usData["todayDeaths"]) + "\n" +
	"US recovered : " + fmt.Sprintf("%v", usData["recovered"]) + "\n")

	//For Debugging
    //fmt.Println("These are the current COVID-19 numbers as of", today)
    //fmt.Println("World Infections -----" )
    //fmt.Println(worldInfections)
    //fmt.Println("US Infections -----")
    //fmt.Println(usInfections)

	return "Infections Printed", nil
}

func main() {
    //start lambda function
    lambda.Start(HandleRequest)
}