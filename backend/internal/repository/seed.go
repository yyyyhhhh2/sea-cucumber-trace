package repository

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"sea-cucumber-trace/backend/internal/fabric"
	"sea-cucumber-trace/backend/internal/model"
)

func (r *Repository) SeedIfEmpty() error {
	var n int64
	if err := r.db.Model(&model.User{}).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hashAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	hashOrg, _ := bcrypt.GenerateFromPassword([]byte("org123"), bcrypt.DefaultCost)

	org := model.Org{
		Name:        "大连旅顺口海洋产业合作社",
		Type:        model.OrgFarm,
		LicenseNo:   "LN-DL-2025-旅顺-018",
		Address:     "辽宁省大连市旅顺口区",
		Contact:     "0411-0000-0000",
		Description: "刺参底播养殖与链上溯源示范主体",
	}
	if err := r.db.Create(&org).Error; err != nil {
		return err
	}
	oid := org.ID

	uAdmin := model.User{
		Username:     "admin",
		PasswordHash: string(hashAdmin),
		DisplayName:  "系统管理员",
		Role:         model.RoleAdmin,
	}
	uOrg := model.User{
		Username:     "orguser",
		PasswordHash: string(hashOrg),
		DisplayName:  "牧场操作员",
		Role:         model.RoleOrg,
		OrgID:        &oid,
	}
	if err := r.db.Create(&uAdmin).Error; err != nil {
		return err
	}
	if err := r.db.Create(&uOrg).Error; err != nil {
		return err
	}

	breedStart := time.Date(2025, 3, 1, 0, 0, 0, 0, time.Local)
	catchDate := time.Date(2025, 6, 2, 10, 0, 0, 0, time.Local)

	batch := model.SeaCucumberBatch{
		BatchNo:        "HSC-2025-DL-PREMIUM-001",
		OrgID:          oid,
		ProductName:    "大连刺参 (精品级)",
		FarmBase:       "大连旅顺口区海参养殖场",
		Quality:        "合格 — 国家A级标准",
		CatchDate:      &catchDate,
		BreedArea:      "旅顺口区底播养殖区 S2",
		BreedStartDate: &breedStart,
		Spec:           "精品级 / 淡干",
		Quantity:       "500 kg",
		ExtraJSON:      `{"inspector":"大连市农产品质检中心","standardRef":"GB/T 33184-2016"}`,
	}
	if err := r.db.Create(&batch).Error; err != nil {
		return err
	}

	detailJSON := `{"quality":"合格 — 国家A级标准","step":"采捕后质检与入链锚定"}`
	loc := "大连旅顺口区海参养殖场"
	op := "质检组·李娜"
	ev := model.TraceEvent{
		BatchID:      batch.ID,
		Stage:        model.StageHarvest,
		Title:        "采捕与质检上链",
		DetailJSON:   detailJSON,
		Location:     loc,
		OperatorName: op,
		EvidenceURLs: "",
		OccurredAt:   catchDate,
		DataHash: fabric.HashTraceEvent(
			batch.BatchNo,
			string(model.StageHarvest),
			detailJSON,
			loc,
			op,
			catchDate,
			"",
		),
		CreatedBy: uOrg.ID,
	}
	if err := r.db.Create(&ev).Error; err != nil {
		return err
	}

	cr := model.ChainRecord{
		RefType:     "event",
		RefID:       ev.ID,
		ChainType:   "fabric",
		TxID:        "seed_anchor_tx_dalian_premium_001",
		BlockNumber: 1048576,
		Status:      model.ChainSuccess,
		PayloadJSON: `{"note":"演示种子数据：区块高度对应 BlockNumber 1048576（展示为 #1,048,576）"}`,
	}
	if err := r.db.Create(&cr).Error; err != nil {
		return err
	}

	return nil
}
