package utils

import (
	"time"

	"github.com/goombaio/namegenerator"
)

func GenerateNewName() string {
	seed := time.Now().UTC().UnixNano()
	generator := namegenerator.NewNameGenerator(seed)
	return generator.Generate()
}
