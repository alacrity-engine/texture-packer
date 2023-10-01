package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

func TextureFilteringByID(id string) uint32 {
	switch id {
	case "nearest":
		return gl.NEAREST

	case "linear":
		return gl.LINEAR

	default:
		return 0
	}
}
