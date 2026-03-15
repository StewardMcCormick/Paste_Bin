package test

import (
	"testing"

	appcache "github.com/StewardMcCormick/Paste_Bin/internal/repository/cache"
	"github.com/StewardMcCormick/Paste_Bin/internal/repository/paste"
	"github.com/stretchr/testify/suite"
)

type PasteRepoIntTestSuite struct {
	suite.Suite
	repo *paste.Repository
}

func TestPasteRepoInt(t *testing.T) {
	suite.Run(t, new(PasteRepoIntTestSuite))
}

func (s *PasteRepoIntTestSuite) SetupSuite() {
	pasteCache := appcache.NewPasteCache(pasteCacheRedisClient)
	s.repo = paste.NewRepository(pool, pasteCache)
}
