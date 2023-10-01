package main

import (
	"flag"
	"fmt"
	"os"

	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

var (
	texturesMetaPath string
	resourceFilePath string
)

func parseFlags() {
	flag.StringVar(&texturesMetaPath, "textures-meta",
		"./textures-meta.yml", "A path to the file with textures metadata")
	flag.StringVar(&resourceFilePath, "out", "./stage.res",
		"Resource file to store animations and spritesheets.")

	flag.Parse()
}

func main() {
	parseFlags()

	data, err := os.ReadFile(texturesMetaPath)
	handleError(err)
	textureMetas, err := ReadTexturesData(data)
	handleError(err)
	textureDatas := make([]codec.TextureData,
		0, len(textureMetas))

	for i := 0; i < len(textureMetas); i++ {
		textureMeta := textureMetas[i]
		textureDatas = append(textureDatas,
			textureMeta.ToTextureData())
	}

	resourceFile, err := bolt.Open(resourceFilePath, 0666, nil)
	handleError(err)
	defer resourceFile.Close()

	for i := 0; i < len(textureDatas); i++ {
		textureMeta := textureMetas[i]
		textureData := textureDatas[i]

		err = resourceFile.Update(func(tx *bolt.Tx) error {
			buck := tx.Bucket([]byte("pictures"))

			if buck == nil {
				return fmt.Errorf("no pictures bucket present")
			}

			picData := buck.Get([]byte(textureData.PictureID))

			if picData == nil {
				return fmt.Errorf(
					"the '%s' picture is absent", textureData.PictureID)
			}

			buck = tx.Bucket([]byte("textures"))

			if buck == nil {
				return fmt.Errorf("no textures bucket present")
			}

			textureBytes, err := textureData.ToBytes()

			if err != nil {
				return err
			}

			err = buck.Put([]byte(textureMeta.Name), textureBytes)

			if err != nil {
				return err
			}

			return nil
		})
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
