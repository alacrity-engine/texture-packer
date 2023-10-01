package main

import "gopkg.in/yaml.v2"

func ReadTexturesData(data []byte) ([]TextureMeta, error) {
	textures := make([]TextureMeta, 0)
	err := yaml.Unmarshal(data, &textures)

	if err != nil {
		return nil, err
	}

	return textures, nil
}
