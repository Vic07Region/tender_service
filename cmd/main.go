package main

import (
	"context"
	"fmt"
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
	router := handles.New(ctx, storage)
	//types := []string{}
	//listtenders, _ := storage.FetchListTenders(ctx, &database.ListTendersParams{
	//	Service_type: types,
	//	Offset:       0,
	//	Limit:        10,
	//})
	//var listTenders []Tender
	//for _, item := range listtenders {
	//	listTenders = append(listTenders, Tender{
	//		ID:          item.ID.String(),
	//		Name:        item.Name,
	//		Description: item.Description,
	//		Status:      item.Status,
	//		ServiceType: item.ServiceType,
	//		Version:     item.Version,
	//		CreatedAt:   item.CreatedAt,
	//	})
	//}
	//result, _ := json.Marshal(listTenders)
	//fmt.Println(string(result))
	http.HandleFunc("/api/ping", router.Ping)
	http.HandleFunc("/api/tenders", router.TenderList)
	http.HandleFunc("/api/tenders/new", router.NewTender)
	//http.HandleFunc("/users/refresh", router.RefreshToken)
	//http.HandleFunc("/users/refresh", router.RefreshToken)
	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
