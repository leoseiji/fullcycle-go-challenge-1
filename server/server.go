package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbDolar *gorm.DB

func main() {
	var dbErr error
	dbDolar, dbErr = gorm.Open(sqlite.Open("dolar.db"), &gorm.Config{})
	if dbErr != nil {
		log.Fatalln("failed to connect database", dbErr)
	}
	dbDolar.AutoMigrate(&DolarExchangeDB{})

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", CotacaoHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln("error starting server", err)
	}
}

type DolarExchange struct {
	DolarExchangeDB DolarExchangeDB `json:"USDBRL"`
}

type DolarExchangeDB struct {
	gorm.Model
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()
	log.Println("CotacaoHandler started")
	defer log.Println("CotacaoHandler ended")

	dolarExchange, err := GetCotacao(ctx)
	if err != nil {
		log.Println("get Cotacao error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = SaveDatabase(ctx, dolarExchange)
	if err != nil {
		log.Println("save Database error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dolarExchange.DolarExchangeDB)
}

func GetCotacao(ctx context.Context) (*DolarExchange, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println("get cotacao endpoint error", err)
		return nil, err
	}

	res, resErr := http.DefaultClient.Do(req)
	if condition := resErr != nil; condition {
		log.Println("response error", resErr)
		return nil, resErr
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Println("read all error", readErr)
		return nil, readErr
	}

	var dolarExchange DolarExchange
	err = json.Unmarshal(body, &dolarExchange)
	if err != nil {
		log.Println("unmarshall error", err)
		return nil, err
	}
	return &dolarExchange, nil
}

func SaveDatabase(ctx context.Context, dolarExchange *DolarExchange) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	log.Println("SaveDatabase started")
	defer log.Println("SaveDatabase ended")

	insertErr := insertDollarExchange(ctx, dolarExchange)
	if insertErr != nil {
		log.Println("insert database error", insertErr)
		return insertErr
	}

	return nil
}

func insertDollarExchange(ctx context.Context, dolarExchange *DolarExchange) error {
	log.Println("insertDollarExchange started")
	defer log.Println("insertDollarExchange ended")

	dbErr := dbDolar.WithContext(ctx).Create(&dolarExchange.DolarExchangeDB).Error
	if dbErr != nil {
		log.Println("error inserting db", dbErr)
		return dbErr
	}
	return nil
}
