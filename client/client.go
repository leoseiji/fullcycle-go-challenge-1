package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type DolarExchange struct {
	Bid string `json:"bid"`
}

func GetCotacao() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	log.Println("Client GetCotacao started")
	defer log.Println("Client GetCotacao ended")

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalln("error to create cotacao endpoint", err)
	}
	client := http.DefaultClient
	res, resErr := client.Do(req)
	if resErr != nil {
		log.Fatalln("error to call cotacao endpoint", resErr)
	}
	defer res.Body.Close()
	resBody, resBodyErr := io.ReadAll(res.Body)
	if resBodyErr != nil {
		log.Fatalln("error to read body cotacao endpoint", resBodyErr)
	}
	var dolarExchange DolarExchange
	if err = json.Unmarshal(resBody, &dolarExchange); err != nil {
		log.Fatalln("error to unmarshall body cotacao endpoint", err)
	}

	log.Printf("Response Bid: %s", dolarExchange.Bid)

	if err := writeCotacaoTxt(dolarExchange); err != nil {
		log.Fatalln("error to write cotacao text", err)
	}
}

func main() {
	GetCotacao()
}

func writeCotacaoTxt(dolarExchange DolarExchange) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println("error to create file", err)
		return err
	}
	defer file.Close()
	_, writeErr := file.WriteString(fmt.Sprintf("Dólar: %s", dolarExchange.Bid))
	if writeErr != nil {
		log.Println("error to write file", writeErr)
		return writeErr
	}
	return nil
}
