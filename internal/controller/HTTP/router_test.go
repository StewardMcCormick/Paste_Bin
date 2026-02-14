package HTTP

import (
	"github.com/StewardMcCormick/Paste_Bin/internal/controller/HTTP/handlers/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MiddlewareTestSuite struct {
	suite.Suite
	handler *mocks.MockHandler
	router  http.Handler
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.handler = mocks.NewMockHandler(s.T())
	s.router = Router(s.handler, zap.L())
}

func (s *MiddlewareTestSuite) TestRouter_OnPanic() {
	s.handler.EXPECT().
		HelloHandler(mock.Anything, mock.Anything).
		Run(func(w http.ResponseWriter, r *http.Request) {
			panic("foo")
		}).Once()
	r := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusInternalServerError, w.Result().StatusCode)
}
