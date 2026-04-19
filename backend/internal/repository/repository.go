package repository

import (
	"errors"

	"gorm.io/gorm"

	"sea-cucumber-trace/backend/internal/model"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UserByUsername(name string) (*model.User, error) {
	var u model.User
	err := r.db.Preload("Org").Where("username = ?", name).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) ListBatches(orgID *uint) ([]model.SeaCucumberBatch, error) {
	q := r.db.Preload("Org").Order("id desc")
	if orgID != nil {
		q = q.Where("org_id = ?", *orgID)
	}
	var list []model.SeaCucumberBatch
	return list, q.Find(&list).Error
}

func (r *Repository) CreateBatch(b *model.SeaCucumberBatch) error {
	return r.db.Create(b).Error
}

func (r *Repository) GetBatchByID(id uint) (*model.SeaCucumberBatch, error) {
	var b model.SeaCucumberBatch
	err := r.db.Preload("Org").First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) GetBatchByNo(no string) (*model.SeaCucumberBatch, error) {
	var b model.SeaCucumberBatch
	err := r.db.Preload("Org").Where("batch_no = ?", no).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) ListEvents(batchID uint) ([]model.TraceEvent, error) {
	var evs []model.TraceEvent
	err := r.db.Where("batch_id = ?", batchID).Order("occurred_at asc").Find(&evs).Error
	return evs, err
}

func (r *Repository) CreateEvent(ev *model.TraceEvent) error {
	return r.db.Create(ev).Error
}

func (r *Repository) CreateChainRecord(cr *model.ChainRecord) error {
	return r.db.Create(cr).Error
}

func (r *Repository) ChainByEventID(eventID uint) (*model.ChainRecord, error) {
	var c model.ChainRecord
	err := r.db.Where("ref_type = ? AND ref_id = ?", "event", eventID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *Repository) UserByID(id uint) (*model.User, error) {
	var u model.User
	err := r.db.Preload("Org").First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

