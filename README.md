# API Race - Consulta de CEP em Go

Este projeto é uma aplicação de linha de comando em Go que demonstra o uso de concorrência (goroutines e channels) para consultar duas APIs de CEP simultaneamente e retornar a resposta da mais rápida.

## 🎯 Objetivo

O desafio consiste em criar uma aplicação que, ao receber um CEP, dispare duas requisições simultâneas para as APIs **BrasilAPI** e **ViaCEP**. A aplicação deve acatar a resposta que chegar primeiro, descartar a outra e exibir o resultado no terminal. Além disso, um timeout global de 1 segundo deve ser implementado.

## ✨ Funcionalidades

-   **Consulta Simultânea:** Utiliza goroutines para fazer requisições HTTP em paralelo, sem que uma espere pela outra.
-   **Resposta Mais Rápida:** Utiliza channels para comunicação entre goroutines. A primeira a obter uma resposta válida "vence" a corrida.
-   **Timeout Global:** Implementa um timeout de 1 segundo. Caso nenhuma das APIs responda a tempo, uma mensagem de erro é exibida.
-   **Exibição no Terminal:** O resultado da API vencedora, incluindo seu nome e os dados do endereço, é exibido de forma clara no terminal.

## 🛠️ Tecnologias Utilizadas

-   **Go (Golang)**
-   **Goroutines** para concorrência.
-   **Channels** para comunicação segura entre goroutines.
-   Pacote `net/http` para realizar as requisições HTTP.
-   Pacote `encoding/json` para formatar a saída JSON.
-   Pacote `time` para controle do timeout.

## 🚀 Como Executar

### Pré-requisitos

-   Você precisa ter o [Go](https://go.dev/doc/install) instalado em sua máquina (versão 1.18 ou superior).

### Passos

1.  **Clone o repositório ou salve o código:**
    Salve o código acima em um arquivo chamado `main.go`.

2.  **Navegue até a pasta do projeto:**
    Abra seu terminal e acesse o diretório onde você salvou o arquivo `main.go`.

3.  **Execute a aplicação:**
    Use o comando `go run` para compilar e executar o arquivo.

    ```bash
    go run main.go
    ```

## 📝 Exemplo de Saída

### Cenário de Sucesso (BrasilAPI foi mais rápida)

```
Iniciando corrida de APIs para o CEP: 01153000

-> Tentando API: BrasilAPI
-> Tentando API: ViaCEP

----------------------------------------------------
🏆 API VENCEDORA: BrasilAPI (em 153.456789ms)
----------------------------------------------------
{
  "cep": "01153000",
  "city": "São Paulo",
  "neighborhood": "Barra Funda",
  "service": "viacep",
  "state": "SP",
  "street": "Rua Vitorino Carmilo"
}
----------------------------------------------------
```

### Cenário de Timeout

```
Iniciando corrida de APIs para o CEP: 01153000

-> Tentando API: BrasilAPI
-> Tentando API: ViaCEP

----------------------------------------------------
❌ ERRO: Timeout. Nenhuma API respondeu em 1 segundo.
----------------------------------------------------
```

## 🧠 Estrutura do Código

-   **`struct APIResult`**: Uma estrutura simples para encapsular a resposta da API, contendo seu nome e os dados retornados.
-   **`func fetchCEP(...)`**: Esta função é o "trabalhador". Ela é executada como uma goroutine, faz a requisição HTTP para uma API específica e, em caso de sucesso, envia o `APIResult` para um channel compartilhado.
-   **`func main()`**: A função principal orquestra a aplicação. Ela cria o channel, inicia as duas goroutines (`fetchCEP`) e utiliza um `select` para aguardar o primeiro evento que ocorrer: ou o recebimento de um `APIResult` no channel ou o acionamento de um `time.After(1 * time.Second)`.