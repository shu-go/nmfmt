package nmfmt

import (
	"fmt"
	"io"
	"sync"
)

type formatterOptions struct {
	cacheResetLimit int
}

type OptionFunc func(*formatterOptions)

// CacheResetLimit sets when to clear cache.
func CacheResetLimit(misses int) OptionFunc {
	return func(f *formatterOptions) {
		f.cacheResetLimit = misses
	}
}

type Formatter struct {
	cache       *cache
	argsBufPool sync.Pool
}

func New(opts ...OptionFunc) Formatter {
	fo := formatterOptions{
		cacheResetLimit: 100,
	}
	for _, o := range opts {
		o(&fo)
	}

	return Formatter{
		cache: newCache(fo.cacheResetLimit),
		argsBufPool: sync.Pool{
			New: func() any {
				return &[]any{}
			},
		},
	}
}

func (f *Formatter) Printf(format string, a ...any) (int, error) {
	if len(a) == 0 {
		return fmt.Printf(format)
	}

	cn := f.cache.get(format)

	aa, err := cn.construct(a, func() *[]any {
		return f.argsBufPool.Get().(*[]any)
	})
	if err != nil {
		return 0, err
	}

	var n int
	if aa == nil {
		n, err = fmt.Printf(cn.format)
	} else {
		n, err = fmt.Printf(cn.format, (*aa)...)
		f.argsBufPool.Put(aa)
	}

	return n, err
}

func (f *Formatter) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	if len(a) == 0 {
		return fmt.Fprintf(w, format)
	}

	cn := f.cache.get(format)

	aa, err := cn.construct(a, func() *[]any {
		return f.argsBufPool.Get().(*[]any)
	})
	if err != nil {
		return 0, err
	}

	var n int
	if aa == nil {
		n, err = fmt.Fprintf(w, cn.format)
	} else {
		n, err = fmt.Fprintf(w, cn.format, (*aa)...)
		f.argsBufPool.Put(aa)
	}

	return n, err
}

func (f *Formatter) Sprintf(format string, a ...any) string {
	if len(a) == 0 {
		return fmt.Sprintf(format)
	}

	cn := f.cache.get(format)

	aa, err := cn.construct(a, func() *[]any {
		return f.argsBufPool.Get().(*[]any)
	})
	if err != nil {
		return ""
	}

	var s string
	if aa == nil {
		s = fmt.Sprintf(cn.format)
	} else {
		s = fmt.Sprintf(cn.format, (*aa)...)
		f.argsBufPool.Put(aa)
	}

	return s
}

func (f *Formatter) Errorf(format string, a ...any) error {
	if len(a) == 0 {
		return fmt.Errorf(format)
	}

	cn := f.cache.get(format)

	aa, err := cn.construct(a, func() *[]any {
		return f.argsBufPool.Get().(*[]any)
	})
	if err != nil {
		return err
	}

	if aa == nil {
		err = fmt.Errorf(cn.format)
	} else {
		err = fmt.Errorf(cn.format, (*aa)...)
		f.argsBufPool.Put(aa)
	}

	return err
}
