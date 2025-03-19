package main

import "testing"

func TestReplaceBadWords(t *testing.T) {
	input := "Hi you are such a kerfuffle Sharbert its crazy. You little Fornax!"
	res := replaceBadWords(input)
	expected := "Hi you are such a **** **** its crazy. You little Fornax!"
	if res != expected {
		t.Errorf("got %s, wanted %s", res, expected)
	}
}
