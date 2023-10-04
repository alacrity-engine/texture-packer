package main

func TextureFilteringByID(id string) uint32 {
	switch id {
	case "nearest":
		return 0x2600

	case "linear":
		return 0x2601

	default:
		return 0
	}
}
