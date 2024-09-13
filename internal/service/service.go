package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"tender_service/internal/database"
	"tender_service/internal/utils"
	"time"
)

var (
	UnknowError           = fmt.Errorf("unknown server error")
	CreateTenderError     = fmt.Errorf("Ошибка в создании тендера")
	UserNotFound          = fmt.Errorf("Пользователь с таким именем не существует")
	IsNotResponsible      = fmt.Errorf("Пользователь не является ответственным в организации")
	IsResponsible         = fmt.Errorf("Пользователь является ответственным в организации")
	TenderNotFound        = fmt.Errorf("Тендер с таким id не существует")
	TenderHistoryNotFound = fmt.Errorf("Версия тендера с таким номером не существует")
	NotAllowValue         = fmt.Errorf("Недопустимое значение поля")
	BidNotFound           = fmt.Errorf("Предложение с таким id не существует")
	IsNotAuthor           = fmt.Errorf("Пользователь не является автором")
	BidCanceled           = fmt.Errorf("Предложение уже закрыто")
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
		log.Println("check IsResponsible error -", err)
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
	new_tender, err := s.query.EditTenderWithTX(ctx, database.EditTenderWithTxParam{
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
		log.Println("EditTenderStatus: EditTenderWithTX err -", err)
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

type RollbackTenderRequest struct {
	Username  string
	Tender_id string
	Version   int32
}

func (s *Service) RollbackTender(ctx context.Context, params RollbackTenderRequest) (*Tender, error) {
	user_id, err := s.query.FetchUserID(ctx, params.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("RollbackTender: FetchUserID err -", err)
		return nil, UnknowError
	}
	tender, err := s.query.GetTender(ctx, params.Tender_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("RollbackTender: GetTender err -", err)
		return nil, UnknowError
	}

	err = s.isResponsibleUser(ctx, tender.OrganizationID.String(), user_id)
	if err != nil {
		return nil, err
	}

	tender_history, err := s.query.GetTenderHistory(ctx, tender.ID.String(), params.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderHistoryNotFound
		}
		log.Println("RollbackTender: GetTenderHistory err -", err)
		return nil, UnknowError
	}

	new_tender, err := s.query.EditTenderWithTX(ctx, database.EditTenderWithTxParam{
		ChangeTenderParam: database.TenderChangeParam{
			Tender_id:    tender.ID.String(),
			Name:         tender_history.Name,
			Description:  tender_history.Description,
			Service_type: tender_history.ServiceType,
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
		log.Println("RollbackTender: EditTenderWithTX err -", err)
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

//offers

type Bid struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	AuthorType string    `json:"authorType"`
	AuthorId   string    `json:"authorId"`
	Version    int32     `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CreateBidParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderId    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorId    string `json:"authorId"`
}

func (s *Service) CreateNewBid(ctx context.Context, param CreateBidParam) (*Bid, error) {
	allowValue := []string{"User", "Organization"}
	if !utils.CheckString(param.AuthorType, allowValue) {
		return nil, NotAllowValue
	}

	user, err := s.query.FetchUser(ctx, param.AuthorId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("CreateNewBid: Fetching User id error -", err)
		return nil, UnknowError
	}

	tender, err := s.query.GetTender(ctx, param.TenderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("CreateNewBid: GetTender err -", err)
		return nil, UnknowError
	}

	err = s.isResponsibleUser(ctx, tender.OrganizationID.String(), param.AuthorId)
	if err == nil {
		return nil, IsResponsible
	}

	org_id, err := s.query.GetUserOrganization(ctx, user.ID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, IsNotResponsible
		}
		return nil, UnknowError
	}

	bid, err := s.query.CreateOffer(ctx, database.CreateOfferParam{
		Name:            param.Name,
		Description:     param.Description,
		TenderId:        param.TenderId,
		AuthorType:      param.AuthorType,
		AuthorId:        param.AuthorId,
		Organization_id: org_id,
	})
	if err != nil {
		log.Println("CreateNewBid: CreateOffer error -", err)
		return nil, UnknowError
	}
	return &Bid{
		ID:         bid.ID.String(),
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId.String(),
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil

}

type ListMyBidsRequest struct {
	Username string
	Offset   int32
	Limit    int32
}

func (s *Service) ListMyBids(ctx context.Context, param ListMyBidsRequest) ([]Bid, error) {
	user_id, err := s.query.FetchUserID(ctx, param.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("CreateNewBid: Fetching User id error -", err)
		return nil, UnknowError
	}

	listoffers, err := s.query.MyListOffers(ctx, &database.MyListOffersParams{
		Creator_id: user_id,
		Offset:     param.Offset,
		Limit:      param.Limit,
	})
	if err != nil {
		log.Println("ListMyBids: MyListOffers err -", err)
		return nil, UnknowError
	}
	var bidslist []Bid
	for _, item := range listoffers {
		bidslist = append(bidslist, Bid{
			ID:         item.ID.String(),
			Name:       item.Name,
			Status:     item.Status,
			AuthorType: item.AuthorType,
			AuthorId:   item.AuthorId.String(),
			Version:    item.Version,
			CreatedAt:  item.CreatedAt,
		})
	}
	return bidslist, nil
}

type TenderListBidsRequest struct {
	Tender_id string
	Username  string
	Offset    int32
	Limit     int32
}

func (s *Service) TenderListBids(ctx context.Context, param TenderListBidsRequest) ([]Bid, error) {
	user_id, err := s.query.FetchUserID(ctx, param.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("TenderListBids: Fetching User id error -", err)
		return nil, UnknowError
	}

	tender, err := s.query.GetTender(ctx, param.Tender_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("TenderListBids: GetTender err -", err)
		return nil, UnknowError
	}
	org_id, _ := s.query.GetUserOrganization(ctx, user_id)

	listoffers, err := s.query.TenderListOffers(ctx, &database.TenderListOffersParams{
		Tender_id:       tender.ID.String(),
		Offset:          param.Offset,
		Limit:           param.Limit,
		Organization_id: org_id,
	})
	if err != nil {
		log.Println("TenderListBids: MyListOffers err -", err)
		return nil, UnknowError
	}
	var bidslist []Bid
	for _, item := range listoffers {
		bidslist = append(bidslist, Bid{
			ID:         item.ID.String(),
			Name:       item.Name,
			Status:     item.Status,
			AuthorType: item.AuthorType,
			AuthorId:   item.AuthorId.String(),
			Version:    item.Version,
			CreatedAt:  item.CreatedAt,
		})
	}
	return bidslist, nil
}

type GetBidStatus struct {
	Username string
	BidID    string
}

func (s *Service) GetBidStatus(ctx context.Context, param GetBidStatus) (string, error) {
	user_id, err := s.query.FetchUserID(ctx, param.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", UserNotFound
		}
		log.Println("GetBidStatus: Fetching User id error -", err)
		return "", UnknowError
	}

	offer, err := s.query.GetOffer(ctx, param.BidID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", BidNotFound
		}
		log.Println("GetBidStatus: GetOffer error -", err)
		return "", UnknowError

	}

	err = s.isResponsibleUser(ctx, offer.Organization_ID.String(), user_id)
	if err != nil {
		return "", err
	}

	return offer.Status, nil
}

type ChangeBidStatus struct {
	Username string
	BidID    string
	Status   string
}

func (s *Service) ChangeBidStatus(ctx context.Context, param ChangeBidStatus) (*Bid, error) {

	user_id, err := s.query.FetchUserID(ctx, param.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("ChangeBidStatus: Fetching User id error -", err)
		return nil, UnknowError
	}

	offer, err := s.query.GetOffer(ctx, param.BidID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, BidNotFound
		}
		log.Println("ChangeBidStatus: GetOffer error -", err)
		return nil, UnknowError

	}

	err = s.isResponsibleUser(ctx, offer.Organization_ID.String(), user_id)
	if err != nil {
		return nil, err
	}

	if offer.Status == "Canceled" {
		return nil, BidCanceled
	}

	bid, err := s.query.ChangeOfferStatus(ctx, offer.ID.String(), param.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, BidNotFound
		}
		log.Println("ChangeBidStatus: ChangeTenderStatus err -", err)
		return nil, UnknowError
	}
	return &Bid{
		ID:         bid.ID.String(),
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId.String(),
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil
}

type EditBidRequest struct {
	Username    string
	Bid_id      string
	Name        string
	Description string
}

func (s *Service) EditBid(ctx context.Context, params EditBidRequest) (*Bid, error) {
	user_id, err := s.query.FetchUserID(ctx, params.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("EditTenderStatus: FetchUserID err -", err)
		return nil, UnknowError
	}
	bid, err := s.query.GetOffer(ctx, params.Bid_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, BidNotFound
		}
		log.Println("EditTenderStatus: GetTender err -", err)
		return nil, UnknowError
	}

	err = s.isResponsibleUser(ctx, bid.Organization_ID.String(), user_id)
	if err != nil {
		return nil, err
	}
	new_bid, err := s.query.EditOfferWithTX(ctx, database.EditOfferWithTxParam{
		ChangeOfferParam: database.OfferChangeParam{
			Bid_id:      bid.ID.String(),
			Name:        params.Name,
			Description: params.Description,
		},
		OfferHistoryParam: database.CreateOfferHistoryParams{
			Offer_id:   bid.ID.String(),
			Creator_id: user_id,

			Name:        bid.Name,
			Description: bid.Description,
			OldVersion:  bid.Version,
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("EditTenderStatus: EditTenderWithTX err -", err)
		return nil, UnknowError
	}
	return &Bid{
		ID:         new_bid.ID.String(),
		Name:       new_bid.Name,
		Status:     new_bid.Status,
		AuthorType: new_bid.AuthorType,
		AuthorId:   new_bid.AuthorId.String(),
		Version:    new_bid.Version,
		CreatedAt:  new_bid.CreatedAt,
	}, nil
}

type RollbackOfferRequest struct {
	Username string
	Offer_id string
	Version  int32
}

func (s *Service) RollbackOffer(ctx context.Context, params RollbackOfferRequest) (*Bid, error) {
	user_id, err := s.query.FetchUserID(ctx, params.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		log.Println("RollbackTender: FetchUserID err -", err)
		return nil, UnknowError
	}
	bid, err := s.query.GetOffer(ctx, params.Offer_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("RollbackTender: GetTender err -", err)
		return nil, UnknowError
	}

	err = s.isResponsibleUser(ctx, bid.Organization_ID.String(), user_id)
	if err != nil {
		return nil, err
	}

	offer_history, err := s.query.GetOfferHistory(ctx, bid.ID.String(), params.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderHistoryNotFound
		}
		log.Println("RollbackTender: GetTenderHistory err -", err)
		return nil, UnknowError
	}

	new_tender, err := s.query.EditOfferWithTX(ctx, database.EditOfferWithTxParam{
		ChangeOfferParam: database.OfferChangeParam{
			Bid_id:      bid.ID.String(),
			Name:        offer_history.Name,
			Description: offer_history.Description,
		},
		OfferHistoryParam: database.CreateOfferHistoryParams{
			Offer_id:    bid.ID.String(),
			Creator_id:  user_id,
			Name:        bid.Name,
			Description: bid.Description,
			OldVersion:  bid.Version,
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TenderNotFound
		}
		log.Println("RollbackTender: EditTenderWithTX err -", err)
		return nil, UnknowError
	}
	return &Bid{
		ID:         new_tender.ID.String(),
		Name:       new_tender.Name,
		Status:     new_tender.Status,
		AuthorType: new_tender.AuthorType,
		AuthorId:   new_tender.AuthorId.String(),
		Version:    new_tender.Version,
		CreatedAt:  new_tender.CreatedAt,
	}, nil
}
