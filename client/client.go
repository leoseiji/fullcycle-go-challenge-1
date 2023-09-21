package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type DolarExchange struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func GetCotacao() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	log.Println("GetCotacao started")
	defer log.Println("GetCotacao ended")

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	var dolarExchange DolarExchange
	err = json.Unmarshal(res, &dolarExchange)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("USD-BRL Bid: " + dolarExchange.Usdbrl.Bid)
}

func main() {
	GetCotacao()
}
