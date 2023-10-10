package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	codec "github.com/alacrity-engine/resource-codec"
	"github.com/golang-collections/collections/queue"
	bolt "go.etcd.io/bbolt"
)

var (
	projectPath      string
	resourceFilePath string
)

func parseFlags() {
	flag.StringVar(&projectPath, "project", ".",
		"Path to the project to pack spritesheets for.")
	flag.StringVar(&resourceFilePath, "out", "./stage.res",
		"Resource file to store animations and spritesheets.")

	flag.Parse()
}

func main() {
	parseFlags()

	resourceFile, err := bolt.Open(resourceFilePath, 0666, nil)
	handleError(err)
	defer resourceFile.Close()

	entries, err := os.ReadDir(projectPath)
	handleError(err)

	traverseQueue := queue.New()

	if len(entries) <= 0 {
		return
	}

	for _, entry := range entries {
		traverseQueue.Enqueue(FileTracker{
			EntryPath: projectPath,
			Entry:     entry,
		})
	}

	for traverseQueue.Len() > 0 {
		fsEntry := traverseQueue.Dequeue().(FileTracker)

		if fsEntry.Entry.IsDir() {
			entries, err = os.ReadDir(path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()))
			handleError(err)

			for _, entry := range entries {
				traverseQueue.Enqueue(FileTracker{
					EntryPath: path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()),
					Entry:     entry,
				})
			}

			continue
		}

		if !strings.HasSuffix(fsEntry.Entry.Name(), ".texture.yml") {
			continue
		}

		data, err := os.ReadFile(path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()))
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

				buck, err = tx.CreateBucketIfNotExists([]byte("textures"))

				if err != nil {
					return err
				}

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
			handleError(err)
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
