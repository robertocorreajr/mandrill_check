package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/rodaine/table"
	"mvdan.cc/xurls/v2"
)

func Search(args []string) {

	days := -7
	dateFrom := (time.Now().AddDate(0, 0, days)).Format("2006-01-02")
	dateTo := (time.Now()).Format("2006-01-02")
	key := os.Getenv("KEY")
	email := args[2]

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

	payload.Key = key
	payload.Query = email
	payload.DateFrom = dateFrom
	payload.DateTo = dateTo
	payload.Limit = limit

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)
	req, _ := http.NewRequest("POST", urlContentSearch, payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	fmt.Println(" response Status:", string("\033[32m"), resp.Status)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)

		var emails = make([]returnContentSearch, 0)
		err = json.Unmarshal(data, &emails)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			os.Exit(1)
		}

		// Start table structure
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		tbl := table.New("ID", "Subject", "Email To", "State", "Opens/Clicks", "Timestamp")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		// Finish table structure

		for i := range emails {
			openClicks := strconv.Itoa(emails[i].Opens) + "/" + strconv.Itoa(emails[i].Clicks)
			timeStamp := time.Unix(emails[i].Ts, 0)
			tbl.AddRow(emails[i].ID, emails[i].Subject, emails[i].Email, emails[i].State, openClicks, timeStamp.Format("2006-01-02 15:04:05"))
		}
		tbl.Print()
	}

}

func Info(args []string) {

	key := os.Getenv("KEY")
	id := args[2]

	payload.Key = key
	payload.ID = id

	payloadBuf := new(bytes.Buffer)

	json.NewEncoder(payloadBuf).Encode(payload)
	req, _ := http.NewRequest("POST", urlContentInfo, payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	fmt.Println(" response Status:", string("\033[32m"), resp.Status)

	data, _ := ioutil.ReadAll(resp.Body)

	var returnContent returnContentInfo

	err = json.Unmarshal(data, &returnContent)
	if err != nil {
		fmt.Printf("Error Unmarshal returnContent with error %s\n", err)
		os.Exit(1)
	}

	text := xurls.Strict().FindAllString(returnContent.Text, -1)
	returnContent.Text = strings.Trim(text[0], " ")

	// Start table structure
	tbl := table.New("ID", "Subject", "Email To", "Timestamp")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	// Finish table structure

	timeStamp := time.Unix(returnContent.TS, 0)
	tbl.AddRow(returnContent.ID, returnContent.Subject, returnContent.To.Email, timeStamp.Format("2006-01-02 15:04:05"))
	tbl.Print()

	fmt.Println(" ")

	fmt.Println(headerFmt("Link                                                                                                             "))
	fmt.Println(columnFmt(strings.Trim(returnContent.Text, " ")))
}

func Setup(args []string) {
	var key string = args[2]
	key = "KEY=\"" + key + "\""
	var file = ".env"
	var i = 0

	for i != 3 { // TODO: Copiar o binário para pasta /usr/bin e deletar do local original ou informar que pode ser apagado sei lá...
		f, err := os.Create(file)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		defer f.Close()

		_, err = f.WriteString(key)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		fmt.Print(".")
		time.Sleep(2 * time.Second)
		i++
	}
	fmt.Println("\nSetup finalizado com sucesso!!!")

}

func main() {

	// Realizar a leitura do arquivo .env e em caso de erro encerra a aplicação...
	// TODO: Adicionar help para criar arquivo na função setup do qual vai criar o arquivo .env e gravar a key lá
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Verificar a quantidade de argumentos, e caso seja menor que 2
	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1) // TODO: Retirar dúvida com o Rodrigo, se é necessário encerrar com código 1 e se não poderia encerrar com 0, já que o erro foi evitado.
	}

	args := os.Args
	action := strings.ToLower(args[1]) // TODO: Verificar a real necessidade de transformar as ações em letras minuculas por default.

	switch action {
	case "search", "s", "-s", "--s":
		fmt.Print("\nSearching")
		Search(args)

	case "info", "i", "-i", "--i":
		fmt.Println("Check Info")
		Info(args)

	case "key", "k", "-k", "--k":
		fmt.Println("Escolhido Key")

	case "setup":
		fmt.Println("Iniciando Setup...")
		Setup(args)
		// return key

	default:
		fmt.Println("Exibir Help")
		os.Exit(1)
	}

}
