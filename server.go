package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
)

type (
	Lead struct {
		Name string `json:"name" validate:"required"`
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

func main() {
	envError := godotenv.Load()

	if envError != nil {
		log.Fatal("Error loading .env file")
	}
	
	router := echo.New()
	router.Validator = &CustomValidator{validator: validator.New()}
	
	router.GET("/healthcheck", func(context echo.Context) error {
		return context.String(http.StatusOK, "ok")
	})
	router.POST("/lead/create", func(context echo.Context) error {
		lead := new(Lead)
		
		if err := context.Bind(lead); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := context.Validate(lead); err != nil {
      return err
    }

		return context.String(http.StatusOK, "created")
	})

	router.Logger.Fatal(router.Start(":9090"))
}