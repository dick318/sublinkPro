package models

import (
	"strconv"
	"testing"
	"time"

	"sublink/cache"
	"sublink/database"
	"sublink/internal/testutil"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func resetSubscriptionShareCacheForTest() {
	subscriptionShareCache = cache.NewMapCache(func(s SubscriptionShare) int { return s.ID })
	subscriptionShareCache.AddIndex("token", func(s SubscriptionShare) string { return s.Token })
	subscriptionShareCache.AddIndex("subscriptionID", func(s SubscriptionShare) string { return strconv.Itoa(s.SubscriptionID) })
}

func setupSubscriptionShareTestDB(t *testing.T) {
	t.Helper()

	oldDB := database.DB
	oldDialect := database.Dialect
	oldInitialized := database.IsInitialized

	db, err := gorm.Open(sqlite.Open(testutil.UniqueMemoryDSN(t, "subscription_share_test")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(&SubscriptionShare{}); err != nil {
		t.Fatalf("auto migrate subscription_shares: %v", err)
	}

	database.DB = db
	database.Dialect = database.DialectSQLite
	database.IsInitialized = false
	resetSubscriptionShareCacheForTest()

	t.Cleanup(func() {
		database.DB = oldDB
		database.Dialect = oldDialect
		database.IsInitialized = oldInitialized
		resetSubscriptionShareCacheForTest()
		testutil.CloseDB(t, db)
	})
}

func setupSubcriptionModelTestDB(t *testing.T) {
	t.Helper()

	oldDB := database.DB
	oldDialect := database.Dialect
	oldInitialized := database.IsInitialized

	db, err := gorm.Open(sqlite.Open(testutil.UniqueMemoryDSN(t, "subcription_model_test")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&Subcription{},
		&Node{},
		&SubcriptionNode{},
		&SubcriptionGroup{},
		&SubcriptionScript{},
		&Script{},
		&SubLogs{},
		&SystemSetting{},
		&SubscriptionShare{},
		&SubscriptionChainRule{},
	); err != nil {
		t.Fatalf("auto migrate subscription model tables: %v", err)
	}

	database.DB = db
	database.Dialect = database.DialectSQLite
	database.IsInitialized = false
	if err := InitNodeCache(); err != nil {
		t.Fatalf("init node cache: %v", err)
	}
	if err := InitSubcriptionCache(); err != nil {
		t.Fatalf("init subcription cache: %v", err)
	}
	if err := InitSubLogsCache(); err != nil {
		t.Fatalf("init sub logs cache: %v", err)
	}
	if err := InitSubscriptionShareCache(); err != nil {
		t.Fatalf("init subscription share cache: %v", err)
	}
	if err := InitChainRuleCache(); err != nil {
		t.Fatalf("init chain rule cache: %v", err)
	}
	if err := InitSettingCache(); err != nil {
		t.Fatalf("init setting cache: %v", err)
	}

	t.Cleanup(func() {
		database.DB = oldDB
		database.Dialect = oldDialect
		database.IsInitialized = oldInitialized
		if oldDB != nil {
			_ = InitNodeCache()
			_ = InitSubcriptionCache()
			_ = InitSubLogsCache()
			_ = InitSubscriptionShareCache()
			_ = InitChainRuleCache()
			_ = InitSettingCache()
		}
		testutil.CloseDB(t, db)
	})
}

func TestSubscriptionShareAddNormalizesOptionalTimestamps(t *testing.T) {
	setupSubscriptionShareTestDB(t)

	zero := time.Time{}
	share := &SubscriptionShare{
		SubscriptionID: 1,
		Name:           "never-expire",
		ExpireType:     ExpireTypeNever,
		ExpireAt:       &zero,
		LastAccessAt:   &zero,
	}

	if err := share.Add(); err != nil {
		t.Fatalf("add share: %v", err)
	}

	var stored SubscriptionShare
	if err := database.DB.First(&stored, share.ID).Error; err != nil {
		t.Fatalf("reload share: %v", err)
	}

	if stored.ExpireAt != nil {
		t.Fatalf("expected expire_at to be nil, got %v", stored.ExpireAt)
	}
	if stored.LastAccessAt != nil {
		t.Fatalf("expected last_access_at to be nil, got %v", stored.LastAccessAt)
	}
}

func TestSubscriptionShareUpdateClearsExpireAtForNonDateTime(t *testing.T) {
	setupSubscriptionShareTestDB(t)

	expireAt := time.Now().Add(24 * time.Hour).Round(time.Second)
	share := &SubscriptionShare{
		SubscriptionID: 1,
		Name:           "datetime-expire",
		ExpireType:     ExpireTypeDateTime,
		ExpireAt:       &expireAt,
	}

	if err := share.Add(); err != nil {
		t.Fatalf("add share: %v", err)
	}

	share.ExpireType = ExpireTypeNever
	share.ExpireAt = &expireAt
	if err := share.Update(); err != nil {
		t.Fatalf("update share: %v", err)
	}

	var stored SubscriptionShare
	if err := database.DB.First(&stored, share.ID).Error; err != nil {
		t.Fatalf("reload share: %v", err)
	}

	if stored.ExpireAt != nil {
		t.Fatalf("expected expire_at to be cleared, got %v", stored.ExpireAt)
	}
}

func TestSubscriptionShareRecordAccessSetsLastAccessAt(t *testing.T) {
	setupSubscriptionShareTestDB(t)

	share := &SubscriptionShare{
		SubscriptionID: 1,
		Name:           "record-access",
		ExpireType:     ExpireTypeNever,
	}

	if err := share.Add(); err != nil {
		t.Fatalf("add share: %v", err)
	}

	share.RecordAccess()

	var stored SubscriptionShare
	if err := database.DB.First(&stored, share.ID).Error; err != nil {
		t.Fatalf("reload share: %v", err)
	}

	if stored.AccessCount != 1 {
		t.Fatalf("expected access_count=1, got %d", stored.AccessCount)
	}
	if stored.LastAccessAt == nil || stored.LastAccessAt.IsZero() {
		t.Fatalf("expected last_access_at to be set, got %v", stored.LastAccessAt)
	}
}

func TestSubcriptionAddNodeSkipsDuplicateNodeIDs(t *testing.T) {
	setupSubcriptionModelTestDB(t)

	sub := &Subcription{Name: "去重测试订阅"}
	if err := sub.Add(); err != nil {
		t.Fatalf("add subscription: %v", err)
	}

	node := Node{
		Name:     "节点A",
		LinkName: "node-a",
		Link:     "ss://YWVzLTEyOC1nY206cGFzc0EyQGV4YW1wbGUuY29tOjQ0Mw==#node-a",
		Protocol: "ss",
		Source:   "manual",
	}
	if err := node.Add(); err != nil {
		t.Fatalf("add node: %v", err)
	}

	sub.Nodes = []Node{node, node}
	if err := sub.AddNode(); err != nil {
		t.Fatalf("expected duplicate node ids to be tolerated, got %v", err)
	}

	var relations []SubcriptionNode
	if err := database.DB.Where("subcription_id = ?", sub.ID).Order("sort asc").Find(&relations).Error; err != nil {
		t.Fatalf("load relations: %v", err)
	}
	if len(relations) != 1 {
		t.Fatalf("expected 1 relation after deduplication, got %d", len(relations))
	}
	if relations[0].NodeID != node.ID || relations[0].Sort != 0 {
		t.Fatalf("unexpected relation after deduplication: %+v", relations[0])
	}
}

func TestSubcriptionUpdateNodesReplacesRelationsWithDedupedIDs(t *testing.T) {
	setupSubcriptionModelTestDB(t)

	sub := &Subcription{Name: "更新去重测试订阅"}
	if err := sub.Add(); err != nil {
		t.Fatalf("add subscription: %v", err)
	}

	first := Node{
		Name:     "节点1",
		LinkName: "node-1",
		Link:     "ss://YWVzLTEyOC1nY206cGFzczEyQGV4YW1wbGUuY29tOjQ0Mw==#node-1",
		Protocol: "ss",
		Source:   "manual",
	}
	second := Node{
		Name:     "节点2",
		LinkName: "node-2",
		Link:     "ss://YWVzLTEyOC1nY206cGFzczIyQGV4YW1wbGUuY29tOjQ0Mw==#node-2",
		Protocol: "ss",
		Source:   "manual",
	}
	if err := first.Add(); err != nil {
		t.Fatalf("add first node: %v", err)
	}
	if err := second.Add(); err != nil {
		t.Fatalf("add second node: %v", err)
	}

	sub.Nodes = []Node{first, first, second}
	if err := sub.UpdateNodes(); err != nil {
		t.Fatalf("expected duplicate node ids to be tolerated on update, got %v", err)
	}

	var relations []SubcriptionNode
	if err := database.DB.Where("subcription_id = ?", sub.ID).Order("sort asc").Find(&relations).Error; err != nil {
		t.Fatalf("load relations: %v", err)
	}
	if len(relations) != 2 {
		t.Fatalf("expected 2 relations after deduplication, got %d", len(relations))
	}
	if relations[0].NodeID != first.ID || relations[0].Sort != 0 {
		t.Fatalf("unexpected first relation: %+v", relations[0])
	}
	if relations[1].NodeID != second.ID || relations[1].Sort != 1 {
		t.Fatalf("unexpected second relation: %+v", relations[1])
	}
}

func TestSubcriptionDelRemovesAssociatedSubLogs(t *testing.T) {
	setupSubcriptionModelTestDB(t)

	sub := &Subcription{Name: "删除日志测试订阅"}
	if err := sub.Add(); err != nil {
		t.Fatalf("add subscription: %v", err)
	}

	log := &SubLogs{
		IP:            "127.0.0.1",
		Date:          "2026-04-06 14:32:55",
		Addr:          "本地",
		Count:         1,
		SubcriptionID: sub.ID,
	}
	if err := log.Add(); err != nil {
		t.Fatalf("add sub log: %v", err)
	}

	if err := sub.Del(); err != nil {
		t.Fatalf("delete subscription: %v", err)
	}

	var count int64
	if err := database.DB.Model(&SubLogs{}).Where("subcription_id = ?", sub.ID).Count(&count).Error; err != nil {
		t.Fatalf("count sub logs: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected associated sub logs to be deleted, got %d", count)
	}
}

func TestSubcriptionGetSubKeepsSameNameNodesWhenOutputDedupDisabled(t *testing.T) {
	setupSubcriptionModelTestDB(t)

	sub := &Subcription{Name: "输出去重测试订阅"}
	if err := sub.Add(); err != nil {
		t.Fatalf("add subscription: %v", err)
	}

	first := Node{
		Name:     "同名节点",
		LinkName: "same-output-a",
		Link:     "ss://YWVzLTEyOC1nY206b3V0cHV0MUFleGFtcGxlLmNvbTo0NDM=#same-output-a",
		Protocol: "ss",
		Source:   "manual",
	}
	second := Node{
		Name:     "同名节点",
		LinkName: "same-output-b",
		Link:     "ss://YWVzLTEyOC1nY206b3V0cHV0MkJleGFtcGxlLmNvbTo0NDM=#same-output-b",
		Protocol: "ss",
		Source:   "manual",
	}
	if err := first.Add(); err != nil {
		t.Fatalf("add first node: %v", err)
	}
	if err := second.Add(); err != nil {
		t.Fatalf("add second node: %v", err)
	}

	sub.Nodes = []Node{first, second}
	if err := sub.AddNode(); err != nil {
		t.Fatalf("add node relations: %v", err)
	}
	if err := SetSetting("subscription_output_name_dedup_enabled", "false"); err != nil {
		t.Fatalf("set subscription output dedup setting: %v", err)
	}

	if err := sub.GetSub("clash"); err != nil {
		t.Fatalf("get subscription output: %v", err)
	}

	if len(sub.Nodes) != 2 {
		t.Fatalf("expected same-name nodes to be preserved when output dedup is disabled, got %d", len(sub.Nodes))
	}
}
