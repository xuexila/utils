package tableRotate

import "time"

type tableSplit struct {
	tableName  string
	createTime time.Time
}

// 定义切片类型
type byCreateTime []tableSplit

// Len 实现 sort.Interface 接口的 Len 方法
func (b byCreateTime) Len() int {
	return len(b)
}

// Less 实现 sort.Interface 接口的 Less 方法
func (b byCreateTime) Less(i, j int) bool {
	return b[i].createTime.After(b[j].createTime) // 按 createTime 降序排序
}

// Swap 实现 sort.Interface 接口的 Swap 方法
func (b byCreateTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
