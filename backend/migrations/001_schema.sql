-- 参考：GORM AutoMigrate 等价结构（SQLite 方言可微调）
-- 实际表结构以 model 为准；索引含：batch_no、org_id、trace_events.batch_id、
-- chain_records(ref_type,ref_id) 复合索引、data_hash 等，由 GORM tag 维护。

-- users, orgs, sea_cucumber_batches, trace_events, chain_records
-- 详见 backend/internal/model/models.go
