// Composition root. The only place that wires concrete implementations
// onto the application's ports.
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lindritprekaj/user-authmanagement/internal/application/admin"
	"github.com/lindritprekaj/user-authmanagement/internal/application/auth"
	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	usercases "github.com/lindritprekaj/user-authmanagement/internal/application/user"
	domainuser "github.com/lindritprekaj/user-authmanagement/internal/domain/user"
	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/config"
	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/crypto"
	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/logging"
	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/mongodb"
	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/outbox"
	"github.com/lindritprekaj/user-authmanagement/internal/infrastructure/rabbitmq"
	httpiface "github.com/lindritprekaj/user-authmanagement/internal/interfaces/http"
	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/handlers"
)

func main() {
	if err := run(); err != nil {
		// Logger may not be set up yet; fall back to stderr.
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error("startup failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logger := logging.New(cfg.Logging.Level, cfg.Logging.Format)
	slog.SetDefault(logger)

	rootCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// --- Mongo ---------------------------------------------------------
	mongoClient, err := mongodb.Connect(rootCtx, cfg.Mongo.URI, cfg.Mongo.ConnectTimeout)
	if err != nil {
		return err
	}
	defer func() {
		dctx, c := context.WithTimeout(context.Background(), 5*time.Second)
		defer c()
		_ = mongoClient.Disconnect(dctx)
	}()
	db := mongoClient.Database(cfg.Mongo.Database)
	if err := mongodb.EnsureIndexes(rootCtx, db); err != nil {
		return err
	}

	userRepo := mongodb.NewUserRepository(db)
	tokenRepo := mongodb.NewRefreshTokenRepository(db)
	resetRepo := mongodb.NewPasswordResetRepository(db)
	outboxRepo := mongodb.NewOutboxRepository(db)
	if err := mongodb.EnsureOutboxIndexes(rootCtx, db); err != nil {
		return err
	}

	// --- RabbitMQ ------------------------------------------------------
	rabbitConn, err := rabbitmq.Dial(rootCtx, cfg.RabbitMQ.URL, cfg.RabbitMQ.Exchange,
		cfg.RabbitMQ.ReconnectBackoff, logger)
	if err != nil {
		return err
	}
	defer rabbitConn.Close()

	// Use cases see only the outbox publisher — they never touch RabbitMQ
	// directly. The Relay below drains the outbox to RabbitMQ asynchronously.
	publisher := outbox.NewPublisher(outboxRepo)
	relay := outbox.NewRelay(outboxRepo, rabbitConn, logger, outbox.Config{})
	go relay.Run(rootCtx)

	// --- Crypto / Tokens ----------------------------------------------
	hasher := crypto.NewBcryptHasher(cfg.Auth.BcryptCost)
	tokenSvc, err := crypto.NewJWTService(cfg.JWT.PrivateKey, cfg.JWT.PublicKey, cfg.JWT.Issuer, cfg.JWT.AccessTTL)
	if err != nil {
		return err
	}
	clock := ports.SystemClock{}
	ids := crypto.UUIDGenerator{}

	policy := auth.Policy{
		PasswordMinLength: cfg.Auth.PasswordMinLength,
		AccessTTL:         cfg.JWT.AccessTTL,
		RefreshTTL:        cfg.JWT.RefreshTTL,
		PasswordResetTTL:  cfg.Auth.PasswordResetTTL,
	}

	// --- Use Cases ----------------------------------------------------
	registerUC := auth.NewRegisterUseCase(userRepo, hasher, publisher, clock, ids, policy)
	loginUC := auth.NewLoginUseCase(userRepo, tokenRepo, hasher, tokenSvc, publisher, clock, ids, policy)
	refreshUC := auth.NewRefreshUseCase(userRepo, tokenRepo, tokenSvc, clock, ids, policy)
	logoutUC := auth.NewLogoutUseCase(tokenRepo, tokenSvc, clock)
	requestResetUC := auth.NewRequestPasswordResetUseCase(userRepo, resetRepo, tokenSvc, publisher, clock, ids, policy)
	confirmResetUC := auth.NewConfirmPasswordResetUseCase(userRepo, resetRepo, tokenRepo, hasher, tokenSvc, clock, policy)

	getProfileUC := usercases.NewGetProfileUseCase(userRepo)
	updateProfileUC := usercases.NewUpdateProfileUseCase(userRepo, publisher, clock, ids)
	deleteAccountUC := usercases.NewDeleteAccountUseCase(userRepo, tokenRepo, publisher, clock, ids)

	listUsersUC := admin.NewListUsersUseCase(userRepo)
	getUserUC := admin.NewGetUserUseCase(userRepo)
	setRolesUC := admin.NewSetRolesUseCase(userRepo, publisher, clock, ids)
	deleteUserUC := admin.NewDeleteUserUseCase(userRepo, tokenRepo, publisher, clock, ids)

	// --- Bootstrap admin ----------------------------------------------
	if err := seedBootstrapAdmin(rootCtx, logger, cfg, userRepo, hasher, publisher, ids, clock); err != nil {
		return err
	}

	// --- HTTP ---------------------------------------------------------
	authH := handlers.NewAuthHandler(logger, registerUC, loginUC, refreshUC, logoutUC, requestResetUC, confirmResetUC)
	userH := handlers.NewUserHandler(logger, getProfileUC, updateProfileUC, deleteAccountUC)
	adminH := handlers.NewAdminHandler(logger, listUsersUC, getUserUC, setRolesUC, deleteUserUC)
	healthH := handlers.NewHealthHandler(mongoClient, rabbitConn)

	router := httpiface.NewRouter(httpiface.RouterDeps{
		Logger:             logger,
		TokenService:       tokenSvc,
		AuthHandler:        authH,
		UserHandler:        userH,
		AdminHandler:       adminH,
		HealthHandler:      healthH,
		CORSAllowedOrigins: cfg.HTTP.CORSAllowedOrigins,
		LoginRateLimit:     cfg.Auth.LoginRateLimitPerMin,
		RegisterRateLimit:  cfg.Auth.RegisterRateLimitPerMin,
	})

	server := httpiface.NewServer(router, cfg.HTTP.Port,
		cfg.HTTP.ReadTimeout, cfg.HTTP.WriteTimeout, cfg.HTTP.ShutdownTimeout, logger)

	// Run server in a goroutine; wait for signal or fatal error.
	errCh := make(chan error, 1)
	go func() { errCh <- server.Start() }()

	select {
	case <-rootCtx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			return err
		}
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer stopCancel()
	if err := server.Stop(stopCtx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("graceful shutdown error", slog.Any("error", err))
	}
	logger.Info("server stopped")
	return nil
}

// seedBootstrapAdmin creates an admin user from ADMIN_BOOTSTRAP_EMAIL/PASSWORD
// if they are configured and no user with that email exists yet. Idempotent.
//
// Publishes user.registered through the outbox on first creation so downstream
// services (booking) can build their customer read-model. Without this, the
// admin can log in but every action that needs a customer projection fails.
func seedBootstrapAdmin(
	ctx context.Context,
	logger *slog.Logger,
	cfg *config.Config,
	users domainuser.Repository,
	hasher ports.PasswordHasher,
	publisher ports.EventPublisher,
	ids ports.IDGenerator,
	clock ports.Clock,
) error {
	if cfg.Auth.AdminBootstrapEmail == "" {
		return nil
	}
	if len(cfg.Auth.AdminBootstrapPassword) < cfg.Auth.PasswordMinLength {
		return errors.New("ADMIN_BOOTSTRAP_PASSWORD does not meet PASSWORD_MIN_LENGTH")
	}
	email := domainuser.NormalizeEmail(cfg.Auth.AdminBootstrapEmail)

	existing, err := users.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		// Already exists. Ensure they have the admin role.
		updated := false
		if !existing.HasRole(domainuser.RoleAdmin) {
			existing.Roles = append(existing.Roles, domainuser.RoleAdmin)
			existing.UpdatedAt = clock.Now()
			if err := users.Update(ctx, existing); err != nil {
				return err
			}
			updated = true
			logger.Info("bootstrap admin: granted admin role to pre-existing user", slog.String("email", email))
		}
		// Republish user.registered so booking's customer projection is healed
		// even when the admin was created before the outbox was wired up.
		// The booking-side EventLog inbox dedupes if the original event did
		// land, so this is safe to run on every startup.
		publishUserRegistered(ctx, publisher, ids, clock, existing.ID, existing.Email, logger)
		_ = updated
		return nil
	} else if err != nil && !errors.Is(err, domainuser.ErrNotFound) {
		return err
	}

	hash, err := hasher.Hash(cfg.Auth.AdminBootstrapPassword)
	if err != nil {
		return err
	}
	now := clock.Now()
	u := domainuser.New(ids.New(), email, hash, now)
	u.Roles = []domainuser.Role{domainuser.RoleUser, domainuser.RoleAdmin}
	u.Verified = true
	if err := users.Create(ctx, u); err != nil {
		return err
	}
	publishUserRegistered(ctx, publisher, ids, clock, u.ID, u.Email, logger)
	logger.Info("bootstrap admin: created", slog.String("email", email), slog.String("user_id", u.ID))
	return nil
}

func publishUserRegistered(
	ctx context.Context,
	publisher ports.EventPublisher,
	ids ports.IDGenerator,
	clock ports.Clock,
	userID, email string,
	logger *slog.Logger,
) {
	err := publisher.Publish(ctx, ports.Event{
		ID:         ids.New(),
		Type:       auth.EventUserRegistered,
		Version:    auth.EventVersion,
		OccurredAt: clock.Now().UTC().Format(time.RFC3339),
		RoutingKey: auth.EventUserRegistered,
		Data:       auth.UserRegisteredData{UserID: userID, Email: email},
	})
	if err != nil {
		logger.Error("bootstrap admin: failed to publish user.registered",
			slog.String("user_id", userID), slog.Any("error", err))
	}
}
