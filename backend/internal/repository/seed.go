package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"sea-cucumber-trace/backend/internal/fabric"
	"sea-cucumber-trace/backend/internal/model"
)

type demoOrgSeed struct {
	Name        string
	Type        model.OrgType
	LicenseNo   string
	Address     string
	Contact     string
	Description string
}

type demoUserSeed struct {
	Username    string
	Password    string
	DisplayName string
	Role        model.UserRole
	OrgLicense  string
}

type demoBatchSeed struct {
	BatchNo        string
	OrgLicense     string
	ProductName    string
	FarmBase       string
	Quality        string
	CatchDate      *time.Time
	BreedArea      string
	BreedStartDate *time.Time
	Spec           string
	Quantity       string
	Extra          map[string]any
}

type demoEventSeed struct {
	BatchNo       string
	Stage         model.TraceStage
	Title         string
	Location      string
	Operator      string
	EvidenceURLs  string
	OccurredAt    time.Time
	CreatedByUser string
	Detail        map[string]any
	TxID          string
	BlockNumber   uint64
}

func (r *Repository) SeedIfEmpty() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.SeaCucumberBatch{}).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
		return seedDemoData(tx)
	})
}

func (r *Repository) ImportDemoData() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return seedDemoData(tx)
	})
}

func seedDemoData(tx *gorm.DB) error {
	orgs, err := seedOrgs(tx)
	if err != nil {
		return err
	}

	users, err := seedUsers(tx, orgs)
	if err != nil {
		return err
	}

	batches, err := seedBatches(tx, orgs)
	if err != nil {
		return err
	}

	return seedEvents(tx, batches, users)
}

func seedOrgs(tx *gorm.DB) (map[string]model.Org, error) {
	seeds := []demoOrgSeed{
		{
			Name:        "大连旅顺海洋牧场合作社",
			Type:        model.OrgFarm,
			LicenseNo:   "LN-DL-2025-LS-018",
			Address:     "辽宁省大连市旅顺口区双岛湾街道 18 号",
			Contact:     "0411-8800-1001",
			Description: "负责海参底播养殖、采捕及原产地质量管理。",
		},
		{
			Name:        "大连深蓝海产品加工有限公司",
			Type:        model.OrgProcess,
			LicenseNo:   "LN-DL-2025-PR-102",
			Address:     "辽宁省大连市甘井子区海泽路 66 号",
			Contact:     "0411-8800-2002",
			Description: "负责原料验收、分级、低温清洗和加工包装。",
		},
		{
			Name:        "渤海冷链物流有限公司",
			Type:        model.OrgLogistics,
			LicenseNo:   "LN-DL-2025-LG-205",
			Address:     "辽宁省大连市金普新区港兴街 9 号",
			Contact:     "0411-8800-3003",
			Description: "负责冷链仓储、干线运输和在途温控记录。",
		},
		{
			Name:        "鲜参优选零售中心",
			Type:        model.OrgRetail,
			LicenseNo:   "LN-DL-2025-RT-309",
			Address:     "辽宁省大连市中山区人民路 128 号",
			Contact:     "0411-8800-4004",
			Description: "负责终端陈列、销售和消费者查询展示。",
		},
	}

	out := make(map[string]model.Org, len(seeds))
	for _, seed := range seeds {
		org := model.Org{
			Name:        seed.Name,
			Type:        seed.Type,
			LicenseNo:   seed.LicenseNo,
			Address:     seed.Address,
			Contact:     seed.Contact,
			Description: seed.Description,
		}
		if err := tx.Where("license_no = ?", seed.LicenseNo).Assign(org).FirstOrCreate(&org).Error; err != nil {
			return nil, err
		}
		out[seed.LicenseNo] = org
	}
	return out, nil
}

func seedUsers(tx *gorm.DB, orgs map[string]model.Org) (map[string]model.User, error) {
	seeds := []demoUserSeed{
		{
			Username:    "admin",
			Password:    "admin123",
			DisplayName: "系统管理员",
			Role:        model.RoleAdmin,
		},
		{
			Username:    "farm_user",
			Password:    "farm123",
			DisplayName: "养殖场操作员",
			Role:        model.RoleOrg,
			OrgLicense:  "LN-DL-2025-LS-018",
		},
		{
			Username:    "process_user",
			Password:    "process123",
			DisplayName: "加工厂质检员",
			Role:        model.RoleOrg,
			OrgLicense:  "LN-DL-2025-PR-102",
		},
		{
			Username:    "logistics_user",
			Password:    "logistics123",
			DisplayName: "冷链调度员",
			Role:        model.RoleOrg,
			OrgLicense:  "LN-DL-2025-LG-205",
		},
		{
			Username:    "retail_user",
			Password:    "retail123",
			DisplayName: "门店管理员",
			Role:        model.RoleOrg,
			OrgLicense:  "LN-DL-2025-RT-309",
		},
		{
			Username:    "orguser",
			Password:    "org123",
			DisplayName: "企业演示账号",
			Role:        model.RoleOrg,
			OrgLicense:  "LN-DL-2025-LS-018",
		},
	}

	out := make(map[string]model.User, len(seeds))
	for _, seed := range seeds {
		hash, err := bcrypt.GenerateFromPassword([]byte(seed.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user := model.User{
			Username:     seed.Username,
			PasswordHash: string(hash),
			DisplayName:  seed.DisplayName,
			Role:         seed.Role,
		}
		if seed.OrgLicense != "" {
			org := orgs[seed.OrgLicense]
			user.OrgID = &org.ID
		}
		if err := tx.Where("username = ?", seed.Username).Assign(user).FirstOrCreate(&user).Error; err != nil {
			return nil, err
		}
		out[seed.Username] = user
	}
	return out, nil
}

func seedBatches(tx *gorm.DB, orgs map[string]model.Org) (map[string]model.SeaCucumberBatch, error) {
	breedStartA := time.Date(2025, 2, 18, 0, 0, 0, 0, time.Local)
	catchDateA := time.Date(2025, 6, 2, 6, 30, 0, 0, time.Local)
	breedStartB := time.Date(2025, 3, 12, 0, 0, 0, 0, time.Local)
	catchDateB := time.Date(2025, 7, 8, 8, 15, 0, 0, time.Local)

	seeds := []demoBatchSeed{
		{
			BatchNo:        "HSC-2025-DL-PREMIUM-001",
			OrgLicense:     "LN-DL-2025-LS-018",
			ProductName:    "大连淡干海参精品装",
			FarmBase:       "旅顺口双岛湾海参养殖基地",
			Quality:        "合格，蛋白及感官指标通过抽检",
			CatchDate:      &catchDateA,
			BreedArea:      "旅顺口区底播养殖区 S2",
			BreedStartDate: &breedStartA,
			Spec:           "60-80 头/500g",
			Quantity:       "500 kg",
			Extra: map[string]any{
				"inspector":   "大连市农产品质量检验中心",
				"standardRef": "GB/T 34747-2017",
				"waterTemp":   "14C",
			},
		},
		{
			BatchNo:        "HSC-2025-DL-FRESH-002",
			OrgLicense:     "LN-DL-2025-LS-018",
			ProductName:    "鲜活即食海参礼盒",
			FarmBase:       "旅顺口近海生态浮筏基地",
			Quality:        "合格，活性与净含量符合企业标准",
			CatchDate:      &catchDateB,
			BreedArea:      "旅顺口区近海养殖区 F7",
			BreedStartDate: &breedStartB,
			Spec:           "单只 80-100g",
			Quantity:       "320 kg",
			Extra: map[string]any{
				"inspector":   "旅顺口区海洋食品检测站",
				"standardRef": "Q/DLSL 0003S-2025",
				"salinity":    "31ppt",
			},
		},
	}

	out := make(map[string]model.SeaCucumberBatch, len(seeds))
	for _, seed := range seeds {
		extraJSON, err := mustJSON(seed.Extra)
		if err != nil {
			return nil, err
		}
		org := orgs[seed.OrgLicense]
		batch := model.SeaCucumberBatch{
			BatchNo:        seed.BatchNo,
			OrgID:          org.ID,
			ProductName:    seed.ProductName,
			FarmBase:       seed.FarmBase,
			Quality:        seed.Quality,
			CatchDate:      seed.CatchDate,
			BreedArea:      seed.BreedArea,
			BreedStartDate: seed.BreedStartDate,
			Spec:           seed.Spec,
			Quantity:       seed.Quantity,
			ExtraJSON:      extraJSON,
		}
		if err := tx.Where("batch_no = ?", seed.BatchNo).Assign(batch).FirstOrCreate(&batch).Error; err != nil {
			return nil, err
		}
		batch.Org = org
		out[seed.BatchNo] = batch
	}
	return out, nil
}

func seedEvents(tx *gorm.DB, batches map[string]model.SeaCucumberBatch, users map[string]model.User) error {
	seeds := []demoEventSeed{
		{
			BatchNo:       "HSC-2025-DL-PREMIUM-001",
			Stage:         model.StageBreeding,
			Title:         "养殖建档",
			Location:      "旅顺口双岛湾海参养殖基地",
			Operator:      "王海波",
			EvidenceURLs:  "https://example.com/evidence/premium-001-breeding",
			OccurredAt:    time.Date(2025, 2, 18, 9, 0, 0, 0, time.Local),
			CreatedByUser: "farm_user",
			Detail: map[string]any{
				"seedSource": "獐子岛苗种中心",
				"waterDepth": "13m",
				"density":    "1800 头/亩",
			},
			TxID:        "seed_anchor_premium_001_breeding",
			BlockNumber: 1048576,
		},
		{
			BatchNo:       "HSC-2025-DL-PREMIUM-001",
			Stage:         model.StageHarvest,
			Title:         "采捕与初检",
			Location:      "旅顺口双岛湾码头",
			Operator:      "李明哲",
			EvidenceURLs:  "https://example.com/evidence/premium-001-harvest",
			OccurredAt:    time.Date(2025, 6, 2, 6, 30, 0, 0, time.Local),
			CreatedByUser: "farm_user",
			Detail: map[string]any{
				"sampling":  "抽检 30 只",
				"result":    "合格",
				"transport": "全程冷藏筐转运",
			},
			TxID:        "seed_anchor_premium_001_harvest",
			BlockNumber: 1048577,
		},
		{
			BatchNo:       "HSC-2025-DL-PREMIUM-001",
			Stage:         model.StageProcessing,
			Title:         "分级与清洗",
			Location:      "大连深蓝海产品加工有限公司 A 车间",
			Operator:      "周倩",
			EvidenceURLs:  "https://example.com/evidence/premium-001-processing",
			OccurredAt:    time.Date(2025, 6, 2, 14, 20, 0, 0, time.Local),
			CreatedByUser: "process_user",
			Detail: map[string]any{
				"grade":       "A",
				"cleanMethod": "低温净化 8 小时",
				"packageTemp": "4C",
			},
			TxID:        "seed_anchor_premium_001_processing",
			BlockNumber: 1048578,
		},
		{
			BatchNo:       "HSC-2025-DL-PREMIUM-001",
			Stage:         model.StageLogistics,
			Title:         "冷链发运",
			Location:      "渤海冷链物流大连分拨中心",
			Operator:      "陈宇",
			EvidenceURLs:  "https://example.com/evidence/premium-001-logistics",
			OccurredAt:    time.Date(2025, 6, 3, 8, 10, 0, 0, time.Local),
			CreatedByUser: "logistics_user",
			Detail: map[string]any{
				"vehicleNo": "辽B CL908",
				"tempRange": "-2C ~ 2C",
				"route":     "大连 -> 沈阳",
			},
			TxID:        "seed_anchor_premium_001_logistics",
			BlockNumber: 1048579,
		},
		{
			BatchNo:       "HSC-2025-DL-PREMIUM-001",
			Stage:         model.StageRetail,
			Title:         "门店上架",
			Location:      "鲜参优选零售中心 中山店",
			Operator:      "林洁",
			EvidenceURLs:  "https://example.com/evidence/premium-001-retail",
			OccurredAt:    time.Date(2025, 6, 4, 10, 45, 0, 0, time.Local),
			CreatedByUser: "retail_user",
			Detail: map[string]any{
				"shelfLife": "30 天",
				"qrStatus":  "已张贴",
				"stock":     "200 盒",
			},
			TxID:        "seed_anchor_premium_001_retail",
			BlockNumber: 1048580,
		},
		{
			BatchNo:       "HSC-2025-DL-FRESH-002",
			Stage:         model.StageBreeding,
			Title:         "苗种投放",
			Location:      "旅顺口近海生态浮筏基地",
			Operator:      "孙鹏",
			EvidenceURLs:  "https://example.com/evidence/fresh-002-breeding",
			OccurredAt:    time.Date(2025, 3, 12, 7, 40, 0, 0, time.Local),
			CreatedByUser: "farm_user",
			Detail: map[string]any{
				"seedSource": "辽宁海洋大学育苗站",
				"raftNo":     "F7-12",
				"feedPlan":   "海藻天然增殖",
			},
			TxID:        "seed_anchor_fresh_002_breeding",
			BlockNumber: 1048581,
		},
		{
			BatchNo:       "HSC-2025-DL-FRESH-002",
			Stage:         model.StagePackaging,
			Title:         "即食包装",
			Location:      "大连深蓝海产品加工有限公司 无菌包装间",
			Operator:      "赵宁",
			EvidenceURLs:  "https://example.com/evidence/fresh-002-packaging",
			OccurredAt:    time.Date(2025, 7, 8, 16, 0, 0, 0, time.Local),
			CreatedByUser: "process_user",
			Detail: map[string]any{
				"packageType": "锁鲜袋 + 礼盒",
				"netWeight":   "1200g",
				"sterility":   "合格",
			},
			TxID:        "seed_anchor_fresh_002_packaging",
			BlockNumber: 1048582,
		},
		{
			BatchNo:       "HSC-2025-DL-FRESH-002",
			Stage:         model.StageLogistics,
			Title:         "同城冷配送",
			Location:      "渤海冷链物流即时配送仓",
			Operator:      "郭峰",
			EvidenceURLs:  "https://example.com/evidence/fresh-002-logistics",
			OccurredAt:    time.Date(2025, 7, 9, 6, 50, 0, 0, time.Local),
			CreatedByUser: "logistics_user",
			Detail: map[string]any{
				"vehicleNo":   "辽B DD512",
				"dispatchETA": "2 小时",
				"tempRange":   "0C ~ 4C",
			},
			TxID:        "seed_anchor_fresh_002_logistics",
			BlockNumber: 1048583,
		},
	}

	for _, seed := range seeds {
		batch := batches[seed.BatchNo]
		user := users[seed.CreatedByUser]
		detailJSON, err := mustJSON(seed.Detail)
		if err != nil {
			return err
		}
		dataHash := fabric.HashTraceEvent(
			batch.BatchNo,
			string(seed.Stage),
			detailJSON,
			seed.Location,
			seed.Operator,
			seed.OccurredAt,
			seed.EvidenceURLs,
		)

		event := model.TraceEvent{
			BatchID:      batch.ID,
			Stage:        seed.Stage,
			Title:        seed.Title,
			DetailJSON:   detailJSON,
			Location:     seed.Location,
			OperatorName: seed.Operator,
			EvidenceURLs: seed.EvidenceURLs,
			OccurredAt:   seed.OccurredAt,
			DataHash:     dataHash,
			CreatedBy:    user.ID,
		}
		if err := tx.Where("data_hash = ?", dataHash).Assign(event).FirstOrCreate(&event).Error; err != nil {
			return err
		}

		payloadJSON, err := mustJSON(map[string]any{
			"batchNo":    batch.BatchNo,
			"eventID":    event.ID,
			"stage":      string(seed.Stage),
			"dataHash":   dataHash,
			"occurredAt": seed.OccurredAt.UTC().Format(time.RFC3339Nano),
			"orgName":    batch.Org.Name,
			"seeded":     true,
		})
		if err != nil {
			return err
		}
		record := model.ChainRecord{
			RefType:     "event",
			RefID:       event.ID,
			ChainType:   "fabric",
			TxID:        seed.TxID,
			BlockNumber: seed.BlockNumber,
			Status:      model.ChainSuccess,
			PayloadJSON: payloadJSON,
		}
		if err := tx.Where("ref_type = ? AND ref_id = ?", "event", event.ID).Assign(record).FirstOrCreate(&record).Error; err != nil {
			return err
		}
	}

	return nil
}

func mustJSON(v any) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal seed json: %w", err)
	}
	return string(raw), nil
}
