package funtranslations

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	// funtranslationsBaseURL is the base URL of the `funtranslations` APIs
	//
	// Reference: https://funtranslations.com/api/
	funtranslationsBaseURL = "https://api.funtranslations.com/translate"

	shakespearePath = "shakespeare.json"
	yodaPath        = "yoda.json"

	// TranslatorYoda passed to FunTranslate will make it output a Yoda translation
	TranslatorYoda = "yoda"

	// TranslatorShakespeare passed to FunTranslate will make it output a Shakespeare translation
	TranslatorShakespeare = "shakespeare"
)

var ErrUnrecognizedTranslator = errors.New("unrecognized translator type")

type Client interface {
	// FunTranslate given a Translator type and a text will output the translation
	// only Yoda and Shakespeare translations are currently supported.
	// Providing an unknown translatorType argument results in an error
	FunTranslate(translatorType, text string) (string, error)
}

type client struct {
	baseURL string
}

func NewClient() Client {
	return &client{funtranslationsBaseURL}
}

func (c *client) FunTranslate(translatorType, text string) (string, error) {
	path, err := mapTranslatorToPath(translatorType)
	if err != nil {
		return "", err
	}
	body := &translateReqBody{text}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(c.baseURL+"/"+path, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBody := &translateRespBody{}
	if err := json.Unmarshal(respBodyBytes, respBody); err != nil {
		return "", err
	}

	return respBody.Contents.Translated, nil
}

type translateReqBody struct {
	Text string `json:"text"`
}

type translateRespBody struct {
	Contents translateRespBodyContents `json:"contents"`
}

type translateRespBodyContents struct {
	Translated string `json:"translated"`
}

func mapTranslatorToPath(translatorType string) (string, error) {
	if translatorType == TranslatorYoda {
		return yodaPath, nil
	}
	if translatorType == TranslatorShakespeare {
		return shakespearePath, nil
	}
	return "", ErrUnrecognizedTranslator
}
