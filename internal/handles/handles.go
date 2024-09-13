package handles

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"tender_service/internal/service"
	"tender_service/internal/utils"
	"time"
)

var (
	InvalidParams    = "Некорректные параметры запроса"
	MethodNotAllowed = "Метод не разрешен"
	IsNotResponsible = "Пользователь не является ответственным в организации"
	UserNotFound     = "Пользователь с таким именем не существует"
	TenderNotFound   = "Тендер с таким id не существует"
	FieldRequired    = " является обязательным для заполнения"
)

type Handle struct {
	srv *service.Service
	ctx context.Context
}

func New(ctx context.Context, s *service.Service) *Handle {
	return &Handle{
		srv: s,
		ctx: ctx,
	}
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

func (h *Handle) Ping(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodGet {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode("ok")
	w.Write([]byte("ok"))
}

func (h *Handle) TenderList(w http.ResponseWriter, r *http.Request) {
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

	if limit_param == "" || !utils.IsNumeric(limit_param) {
		limit = 5
	} else {
		tl, _ := strconv.Atoi(limit_param)
		limit = int32(tl)
	}

	if offset_param == "" || !utils.IsNumeric(offset_param) {
		offset = 0
	} else {
		tl, _ := strconv.Atoi(offset_param)
		offset = int32(tl)
	}
	if type_param != "" {
		service_type = strings.Split(type_param, ",")
	}
	tender_list_request := service.ListTendersRequest{
		Service_type: service_type,
		Offset:       offset,
		Limit:        limit,
	}

	listTenders, err := h.srv.FetchPublishedTenders(h.ctx, tender_list_request)
	if err != nil {
		if err == service.CreateTenderError {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listTenders)
}

func (h *Handle) NewTender(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodPost {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	var params service.TenderParams

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
	_, err = uuid.Parse(params.OrganizationId)
	if err != nil {
		err_response["reason"] = InvalidParams + ": неверный формат поля organizationId"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	tenderResponse, err := h.srv.CreateNewTender(h.ctx, params)
	if err != nil {
		if err == service.IsNotResponsible {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = service.UnknowError
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenderResponse)
}

func (h *Handle) TenderMyList(w http.ResponseWriter, r *http.Request) {
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
	limit_param := queryParams.Get("limit")
	offset_param := queryParams.Get("offset")
	username := queryParams.Get("username")

	if limit_param == "" || !utils.IsNumeric(limit_param) {
		limit = 5
	} else {
		tl, _ := strconv.Atoi(limit_param)
		limit = int32(tl)
	}

	if offset_param == "" || !utils.IsNumeric(offset_param) {
		offset = 0
	} else {
		tl, _ := strconv.Atoi(offset_param)
		offset = int32(tl)
	}

	listTenders, err := h.srv.FetchMyTenders(h.ctx, service.ListMyTendersRequest{
		Username: username,
		Offset:   offset,
		Limit:    limit,
	})
	if err != nil {
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listTenders)
}

func (h *Handle) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
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
	username := queryParams.Get("username")

	pathParts := r.URL.Path[len("/api/tenders/"):]
	tender_id := strings.Split(pathParts, "/")[0]

	_, err := uuid.Parse(tender_id)
	if err != nil {
		err_response["reason"] = InvalidParams + ": некорректный формат id тендера"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	tender_status, err := h.srv.FetchTenderStatus(h.ctx, username, tender_id)
	if err != nil {
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.IsNotResponsible {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.TenderNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender_status)
}

func (h *Handle) ChangeTenderStatus(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodPut {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	allowedValues := []string{"Created", "Published", "Closed"}

	pathParts := r.URL.Path[len("/api/tenders/"):]

	tender_id := strings.Split(pathParts, "/")[0]
	_, err := uuid.Parse(tender_id)
	if err != nil {
		err_response["reason"] = InvalidParams + ": некорректный формат id тендера"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	queryParams := r.URL.Query()

	username := queryParams.Get("username")
	newstatus := queryParams.Get("status")

	if newstatus == "" {
		err_response["reason"] = "status (query param)" + FieldRequired
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	if !utils.CheckString(newstatus, allowedValues) {
		err_response["reason"] = InvalidParams + ": неверное значение status"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	tender, err := h.srv.EditTenderStatus(h.ctx, service.EditTenderStatusRequest{
		Username:   username,
		Tender_id:  tender_id,
		New_status: newstatus,
	})
	if err != nil {
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.IsNotResponsible {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.TenderNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

type TenderChangeRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Service_type string `json:"serviceType"`
}

func (h *Handle) ChangeTender(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodPatch {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	pathParts := r.URL.Path[len("/api/tenders/"):]

	tender_id := strings.Split(pathParts, "/")[0]
	_, err := uuid.Parse(tender_id)
	if err != nil {
		err_response["reason"] = InvalidParams + ": некорректный формат id тендера"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	queryParams := r.URL.Query()

	username := queryParams.Get("username")

	var param TenderChangeRequest
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		err_response["reason"] = InvalidParams + err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	new_tender, err := h.srv.EditTender(h.ctx, service.EditTenderRequest{
		Username:     username,
		Tender_id:    tender_id,
		Name:         param.Name,
		Description:  param.Description,
		Service_type: param.Service_type,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(new_tender)
}

func (h *Handle) RollbackTender(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodPut {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	pathParts := r.URL.Path[len("/api/tenders/"):]

	tender_id := strings.Split(pathParts, "/")[0]
	version_param := strings.Split(pathParts, "/")[2]

	if !utils.IsNumeric(version_param) {
		err_response["reason"] = InvalidParams + ": некорректный формат версии тендера"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	version, _ := strconv.Atoi(version_param)

	_, err := uuid.Parse(tender_id)
	if err != nil {
		err_response["reason"] = InvalidParams + ": некорректный формат id тендера"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	queryParams := r.URL.Query()

	username := queryParams.Get("username")
	tender, err := h.srv.RollbackTender(h.ctx, service.RollbackTenderRequest{
		Username:  username,
		Tender_id: tender_id,
		Version:   int32(version),
	})
	if err != nil {
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.IsNotResponsible {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.TenderNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.TenderHistoryNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

type NewBidParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderId    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorId    string `json:"authorId"`
}

func (h *Handle) BidNew(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodPost {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	var params NewBidParam

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
	if params.Description == "" {
		err_response["reason"] = "Description" + FieldRequired
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	if params.TenderId == "" {
		err_response["reason"] = "tenderId" + FieldRequired
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	if params.AuthorType == "" {
		err_response["reason"] = "authorType" + FieldRequired
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	if params.AuthorId == "" {
		err_response["reason"] = "authorId" + FieldRequired
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	offer, err := h.srv.CreateNewBid(h.ctx, service.CreateBidParam{
		Name:        params.Name,
		Description: params.Description,
		TenderId:    params.TenderId,
		AuthorType:  params.AuthorType,
		AuthorId:    params.AuthorId,
	})
	if err != nil {
		if err == service.NotAllowValue {
			err_response["reason"] = err.Error() + " authorType"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.TenderNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.IsResponsible {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = service.UnknowError
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(offer)
}

func (h *Handle) MyBids(w http.ResponseWriter, r *http.Request) {
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
	limit_param := queryParams.Get("limit")
	offset_param := queryParams.Get("offset")
	username := queryParams.Get("username")

	if limit_param == "" || !utils.IsNumeric(limit_param) {
		limit = 5
	} else {
		tl, _ := strconv.Atoi(limit_param)
		limit = int32(tl)
	}

	if offset_param == "" || !utils.IsNumeric(offset_param) {
		offset = 0
	} else {
		tl, _ := strconv.Atoi(offset_param)
		offset = int32(tl)
	}
	listoffer, err := h.srv.ListMyBids(h.ctx, service.ListMyBidsRequest{
		Username: username,
		Offset:   offset,
		Limit:    limit,
	})
	if err != nil {
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listoffer)
}

func (h *Handle) BidsTender(w http.ResponseWriter, r *http.Request) {
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
	limit_param := queryParams.Get("limit")
	offset_param := queryParams.Get("offset")
	username := queryParams.Get("username")

	pathParts := r.URL.Path[len("/api/bids/"):]
	tender_id := strings.Split(pathParts, "/")[0]

	_, err := uuid.Parse(tender_id)
	if err != nil {
		err_response["reason"] = InvalidParams + ": некорректный формат id тендера"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	if limit_param == "" || !utils.IsNumeric(limit_param) {
		limit = 5
	} else {
		tl, _ := strconv.Atoi(limit_param)
		limit = int32(tl)
	}

	if offset_param == "" || !utils.IsNumeric(offset_param) {
		offset = 0
	} else {
		tl, _ := strconv.Atoi(offset_param)
		offset = int32(tl)
	}
	listoffer, err := h.srv.TenderListBids(h.ctx, service.TenderListBidsRequest{
		Tender_id: tender_id,
		Username:  username,
		Offset:    offset,
		Limit:     limit,
	})
	if err != nil {
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.TenderNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.IsNotResponsible {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listoffer)
}

func (h *Handle) BidStatus(w http.ResponseWriter, r *http.Request) {
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
	username := queryParams.Get("username")

	pathParts := r.URL.Path[len("/api/bids/"):]
	bid_id := strings.Split(pathParts, "/")[0]

	_, err := uuid.Parse(bid_id)
	if err != nil {
		err_response["reason"] = InvalidParams + ": некорректный формат id тендера"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	bid_status, err := h.srv.GetBidStatus(h.ctx, service.GetBidStatus{
		Username: username,
		BidID:    bid_id,
	})
	if err != nil {
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.BidNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.IsNotResponsible {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid_status)
}

func (h *Handle) ChangeBidStatus(w http.ResponseWriter, r *http.Request) {
	err_response := map[string]interface{}{
		"reason": "",
	}
	if r.Method != http.MethodPut {
		err_response["reason"] = MethodNotAllowed
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	allowed_value := []string{"Created", "Published", "Canceled"}

	queryParams := r.URL.Query()
	username := queryParams.Get("username")
	status := queryParams.Get("status")

	pathParts := r.URL.Path[len("/api/bids/"):]
	bid_id := strings.Split(pathParts, "/")[0]

	_, err := uuid.Parse(bid_id)
	if err != nil {
		err_response["reason"] = InvalidParams + ": некорректный формат id тендера"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	if !utils.CheckString(status, allowed_value) {
		err_response["reason"] = InvalidParams + ": неверное значение status"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}

	bid, err := h.srv.ChangeBidStatus(h.ctx, service.ChangeBidStatus{
		Username: username,
		BidID:    bid_id,
		Status:   status,
	})
	if err != nil {
		if err == service.UserNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.BidNotFound {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		if err == service.IsNotResponsible {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err_response)
			return
		}

		if err == service.BidCanceled {
			err_response["reason"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err_response)
			return
		}
		err_response["reason"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}
