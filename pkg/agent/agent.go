package agent

import (
	"fmt"

	"svrn/internal/config"
	"svrn/internal/logging"

	"go.uber.org/zap"
)

type Agent struct {
	log *zap.SugaredLogger
	cfg *config.Config
}

func New(cfg *config.Config) (*Agent, error) {
	if len(cfg.Roles) == 0 {
		// Default consumer-only if not specified
		cfg.Roles = []string{"consumer"}
	}

	log := logging.New()
	log.Infow("agent init", "roles", cfg.Roles, "services", cfg.Services, "router", cfg.Router)

	return &Agent{log: log, cfg: cfg}, nil
}

// Start launches all configured subsystems.
// Phase 1: only logs startup; future phases will init router, DHT, services, etc.
func (a *Agent) Start() error {
	a.log.Infow("svrn starting", "roles", a.cfg.Roles, "services", a.cfg.Services)

	// TODO: Phase 2 — init I2P router
	// TODO: Phase 3 — init DHT
	// TODO: Phase 4 — init Blob/CRDT services per RFC

	a.log.Info("svrn initialized (stub mode)")
	return nil
}

func (a *Agent) Stop() error {
	a.log.Info("svrn stopping")
	// TODO: graceful shutdown of router, services, etc.
	return nil
}

func (a *Agent) DebugString() string {
	return fmt.Sprintf("roles=%v services=%v router=%s community=%s", a.cfg.Roles, a.cfg.Services, a.cfg.Router, a.cfg.Community)
}
