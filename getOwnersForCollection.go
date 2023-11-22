package main

import (
	"fmt"
	// "github.com/schollz/progressbar/v3"
)

const API_KEY = "MKQFVASMxvydyTMLjfvYQoAyhiqTdcZL"
const CURRENT_NT = "0xB9951B43802dCF3ef5b14567cb17adF367ed1c0F"
const LEGACY_NT = "0xb668beB1Fa440F6cF2Da0399f8C28caB993Bdd65"
const BURN = "0xB9951B43802dCF3ef5b14567cb17adF367ed1c0F"

func difference(a, b []token) []token {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x.TokenID] = struct{}{}
	}
	var diff []token
	for _, x := range a {
		if _, found := mb[x.TokenID]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func removeTheNameless() {

}

func main() {
	//Requests
	getBurnedNfts(true)
	getLegacyNftsInCollection("legacyNT.json", 2082, true)
	//Reading the files
	legacyTokens := readJson("legacyNT.json")
	burnedTokens := readJson("burnedNT.json")
	//diffing
	fmt.Println(len(difference(burnedTokens, legacyTokens)))

}
