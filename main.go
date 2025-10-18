package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIResult armazena o resultado de uma chamada de API,
// incluindo o nome da API e o corpo da resposta.
type APIResult struct {
	APIName string
	Data    []byte // Usamos []byte para o corpo da resposta
}

// fetchCEP é a função que será executada em uma goroutine.
// Ela faz a requisição HTTP para a URL fornecida e envia o resultado
// para o channel (ch).
func fetchCEP(apiName, url string, ch chan<- APIResult) {
	fmt.Printf("-> Tentando API: %s\n", apiName)

	// Cria um cliente HTTP com um timeout um pouco menor que o global
	// para evitar que a requisição fique presa indefinidamente.
	client := http.Client{
		Timeout: 950 * time.Millisecond,
	}

	// Faz a requisição GET
	resp, err := client.Get(url)
	if err != nil {
		// Se houver erro (timeout, DNS, etc.), apenas informa e encerra a goroutine.
		// Não envia nada para o channel.
		fmt.Printf("Erro ao consultar %s: %v\n", apiName, err)
		return
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler resposta de %s: %v\n", apiName, err)
		return
	}

	// Se tudo deu certo, envia o resultado para o channel.
	// A primeira goroutine a chegar aqui "vencerá" a corrida.
	ch <- APIResult{
		APIName: apiName,
		Data:    body,
	}
}

func main() {
	cep := "01153000"

	// O channel que receberá o resultado da API mais rápida.
	// Um channel não bufferizado é suficiente aqui.
	resultChannel := make(chan APIResult)

	// URLs das APIs
	urlBrasilAPI := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	urlViaCEP := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	// Inicia as duas goroutines. Elas executarão em paralelo.
	go fetchCEP("BrasilAPI", urlBrasilAPI, resultChannel)
	go fetchCEP("ViaCEP", urlViaCEP, resultChannel)

	fmt.Printf("Iniciando corrida de APIs para o CEP: %s\n\n", cep)
	startTime := time.Now()

	// O `select` é um controle de concorrência poderoso em Go.
	// Ele vai esperar até que um dos seus "cases" esteja pronto para ser executado.
	select {
	// Caso 1: Um resultado foi recebido no channel.
	case winner := <-resultChannel:
		elapsedTime := time.Since(startTime)

		// Para imprimir o JSON de forma legível (pretty-print)
		var prettyData map[string]interface{}
		json.Unmarshal(winner.Data, &prettyData)
		prettyJSON, _ := json.MarshalIndent(prettyData, "", "  ")

		// --- RESULTADO ---
		fmt.Println("\n----------------------------------------------------")
		fmt.Printf("🏆 API VENCEDORA: %s (em %s)\n", winner.APIName, elapsedTime)
		fmt.Println("----------------------------------------------------")
		fmt.Println(string(prettyJSON))
		fmt.Println("----------------------------------------------------")

	// Caso 2: O tempo de 1 segundo se esgotou.
	case <-time.After(1 * time.Second):
		// --- TIMEOUT ---
		fmt.Println("\n----------------------------------------------------")
		fmt.Println("❌ ERRO: Timeout. Nenhuma API respondeu em 1 segundo.")
		fmt.Println("----------------------------------------------------")
	}
}
