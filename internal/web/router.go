package handler

import (
	"net/http"
	"tic-tac-toe/internal/domain"
)

func RegisterRouter(
	mux *http.ServeMux,
	gameHandler gameHandler,
	authHandler authHandler,
	userHandler userHandler,
	jwtService *domain.JwtService,
	userService *domain.UserService,
) {
	mux.HandleFunc(
		"POST /auth/register",
		authHandler.Register,
	)

	mux.HandleFunc(
		"POST /auth/login",
		authHandler.Login,
	)

	mux.HandleFunc(
		"POST /refreshAccessToken",
		authHandler.RefreshAccessToken,
	)

	mux.HandleFunc(
		"POST /refreshRefreshToken",
		authHandler.RefreshRefreshToken,
	)

	authMiddleware := AuthMiddleware(jwtService, userService)

	mux.Handle(
		"GET /game/new",
		authMiddleware(
			http.HandlerFunc(
				gameHandler.CreateNewGame,
			),
		),
	)

	mux.Handle(
		"GET /game/info/{UUID}",
		authMiddleware(
			http.HandlerFunc(
				gameHandler.GameInfo,
			),
		),
	)

	mux.Handle(
		"POST /game/join/{UUID}",
		authMiddleware(
			http.HandlerFunc(
				gameHandler.JoinGame,
			),
		),
	)

	mux.Handle(
		"POST /game/move/{UUID}",
		authMiddleware(
			http.HandlerFunc(
				gameHandler.MakeMove,
			),
		),
	)

	mux.Handle(
		"GET /game/available",
		authMiddleware(
			http.HandlerFunc(
				gameHandler.GetAvailableGames,
			),
		),
	)

	mux.Handle(
		"GET /game/history",
		authMiddleware(
			http.HandlerFunc(
				gameHandler.GetCompletedGames,
			),
		),
	)

	mux.Handle(
		"GET /user/leaderboard",
		authMiddleware(
			http.HandlerFunc(
				gameHandler.GetLeaderboard,
			),
		),
	)

	mux.Handle(
		"GET /user/{UUID}",
		authMiddleware(
			http.HandlerFunc(
				userHandler.GetUser,
			),
		),
	)

	mux.Handle(
		"GET /user/CurrentUser",
		authMiddleware(
			http.HandlerFunc(
				userHandler.GetCurrentUser,
			),
		),
	)

}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Access-Control-Allow-Origin",
			"*",
		)

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, OPTIONS",
		)

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
