package nmfmt

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

/*
var placeholderRE = regexp.MustCompile(`{([^:{}}]+)(?::([^{}}]+))?}`)
var extract = func(format string, index []int) (string, string) {
	name := strings.TrimSpace(format[index[2]:index[3]])
	verb := ""
	if index[4] != -1 {
		verb = strings.TrimSpace(format[index[4]:index[5]])
	}
	return name, verb
}
*/

var placeholderRE = regexp.MustCompile(`(?:\$(\w+)(?::(\w+))?)|(?:\${([^:{}}]+)(?::([^{}}]+))?})`)
var extract = func(format string, index []int) (string, string) {
	if index[2] != -1 {
		name := strings.TrimSpace(format[index[2]:index[3]])
		verb := ""
		if index[4] != -1 {
			verb = strings.TrimSpace(format[index[4]:index[5]])
		}
		return name, verb
	}

	name := strings.TrimSpace(format[index[6]:index[7]])
	verb := ""
	if index[8] != -1 {
		verb = strings.TrimSpace(format[index[8]:index[9]])
	}
	return name, verb
}

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
	cache *cache
	aPool sync.Pool
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
		aPool: sync.Pool{
			New: func() any {
				return &[]any{}
			},
		},
	}
}

func (f *Formatter) Printf(format string, m map[string]any) (int, error) {
	if len(m) == 0 {
		return fmt.Printf(format)
	}

	cn := f.cache.get(format)

	aa, err := cn.construct(m, &f.aPool)
	if err != nil {
		return 0, err
	}

	var n int
	if aa == nil {
		n, err = fmt.Printf(cn.format)
	} else {
		n, err = fmt.Printf(cn.format, (*aa)...)
		f.aPool.Put(aa)
	}

	return n, err
}

func (f *Formatter) Fprintf(w io.Writer, format string, m map[string]any) (int, error) {
	if len(m) == 0 {
		return fmt.Fprintf(w, format)
	}

	cn := f.cache.get(format)

	aa, err := cn.construct(m, &f.aPool)
	if err != nil {
		return 0, err
	}

	var n int
	if aa == nil {
		n, err = fmt.Fprintf(w, cn.format)
	} else {
		n, err = fmt.Fprintf(w, cn.format, (*aa)...)
		f.aPool.Put(aa)
	}

	return n, err
}

func (f *Formatter) Sprintf(format string, m map[string]any) string {
	if len(m) == 0 {
		return fmt.Sprintf(format)
	}

	cn := f.cache.get(format)

	aa, err := cn.construct(m, &f.aPool)
	if err != nil {
		return ""
	}

	var s string
	if aa == nil {
		s = fmt.Sprintf(cn.format)
	} else {
		s = fmt.Sprintf(cn.format, (*aa)...)
		f.aPool.Put(aa)
	}

	return s
}

func (f *Formatter) Errorf(format string, m map[string]any) error {
	if len(m) == 0 {
		return fmt.Errorf(format)
	}

	cn := f.cache.get(format)

	aa, err := cn.construct(m, &f.aPool)
	if err != nil {
		return err
	}

	if aa == nil {
		err = fmt.Errorf(cn.format)
	} else {
		err = fmt.Errorf(cn.format, (*aa)...)
		f.aPool.Put(aa)
	}

	return err
}
