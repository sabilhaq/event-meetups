package battlestrg

import (
	"context"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/Haraj-backend/hex-pokebattle/internal/shared/telemetry"
)

type Storage struct {
	data map[string]battle.Battle
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	tr := telemetry.GetTracer()
	span := tr.Trace(ctx, "GetBattle BattleStorage")
	defer span.End()

	b, ok := s.data[gameID]
	if !ok {
		return nil, nil
	}
	return &b, nil
}

func (s *Storage) SaveBattle(ctx context.Context, b battle.Battle) error {
	tr := telemetry.GetTracer()
	span := tr.Trace(ctx, "SaveBattle BattleStorage")
	defer span.End()

	s.data[b.GameID] = b
	return nil
}

func New() *Storage {
	return &Storage{data: make(map[string]battle.Battle)}
}
