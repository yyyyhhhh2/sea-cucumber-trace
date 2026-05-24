package model

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleOrg   UserRole = "org"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Username     string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	DisplayName  string         `gorm:"size:128;not null" json:"displayName"`
	Role         UserRole       `gorm:"size:32;not null;index" json:"role"`
	OrgID        *uint          `gorm:"index" json:"orgId,omitempty"`
	Org          *Org           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"org,omitempty"`
}

type OrgType string

const (
	OrgFarm      OrgType = "breeding"
	OrgProcess   OrgType = "processing"
	OrgLogistics OrgType = "logistics"
	OrgRetail    OrgType = "retail"
)

type Org struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Type        OrgType        `gorm:"size:32;not null;index" json:"type"`
	LicenseNo   string         `gorm:"size:128;not null;uniqueIndex" json:"licenseNo"`
	Address     string         `gorm:"size:512;not null" json:"address"`
	Contact     string         `gorm:"size:128;not null" json:"contact"`
	Description string         `gorm:"size:1024;not null" json:"description"`
}

type SeaCucumberBatch struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	BatchNo        string         `gorm:"size:64;not null;uniqueIndex" json:"batchNo"`
	OrgID          uint           `gorm:"not null;index:idx_batch_org_created,priority:1" json:"orgId"`
	Org            Org            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"org,omitempty"`
	ProductName    string         `gorm:"size:255;not null" json:"productName,omitempty"`
	FarmBase       string         `gorm:"size:255;not null" json:"farmBase,omitempty"`
	Quality        string         `gorm:"size:255;not null" json:"quality,omitempty"`
	CatchDate      *time.Time     `gorm:"index" json:"catchDate,omitempty"`
	BreedArea      string         `gorm:"size:255;not null" json:"breedArea"`
	BreedStartDate *time.Time     `gorm:"index" json:"breedStartDate,omitempty"`
	Spec           string         `gorm:"size:128;not null" json:"spec"`
	Quantity       string         `gorm:"size:64;not null" json:"quantity"`
	ExtraJSON      string         `gorm:"type:text" json:"extraJson,omitempty"`
}

type TraceStage string

const (
	StageBreeding   TraceStage = "breeding"
	StageHarvest    TraceStage = "harvest"
	StageProcessing TraceStage = "processing"
	StagePackaging  TraceStage = "packaging"
	StageLogistics  TraceStage = "logistics"
	StageRetail     TraceStage = "retail"
)

type TraceEvent struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
	BatchID      uint             `gorm:"not null;index:idx_event_batch_time,priority:1;index:idx_event_batch_stage,priority:1" json:"batchId"`
	Batch        SeaCucumberBatch `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Stage        TraceStage       `gorm:"size:32;not null;index:idx_event_batch_stage,priority:2" json:"stage"`
	Title        string           `gorm:"size:255;not null" json:"title"`
	DetailJSON   string           `gorm:"type:text" json:"detailJson"`
	Location     string           `gorm:"size:255;not null" json:"location"`
	OperatorName string           `gorm:"size:128;not null" json:"operatorName"`
	EvidenceURLs string           `gorm:"type:text" json:"evidenceUrls"`
	OccurredAt   time.Time        `gorm:"not null;index:idx_event_batch_time,priority:2" json:"occurredAt"`
	DataHash     string           `gorm:"size:128;not null;index" json:"dataHash"`
	CreatedBy    uint             `gorm:"not null;index" json:"createdBy"`
}

type ChainRecordStatus string

const (
	ChainPending ChainRecordStatus = "pending"
	ChainSuccess ChainRecordStatus = "success"
	ChainFailed  ChainRecordStatus = "failed"
)

type ChainRecord struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
	RefType     string            `gorm:"size:32;not null;uniqueIndex:idx_chain_ref,priority:1" json:"refType"`
	RefID       uint              `gorm:"not null;uniqueIndex:idx_chain_ref,priority:2" json:"refId"`
	ChainType   string            `gorm:"size:32;not null" json:"chainType"`
	TxID        string            `gorm:"size:256;index" json:"txId"`
	BlockNumber uint64            `json:"blockNumber,omitempty"`
	Status      ChainRecordStatus `gorm:"size:32;not null;index" json:"status"`
	PayloadJSON string            `gorm:"type:text" json:"payloadJson,omitempty"`
}
