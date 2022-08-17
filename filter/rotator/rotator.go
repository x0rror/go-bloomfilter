package rotator

import (
	"context"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/filter"
	"sync"
	"time"
)

type newFilterFunc func(ctx context.Context) (filter.Filter, error)

type Rotator struct {
	ctx           context.Context
	cfg           config.RotatorConfig
	mutex         *sync.RWMutex
	newFilter     newFilterFunc
	current, next filter.Filter
}

func (r *Rotator) GetBitmap() bitmap.Bitmap {
	return r.current.GetBitmap()
}

func (r *Rotator) handleRotating(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			r.rotate()
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Rotator) rotate() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.current = r.next
	r.next, _ = r.newFilter(r.ctx)
	setExpireTTL(r.next, r.cfg.Freq*2)
	return nil
}

func (r *Rotator) Exist(data string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.current.Exist(data)
}

func (r *Rotator) Add(data string) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	err := r.current.Add(data)
	if err != nil {
		return err
	}
	return r.next.Add(data)
}

// setExpireTTL sets Expire TTL if filter.bitmap is Redis
func setExpireTTL(filter filter.Filter, d time.Duration) error {
	gracefulTTL := d + time.Minute
	if bm, ok := filter.GetBitmap().(*bitmap.Redis); ok {
		// set expiry if bitmap is Redis to avoid wasting capacity of Redis
		err := bm.SetExpireTTL(gracefulTTL)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewRotator(ctx context.Context, cfg config.RotatorConfig, newFilter newFilterFunc) (*Rotator, error) {
	current, err := newFilter(ctx)
	if err != nil {
		return nil, err
	}
	next, err := newFilter(ctx)
	if err != nil {
		return nil, err
	}

	err = setExpireTTL(current, cfg.Freq)
	if err != nil {
		return nil, err
	}
	err = setExpireTTL(next, cfg.Freq*2)
	if err != nil {
		return nil, err
	}

	r := &Rotator{
		ctx:       ctx,
		cfg:       cfg,
		mutex:     &sync.RWMutex{},
		newFilter: newFilter,
		current:   current,
		next:      next,
	}

	go r.handleRotating(time.NewTicker(cfg.Freq))

	return r, nil
}
