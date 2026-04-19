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
	Username     string         `gorm:"uniqueIndex;size:64" json:"username"`
	PasswordHash string         `gorm:"size:255" json:"-"`
	DisplayName  string         `gorm:"size:128" json:"displayName"`
	Role         UserRole       `gorm:"size:32" json:"role"`
	OrgID        *uint          `json:"orgId,omitempty"`
	Org          *Org           `json:"org,omitempty"`
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
	Name        string         `gorm:"size:255" json:"name"`
	Type        OrgType        `gorm:"size:32" json:"type"`
	LicenseNo   string         `gorm:"size:128" json:"licenseNo"`
	Address     string         `gorm:"size:512" json:"address"`
	Contact     string         `gorm:"size:128" json:"contact"`
	Description string         `gorm:"size:1024" json:"description"`
}

type SeaCucumberBatch struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	BatchNo        string         `gorm:"uniqueIndex;size:64" json:"batchNo"`
	OrgID          uint           `gorm:"index" json:"orgId"`
	Org            Org            `json:"org,omitempty"`
	ProductName    string         `gorm:"size:255" json:"productName,omitempty"`
	FarmBase       string         `gorm:"size:255" json:"farmBase,omitempty"`
	Quality        string         `gorm:"size:255" json:"quality,omitempty"`
	CatchDate      *time.Time     `json:"catchDate,omitempty"`
	BreedArea      string         `gorm:"size:255" json:"breedArea"`
	BreedStartDate *time.Time     `json:"breedStartDate,omitempty"`
	Spec           string         `gorm:"size:128" json:"spec"`
	Quantity       string         `gorm:"size:64" json:"quantity"`
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
	BatchID      uint             `gorm:"index" json:"batchId"`
	Batch        SeaCucumberBatch `json:"batch,omitempty"`
	Stage        TraceStage       `gorm:"size:32" json:"stage"`
	Title        string           `gorm:"size:255" json:"title"`
	DetailJSON   string           `gorm:"type:text" json:"detailJson"`
	Location     string           `gorm:"size:255" json:"location"`
	OperatorName string           `gorm:"size:128" json:"operatorName"`
	EvidenceURLs string           `gorm:"type:text" json:"evidenceUrls"`
	OccurredAt   time.Time        `json:"occurredAt"`
	DataHash     string           `gorm:"size:128;index" json:"dataHash"`
	CreatedBy    uint             `json:"createdBy"`
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
	RefType     string            `gorm:"size:32;index:idx_chain_ref,priority:1" json:"refType"`
	RefID       uint              `gorm:"index:idx_chain_ref,priority:2" json:"refId"`
	ChainType   string            `gorm:"size:32" json:"chainType"`
	TxID        string            `gorm:"size:256;index" json:"txId"`
	BlockNumber uint64            `json:"blockNumber,omitempty"`
	Status      ChainRecordStatus `gorm:"size:32" json:"status"`
	PayloadJSON string            `gorm:"type:text" json:"payloadJson,omitempty"`
}
