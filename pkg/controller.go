package pkg



import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Controller interface {
	Init(logger *zap.Logger) error
	SetupRouter(router chi.Router) error
	Terminate() error
}