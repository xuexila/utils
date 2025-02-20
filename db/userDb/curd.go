package userDb

import (
	"errors"
	"fmt"
	"github.com/helays/utils/config"
	"gorm.io/gorm"
	"regexp"
)

type Model interface {
	Valid() error // 验证
}

// Create 通用创建函数（使用泛型）
func Create[T Model](tx *gorm.DB, src T) error {
	// 调用 Valid 方法进行验证
	if err := src.Valid(); err != nil {
		return fmt.Errorf("验证失败：%s", err.Error())
	}
	if config.Dbg {
		tx = tx.Debug()
	}
	// 使用 GORM 的 Create 方法插入数据
	if err := tx.Create(src).Error; err != nil {
		return fmt.Errorf("插入数据失败：%s", err.Error())
	}

	return nil
}

// QueryConfig 通用curd函数，配置结构
type QueryConfig struct {
	SelectQuery any
	SelectArgs  []any
	Omit        []string
	Query       any
	Args        []any
	Pk          string                    // 主键字段 id row_Id
	MustField   map[string]*regexp.Regexp // 必填字段，正则校验
}

// Update 通用更新函数，使用泛型
func Update[T Model](tx *gorm.DB, src T, c QueryConfig) error {
	// 调用 Valid 方法进行验证
	err := src.Valid()
	if err != nil {
		return fmt.Errorf("验证失败：%s", err.Error())
	}
	if config.Dbg {
		tx = tx.Debug()
	}
	// 搜索数据是否存在
	var old T

	if err = tx.Where(c.Query, c.Args...).Take(&old).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("数据不存在")
		}
		return err
	}

	_tx := tx.Where(c.Query, c.Args...)
	if c.Omit != nil && len(c.Omit) > 0 {
		_tx.Omit(c.Omit...)
	}
	if c.SelectQuery != nil {
		_tx.Select(c.SelectQuery, c.SelectArgs...)
	}
	// 使用 GORM 的 Update 方法更新数据
	if err = _tx.Save(src).Error; err != nil {
		return fmt.Errorf("更新数据失败：%s", err.Error())
	}
	return nil
}

func FindOne[T any](tx *gorm.DB, query any, args ...any) (T, error) {
	_tx := tx.Where(query, args...)
	var data T
	err := _tx.Take(&data).Error
	return data, err
}
