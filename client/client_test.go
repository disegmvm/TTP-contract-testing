package client

import (
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"net/http"
	"os"
	"testing"
)

var pact dsl.Pact
var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/pacts", dir)
var logDir = fmt.Sprintf("%s/log", dir)

func createPact() dsl.Pact {
	return dsl.Pact{
		Consumer: "Sample Consumer",
		Provider: "Sample Provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	}
}

func TestMain(m *testing.M) {
	// Setup Pact and related test stuff
	pact = createPact()

	// Proactively start service to get access to the port
	pact.Setup(true)

	// Run all the tests
	code := m.Run()

	// Shutdown the Mock Service and Write pact files to disk
	pact.WritePact()
	pact.Teardown()

	os.Exit(code)
}

func validate() (err error) {
	url := fmt.Sprintf("http://localhost:%d/cars", pact.Server.Port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	return
}

func Test_SomeValuesAndKeysPresence(t *testing.T) {
	ma2 := make(map[string]string)
	ma2["title"] = "BMW"

	pact.
		AddInteraction().
		Given("Validate SOME values and keys are present").
		UponReceiving("A GET request").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars", "/cars"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body: []map[string]string{
				ma2, {},
			},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(validate)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func Test_SomeKeysPresence(t *testing.T) {
	ma2 := make(map[string]dsl.Matcher)
	ma2["title"] = dsl.Like("any string")
	ma2["id"] = dsl.Like("any string")
	ma2["color"] = dsl.Like("any string")

	pact.
		AddInteraction().
		Given("Validate SOME keys are present").
		UponReceiving("A GET request").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars", "/cars"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body: []map[string]dsl.Matcher{
				ma2, ma2,
			},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(validate)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func Test_WholeResponseBody(t *testing.T) {
	type Car struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Color string `json:"color"`
	}

	cars := []Car{
		{ID: "1", Title: "BMW", Color: "Black"},
		{ID: "2", Title: "Tesla", Color: "Red"},
	}

	pact.AddInteraction().
		Given("Match the whole response body").
		UponReceiving("A a GET request").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars", "/cars"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   cars,
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(validate)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func Test_AllKeysPresence(t *testing.T) {
	type car struct {
		ID    dsl.Matcher `json:"id"`
		Title dsl.Matcher `json:"title"`
		Color dsl.Matcher `json:"color"`
	}

	singleCar := car{
		ID:    dsl.Term("any string", `\w+`),
		Title: dsl.Term("any string", `\w+`),
		Color: dsl.Term("any color", `\w+`),
	}

	pact.
		AddInteraction().
		Given("Validate ALL keys are present").
		UponReceiving("A GET request").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/cars", "/cars"),
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
			Body:   []car{singleCar, singleCar},
			Headers: dsl.MapMatcher{
				"Content-Type": dsl.Term("application/json; charset=utf-8", `application\/json`),
			},
		})

	err := pact.Verify(validate)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}

	// specify PACT publisher
	/*publisher := dsl.Publisher{}
	err = publisher.Publish(types.PublishRequest{
		PactURLs: []string{"../client/pacts/consumer_name-provider_name.json"},
		PactBroker:      "https://pen.pactflow.io/", //link to your remote Contract broker
		BrokerToken:     "jEQnxw7xWgYRv-3-G7Cx-g",   //your PactFlow token
		ConsumerVersion: "2.0.2",
		Tags:            []string{"2.0.2", "latest"},
	})*/
}
