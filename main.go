package main

import (
	"fmt"
	"log"
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
			limit, _ = strconv.Atoi(args[3])
		default:
			limit = 100
		}
	}

	var payload Payload

	payload.Key = key           // ou os.Getenv("KEY")
	payload.Query = email       // ou args[2]
	payload.DateFrom = dateFrom // Data inicial da busca
	payload.DateTo = dateTo     // Data final da busca
	payload.Limit = limit       // limit de resultados

	fmt.Println(payload)
	/*
		package main

		func fetchResponse(url string) string{
			resp, _ := http.Get(url)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			return string(body)
		}

		func main() {
			resp := fetchResponse("http://someurl.com")
			fmt.Println(resp)
		}
	*/

	/*
		func createBook(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			mu.Lock()
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)

			book.ID = strconv.Itoa(newId)
			books = append(books, book)
			newId++
			json.NewEncoder(w).Encode(book)
			mu.Unlock()
		}
	*/

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
