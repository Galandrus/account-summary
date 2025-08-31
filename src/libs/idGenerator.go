package libs

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

type IdGenerator interface {
	Generate(prefix string) string
}

type idGenerator struct {
}

func NewIdGenerator() IdGenerator {
	return &idGenerator{}
}

func (i *idGenerator) Generate(prefix string) string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	id, err := ulid.New(ms, entropy)
	if err != nil {
		return ""
	}

	if prefix != "" {
		return fmt.Sprintf("%s-%s", prefix, id.String())
	}

	return id.String()
}
