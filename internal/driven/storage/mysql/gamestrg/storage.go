package gamestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	db "github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	var game entity.Game
	var pokemon entity.Pokemon

	game.Partner = &pokemon

	query := `SELECT g.id, player_name, created_at, battle_won, scenario,
		p.id, p.name, p.avatar_url,
		p.max_health, p.attack, p.defense, p.speed
		FROM games g
		LEFT JOIN pokemons p on partner_id = p.id
		WHERE g.id = ?`

	if err := mappingGame(s.db.QueryRowContext(ctx, query, gameID), &game); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find game with id %s", gameID)
		}
		return nil, fmt.Errorf("unable to find game with id %s: %v", gameID, err)
	}

	return &game, nil
}

func (s *Storage) SaveGame(ctx context.Context, game entity.Game) error {
	queryGame := `
		INSERT INTO games (id, player_name, created_at, battle_won, scenario, partner_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, queryGame, game.ID, game.PlayerName, game.CreatedAt, game.BattleWon, game.Scenario, game.Partner.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func mappingGame(row db.RowResultInterface, g *entity.Game) error {
	return row.Scan(
		&g.ID, &g.PlayerName, &g.CreatedAt,
		&g.BattleWon, &g.Scenario,
		&g.Partner.ID, &g.Partner.Name, &g.Partner.AvatarURL,
		&g.Partner.BattleStats.MaxHealth,
		&g.Partner.BattleStats.Attack,
		&g.Partner.BattleStats.Defense,
		&g.Partner.BattleStats.Speed,
	)
}
