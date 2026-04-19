package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"sea-cucumber-trace/backend/internal/config"
	"sea-cucumber-trace/backend/internal/fabric"
	"sea-cucumber-trace/backend/internal/model"
	"sea-cucumber-trace/backend/internal/repository"
)

type Service struct {
	cfg    *config.Config
	repo   *repository.Repository
	ledger fabric.Ledger
}

// ErrBatchNoRequired is returned when batch number is empty after trim.
var ErrBatchNoRequired = errors.New("batch number required")

func New(cfg *config.Config, repo *repository.Repository) *Service {
	gw := fabric.NewGateway(cfg)
	led := fabric.ResolveLedger(cfg.FabricEnabled, gw)
	return &Service{cfg: cfg, repo: repo, ledger: led}
}

type Claims struct {
	UserID uint           `json:"uid"`
	Role   model.UserRole `json:"role"`
	OrgID  *uint          `json:"oid,omitempty"`
	jwt.RegisteredClaims
}



func (s *Service) GetUser(id uint) (*model.User, error) {
	return s.repo.UserByID(id)
}
func (s *Service) Login(username, password string) (string, *model.User, error) {
	u, err := s.repo.UserByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", nil, errors.New("invalid credentials")
	}
	tok, err := s.issueToken(u)
	return tok, u, err
}

func (s *Service) issueToken(u *model.User) (string, error) {
	claims := Claims{
		UserID: u.ID,
		Role:   u.Role,
		OrgID:  u.OrgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.cfg.JWTSecret))
}

func ParseToken(cfg *config.Config, token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

type CreateBatchInput struct {
	BatchNo        string
	ProductName    string
	FarmBase       string
	Quality        string
	CatchDate      *time.Time
	BreedArea      string
	BreedStartDate *time.Time
	Spec           string
	Quantity       string
	ExtraJSON      string
	OrgID          uint
}

func (s *Service) CreateBatch(actor *Claims, in CreateBatchInput) (*model.SeaCucumberBatch, error) {
	if actor.Role != model.RoleAdmin && (actor.OrgID == nil || *actor.OrgID != in.OrgID) {
		return nil, errors.New("forbidden")
	}
	in.BatchNo = strings.TrimSpace(in.BatchNo)
	if in.BatchNo == "" {
		return nil, ErrBatchNoRequired
	}
	b := &model.SeaCucumberBatch{
		BatchNo:        in.BatchNo,
		OrgID:          in.OrgID,
		ProductName:    strings.TrimSpace(in.ProductName),
		FarmBase:       strings.TrimSpace(in.FarmBase),
		Quality:        strings.TrimSpace(in.Quality),
		CatchDate:      in.CatchDate,
		BreedArea:      strings.TrimSpace(in.BreedArea),
		BreedStartDate: in.BreedStartDate,
		Spec:           strings.TrimSpace(in.Spec),
		Quantity:       strings.TrimSpace(in.Quantity),
		ExtraJSON:      strings.TrimSpace(in.ExtraJSON),
	}
	if err := s.repo.CreateBatch(b); err != nil {
		return nil, err
	}
	return s.repo.GetBatchByID(b.ID)
}

func (s *Service) ListBatches(actor *Claims) ([]model.SeaCucumberBatch, error) {
	if actor.Role == model.RoleAdmin {
		return s.repo.ListBatches(nil)
	}
	if actor.OrgID == nil {
		return []model.SeaCucumberBatch{}, nil
	}
	return s.repo.ListBatches(actor.OrgID)
}

func (s *Service) GetBatch(actor *Claims, id uint) (*model.SeaCucumberBatch, error) {
	b, err := s.repo.GetBatchByID(id)
	if err != nil {
		return nil, err
	}
	if actor.Role != model.RoleAdmin && (actor.OrgID == nil || b.OrgID != *actor.OrgID) {
		return nil, errors.New("forbidden")
	}
	return b, nil
}

type AddEventInput struct {
	Stage        model.TraceStage
	Title        string
	DetailJSON   string
	Location     string
	OperatorName string
	EvidenceURLs string
	OccurredAt   time.Time
}

func (s *Service) AddEvent(actor *Claims, batchID uint, in AddEventInput) (*model.TraceEvent, *model.ChainRecord, error) {
	b, err := s.repo.GetBatchByID(batchID)
	if err != nil {
		return nil, nil, err
	}
	if actor.Role != model.RoleAdmin && (actor.OrgID == nil || b.OrgID != *actor.OrgID) {
		return nil, nil, errors.New("forbidden")
	}

	hash := fabric.HashTraceEvent(b.BatchNo, string(in.Stage), in.DetailJSON, in.Location, in.OperatorName, in.OccurredAt, in.EvidenceURLs)
	ev := &model.TraceEvent{
		BatchID:      batchID,
		Stage:        in.Stage,
		Title:        in.Title,
		DetailJSON:   in.DetailJSON,
		Location:     in.Location,
		OperatorName: in.OperatorName,
		EvidenceURLs: in.EvidenceURLs,
		OccurredAt:   in.OccurredAt,
		DataHash:     hash,
		CreatedBy:    actor.UserID,
	}
	if err := s.repo.CreateEvent(ev); err != nil {
		return nil, nil, err
	}

	req := fabric.RecordRequest{
		BatchNo:    b.BatchNo,
		EventID:    ev.ID,
		Stage:      string(ev.Stage),
		DataHash:   hash,
		OccurredAt: ev.OccurredAt.UTC().Format(time.RFC3339Nano),
		OrgName:    b.Org.Name,
	}
	cr := &model.ChainRecord{
		RefType:   "event",
		RefID:     ev.ID,
		ChainType: "fabric",
		Status:    model.ChainPending,
	}
	payload, _ := json.Marshal(req)
	cr.PayloadJSON = string(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := s.ledger.RecordTrace(ctx, req)
	if err != nil {
		cr.Status = model.ChainFailed
		_ = s.repo.CreateChainRecord(cr)
		return ev, cr, nil
	}
	cr.Status = model.ChainSuccess
	cr.TxID = res.TxID
	cr.BlockNumber = res.BlockNumber
	_ = s.repo.CreateChainRecord(cr)
	return ev, cr, nil
}

type TimelineItem struct {
	Event model.TraceEvent  `json:"event"`
	Chain *model.ChainRecord `json:"chain,omitempty"`
}

func (s *Service) TimelineByBatchID(actor *Claims, batchID uint) ([]TimelineItem, error) {
	b, err := s.repo.GetBatchByID(batchID)
	if err != nil {
		return nil, err
	}
	if actor != nil {
		if actor.Role != model.RoleAdmin && (actor.OrgID == nil || b.OrgID != *actor.OrgID) {
			return nil, errors.New("forbidden")
		}
	}
	evs, err := s.repo.ListEvents(batchID)
	if err != nil {
		return nil, err
	}
	out := make([]TimelineItem, 0, len(evs))
	for _, e := range evs {
		item := TimelineItem{Event: e}
		if c, err := s.repo.ChainByEventID(e.ID); err == nil {
			item.Chain = c
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *Service) PublicTimeline(batchNo string) ([]TimelineItem, *model.SeaCucumberBatch, error) {
	b, err := s.repo.GetBatchByNo(batchNo)
	if err != nil {
		return nil, nil, err
	}
	items, err := s.TimelineByBatchID(nil, b.ID)
	return items, b, err
}

type VerifyResult struct {
	BatchNo   string `json:"batchNo"`
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	EventID   uint   `json:"eventId,omitempty"`
	Expected  string `json:"expectedHash,omitempty"`
	Actual    string `json:"actualHash,omitempty"`
	TxID      string `json:"txId,omitempty"`
}

func (s *Service) Verify(batchNo string) ([]VerifyResult, error) {
	b, err := s.repo.GetBatchByNo(batchNo)
	if err != nil {
		return nil, err
	}
	evs, err := s.repo.ListEvents(b.ID)
	if err != nil {
		return nil, err
	}
	res := make([]VerifyResult, 0, len(evs))
	for _, e := range evs {
		want := fabric.HashTraceEvent(b.BatchNo, string(e.Stage), e.DetailJSON, e.Location, e.OperatorName, e.OccurredAt, e.EvidenceURLs)
		v := VerifyResult{
			BatchNo: b.BatchNo,
			EventID: e.ID,
			Expected: want,
			Actual:   e.DataHash,
			OK:       want == e.DataHash,
		}
		if v.OK {
			v.Message = "链下数据与哈希一致"
		} else {
			v.Message = "数据可能被篡改或字段不一致"
		}
		if c, err := s.repo.ChainByEventID(e.ID); err == nil {
			v.TxID = c.TxID
		}
		res = append(res, v)
	}
	return res, nil
}
