package gosync

import "sync"

type CallGroup struct {
	wg  sync.WaitGroup
	mx  sync.Mutex
	err error
}

func (c *CallGroup) Add(fn func() error) *CallGroup {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		if err := fn(); err != nil {
			c.mx.Lock()
			defer c.mx.Unlock()
			if c.err == nil {
				c.err = err
			}
		}
	}()
	return c
}

func (c *CallGroup) Wait() error {
	c.wg.Wait()
	return c.err
}

func Call(fn ...func() error) error {
	var c CallGroup
	for _, f := range fn {
		c.Add(f)
	}
	return c.Wait()
}

func Do(fn ...func()) {
	var wg sync.WaitGroup
	for _, f := range fn {
		wg.Add(1)
		go func(fn func()) {
			defer wg.Done()
			fn()
		}(f)
	}
	wg.Wait()
}
