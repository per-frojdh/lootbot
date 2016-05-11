package lib

import (
    "errors"
    models "github.com/per-frojdh/lootbot/models"
)

type Parameter struct{
    Key string                  `json:"key"`
    Value string                `json:"value"`
}

type APIError struct {
    Errors []ResponseMessage    `json:"errors"`
    Request RequestData         `json:"request,omitempty"`
}

type RequestData struct {
    ContentType string          `json:"contentType,omitempty"`
    Params []Parameter          `json:"parameters,omitempty"`
}

type ResponseMessage struct{
    StatusCode int              `json:"statusCode"`
    Message string              `json:"message"`
    Error string                `json:"errorMessage"`
}

func CreatePanicResponse(errorString string) error {
    return errors.New(models.ErrorMessages[errorString])
}

func CreateErrorResponse(statusCode int, errorString string) interface{} {
    var response ResponseMessage
    response.StatusCode = statusCode
    response.Message = models.ErrorMessages[errorString]
    response.Error = errorString
    return response
}