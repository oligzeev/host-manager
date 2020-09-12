package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/oligzeev/host-manager/internal/domain"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	cfg        domain.ServerRestConfig
	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(cfg domain.ServerRestConfig, handlers []domain.RestHandler) *Server {
	router := gin.New()
	for _, handler := range handlers {
		handler.Register(router)
	}
	httpServer := &http.Server{
		ReadTimeout:  cfg.ReadTimeoutSec * time.Second,
		WriteTimeout: cfg.WriteTimeoutSec * time.Second,
		Addr:         cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Handler:      router,
	}
	return &Server{cfg: cfg, httpServer: httpServer, router: router}
}

func (s Server) Router() *gin.Engine {
	return s.router
}

func (s Server) Start(ctx context.Context) error {
	const op = "RestServer.Start"

	log.Tracef("%s: %s", op, s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return domain.E(op, err)
	}
	log.Tracef("%s: exit", op)
	return ctx.Err()
}

func (s Server) Stop(ctx context.Context) error {
	const op = "RestServer.Stop"

	log.Tracef("%s: in progress", op)
	timeoutCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeoutSec*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(timeoutCtx); err != nil {
		return domain.E(op, err)
	}
	log.Tracef("%s: finished", op)
	return timeoutCtx.Err()
}
