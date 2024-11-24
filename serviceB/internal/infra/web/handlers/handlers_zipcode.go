package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type CepResponse struct {
	Localidade string `json:"localidade"`
}

func GetLocalidade(ctx context.Context, cep string) (string, error) {

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	resp, err := client.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid cep")
	}

	var cepResp CepResponse
	if err := json.NewDecoder(resp.Body).Decode(&cepResp); err != nil {
		return "", err
	}

	return cepResp.Localidade, nil
}

func IsValidCep(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}
