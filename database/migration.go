package database

import (
	"sublink/utils"
	"time"
)

type Migration struct {
	ID        string `gorm:"primaryKey;size:191"`
	CreatedAt time.Time
}

func EnsureMigrationTable() error {
	if !DB.Migrator().HasTable(&Migration{}) {
		if err := DB.AutoMigrate(&Migration{}); err != nil {
			return err
		}
	}
	return nil
}

func ListMigrationIDs() (map[string]struct{}, error) {
	if err := EnsureMigrationTable(); err != nil {
		return nil, err
	}

	var ids []string
	if err := DB.Model(&Migration{}).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	result := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		result[id] = struct{}{}
	}
	return result, nil
}

func HasMigration(migrationID string) (bool, error) {
	if err := EnsureMigrationTable(); err != nil {
		return false, err
	}

	var count int64
	if err := DB.Model(&Migration{}).Where("id = ?", migrationID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func RecordMigration(migrationID string) error {
	return DB.Create(&Migration{
		ID:        migrationID,
		CreatedAt: time.Now(),
	}).Error
}

// RunAutoMigrate 执行自动迁移，如果 migrationID 已存在则跳过
func RunAutoMigrate(migrationID string, dst ...interface{}) error {
	// 确保 Migration 表存在
	if err := EnsureMigrationTable(); err != nil {
		return err
	}

	exists, err := HasMigration(migrationID)
	if err != nil {
		return err
	}
	if exists {
		// 已经执行过，跳过
		return nil
	}
	utils.Info("执行数据库升级任务：%s", migrationID)
	// 执行迁移
	if err := DB.AutoMigrate(dst...); err != nil {
		return err
	}

	// 记录迁移
	return RecordMigration(migrationID)
}

// RunCustomMigration 执行自定义迁移逻辑，如果 migrationID 已存在则跳过
func RunCustomMigration(migrationID string, action func() error) error {
	// 确保 Migration 表存在
	if err := EnsureMigrationTable(); err != nil {
		return err
	}

	exists, err := HasMigration(migrationID)
	if err != nil {
		return err
	}
	if exists {
		// 已经执行过，跳过
		return nil
	}
	utils.Info("执行数据库升级任务：%s", migrationID)

	// 执行自定义迁移逻辑
	if err := action(); err != nil {
		return err
	}

	// 记录迁移
	return RecordMigration(migrationID)
}
