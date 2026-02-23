package regular_expression

import (
	"regexp"
	"sync"
)

type LazyRegexp struct {
	pattern string

	once sync.Once
	re   *regexp.Regexp
	err  error
}

func NewLazyRegexp(pattern string) *LazyRegexp {
	return &LazyRegexp{pattern: pattern}
}

func (lr *LazyRegexp) Get() (*regexp.Regexp, error) {
	lr.once.Do(func() {
		lr.re, lr.err = regexp.Compile(lr.pattern)
	})
	return lr.re, lr.err
}
