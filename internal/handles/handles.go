package handles

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"tender_service/internal/database"
	"time"
)

var (
	InvalidParams    = "Некорректные параметры запроса"
	MethodNotAllowed = "Метод не разрешен"
	IpParseError     = "Не удалось определить IP-адрес"
	IsNotResponsible = "Пользователь не является ответственным в организации"
	UserNotFound     = "Пользователь с таким именем не существует"
	FieldRequired    = " является обязательным полем"
)

type Routes struct {
	dbq *database.Queries
	ctx context.Context
	mu  sync.Mutex
}

func New(ctx context.Context, dbq *database.Queries) *Routes {
	return &Routes{
		dbq: dbq,
		ctx: ctx,
	}
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (rts *Routes) Ping(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodGet {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode("ok")
	w.Write([]byte("ok"))
}

type Tender struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ServiceType string    `json:"serviceType"`
	Version     int32     `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (rts *Routes) TenderList(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodGet {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	queryParams := r.URL.Query()
	var limit, offset int32
	var service_type []string
	limit_param := queryParams.Get("limit")
	offset_param := queryParams.Get("offset")
	type_param := queryParams.Get("service_type")

	if limit_param == "" || !isNumeric(limit_param) {
		limit = 5
	} else {
		tl, _ := strconv.Atoi(limit_param)
		limit = int32(tl)
	}

	if offset_param == "" || !isNumeric(offset_param) {
		offset = 0
	} else {
		tl, _ := strconv.Atoi(offset_param)
		offset = int32(tl)
	}
	if type_param != "" {
		service_type = strings.Split(type_param, ",")
	}

	listtenders, err := rts.dbq.PublishedListTenders(rts.ctx, &database.ListTendersParams{
		Service_type: service_type,
		Offset:       offset,
		Limit:        limit,
	})
	if err != nil {
		fmt.Println(err)
	}
	var listTenders []Tender
	for _, item := range listtenders {
		listTenders = append(listTenders, Tender{
			ID:          item.ID.String(),
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
			ServiceType: item.ServiceType,
			Version:     item.Version,
			CreatedAt:   item.CreatedAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listTenders)
}

type TenderParams struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	Status          string `json:"status"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

func (rts *Routes) NewTender(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodPost {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	var params TenderParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		err_response["reason"] = InvalidParams
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	if params.Name == "" {
		err_response["reason"] = "name" + FieldRequired
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	if params.ServiceType == "" {
		err_response["reason"] = "serviceType" + FieldRequired
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	if params.Status == "" {
		err_response["reason"] = "status" + FieldRequired
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	//todo checking uuid format
	user_id, err := rts.dbq.FetchUserID(rts.ctx, params.CreatorUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			err_response["reason"] = UserNotFound
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		http.Error(w, err.Error()+"u", http.StatusInternalServerError)
		return
	}
	valid, err := rts.dbq.IsResponsible(rts.ctx, params.OrganizationId, user_id)
	if err != nil {
		http.Error(w, err.Error()+"s", http.StatusInternalServerError)
		return
	}
	if !valid {
		err_response["reason"] = IsNotResponsible
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	tender, err := rts.dbq.CreateTender(rts.ctx, database.CreateTenderParams{
		OrganizationID: params.OrganizationId,
		CreatorID:      user_id,
		Status:         params.Status,
		ServiceType:    params.ServiceType,
		Name:           params.Name,
		Description:    params.Description,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tenderResponse := Tender{
		ID:          tender.ID,
		Name:        params.Name,
		Description: params.Description,
		Status:      params.Status,
		ServiceType: params.ServiceType,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenderResponse)
}
