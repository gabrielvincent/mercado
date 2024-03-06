package utils

import (
	"net/http"
)

func Location(header http.Header, location string) {
	header.Set("HX-Location", location)
}

func Reswap(header http.Header, swapMethod string) {
	header.Set("HX-Reswap", swapMethod)
}

func Retarget(header http.Header, selector string) {
	header.Set("HX-Retarget", selector)
}
