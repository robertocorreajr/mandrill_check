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

const URL string = "https://mandrillapp.com/api/1.0/"

type Payload struct {
	Key      string `json:"key"`
	ID       string `json:"id"`
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
	// fmt.Println("Date From:", dateFrom)
	// fmt.Println("Date To:  ", dateTo)

	key := os.Getenv("KEY")
	email := args[2] // os.Getenv("EMAIL_TEST") // args[2]

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

	payload := Payload{
		Key:      key,
		Query:    email,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Limit:    limit,
	}

	// fmt.Println(payload)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)
	req, _ := http.NewRequest("POST", url, payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(data))

		type Return struct {
			ID      string `json:"_id"`
			Subject string `json:"subject"`
			Email   string `json:"email"`
			State   string `json:"state"`
			Opens   int    `json:"opens"`
			Clicks  int    `json:"clicks"`
			Ts      int64  `json:"ts"`
		}

		var emails = make([]Return, 0)
		err = json.Unmarshal(data, &emails)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			os.Exit(1)
		}
		// fmt.Println("#####################################################")
		// Start table structure
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		tbl := table.New("ID", "Subject", "Email To", "State", "Opens/Clicks", "Timestamp")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		// Finish table structure

		// for i := range emails {
		// 	fmt.Println(emails[i].ID, emails[i].Subject, emails[i].Email, emails[i].State, emails[i].Opens, emails[i].Clicks, time.Unix(emails[i].Ts, 0))
		// }

		for i := range emails {
			openClicks := strconv.Itoa(emails[i].Opens) + "/" + strconv.Itoa(emails[i].Clicks)
			timeStamp := time.Unix(emails[i].Ts, 0)
			tbl.AddRow(emails[i].ID, emails[i].Subject, emails[i].Email, emails[i].State, openClicks, timeStamp.Format("2006-01-02 15:03:04"))
		}

		tbl.Print()

		// fmt.Println(emails)
	}
	// Print the body to the stdout
	// io.Copy(os.Stdout, resp.Body)

	// jsonData := make([]Payload, 0)
	// jsonValue, _ := json.Marshal(jsonData)
	// resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	// if err != nil {
	// 	fmt.Printf("The HTTP request failed with error %s\n", err)
	// } else {
	// 	data, _ := ioutil.ReadAll(resp.Body)
	// 	fmt.Println(string(data))
	// }

}

// func GetMessageContent(x string, y string) []string{

// }

func Info(args []string) {

	// urlInfo := URL + "messages/info.json"
	urlContent := URL + "messages/content.json"
	key := os.Getenv("KEY")
	id := args[2]

	payload := Payload{
		Key: key,
		ID:  id,
	}

	// fmt.Println(payload)

	payloadBuf := new(bytes.Buffer)

	json.NewEncoder(payloadBuf).Encode(payload)
	req, _ := http.NewRequest("POST", urlContent, payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	type ReturnContent struct {
		ID      string `json:"_id"`
		Subject string `json:"subject"`
		To      struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"to"`
		TS   int64  `json:"ts"`
		Text string `json:"text"`
	}

	data, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(data))
	var returnContent ReturnContent
	err = json.Unmarshal(data, &returnContent)
	if err != nil {
		fmt.Printf("Error Unmarshal returnContent with error %s\n", err)
		os.Exit(1)
	}

	text := xurls.Strict().FindAllString(returnContent.Text, -1)
	returnContent.Text = strings.Trim(text[0], " ")
	// fmt.Println(returnContent.Text)

	// Start table structure
	tbl := table.New("ID", "Subject", "Email To", "Timestamp")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	// Finish table structure

	timeStamp := time.Unix(returnContent.TS, 0)
	tbl.AddRow(returnContent.ID, returnContent.Subject, returnContent.To.Email, timeStamp.Format("2006-01-02 15:03:04"))
	tbl.Print()

	fmt.Println(" ")

	fmt.Println(headerFmt("Link                                                                                                             "))
	fmt.Println(columnFmt(strings.Trim(returnContent.Text, " ")))
}

func Info2(args []string) {
	// urlContent := URL + "messages/content.json"
	urlInfo := URL + "messages/info.json"
	key := os.Getenv("KEY")
	id := args[2]

	payload := Payload{
		Key: key,
		ID:  id,
	}

	fmt.Println(payload)

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)
	req, _ := http.NewRequest("POST", urlInfo, payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))
		type ReturnInfo struct {
			SMTPEvents []struct {
				Ts            int    `json:"ts"`
				Type          string `json:"type"`
				Diag          string `json:"diag"`
				SourceIP      string `json:"source_ip"`
				DestinationIP string `json:"destination_ip"`
				Size          int    `json:"size"`
			} `json:"smtp_events"`
		}

		type ReturnContent struct {
			Subject string `json:"subject"`
			To      string `json:"To"`
			Ts      int64  `json:"ts"`
		}
		// var returnInfo = make([]ReturnInfo, 0)
		// var returnContent = make([]ReturnContent, 0)

		// for smtpEvents in

		// err = json.Unmarshal(data, &emails)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			os.Exit(1)
		}
		// fmt.Println("#####################################################")
		// Start table structure
		// headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		// columnFmt := color.New(color.FgYellow).SprintfFunc()
		// tbl := table.New("ID", "Subject", "Email To", "State", "Opens/Clicks", "Timestamp")
		// tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		// tbl.WithFirstColumnFormatter(columnFmt).WithWriter()
		// Finish table structure

		// for i := range emails {
		// 	fmt.Println(emails[i].ID, emails[i].Subject, emails[i].Email, emails[i].State, emails[i].Opens, emails[i].Clicks, time.Unix(emails[i].Ts, 0))
		// }

		// for i := range emails {
		// 	openClicks := strconv.Itoa(emails[i].Opens) + "/" + strconv.Itoa(emails[i].Clicks)
		// 	timeStamp := time.Unix(emails[i].Ts, 0)
		// 	tbl.AddRow(emails[i].ID, emails[i].Subject, emails[i].Email, emails[i].State, openClicks, timeStamp.Format("2006-01-02 15:03:04"))
		// }

		// tbl.Print()

		// fmt.Println(emails)
	}

}

func Setup(args []string) {
	var key string = args[2]
	key = "KEY=\"" + key + "\""
	var file = ".env2" // TODO: Alterar para .env antes de colocar em produção. Acrescentar a criação de um diretório no home do usuário ~/.mandrilc/.env
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

	// key := os.Getenv("KEY")

	// Verificar a quantidade de argumentos, e caso seja menor que 2
	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1) // TODO: Retirar dúvida com o Rodrigo, se é necessário encerrar com código 1 e se não poderia encerrar com 0, já que o erro foi evitado.
	}

	args := os.Args
	action := strings.ToLower(args[1]) // TODO: Verificar a real necessidade de transformar as ações em letras minuculas por default.

	switch action {
	case "search", "s", "-s", "--s":
		// argsSplit := strings.Split(args[2], " ")
		// email := args[2]
		fmt.Printf("Escolhido Search:")

		Search(args)

	case "info", "i", "-i", "--i":
		fmt.Println("Escolhido Info")
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
