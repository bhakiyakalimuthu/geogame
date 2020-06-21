package pkg

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/voi-oss/svc"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
)


var _ svc.Worker = (*ChiWorker)(nil)

type ChiWorker struct {
	port       int
	logger     *zap.Logger
	router     chi.Router
	controller Controller
	server     *http.Server
}

func NewWorker(ctrl Controller)  *ChiWorker {
	return &ChiWorker{
		port:       8080,
		controller: ctrl,
	}
}
func (c *ChiWorker) Init(logger *zap.Logger) error {
	c.logger = logger
	c.router = chi.NewRouter()
	if err := c.controller.Init(c.logger); err !=nil {
		logger.Error("failed to init controller")
		return err
	}
	c.server = &http.Server{
		Addr: net.JoinHostPort("", strconv.Itoa(c.port)),
		Handler:     &ochttp.Handler{
			Handler:   c.router,
		}      ,
	}
	if err:=c.controller.Init(c.logger);err !=nil {
		logger.Error("failed to init controller")
		return err
	}
	r := chi.NewRouter()
	if err:= c.controller.SetupRouter(r);err !=nil {
		logger.Error("failed to setup router")
		return err
	}
	c.router.Mount("/v1", r)
	return nil
}

func (c *ChiWorker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.router.ServeHTTP(rw, req)
}

func (c *ChiWorker) Run() error {
	c.logger.Info("HTTP server running",zap.String("host name",c.server.Addr))
	if err := c.server.ListenAndServe();err!=nil && err != http.ErrServerClosed {
		c.logger.Error("HTTP server failed to RUN")
		return err
	}
	return nil
}

func (c *ChiWorker) Terminate() error {
	if err := c.server.Shutdown(context.Background());err != nil {
		c.logger.Warn("server shutdown failed")
	}
	return c.controller.Terminate()
}

func (c *ChiWorker) Healthy() error {
	if h, ok := c.controller.(svc.Healther);ok {
		h.Healthy()
	}
	return nil
}
