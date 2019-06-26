package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	d := json.NewDecoder(strings.NewReader(request.Body))
	p := new(Payload)

	if err := d.Decode(&p); err != nil {
		log.Print(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, err
	}

	if p.Action != Opened {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	}

	requestPayload := &RequestPayload{
		ChatId: ChatId,
		Text:   p.PullRequest.HTMLURL,
	}

	req, err := http.NewRequest(
		http.MethodPost,
		URL,
		strings.NewReader(requestPayload.String()),
	)
	if err != nil {
		log.Print(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r, _ := ioutil.ReadAll(resp.Body)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       string(r),
		}, errors.New("the response from the Telegram API was: " + string(r))
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
