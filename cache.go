package nmfmt

import (
	"sync"
)

type cache struct {
	m     sync.Mutex
	nodes map[string]cachenode

	cacheResetLimit int
	cachemisses     int
}

func newCache(refreshRate int) *cache {
	return &cache{
		nodes:           make(map[string]cachenode),
		cacheResetLimit: refreshRate,
	}
}

func (c *cache) get(format string, a map[string]any) cachenode {
	c.m.Lock()
	if cn, found := c.nodes[format]; found {
		c.m.Unlock()
		return cn
	}

	cn := newCacheNode(format, a)
	c.cachemisses++
	if c.cachemisses >= c.cacheResetLimit {
		c.cachemisses = 0
		c.nodes = make(map[string]cachenode)
	}
	c.nodes[format] = cn
	c.m.Unlock()
	return cn
}

type cachenode struct {
	format string
	aorder []string
}

func newCacheNode(format string, a map[string]any) cachenode {
	indices := placeholderRE.FindAllStringSubmatchIndex(format, -1)
	if len(indices) == 0 {
		return cachenode{format: format}
	}

	var cformat string
	var caorder []string

	last := 0
	for i := range indices {
		cformat += format[last:indices[i][0]]

		index := indices[i]

		name, verb := extract(format, index)
		if verb == "" { // not found
			cformat += "%v"
		} else {
			cformat += "%" + verb
		}
		caorder = append(caorder, name)

		last = index[1]
	}
	cformat += format[last:]

	return cachenode{
		format: cformat,
		aorder: caorder,
	}
}

func (c cachenode) construct(a map[string]any, aPool *sync.Pool) (*[]any, error) {
	if len(c.aorder) == 0 {
		return nil, nil
	}

	//aa := make([]any, 0, len(c.aorder))
	aa := aPool.Get().(*[]any)
	*aa = (*aa)[:0]

	for _, name := range c.aorder {
		if v, found := a[name]; found {
			*aa = append(*aa, v)
		} else {
			*aa = append(*aa, nil)
		}
	}

	return aa, nil
}
