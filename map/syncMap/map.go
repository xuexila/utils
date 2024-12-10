package syncMap

import (
	"fmt"
	"github.com/cespare/xxhash/v2"
	"reflect"
	"sync"
	"time"
)

const defaultNumSegments = 32
const defaultCleanInterval = 10 * time.Second

// Map 使用分段锁实现线程安全的泛型 map
type Map[K comparable, V any] struct {
	segments     []*segment[K, V]
	wg           sync.WaitGroup
	stopCh       chan struct{}
	closed       bool
	numSegments  int
	segmentMutex sync.Mutex
}

type segment[K comparable, V any] struct {
	mu          sync.RWMutex
	m           map[K]V
	cleaner     *cleaner[K]
	initialized bool
	initOnce    sync.Once
}

type cleaner[K comparable] struct {
	mu      sync.Mutex
	pending map[K]struct{}
}

// Config 用于配置 Map 的选项
type Config[K comparable, V any] struct {
	NumSegments   int
	CleanInterval time.Duration
}

// NewMap 创建一个新的 Map
func NewMap[K comparable, V any](cleanInterval ...time.Duration) *Map[K, V] {
	interval := defaultCleanInterval
	if len(cleanInterval) > 0 {
		interval = cleanInterval[0]
	}
	return newMap[K, V](defaultNumSegments, interval)
}

// NewMapWithConfig 使用配置选项创建一个新的 Map
func NewMapWithConfig[K comparable, V any](config Config[K, V]) *Map[K, V] {
	interval := config.CleanInterval
	if interval == 0 {
		interval = defaultCleanInterval
	}
	return newMap[K, V](config.NumSegments, interval)
}

// Close 关闭 Map 并等待所有后台任务完成
func (sm *Map[K, V]) Close() {
	close(sm.stopCh)
	sm.wg.Wait()
	sm.closed = true
}

// IsClosed 检查 Map 是否已关闭
func (sm *Map[K, V]) IsClosed() bool {
	return sm.closed
}

// hash 计算哈希值，用于选择段
func (sm *Map[K, V]) hash(key K) int {
	v := reflect.ValueOf(key)
	switch v.Kind() {
	case reflect.String:
		return int(xxhash.Sum64String(v.String()) % uint64(sm.numSegments))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(uint64(v.Int()) % uint64(sm.numSegments))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(v.Uint() % uint64(sm.numSegments))
	default:
		return int(xxhash.Sum64String(fmt.Sprintf("%v", key)) % uint64(sm.numSegments))
	}
}

// Load 获取键对应的值
func (sm *Map[K, V]) Load(key K) (V, bool) {
	if err := sm.checkClosed(); err != nil {
		return *new(V), false
	}
	seg := sm.segments[sm.hash(key)]
	seg.mu.RLock()
	defer seg.mu.RUnlock()
	value, ok := seg.m[key]
	return value, ok
}

// Store 设置键对应的值
func (sm *Map[K, V]) Store(key K, value V) error {
	if err := sm.checkClosed(); err != nil {
		return err
	}
	seg := sm.segments[sm.hash(key)]
	seg.mu.Lock()
	defer seg.mu.Unlock()

	seg.initOnce.Do(func() {
		seg.initCleaner()
	})

	seg.m[key] = value
	if seg.initialized {
		seg.cleaner.mu.Lock()
		delete(seg.cleaner.pending, key)
		seg.cleaner.mu.Unlock()
	}
	return nil
}

// Delete 删除键对应的值
func (sm *Map[K, V]) Delete(key K) {
	if err := sm.checkClosed(); err != nil {
		return
	}
	seg := sm.segments[sm.hash(key)]
	seg.mu.Lock()
	defer seg.mu.Unlock()

	if seg.initialized {
		seg.cleaner.mu.Lock()
		seg.cleaner.pending[key] = struct{}{}
		seg.cleaner.mu.Unlock()
	}
	delete(seg.m, key)
}

// initCleaner 初始化 cleaner，确保线程安全
func (seg *segment[K, V]) initCleaner() {
	seg.initialized = true
	seg.cleaner = &cleaner[K]{pending: make(map[K]struct{})}
}

// clean 清理已删除的条目
func (seg *segment[K, V]) clean() {
	if seg.cleaner == nil {
		return // 或者返回一个错误
	}

	seg.cleaner.mu.Lock()
	pending := make(map[K]struct{})
	for k := range seg.cleaner.pending {
		pending[k] = struct{}{}
	}
	seg.cleaner.mu.Unlock()

	if seg.m == nil {
		return // 或者返回一个错误
	}

	seg.mu.RLock()
	for key := range pending {
		_, exists := seg.m[key]
		if !exists {
			seg.cleaner.mu.Lock()
			delete(seg.cleaner.pending, key)
			seg.cleaner.mu.Unlock()
		}
	}
	seg.mu.RUnlock()
}

// startCleaner 启动后台清理 goroutine
func (sm *Map[K, V]) startCleaner(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer sm.wg.Done()

	var cleanWg sync.WaitGroup

	for {
		select {
		case <-ticker.C:
			cleanWg.Add(len(sm.segments))
			for _, seg := range sm.segments {
				seg := seg
				go func() {
					defer cleanWg.Done()
					seg.clean()
				}()
			}
			cleanWg.Wait()
		case <-sm.stopCh:
			return
		}
	}
}

// Range 遍历 map 中的所有键值对
func (sm *Map[K, V]) Range(f func(key K, value V) bool, batchSize ...int) error {
	if err := sm.checkClosed(); err != nil {
		return err
	}
	batch := 100
	if len(batchSize) > 0 {
		batch = batchSize[0]
	}

	var wg sync.WaitGroup
	errCh := make(chan error, sm.numSegments)

	// 复用切片
	keysPool := sync.Pool{
		New: func() interface{} {
			return make([]K, 0, 100)
		},
	}

	for _, seg := range sm.segments {
		seg := seg
		wg.Add(1)
		go func() {
			defer wg.Done()
			keys := keysPool.Get().([]K)
			defer keysPool.Put(keys)

			seg.mu.RLock()
			keys = keys[:0]
			for k := range seg.m {
				keys = append(keys, k)
			}
			seg.mu.RUnlock()

			for i := 0; i < len(keys); i += batch {
				end := i + batch
				if end > len(keys) {
					end = len(keys)
				}
				for _, k := range keys[i:end] {
					seg.mu.RLock()
					v, exists := seg.m[k]
					seg.mu.RUnlock()
					if exists && !f(k, v) {
						errCh <- nil
						return
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

// initSegments 初始化段
func (sm *Map[K, V]) initSegments() {
	sm.segmentMutex.Lock()
	defer sm.segmentMutex.Unlock()
	sm.segments = make([]*segment[K, V], sm.numSegments)
	for i := range sm.segments {
		sm.segments[i] = &segment[K, V]{m: make(map[K]V)}
	}
}

// resizeSegments 动态调整段的数量
func (sm *Map[K, V]) resizeSegments(newNumSegments int) {
	sm.segmentMutex.Lock()
	defer sm.segmentMutex.Unlock()

	newSegments := make([]*segment[K, V], newNumSegments)
	for i := range newSegments {
		newSegments[i] = &segment[K, V]{m: make(map[K]V)}
	}

	var wg sync.WaitGroup
	for _, seg := range sm.segments {
		seg := seg
		wg.Add(1)
		go func() {
			defer wg.Done()
			seg.mu.RLock()
			for k, v := range seg.m {
				newHash := sm.hash(k)
				newSegments[newHash].m[k] = v
			}
			seg.mu.RUnlock()
		}()
	}
	wg.Wait()

	sm.segments = newSegments
	sm.numSegments = newNumSegments
}

// checkClosed 检查 Map 是否已关闭
func (sm *Map[K, V]) checkClosed() error {
	if sm.closed {
		return fmt.Errorf("map is already closed")
	}
	return nil
}

// 提取公共部分到一个私有方法中
func newMap[K comparable, V any](numSegments int, cleanInterval time.Duration) *Map[K, V] {
	sm := &Map[K, V]{numSegments: numSegments, stopCh: make(chan struct{})}
	sm.initSegments()
	sm.wg.Add(1)
	go sm.startCleaner(cleanInterval)
	return sm
}
