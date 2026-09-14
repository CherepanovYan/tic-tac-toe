package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"tic-tac-toe/internal/domain/model"
	"tic-tac-toe/internal/repository/mappers"
	data "tic-tac-toe/internal/repository/model"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetGameByUUID(id uuid.UUID) (*model.Game, error) {
	if id == uuid.Nil {
		return nil, model.ErrorInvalidUUID
	}

	const query = `
		SELECT
			id,
			board,
			status,
			game_type,
			userX,
			userO,
			CurrentPlayer,
			created_at
		FROM games
		WHERE id = $1
`

	var (
		gameDTO  data.GameDTO
		boardRaw []byte
	)

	if err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&gameDTO.ID,
		&boardRaw,
		&gameDTO.StatusDTO,
		&gameDTO.GameTypeDTO,
		&gameDTO.UserX,
		&gameDTO.UserO,
		&gameDTO.CurrentPlayer,
		&gameDTO.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf(
			"scan game by uuid: %w",
			err,
		)
	}

	if err := json.Unmarshal(
		boardRaw,
		&gameDTO.BoardDTO,
	); err != nil {
		return nil, fmt.Errorf(
			"unmarshal game board: %w",
			err,
		)
	}

	return mappers.ToGame(&gameDTO), nil
}

func (r *Repository) SaveGame(game *model.Game) error {
	if game == nil {
		return model.ErrorInvalidBoard
	}
	dto := mappers.ToGameDTO(game)

	boardJSON, err := json.Marshal(dto.BoardDTO)
	if err != nil {
		return fmt.Errorf(
			"marshal game board: %w",
			err,
		)
	}

	const query = `
		INSERT INTO games (
			id,
			board,
			status,
			game_type,
			userX,
			userO,
			CurrentPlayer,
			created_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		)
		ON CONFLICT (id)
		DO UPDATE SET
			board = EXCLUDED.board,
			status = EXCLUDED.status,
			game_type = EXCLUDED.game_type,
			userX = EXCLUDED.userX,
			userO = EXCLUDED.userO,
			CurrentPlayer = EXCLUDED.CurrentPlayer,
			created_at = EXCLUDED.created_at
	`

	_, err = r.db.Exec(
		context.Background(),
		query,
		dto.ID,
		boardJSON,
		dto.StatusDTO,
		dto.GameTypeDTO,
		dto.UserX,
		dto.UserO,
		dto.CurrentPlayer,
		dto.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"save game: %w",
			err,
		)
	}

	return nil
}

func (r *Repository) GetAvailableGames() ([]*model.Game, error) {
	const query = `
		SELECT
			id,
			board,
			status,
			game_type,
			userX,
			userO,
			CurrentPlayer,
			created_at
		FROM games
		WHERE game_type = $1
		  AND status = $2
		ORDER BY id
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
		string(model.GameTypePlayer),
		string(model.StatusWaiting),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get available games: %w",
			err,
		)
	}

	defer rows.Close()

	var games []*model.Game

	for rows.Next() {
		var (
			gameDTO  data.GameDTO
			boardRaw []byte
		)

		err := rows.Scan(
			&gameDTO.ID,
			&boardRaw,
			&gameDTO.StatusDTO,
			&gameDTO.GameTypeDTO,
			&gameDTO.UserX,
			&gameDTO.UserO,
			&gameDTO.CurrentPlayer,
			&gameDTO.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan available game: %w",
				err,
			)
		}

		if err := json.Unmarshal(
			boardRaw,
			&gameDTO.BoardDTO,
		); err != nil {
			return nil, fmt.Errorf(
				"unmarshal available game board: %w",
				err,
			)
		}

		games = append(
			games,
			mappers.ToGame(&gameDTO),
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate available games: %w",
			err,
		)
	}

	return games, nil
}

func (r *Repository) GetCompletedGamesByUserUUID(
	userUUID uuid.UUID,
) ([]*model.Game, error) {
	if userUUID == uuid.Nil {
		return nil, model.ErrorInvalidUUID
	}

	const query = `
		SELECT
			id,
			board,
			status,
			game_type,
			userX,
			userO,
			CurrentPlayer,
			created_at
		FROM games
		WHERE (
			(status = $1 AND userX = $3)
			OR (status = $2 AND userO = $3)
			OR (status = $4 AND (userX = $3 OR userO = $3))
		)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
		string(model.StatusXWon),
		string(model.StatusOWon),
		userUUID,
		string(model.StatusDraw),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get completed games: %w",
			err,
		)
	}
	defer rows.Close()

	var games []*model.Game

	for rows.Next() {
		var (
			gameDTO  data.GameDTO
			boardRaw []byte
		)

		err := rows.Scan(
			&gameDTO.ID,
			&boardRaw,
			&gameDTO.StatusDTO,
			&gameDTO.GameTypeDTO,
			&gameDTO.UserX,
			&gameDTO.UserO,
			&gameDTO.CurrentPlayer,
			&gameDTO.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan completed game: %w",
				err,
			)
		}

		if err := json.Unmarshal(
			boardRaw,
			&gameDTO.BoardDTO,
		); err != nil {
			return nil, fmt.Errorf(
				"unmarshal completed game board: %w",
				err,
			)
		}

		games = append(
			games,
			mappers.ToGame(&gameDTO),
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate completed games: %w",
			err,
		)
	}

	return games, nil
}

func (r *Repository) GetLeaderboard(
	limit int,
) ([]*model.LeaderboardEntry, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive")
	}

	const query = `
		SELECT
			stats.id,
			stats.login,
			CASE
				WHEN stats.nonwins = 0 THEN 0
				ELSE stats.wins::float / stats.nonwins
			END AS win_ratio
		FROM (
			SELECT
				u.id,
				u.login,
				COALESCE(
					SUM(
						CASE
							WHEN (g.status = 'X_WON' AND g.userX = u.id)
								OR (g.status = 'O_WON' AND g.userO = u.id)
							THEN 1
							ELSE 0
						END
					), 0
				) AS wins,
				COALESCE(
					SUM(
						CASE
							WHEN (g.status = 'X_WON' AND g.userO = u.id)
								OR (g.status = 'O_WON' AND g.userX = u.id)
							THEN 1
							ELSE 0
						END
					), 0
				) + COALESCE(
					SUM(
						CASE
							WHEN g.status = 'DRAW' AND (g.userX = u.id OR g.userO = u.id)
							THEN 1
							ELSE 0
						END
					), 0
				) AS nonwins
			FROM users u
			LEFT JOIN games g ON g.userX = u.id OR g.userO = u.id
			GROUP BY u.id, u.login
		) AS stats
		ORDER BY win_ratio DESC, stats.wins DESC
		LIMIT $1
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get leaderboard: %w",
			err,
		)
	}
	defer rows.Close()

	var result []*model.LeaderboardEntry

	for rows.Next() {
		var entry model.LeaderboardEntry

		err := rows.Scan(
			&entry.UUID,
			&entry.Login,
			&entry.WinRatio,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan leaderboard entry: %w",
				err,
			)
		}

		result = append(result, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate leaderboard: %w",
			err,
		)
	}

	return result, nil
}

func (r *Repository) GetUserByLogin(
	login string,
) (*model.User, error) {
	const query = `
		SELECT
			id,
			login,
			password
		FROM users
		WHERE login = $1
	`

	var userDTO data.UserDTO

	err := r.db.QueryRow(
		context.Background(),
		query,
		login,
	).Scan(
		&userDTO.ID,
		&userDTO.Login,
		&userDTO.Password,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get user by login: %w",
			err,
		)
	}

	return mappers.ToUser(&userDTO), nil
}

func (r *Repository) CreateUser(user *model.User) error {
	if user == nil {
		return fmt.Errorf("create user: %w", model.ErrorInvalidBoard)
	}

	const query = `
		INSERT INTO users (
			id,
			login,
			password
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		user.UUID,
		user.Login,
		user.Password,
	)

	if err != nil {
		return fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	return nil
}

func (r *Repository) GetUserByUUID(
	id uuid.UUID,
) (*model.User, error) {
	if id == uuid.Nil {
		return nil, model.ErrorInvalidUUID
	}

	const query = `
		SELECT
			id,
			login,
			password
		FROM users
		WHERE id = $1
	`

	var userDTO data.UserDTO

	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&userDTO.ID,
		&userDTO.Login,
		&userDTO.Password,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get user by uuid: %w",
			err,
		)
	}

	return mappers.ToUser(&userDTO), nil
}
