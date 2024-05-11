package retry

import (
	"context"
	"errors"
	"time"
)

var (
	panicErr    = errors.New("Redundant.Do execute fn() panic, but recover")
	canceledErr = errors.New("redundant.Do canceled by external context, context canceled")
)

type Redundant[Resp any] struct {
	fn func() (Resp, error)

	count    int
	nextTime int // 单位 ms

	errCount int
	respChan chan Resp
	errChan  chan error
}

func NewRedundant[Resp any](count, nextTime int, fn func() (Resp, error)) *Redundant[Resp] {
	if count <= 0 && nextTime <= 0 {
		count = 3
	}
	if count <= 0 && nextTime > 0 {
		count = 2
	}
	return &Redundant[Resp]{
		fn:       fn,
		count:    count,
		nextTime: nextTime,
		respChan: make(chan Resp, count),
		errChan:  make(chan error, count),
	}
}

// Do 每一个 G 并行执行冗余请求，有一个执行成功 Do 即刻返回，都执行失败（即 fn 返回 error）或内部 panic 返回 error，注意：如果全部失败，只会返回最后一个请求的 error
func (r *Redundant[Resp]) Do(ctx context.Context) (resp Resp, err error) {
	for i := 0; i < r.count; i++ {
		if i != 0 && r.nextTime != 0 {
			select {
			case <-ctx.Done():
				err = canceledErr
				return
			case resp = <-r.respChan:
				return
			case err = <-r.errChan:
				r.errCount++
			case <-time.After(time.Duration(r.nextTime) * time.Millisecond):
			}
		}
		go func() {
			defer func() {
				if re := recover(); re != nil {
					r.errChan <- panicErr
				}
			}()
			if resp, err := r.fn(); err != nil {
				r.errChan <- err
			} else {
				r.respChan <- resp
			}
		}()
	}
	for {
		select {
		case <-ctx.Done():
			err = canceledErr
			return
		case resp = <-r.respChan:
			return
		case err = <-r.errChan:
			r.errCount++
			if r.errCount == r.count {
				return
			}
		}
	}
}
