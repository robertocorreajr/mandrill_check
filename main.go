package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const URL = "https://mandrillapp.com/api/1.0/"

type Payload struct {
	Key      string `json:"key"`
	Query    string `json:"query"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Limit    int    `json:"limit"`
}

// var payloads = make([]Payload, 0)

func Search(args []string) {

	url := URL + "messages/search.json"

	days := -7
	dateFrom := (time.Now().AddDate(0, 0, days)).Format("2006-01-02")
	dateTo := (time.Now()).Format("2006-01-02")
	fmt.Println("Date From:", dateFrom)
	fmt.Println("Date To:  ", dateTo)

	key := os.Getenv("KEY")
	email := os.Getenv("EMAIL_TEST") // args[2]

	var limit int

	if args[3] != "" && args[4] != "" {
		action := args[3]
		switch action {
		case "limit", "l", "-l", "--l":

			var err error
			limit, err = strconv.Atoi(args[4])

			if err != nil {
				fmt.Println("First input parameter must be integer")
				os.Exit(1)
			}

		default:
			limit = 100
		}
	}

	payload := &Payload{
		Key:      key,
		Query:    email,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Limit:    limit,
	}

	fmt.Println(payload)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)
	req, _ := http.NewRequest("POST", url, payloadBuf)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// Print the body to the stdout
	io.Copy(os.Stdout, resp.Body)

}

func Setup() {

}

func main() {

	// Realizar a leitura do arquivo .env e em caso de erro encerra a aplicação...
	// @TODO: Adicionar help para criar arquivo na função setup do qual vai criar o arquivo .env e gravar a key lá
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("KEY")

	// Verificar a quantidade de argumentos, e caso seja menor que 2
	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1) // @TODO: Retirar dúvida com o Rodrigo, se é necessário encerrar com código 1 e se não poderia encerrar com 0, já que o erro foi evitado.
	}

	args := os.Args
	action := strings.ToLower(args[1]) // @TODO: Verificar a real necessidade de transformar as ações em letras minuculas por default.

	switch action {
	case "search", "s", "-s", "--s":
		// argsSplit := strings.Split(args[2], " ")
		email := args[2]
		fmt.Printf("Escolhido Search:\nKey: %s\nemail: %s\n", key, email)

		Search(args)

	case "info", "i", "-i", "--i":
		fmt.Println("Escolhido Info")

	case "key", "k", "-k", "--k":
		fmt.Println("Escolhido Key")

	case "setup":
		fmt.Println("Iniciando Setup...")

	default:
		fmt.Println("Exibir Help")
		os.Exit(1)
	}

}
