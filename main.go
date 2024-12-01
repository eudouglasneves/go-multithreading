package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Address struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"localidade"`
	UF         string `json:"uf"`
}

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func getFromBrasilAPI(cep string, ch chan<- string) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Erro na API BrasilAPI: %v", err)
		return
	}
	defer resp.Body.Close()

	var apiResp BrasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		ch <- fmt.Sprintf("Erro ao decodificar resposta da API BrasilAPI: %v", err)
		return
	}

	address := Address{
		Logradouro:  apiResp.Street,
		Bairro:      apiResp.Neighborhood,
		Cidade:      apiResp.City,
		UF:          apiResp.State,
	}

	ch <- fmt.Sprintf("BrasilAPI: %s, %s, %s, %s", address.Logradouro, address.Bairro, address.Cidade, address.UF)
}

func getFromViaCEP(cep string, ch chan<- string) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Erro na API ViaCEP: %v", err)
		return
	}
	defer resp.Body.Close()

	var apiResp Address
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		ch <- fmt.Sprintf("Erro ao decodificar resposta da API ViaCEP: %v", err)
		return
	}

	ch <- fmt.Sprintf("ViaCEP: %s, %s, %s, %s", apiResp.Logradouro, apiResp.Bairro, apiResp.Cidade, apiResp.UF)
}

func main() {
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scan(&cep)

	ch := make(chan string, 2)
	timeout := time.After(1 * time.Second)

	go getFromBrasilAPI(cep, ch)
	go getFromViaCEP(cep, ch)

	select {
	case result := <-ch:
		fmt.Println(result)
	case <-timeout:
		fmt.Println("Erro: Tempo de resposta excedido.")
	}
}
