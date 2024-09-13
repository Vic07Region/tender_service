package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"tender_service/internal/database"
	"tender_service/internal/handles"
	"tender_service/internal/service"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Ошибка загрузки файла .env:", err)
		return
	}

	//todo убрать godotenv
	dbhost := os.Getenv("POSTGRES_HOST")
	dbname := os.Getenv("POSTGRES_DATABASE")
	dbport := os.Getenv("POSTGRES_PORT")
	dbusername := os.Getenv("POSTGRES_USERNAME")
	dbpassword := os.Getenv("POSTGRES_PASSWORD")
	server_addres := os.Getenv("SERVER_ADDRESS")
	conn_str := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=verify-full sslrootcert=./root.crt",
		dbname,
		dbusername,
		dbpassword,
		dbhost,
		dbport,
	)
	fmt.Println(conn_str)
	db, err := database.New("user=postgres dbname=tenders password=85428542 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	storage := database.NewService(db)
	srv := service.New(storage)
	handle := handles.New(ctx, srv)
	router := mux.NewRouter()
	router.HandleFunc("/api/ping", handle.Ping).Methods("GET")
	router.HandleFunc("/api/tenders", handle.TenderList)
	router.HandleFunc("/api/tenders/new", handle.NewTender)
	router.HandleFunc("/api/tenders/my", handle.TenderMyList)
	router.HandleFunc("/api/tenders/{id}/status", handle.GetTenderStatus).Methods("GET")
	router.HandleFunc("/api/tenders/{id}/status", handle.ChangeTenderStatus).Methods("PUT")
	router.HandleFunc("/api/tenders/{id}/edit", handle.ChangeTender).Methods("PATCH")
	router.HandleFunc("/api/tenders/{id}/rollback/{version}", handle.RollbackTender).Methods("PUT")
	router.HandleFunc("/api/bids/new", handle.BidNew).Methods("POST")
	router.HandleFunc("/api/bids/{tenderID}/list", handle.BidsTender).Methods("GET")
	router.HandleFunc("/api/bids/my", handle.MyBids).Methods("GET")
	router.HandleFunc("/api/bids/{tenderID}/status", handle.BidStatus).Methods("GET")
	router.HandleFunc("/api/bids/{tenderID}/status", handle.BidStatus).Methods("PUT")
	fmt.Println("Сервер запущен на ", server_addres)
	log.Fatal(http.ListenAndServe(server_addres, router))
}
