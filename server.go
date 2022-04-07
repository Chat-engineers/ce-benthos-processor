package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Jeffail/gabs"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
)

type (
	Lead struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName string `json:"lastName" validate:"required"`
		Email string `json:"email" validate:"required"`
		Phone string `json:"phone"`
		Location string `json:"location"`
	}

	CustomValidator struct {
    validator *validator.Validate
  }
)

func (cv *CustomValidator) Validate(i interface{}) error {
  if err := cv.validator.Struct(i); err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  }
  return nil
}

func handleHealthcheck(context echo.Context) error {
	return context.String(http.StatusOK, "ok")
}

func createLead(lead *Lead) ([]byte, error) {
	client := http.Client{}

	postBody, err := gabs.ParseJSON([]byte(`{
		"party": {
			"type": "person",
			"firstName": "First name",
			"lastName": "LastName",
			"emailAddresses": [{ "address": "test@domain.com" }],
			"tags": [{ "name": "lead" }]
		}
	}`))

	if err != nil {
		log.Fatal("postBody", err)
	}

	emailAddress := gabs.New()
	emailAddress.Set(lead.Email, "address") 

	postBody.Set(lead.FirstName, "party", "firstName")
	postBody.Set(lead.LastName, "party", "lastName")
	// postBody.Array("party", "emailAddresses")

	log.Print(emailAddress.String())

	// TODO: how to inject object to array?????
	postBody.ArrayAppend(emailAddress, "party", "emailAddresses")

	log.Print(postBody.String())

	responseBody := bytes.NewBuffer(postBody.Bytes())

	req, err := http.NewRequest("POST", "https://api.capsulecrm.com/api/v2/parties", responseBody)
	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"Authorization": []string{"Bearer " + os.Getenv("CAPSULE_API_TOKEN")},
	}

	if err != nil {
		log.Fatal("Capsule error", err)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal("execute request", err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Read response", err)
	}

	return body, err
}

func handleCreateLead(context echo.Context) error {
	lead := new(Lead)
		
	if err := context.Bind(lead); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := context.Validate(lead); err != nil {
		return err
	}

	body, err := createLead(lead)

	log.Print(string(body))

	if err != nil {
		log.Fatal("Response ", err)
	}


	return context.String(http.StatusOK, "ok")
}

func main() {
	envError := godotenv.Load()

	if envError != nil {
		log.Fatal("Error loading .env file")
	}
	
	router := echo.New()
	router.Validator = &CustomValidator{validator: validator.New()}

	router.GET("/healthcheck", handleHealthcheck)
	router.POST("/lead/create", handleCreateLead)

	router.Logger.Fatal(router.Start(":9090"))
}