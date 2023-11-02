package nmfmt

import (
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

var placeholderRE = regexp.MustCompile(`(?:\$(=?\w+)(?::(\w+))?)|(?:\${(=?[^:{}}]+)(?::([^{}}]+))?})`)
var extract = func(format string, index []int) (string, string, bool) {
	var eq bool

	if index[2] != -1 {
		name := strings.TrimSpace(format[index[2]:index[3]])
		if strings.HasPrefix(name, "=") {
			name = name[1:]
			eq = true
		}
		verb := ""
		if index[4] != -1 {
			verb = strings.TrimSpace(format[index[4]:index[5]])
		}
		return name, verb, eq
	}

	name := strings.TrimSpace(format[index[6]:index[7]])
	if strings.HasPrefix(name, "=") {
		name = name[1:]
		eq = true
	}
	verb := ""
	if index[8] != -1 {
		verb = strings.TrimSpace(format[index[8]:index[9]])
	}
	return name, verb, eq
}

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

func (c *cache) get(format string) (cn cachenode) {
	c.m.Lock()

	var found bool
	if cn, found = c.nodes[format]; found {
		c.m.Unlock()
		return
	}

	cn = newCacheNode(format)
	c.cachemisses++
	if c.cachemisses >= c.cacheResetLimit {
		c.cachemisses = 0
		c.nodes = make(map[string]cachenode)
	}
	c.nodes[format] = cn
	c.m.Unlock()
	return
}

type cachenode struct {
	format string
	aorder []string
}

func newCacheNode(format string) cachenode {
	indices := placeholderRE.FindAllStringSubmatchIndex(format, -1)
	if len(indices) == 0 {
		return cachenode{format: format}
	}

	var cformat string
	var caorder []string

	last := 0
	for i := 0; i < len(indices); i++ {
		cformat += format[last:indices[i][0]]

		index := indices[i]

		name, verb, eq := extract(format, index)
		if eq {
			cformat += name + "="
		}
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

	for i := 0; i < len(c.aorder); i++ {
		*aa = append(*aa, a[c.aorder[i]])
	}

	return aa, nil
}
