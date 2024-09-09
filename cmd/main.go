package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tender_service/internal/database"
	"tender_service/internal/handles"
	"time"
)

type Tender struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ServiceType string    `json:"serviceType"`
	Version     int32     `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

func main() {
	db, err := database.New("user=postgres dbname=tenders password=85428542 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	storage := database.NewService(db)
	handle := handles.New(ctx, storage)
	router := mux.NewRouter()
	router.HandleFunc("/api/ping", handle.Ping)
	router.HandleFunc("/api/tenders", handle.TenderList)
	router.HandleFunc("/api/tenders/new", handle.NewTender)
	router.HandleFunc("/api/tenders/my", handle.TenderMyList)
	router.HandleFunc("/api/tenders/{id}/status", handle.GetTenderStatus).Methods("GET")
	router.HandleFunc("/api/tenders/{id}/status", handle.ChangeTenderStatus).Methods("PUT")
	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
