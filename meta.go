package main

import codec "github.com/alacrity-engine/resource-codec"

type TextureMeta struct {
	Name             string `yaml:"name"`
	PictureID        string `yaml:"pictureID"`
	TextureFiltering string `yaml:"textureFiltering"`
}

func (meta TextureMeta) ToTextureData() codec.TextureData {
	return codec.TextureData{
		PictureID: meta.PictureID,
		Filtering: TextureFilteringByID(meta.TextureFiltering),
	}
}
