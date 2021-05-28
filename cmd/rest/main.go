package main

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

var validate = validator.New()

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/numbers", getNumbers)
	e.Logger.Fatal(e.Start(":8080"))
}

type getNumbersInput struct {
	From int `query:"from" validate:"required"`
	To   int `query:"to" validate:"required"`
}

type getNumbersOutput struct {
	Numbers []int `json:"numbers"`
}

func getNumbers(c echo.Context) error {
	var query getNumbersInput
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := validate.Struct(query); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if query.To < query.From || query.To-query.From > 10000000 {
		return c.JSON(http.StatusBadRequest, errors.New("invalid query"))
	}

	var resp getNumbersOutput
	for i := query.From; i <= query.To; i++ {
		resp.Numbers = append(resp.Numbers, i)
	}

	return c.JSON(http.StatusOK, resp)
}
