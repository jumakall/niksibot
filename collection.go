package main

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

// SoundCollection contains sounds and the commands associated with the collection
type SoundCollection struct {
	Prefix   string
	Commands []string
	Sounds   []*Sound

	soundRange int
}

// Create a collection from each directory inside the given path
func discoverSounds(path string) []*SoundCollection {
	collections := []*SoundCollection{}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() != path {
			sc := createCollection(info.Name(), path)
			if sc != nil {
				collections = append(collections, sc)
				SoundCount += len(sc.Sounds)
			}

			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"count": len(collections),
	}).Info("Collections discovered")
	return collections
}

// Create a collection from the given path
func createCollection(name string, path string) *SoundCollection {
	sc := SoundCollection{
		Prefix: name,
		Commands: []string{
			"!" + name,
		},
		Sounds: []*Sound{},
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() != name {
			return filepath.SkipDir
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".dca" {
			name := info.Name()
			extension := filepath.Ext(name)

			sc.Sounds = append(sc.Sounds, createSound(name[0:len(name)-len(extension)], 1, 100, &sc))
		}

		return nil
	})

	if err != nil || len(sc.Sounds) <= 0 {
		return nil
	}

	log.WithFields(log.Fields{
		"name":   name,
		"length": len(sc.Sounds),
		"path":   path,
	}).Debug("Collection created")
	return &sc
}
