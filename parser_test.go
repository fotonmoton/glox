package main

import "testing"

func TestSimpleParser(t *testing.T) {
	s := newScanner([]byte("print 1;"))
	tokens, _ := s.scan()
	p, _ := newParser(tokens).parse()

	if p == nil {
		t.Fatal("cant parse")
	}
}
