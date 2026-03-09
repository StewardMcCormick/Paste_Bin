package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StewardMcCormick/Paste_Bin/config"
	"github.com/StewardMcCormick/Paste_Bin/internal/adapter/postgres"
	"github.com/StewardMcCormick/Paste_Bin/internal/adapter/redis"
	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	"github.com/StewardMcCormick/Paste_Bin/internal/handler"
	"github.com/StewardMcCormick/Paste_Bin/internal/handler/middleware"
	pasteH "github.com/StewardMcCormick/Paste_Bin/internal/handler/paste"
	userH "github.com/StewardMcCormick/Paste_Bin/internal/handler/user"
	"github.com/StewardMcCormick/Paste_Bin/internal/repository"
	appcache "github.com/StewardMcCormick/Paste_Bin/internal/repository/cache"
	"github.com/StewardMcCormick/Paste_Bin/internal/repository/paste"
	userUseCase "github.com/StewardMcCormick/Paste_Bin/internal/usecase/auth"
	pasteUseCase "github.com/StewardMcCormick/Paste_Bin/internal/usecase/paste"
	"github.com/StewardMcCormick/Paste_Bin/internal/util/security"
	"github.com/StewardMcCormick/Paste_Bin/internal/util/validation"
	views "github.com/StewardMcCormick/Paste_Bin/internal/util/views_worker"
	"github.com/StewardMcCormick/Paste_Bin/pkg/httpserver"
	"github.com/StewardMcCormick/Paste_Bin/pkg/logging"
	"github.com/StewardMcCormick/Paste_Bin/pkg/migrations"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	AppRun(ctx, cfg)

	cancel()
}

func AppRun(ctx context.Context, cfg *config.Config) {
	log, err := logging.NewLogger(cfg.Logger, cfg.App.Env, cfg.App.Name, cfg.App.Version)
	if err != nil {
		panic(err)
	}
	log.Info("[START] Logger initialization completed")

	log.Info("[START] PGX pool initialization...")
	pool, err := postgres.NewPool(ctx, &cfg.Postgres)
	if err != nil {
		panic(err)
	}
	log.Info("[START] PGX initialization completed")

	log.Info("[START] DataBase migrations executing...")
	err = migrations.Exec(cfg.Postgres.DbUrl, cfg.Postgres.MigrationsPath)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("[START] Migrations - nothing to change")
		} else {
			panic(err)
		}
	}
	log.Info("[START] DataBase migrations executing completed")

	log.Info("[START] Redis client initialization...")
	redisManager, err := redis.NewManager(cfg.Redis)
	if err != nil {
		panic(err)
	}
	log.Info("[START] Redis client initialization completed")

	pasteCache := appcache.NewPasteCache(redisManager.GetPasteCacheClient())
	apiKeyCache := appcache.NewAPIKeyCache(redisManager.GetAPIKeyCacheClient())

	uowFactory := repository.NewUWFactory(pool, apiKeyCache)
	pasteRepo := paste.NewRepository(pool, pasteCache)
	securityUtil := security.NewUtil()

	userValid := validation.NewValidator[*dto.UserRequest](validator.New(validator.WithRequiredStructEnabled()))
	pasteValid := validation.NewValidator[*dto.PasteRequest](validator.New(validator.WithRequiredStructEnabled()))

	viewWorker := views.NewViewsWorker(pool, 10, 10*time.Millisecond)
	viewWorker.Start(ctx)
	log.Info("[START] View Worker started")

	authUc := userUseCase.NewUseCase(uowFactory, securityUtil, userValid, cfg.Auth)
	pasteUc := pasteUseCase.NewUseCase(cfg.Paste, pasteRepo, pasteValid, securityUtil, viewWorker)

	log.Info("[START] Server initialization...")

	logMid := middleware.NewLogging(log)
	recoverMid := middleware.NewRecoverer()
	envMid := middleware.NewEnv(cfg.App.Env)
	validMid := middleware.NewJSONValidation()
	authMid := middleware.NewAuth(authUc)

	userHandler := userH.NewHandler(authUc)
	pasteHandler := pasteH.NewHandlers(pasteUc)
	router := handler.NewRouter(
		userHandler,
		pasteHandler,
		logMid,
		recoverMid,
		envMid,
		validMid,
		authMid,
	)
	server := httpserver.New(router, &cfg.Server)

	go func() {
		log.Info(fmt.Sprintf("[START] Server starts on %s:%s", cfg.Server.Host, cfg.Server.Port))
		err = server.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	log.Info("[SHUTDOWN] Start shutting down...")

	viewWorker.Close(ctx)
	log.Info("[SHUTDOWN] View Worker closed")

	err = redisManager.Close()
	if err != nil {
		log.Error(fmt.Sprintf("[SHUTDOWN] Redis client close error - %v", err))
	} else {
		log.Info("[SHUTDOWN] Redis client closed")
	}

	pool.Close()
	log.Info("[SHUTDOWN] PGX close completed")

	err = server.Close()
	if err != nil {
		log.Panic(fmt.Sprintf("[SHUTDOWN] Server closing error: %v", err))
	}
	log.Info("[SHUTDOWN] Server close completed")

	err = log.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EINVAL) && !errors.Is(err, syscall.EBADF) {
		log.Panic(fmt.Sprintf("[SHUTDOWN] Log sync error: %v", err))
	}
	log.Info("[SHUTDOWN] Logger sync completed")

	log.Info("[SHUTDOWN] Shutdown completed")
}
