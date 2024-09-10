package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"tender_service/internal/database"
	"time"
)

var (
	UnknowError       = fmt.Errorf("unknown server error")
	CreateTenderError = fmt.Errorf("Ошибка в создании тендера")
	UserNotFound      = fmt.Errorf("Пользователь с таким именем не существует")
	IsNotResponsible  = fmt.Errorf("Пользователь не является ответственным в организации")
	TenderNotFound    = fmt.Errorf("Тендер с таким id не существует")
)

type Service struct {
	query *database.Queries
	mu    sync.Mutex
}

func New(query *database.Queries) *Service {
	return &Service{
		query: query,
	}
}

func (s *Service) isResponsibleUser(ctx context.Context, org_id, user_id string) error {
	valid, err := s.query.IsResponsible(ctx, org_id, user_id)
	if err != nil {
		log.Println("CreateNewTender: check IsResponsible error -", err)
		return UnknowError
	}
	if !valid {
		return IsNotResponsible
	}
	return nil
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

type ListTendersRequest struct {
	Service_type []string
	Offset       int32
	Limit        int32
}

func (s *Service) FetchPublishedTenders(ctx context.Context, params ListTendersRequest) ([]Tender, error) {
	listtenders, err := s.query.PublishedListTenders(ctx, database.ListTendersParams{
		Service_type: params.Service_type,
		Offset:       params.Offset,
		Limit:        params.Limit,
	})
	if err != nil {
		log.Printf("FetchPublisgTender error: %s", err)
		return nil, CreateTenderError
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
	return listTenders, err
}

type TenderParams struct {
	Name            string
	Description     string
	ServiceType     string
	Status          string
	OrganizationId  string
	CreatorUsername string
}

func (s *Service) CreateNewTender(ctx context.Context, params TenderParams) (*Tender, error) {

	user_id, err := s.query.FetchUserID(ctx, params.CreatorUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("CreateNewTender: Fetching User id error -", err)
		return nil, UnknowError
	}

	err = s.isResponsibleUser(ctx, params.OrganizationId, user_id)
	if err != nil {
		return nil, err
	}

	tender, err := s.query.CreateTender(ctx, database.CreateTenderParams{
		OrganizationID: params.OrganizationId,
		CreatorID:      user_id,
		Status:         params.Status,
		ServiceType:    params.ServiceType,
		Name:           params.Name,
		Description:    params.Description,
	})
	if err != nil {
		log.Println("CreateNewTender: CreateTender error -", err)
		return nil, UnknowError
	}

	return &Tender{
		ID:          tender.ID,
		Name:        params.Name,
		Description: params.Description,
		Status:      params.Status,
		ServiceType: params.ServiceType,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt,
	}, nil
}

type ListMyTendersRequest struct {
	Username string
	Offset   int32
	Limit    int32
}

func (s *Service) FetchMyTenders(ctx context.Context, params ListMyTendersRequest) ([]Tender, error) {
	user_id, err := s.query.FetchUserID(ctx, params.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("FetchMyTenders: fetching user id err -", err)
		return nil, UnknowError
	}

	listtenders, err := s.query.MyListTenders(ctx, &database.MyListTendersParams{
		User_id: user_id,
		Offset:  params.Offset,
		Limit:   params.Limit,
	})
	if err != nil {
		log.Println("FetchMyTenders: MyListTenders err -", err)
		return nil, UnknowError
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
	return listTenders, nil
}

func (s *Service) FetchTenderStatus(ctx context.Context, username, tender_id string) (string, error) {
	user_id, err := s.query.FetchUserID(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", UserNotFound
		}
		log.Println("FetchTenderStatus: FetchUserID err -", err)
		return "", UnknowError
	}
	tender, err := s.query.GetTender(ctx, tender_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", TenderNotFound
		}
		log.Println("FetchTenderStatus: GetTender err -", err)
		return "", UnknowError
	}

	err = s.isResponsibleUser(ctx, tender.OrganizationID.String(), user_id)
	if err != nil {
		return "", err
	}

	tender_status, err := s.query.CheckTenderStatus(ctx, tender_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", TenderNotFound
		}
		log.Println("FetchTenderStatus: CheckTenderStatus err -", err)
		return "", UnknowError
	}
	return tender_status, nil
}

type EditTenderStatusRequest struct {
	Username   string
	Tender_id  string
	New_status string
}

func (s *Service) EditTenderStatus(ctx context.Context, param EditTenderStatusRequest) (*Tender, error) {
	user_id, err := s.query.FetchUserID(ctx, param.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("EditTenderStatus: FetchUserID err -", err)
		return nil, UnknowError
	}
	tender, err := s.query.GetTender(ctx, param.Tender_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("EditTenderStatus: GetTender err -", err)
		return nil, UnknowError
	}

	err = s.isResponsibleUser(ctx, tender.OrganizationID.String(), user_id)
	if err != nil {
		return nil, err
	}

	new_tender, err := s.query.ChangeTenderStatus(ctx, tender.ID.String(), param.New_status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("EditTenderStatus: ChangeTenderStatus err -", err)
		return nil, UnknowError
	}

	return &Tender{
		ID:          new_tender.ID.String(),
		Name:        new_tender.Name,
		Description: new_tender.Description,
		Status:      new_tender.Status,
		ServiceType: new_tender.ServiceType,
		Version:     new_tender.Version,
		CreatedAt:   new_tender.CreatedAt,
	}, nil
}

type EditTenderRequest struct {
	Username     string
	Tender_id    string
	Name         string
	Description  string
	Service_type string
}

func (s *Service) EditTender(ctx context.Context, params EditTenderRequest) (*Tender, error) {
	user_id, err := s.query.FetchUserID(ctx, params.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("EditTenderStatus: FetchUserID err -", err)
		return nil, UnknowError
	}
	tender, err := s.query.GetTender(ctx, params.Tender_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("EditTenderStatus: GetTender err -", err)
		return nil, UnknowError
	}

	err = s.isResponsibleUser(ctx, tender.OrganizationID.String(), user_id)
	if err != nil {
		return nil, err
	}
	new_tender, err := s.query.CreateTenderWithTX(ctx, database.CreateTenderWithTxParam{
		ChangeTenderParam: database.TenderChangeParam{
			Tender_id:    tender.ID.String(),
			Name:         params.Name,
			Description:  params.Description,
			Service_type: params.Service_type,
		},
		TenderHistoryParam: database.CreateTenderHistoryParams{
			Tender_id:   tender.ID.String(),
			Creator_id:  user_id,
			ServiceType: tender.ServiceType,
			Name:        tender.Name,
			Description: tender.Description,
			OldVersion:  tender.Version,
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("EditTenderStatus: CreateTenderWithTX err -", err)
		return nil, UnknowError
	}
	return &Tender{
		ID:          new_tender.ID.String(),
		Name:        new_tender.Name,
		Description: new_tender.Description,
		Status:      new_tender.Status,
		ServiceType: new_tender.ServiceType,
		Version:     new_tender.Version,
		CreatedAt:   new_tender.CreatedAt,
	}, nil
}
