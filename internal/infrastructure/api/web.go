package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yvv4git/task-ef/internal/config"
	"github.com/yvv4git/task-ef/internal/service"
)

// @title ETH Transaction Analysis API
// @version 1.0
// @description Service for analyzing Ethereum transactions and finding addresses with maximum balance differences
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
// @schemes http https
type Web struct {
	cfg        *config.API
	logger     *slog.Logger
	health     *service.Health
	analyzeETH *service.AnalyzeETH
}

func NewWeb(cfg *config.API, log *slog.Logger, health *service.Health, analyzeETH *service.AnalyzeETH) *Web {
	return &Web{
		cfg:        cfg,
		logger:     log,
		health:     health,
		analyzeETH: analyzeETH,
	}
}

func (w *Web) Run(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/address-eth-diff", w.AddressWithMaxBalanceDiffHandler)
	router.GET("/health", w.HealthCheckHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", w.cfg.Host, w.cfg.Port),
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			w.logger.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(w.cfg.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		w.logger.Error("HTTP server shutdown error", slog.String("error", err.Error()))
		return err
	}

	w.logger.Info("HTTP server stopped")
	return nil
}

// AddressWithMaxBalanceDiffHandler godoc
// @Summary Find address with maximum balance difference
// @Description Returns the Ethereum address with the largest balance change over the last N blocks.
// @Tags analyze
// @Accept json
// @Produce json
// @Param n query int true "Number of blocks to analyze"
// @Success 200 {object} map[string]string "address: The Ethereum address"
// @Failure 500 {object} map[string]string "error: Error message"
// @Router /address-eth-diff [get]
func (w *Web) AddressWithMaxBalanceDiffHandler(c *gin.Context) {
	address, err := w.analyzeETH.FindAddressWithMaxDiffBalance(c.Request.Context(), uint32(w.cfg.AnalyzeTxCountETH)) // TODO: set n value from config
	if err != nil {
		w.logger.Error("failed to get address with max balance diff", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

// HealthCheckHandler godoc
// @Summary Health check
// @Description Checks the health of the service.
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Service is unhealthy"
// @Router /health [get]
func (w *Web) HealthCheckHandler(c *gin.Context) {
	if err := w.health.Check(c.Request.Context()); err != nil {
		w.logger.Error("health check failed", slog.String("error", err.Error()))
		c.String(http.StatusInternalServerError, "Service is unhealthy")
		return
	}

	c.String(http.StatusOK, "OK")
}
