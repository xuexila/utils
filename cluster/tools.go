package cluster

import (
	"context"
	"sync"
)

var (
	EnableCluster, isLeader bool
	lock                    sync.Mutex
)

// RunWithLeader 运行函数，当有master和slave的时候，只有master才能运行函数
func RunWithLeader(ch chan bool, ctx context.Context, call func(ctx context.Context)) {
	// 当未开启多节点模式的时候，只有一个节点，不区分master和slave
	if !EnableCluster {
		call(ctx)
		return
	}
	_ctx, cancel := context.WithCancel(context.Background())
	for {
		select {
		case <-ctx.Done():
			cancel()
			return
		case leader := <-ch:
			if !leader {
				cancel()
				setLeader(false) // 更新当前节点类型，
				_ctx, cancel = context.WithCancel(context.Background())
				continue
			}
			setLeader(true) // 更新当前节点类型，
			call(_ctx)
		}
	}
}

func setLeader(leader bool) {
	lock.Lock()
	defer lock.Unlock()
	isLeader = leader
}

func IsLeader() bool {
	lock.Lock()
	defer lock.Unlock()
	return isLeader
}
