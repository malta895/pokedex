package funtranslations

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
	return "", nil
}
