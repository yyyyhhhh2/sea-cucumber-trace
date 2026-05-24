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

func (r *Repository) Ping() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func PrepareForMigration(db *gorm.DB) error {
	updates := map[string][]string{
		"users": {
			"UPDATE users SET username = '' WHERE username IS NULL",
			"UPDATE users SET password_hash = '' WHERE password_hash IS NULL",
			"UPDATE users SET display_name = '' WHERE display_name IS NULL",
			"UPDATE users SET role = '' WHERE role IS NULL",
		},
		"orgs": {
			"UPDATE orgs SET name = '' WHERE name IS NULL",
			"UPDATE orgs SET type = '' WHERE type IS NULL",
			"UPDATE orgs SET license_no = '' WHERE license_no IS NULL",
			"UPDATE orgs SET address = '' WHERE address IS NULL",
			"UPDATE orgs SET contact = '' WHERE contact IS NULL",
			"UPDATE orgs SET description = '' WHERE description IS NULL",
		},
		"sea_cucumber_batches": {
			"UPDATE sea_cucumber_batches SET batch_no = '' WHERE batch_no IS NULL",
			"UPDATE sea_cucumber_batches SET product_name = '' WHERE product_name IS NULL",
			"UPDATE sea_cucumber_batches SET farm_base = '' WHERE farm_base IS NULL",
			"UPDATE sea_cucumber_batches SET quality = '' WHERE quality IS NULL",
			"UPDATE sea_cucumber_batches SET breed_area = '' WHERE breed_area IS NULL",
			"UPDATE sea_cucumber_batches SET spec = '' WHERE spec IS NULL",
			"UPDATE sea_cucumber_batches SET quantity = '' WHERE quantity IS NULL",
			"UPDATE sea_cucumber_batches SET extra_json = '' WHERE extra_json IS NULL",
		},
		"trace_events": {
			"UPDATE trace_events SET stage = '' WHERE stage IS NULL",
			"UPDATE trace_events SET title = '' WHERE title IS NULL",
			"UPDATE trace_events SET detail_json = '' WHERE detail_json IS NULL",
			"UPDATE trace_events SET location = '' WHERE location IS NULL",
			"UPDATE trace_events SET operator_name = '' WHERE operator_name IS NULL",
			"UPDATE trace_events SET evidence_urls = '' WHERE evidence_urls IS NULL",
			"UPDATE trace_events SET data_hash = '' WHERE data_hash IS NULL",
		},
		"chain_records": {
			"UPDATE chain_records SET ref_type = '' WHERE ref_type IS NULL",
			"UPDATE chain_records SET chain_type = '' WHERE chain_type IS NULL",
			"UPDATE chain_records SET tx_id = '' WHERE tx_id IS NULL",
			"UPDATE chain_records SET status = '' WHERE status IS NULL",
			"UPDATE chain_records SET payload_json = '' WHERE payload_json IS NULL",
		},
	}
	for table, statements := range updates {
		if !db.Migrator().HasTable(table) {
			continue
		}
		for _, stmt := range statements {
			if err := db.Exec(stmt).Error; err != nil {
				return err
			}
		}
	}
	return nil
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

func (r *Repository) ListOrgs() ([]model.Org, error) {
	var list []model.Org
	return list, r.db.Order("id asc").Find(&list).Error
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

type DashboardSummary struct {
	OrgCount           int64
	BatchCount         int64
	EventCount         int64
	AnchoredEventCount int64
	PendingChainCount  int64
}

func (r *Repository) DashboardSummary(orgID *uint) (*DashboardSummary, error) {
	out := &DashboardSummary{}

	orgQuery := r.db.Model(&model.Org{})
	if orgID != nil {
		orgQuery = orgQuery.Where("id = ?", *orgID)
	}
	if err := orgQuery.Count(&out.OrgCount).Error; err != nil {
		return nil, err
	}

	batchQuery := r.db.Model(&model.SeaCucumberBatch{})
	if orgID != nil {
		batchQuery = batchQuery.Where("org_id = ?", *orgID)
	}
	if err := batchQuery.Count(&out.BatchCount).Error; err != nil {
		return nil, err
	}

	eventQuery := r.db.Model(&model.TraceEvent{})
	if orgID != nil {
		eventQuery = eventQuery.Joins("JOIN sea_cucumber_batches ON sea_cucumber_batches.id = trace_events.batch_id").
			Where("sea_cucumber_batches.org_id = ?", *orgID)
	}
	if err := eventQuery.Count(&out.EventCount).Error; err != nil {
		return nil, err
	}

	anchoredQuery := r.db.Model(&model.ChainRecord{}).Where("ref_type = ? AND status = ?", "event", model.ChainSuccess)
	if orgID != nil {
		anchoredQuery = anchoredQuery.Joins("JOIN trace_events ON trace_events.id = chain_records.ref_id").
			Joins("JOIN sea_cucumber_batches ON sea_cucumber_batches.id = trace_events.batch_id").
			Where("sea_cucumber_batches.org_id = ?", *orgID)
	}
	if err := anchoredQuery.Count(&out.AnchoredEventCount).Error; err != nil {
		return nil, err
	}

	pendingQuery := r.db.Model(&model.ChainRecord{}).Where("ref_type = ? AND status = ?", "event", model.ChainPending)
	if orgID != nil {
		pendingQuery = pendingQuery.Joins("JOIN trace_events ON trace_events.id = chain_records.ref_id").
			Joins("JOIN sea_cucumber_batches ON sea_cucumber_batches.id = trace_events.batch_id").
			Where("sea_cucumber_batches.org_id = ?", *orgID)
	}
	if err := pendingQuery.Count(&out.PendingChainCount).Error; err != nil {
		return nil, err
	}

	return out, nil
}
