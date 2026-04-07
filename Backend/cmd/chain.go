package main

import (
	"net/http"
	"slices"
)

type chain []func(http.Handler) http.Handler

func (c chain) Then(h http.Handler) http.Handler {
	for _, fn := range slices.Backward(c) {
		h = fn(h)
	}
	return h
}

func (c chain) ThenFunc(h http.HandlerFunc) http.Handler {
	return c.Then(h)
}
