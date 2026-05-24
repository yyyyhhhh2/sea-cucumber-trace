package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	appcache "sea-cucumber-trace/backend/internal/cache"
	"sea-cucumber-trace/backend/internal/config"
	"sea-cucumber-trace/backend/internal/fabric"
	"sea-cucumber-trace/backend/internal/model"
	"sea-cucumber-trace/backend/internal/repository"
)

type Service struct {
	cfg    *config.Config
	repo   *repository.Repository
	ledger fabric.Ledger
	cache  *appcache.Client
}

var ErrBatchNoRequired = errors.New("batch number required")
var ErrStageRequired = errors.New("stage required")
var ErrTitleRequired = errors.New("title required")

func New(cfg *config.Config, repo *repository.Repository, cacheClient *appcache.Client) *Service {
	gw := fabric.NewGateway(cfg)
	led := fabric.ResolveLedger(cfg.FabricEnabled, gw)
	return &Service{cfg: cfg, repo: repo, ledger: led, cache: cacheClient}
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

type HealthStatus struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	DB        string `json:"db"`
	Cache     string `json:"cache"`
	Ledger    string `json:"ledger"`
	Timestamp string `json:"timestamp"`
}

func (s *Service) Health(ctx context.Context) *HealthStatus {
	status := &HealthStatus{
		Status:    "ok",
		Service:   "sea-cucumber-trace",
		DB:        "ok",
		Cache:     "disabled",
		Ledger:    "mock",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.repo.Ping(); err != nil {
		status.Status = "degraded"
		status.DB = "error"
	}

	if s.cache != nil {
		status.Cache = "ok"
		if err := s.cache.Ping(ctx); err != nil {
			status.Status = "degraded"
			status.Cache = "error"
		}
	}

	if s.cfg.FabricEnabled {
		status.Ledger = "fabric"
	}

	return status
}

func (s *Service) ListOrgs(actor *Claims) ([]model.Org, error) {
	if actor == nil {
		return nil, errors.New("unauthorized")
	}
	list, err := s.repo.ListOrgs()
	if err != nil {
		return nil, err
	}
	if actor.Role == model.RoleAdmin {
		return list, nil
	}
	if actor.OrgID == nil {
		return []model.Org{}, nil
	}
	filtered := make([]model.Org, 0, 1)
	for _, item := range list {
		if item.ID == *actor.OrgID {
			filtered = append(filtered, item)
			break
		}
	}
	return filtered, nil
}

type DashboardSummary struct {
	OrgCount           int64                    `json:"orgCount"`
	BatchCount         int64                    `json:"batchCount"`
	EventCount         int64                    `json:"eventCount"`
	AnchoredEventCount int64                    `json:"anchoredEventCount"`
	PendingChainCount  int64                    `json:"pendingChainCount"`
	RecentBatches      []model.SeaCucumberBatch `json:"recentBatches"`
}

func (s *Service) DashboardSummary(actor *Claims) (*DashboardSummary, error) {
	if actor == nil {
		return nil, errors.New("unauthorized")
	}

	var scopeOrgID *uint
	if actor.Role != model.RoleAdmin {
		scopeOrgID = actor.OrgID
	}

	raw, err := s.repo.DashboardSummary(scopeOrgID)
	if err != nil {
		return nil, err
	}
	batches, err := s.repo.ListBatches(scopeOrgID)
	if err != nil {
		return nil, err
	}
	if len(batches) > 5 {
		batches = batches[:5]
	}

	return &DashboardSummary{
		OrgCount:           raw.OrgCount,
		BatchCount:         raw.BatchCount,
		EventCount:         raw.EventCount,
		AnchoredEventCount: raw.AnchoredEventCount,
		PendingChainCount:  raw.PendingChainCount,
		RecentBatches:      batches,
	}, nil
}

func (s *Service) ImportDemoData(actor *Claims) error {
	if actor == nil || actor.Role != model.RoleAdmin {
		return errors.New("forbidden")
	}
	if err := s.repo.ImportDemoData(); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.FlushPrefix(context.Background())
	}
	return nil
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
	if actor == nil {
		return nil, errors.New("unauthorized")
	}
	if actor.Role == model.RoleAdmin {
		if in.OrgID == 0 {
			return nil, errors.New("org id required")
		}
	} else {
		if actor.OrgID == nil {
			return nil, errors.New("forbidden")
		}
		if in.OrgID == 0 {
			in.OrgID = *actor.OrgID
		}
		if *actor.OrgID != in.OrgID {
			return nil, errors.New("forbidden")
		}
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
	s.invalidateBatchCaches(context.Background(), in.BatchNo, in.OrgID)
	return s.repo.GetBatchByID(b.ID)
}

func (s *Service) ListBatches(actor *Claims) ([]model.SeaCucumberBatch, error) {
	key := "batches:admin"
	if actor.Role == model.RoleAdmin {
		var cached []model.SeaCucumberBatch
		if ok, err := s.cacheGet(context.Background(), key, &cached); ok && err == nil {
			return cached, nil
		}
		list, err := s.repo.ListBatches(nil)
		if err == nil {
			_ = s.cacheSet(context.Background(), key, list)
		}
		return list, err
	}
	if actor.OrgID == nil {
		return []model.SeaCucumberBatch{}, nil
	}
	key = fmt.Sprintf("batches:org:%d", *actor.OrgID)
	var cached []model.SeaCucumberBatch
	if ok, err := s.cacheGet(context.Background(), key, &cached); ok && err == nil {
		return cached, nil
	}
	list, err := s.repo.ListBatches(actor.OrgID)
	if err == nil {
		_ = s.cacheSet(context.Background(), key, list)
	}
	return list, err
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
	if strings.TrimSpace(string(in.Stage)) == "" {
		return nil, nil, ErrStageRequired
	}
	if strings.TrimSpace(in.Title) == "" {
		return nil, nil, ErrTitleRequired
	}

	in.Title = strings.TrimSpace(in.Title)
	in.DetailJSON = strings.TrimSpace(in.DetailJSON)
	in.Location = strings.TrimSpace(in.Location)
	in.OperatorName = strings.TrimSpace(in.OperatorName)
	in.EvidenceURLs = strings.TrimSpace(in.EvidenceURLs)

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
	defer s.invalidateBatchCaches(context.Background(), b.BatchNo, b.OrgID)

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
	Event model.TraceEvent   `json:"event"`
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
	key := "trace:" + batchNo
	var cached struct {
		Items []TimelineItem          `json:"items"`
		Batch *model.SeaCucumberBatch `json:"batch"`
	}
	if ok, err := s.cacheGet(context.Background(), key, &cached); ok && err == nil && cached.Batch != nil {
		return cached.Items, cached.Batch, nil
	}
	b, err := s.repo.GetBatchByNo(batchNo)
	if err != nil {
		return nil, nil, err
	}
	items, err := s.TimelineByBatchID(nil, b.ID)
	if err == nil {
		_ = s.cacheSet(context.Background(), key, struct {
			Items []TimelineItem          `json:"items"`
			Batch *model.SeaCucumberBatch `json:"batch"`
		}{Items: items, Batch: b})
	}
	return items, b, err
}

type VerifyResult struct {
	BatchNo  string `json:"batchNo"`
	OK       bool   `json:"ok"`
	Message  string `json:"message"`
	EventID  uint   `json:"eventId,omitempty"`
	Expected string `json:"expectedHash,omitempty"`
	Actual   string `json:"actualHash,omitempty"`
	TxID     string `json:"txId,omitempty"`
}

func (s *Service) Verify(batchNo string) ([]VerifyResult, error) {
	key := "verify:" + batchNo
	var cached []VerifyResult
	if ok, err := s.cacheGet(context.Background(), key, &cached); ok && err == nil {
		return cached, nil
	}
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
			BatchNo:  b.BatchNo,
			EventID:  e.ID,
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
	_ = s.cacheSet(context.Background(), key, res)
	return res, nil
}

func (s *Service) cacheGet(ctx context.Context, key string, dest any) (bool, error) {
	if s.cache == nil {
		return false, nil
	}
	return s.cache.GetJSON(ctx, key, dest)
}

func (s *Service) cacheSet(ctx context.Context, key string, value any) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.SetJSON(ctx, key, value)
}

func (s *Service) invalidateBatchCaches(ctx context.Context, batchNo string, orgID uint) {
	if s.cache == nil {
		return
	}
	keys := []string{"batches:admin", "trace:" + batchNo, "verify:" + batchNo}
	if orgID > 0 {
		keys = append(keys, fmt.Sprintf("batches:org:%d", orgID))
	}
	if batchNo == "" {
		keys = []string{"batches:admin"}
		if orgID > 0 {
			keys = append(keys, fmt.Sprintf("batches:org:%d", orgID))
		}
	}
	_ = s.cache.Delete(ctx, keys...)
}
