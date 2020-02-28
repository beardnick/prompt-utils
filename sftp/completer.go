package main

import "regexp"

//type Completer interface {
//}

type Completer struct {
	Source  ISource
	pattern *regexp.Regexp
}

func (c Completer) Of(r string) Completer {
	var err error
	c.pattern, err = regexp.Compile(r)
	if err != nil {
		panic(err)
	}
	return c
}

func (c Completer) Match(s string) bool {
	return c.pattern.MatchString(s)
}
