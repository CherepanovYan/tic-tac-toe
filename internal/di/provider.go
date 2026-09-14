package di

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"tic-tac-toe/config"
	"tic-tac-toe/internal/domain"
	"tic-tac-toe/internal/repository"
	web "tic-tac-toe/internal/web"
	"tic-tac-toe/internal/web/handler"
	"tic-tac-toe/internal/web/mappers"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"tic-tac-toe",
		provideConfig(),
		provideRepository(),
		provideDomain(),
		provideWeb(),
		provideHTTP(),
		fx.Invoke(StartServer),
	)
}

func provideConfig() fx.Option {
	return fx.Provide(
		config.Load,
	)
}

func provideRepository() fx.Option {
	return fx.Options(
		fx.Provide(
			repository.NewPostgresPool,
		),
		fx.Provide(
			fx.Annotate(
				repository.NewRepository,
				fx.As(new(domain.Repository)),
			),
		),
	)
}

func provideDomain() fx.Option {
	return fx.Options(
		fx.Provide(
			domain.NewUserService,
			func(cfg config.Config) *domain.JwtService {
				return domain.NewJwtService(cfg.JwtSecret)
			},
			domain.NewAuthService,
		),
		fx.Provide(
			fx.Annotate(
				domain.NewGameService,
				fx.As(new(handler.GameService)),
			),
		),
	)
}

func provideWeb() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				mappers.NewWebMappers,
				fx.As(new(handler.WebMappers)),
			),
		),
		fx.Provide(
			handler.NewGameHandler,
			handler.NewAuthHandler,
			handler.NewUserHandler,
		),
	)
}

func provideHTTP() fx.Option {
	return fx.Options(
		fx.Provide(
			NewHTTPHandler,
			NewServer,
		),
	)
}

func NewHTTPHandler(
	gameHandler *handler.GameHandler,
	authHandler *handler.AuthHandler,
	jwtService *domain.JwtService,
	userService *domain.UserService,
	userHandler *handler.UserHandler,
) http.Handler {
	mux := http.NewServeMux()

	web.RegisterRouter(
		mux,
		gameHandler,
		authHandler,
		userHandler,
		jwtService,
		userService,
	)

	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	return web.CORS(mux)
}

func NewServer(
	cfg config.Config,
	httpHandler http.Handler,
) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httpHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func StartServer(
	lc fx.Lifecycle,
	server *http.Server,
	cfg config.Config,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Printf(
				"Server started at http://localhost%s\n",
				server.Addr,
			)

			go func() {
				if err := server.ListenAndServe(); err != nil &&
					err != http.ErrServerClosed {
					fmt.Printf("Server error: %v\n", err)
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}
