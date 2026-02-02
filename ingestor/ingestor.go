package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func matchValid(md *MatchMessage) error {
	msg := []string{}
	if md.SeriesID == nil || *md.SeriesID == 0 {
		msg = append(msg, "invalid series_id")
	}
	if md.RadiantTeamID == 0 {
		msg = append(msg, "invalid radiant_team_id")
	}
	if md.DireTeamID == 0 {
		msg = append(msg, "invalid dire_team_id")
	}
	if len(msg) > 0 {
		return fmt.Errorf("invalid match. errors: {%v}", msg)
	}
	return nil
}

// processMessage unmarshals the SQS message body and orchestrates database insertions.
func processMessage(ctx context.Context, pool *pgxpool.Pool, messageBody string) (err error) {
	matchData, err := parseMatchMessage(messageBody)
	if err != nil {
		return err
	}
	if err := matchValid(&matchData); err != nil {
		return err
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		if commitErr := tx.Commit(ctx); commitErr != nil {
			err = fmt.Errorf("failed to commit transaction: %w", commitErr)
		}
	}()
	r := pool.QueryRow(ctx, `SELECT count(*) FROM matches WHERE match_id = $1`, matchData.MatchID)
	var exists int
	if err := r.Scan(&exists); err != nil {
		return err
	}
	if exists == 1 {
		log.Info().
			Int64("match_id", matchData.MatchID).
			Msg("Match already exists, skipping.")
		return nil
	}

	if err = upsertReferenceData(ctx, tx, &matchData); err != nil {
		return err
	}

	radiantHeroes, direHeroes, radiantPlayers, direPlayers := extractIdsForMatch(&matchData)

	if err = insertMatch(ctx, tx, &matchData, radiantHeroes, direHeroes, radiantPlayers, direPlayers); err != nil {
		return fmt.Errorf("failed to insert match: %w", err)
	}

	if err = upsertPlayers(ctx, tx, &matchData); err != nil {
		return err
	}

	if err = insertMatchMetadata(ctx, tx, &matchData); err != nil {
		return fmt.Errorf("failed to insert match metadata: %w", err)
	}

	return nil
}

func parseMatchMessage(messageBody string) (MatchMessage, error) {
	var matchData MatchMessage
	if err := json.Unmarshal([]byte(messageBody), &matchData); err != nil {
		return MatchMessage{}, fmt.Errorf("failed to unmarshal match: %w", err)
	}
	return matchData, nil
}

func upsertReferenceData(ctx context.Context, tx pgx.Tx, m *MatchMessage) error {
	if err := insertLeague(ctx, tx, m.LeagueID, m.LeagueName, m.LeagueTier); err != nil {
		return fmt.Errorf("failed to insert league: %w", err)
	}

	if err := insertTeam(ctx, tx, m.RadiantTeamID, m.RadiantTeamName, m.RadiantTeamTag, m.RadiantTeamLogoURL); err != nil {
		return fmt.Errorf("failed to insert radiant team: %w", err)
	}

	if err := insertTeam(ctx, tx, m.DireTeamID, m.DireTeamName, m.DireTeamTag, m.DireTeamLogoURL); err != nil {
		return fmt.Errorf("failed to insert dire team: %w", err)
	}

	return nil
}

func extractIdsForMatch(m *MatchMessage) (radiantHeroes, direHeroes, radiantPlayers, direPlayers []int64) {
	for _, p := range m.Players {
		if p.IsRadiant {
			radiantHeroes = append(radiantHeroes, int64(p.HeroID))
			if p.AccountID != nil {
				radiantPlayers = append(radiantPlayers, *p.AccountID)
			}
		} else {
			direHeroes = append(direHeroes, int64(p.HeroID))
			if p.AccountID != nil {
				direPlayers = append(direPlayers, *p.AccountID)
			}
		}
	}
	return
}

func upsertPlayers(ctx context.Context, tx pgx.Tx, m *MatchMessage) error {
	if m.RadiantCaptain != nil {
		if err := insertPlayer(ctx, tx, *m.RadiantCaptain); err != nil {
			log.Warn().
				Err(err).
				Int64("account_id", *m.RadiantCaptain).
				Msg("Failed to insert radiant captain.")
		}
	}
	if m.DireCaptain != nil {
		if err := insertPlayer(ctx, tx, *m.DireCaptain); err != nil {
			log.Warn().
				Err(err).
				Int64("account_id", *m.DireCaptain).
				Msg("Failed to insert dire captain.")
		}
	}
	for _, player := range m.Players {
		if player.AccountID != nil {
			if err := insertPlayerWithDetails(ctx, tx, *player.AccountID, player.Name); err != nil {
				log.Warn().
					Err(err).
					Int64("account_id", *player.AccountID).
					Msg("Failed to insert player details.")
			}
		}
	}
	return nil
}

// insertPlayer inserts a player into the players table if they don't already exist.
func insertPlayer(ctx context.Context, tx pgx.Tx, accountID int64) error {

	query := `
		INSERT INTO players (account_id)
		VALUES ($1)
		ON CONFLICT (account_id) DO NOTHING;
	`
	_, err := tx.Exec(ctx, query, accountID)
	return err
}

// insertPlayerWithDetails inserts a player with their details into the players table.
func insertPlayerWithDetails(ctx context.Context, tx pgx.Tx, accountID int64, name *string) error {
	query := `
		INSERT INTO players (account_id, name)
		VALUES ($1, $2)
		ON CONFLICT (account_id) DO UPDATE SET
			name = COALESCE(EXCLUDED.name, players.name);
	`
	_, err := tx.Exec(ctx, query, accountID, name)
	return err
}

// insertMatchMetadata inserts match metadata into the matches_metadata table.
func insertMatchMetadata(ctx context.Context, tx pgx.Tx, match *MatchMessage) error {
	// Marshal complex fields to JSON for JSONB columns
	playersJSON, err := json.Marshal(match.Players)
	if err != nil {
		return fmt.Errorf("failed to marshal players to JSON: %w", err)
	}
	picksBansJSON, err := json.Marshal(match.PicksBans)
	if err != nil {
		return fmt.Errorf("failed to marshal picks_bans to JSON: %w", err)
	}

	query := `
        INSERT INTO matches_metadata (
            match_id, series_id, radiant_captain, dire_captain,
		picks_bans, players_data, radiant_gold_adv, radiant_xp_adv, radiant_score, dire_score
        ) VALUES (
            $1, $2, $3, $4,
            $5, $6, $7, $8, $9, $10
        )
        ON CONFLICT (match_id) DO UPDATE SET
            series_id = EXCLUDED.series_id,
            radiant_captain = EXCLUDED.radiant_captain,
            dire_captain = EXCLUDED.dire_captain,
            picks_bans = EXCLUDED.picks_bans,
            players_data = EXCLUDED.players_data,
            radiant_gold_adv = EXCLUDED.radiant_gold_adv,
            radiant_xp_adv = EXCLUDED.radiant_xp_adv;
    `

	_, err = tx.Exec(ctx, query,
		match.MatchID,
		match.SeriesID,
		match.RadiantCaptain,
		match.DireCaptain,
		picksBansJSON,
		playersJSON,
		toInt32Slice(match.RadiantGoldAdv),
		toInt32Slice(match.RadiantXPAdv),
		match.RadiantScore,
		match.DireScore,
	)
	return err
}

// insertMatch inserts a match into the matches table.
func insertMatch(ctx context.Context, tx pgx.Tx, match *MatchMessage, radiantHeroes, direHeroes, radiantPlayers, direPlayers []int64) error {
	query := `
		INSERT INTO matches (
			match_id, league_id, series_id, radiant_team_id,
			dire_team_id, radiant_heroes, dire_heroes, radiant_players, dire_players,
			duration, start_time, patch, version, radiant_win
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14
		) ON CONFLICT (match_id) DO NOTHING;
	`

	_, err := tx.Exec(ctx, query,
		match.MatchID, match.LeagueID, match.SeriesID,
		match.RadiantTeamID,
		match.DireTeamID,
		radiantHeroes, direHeroes, radiantPlayers, direPlayers,
		match.Duration, time.Unix(match.StartTime, 0), match.Patch, match.Version, match.RadiantWin,
	)
	return err
}

// toInt32Slice converts a []int (from JSON) to []int32 for INTEGER[] columns.
func toInt32Slice(input []int) []int32 {
	if input == nil {
		return nil
	}
	out := make([]int32, len(input))
	for i, v := range input {
		out[i] = int32(v)
	}
	return out
}

func insertLeague(ctx context.Context, tx pgx.Tx, leagueID int, name, tier string) error {
	query := `
		INSERT INTO leagues (league_id, name, tier)
		VALUES ($1, $2, $3)
		ON CONFLICT (league_id) DO UPDATE SET
			name = EXCLUDED.name,
			tier = EXCLUDED.tier;
	`
	_, err := tx.Exec(ctx, query, leagueID, name, tier)
	return err
}

func insertTeam(ctx context.Context, tx pgx.Tx, teamID int64, name string, tag, logoURL *string) error {
	query := `
		INSERT INTO teams (team_id, name, tag, logo_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (team_id) DO UPDATE SET
			name = EXCLUDED.name,
			tag = EXCLUDED.tag,
			logo_url = EXCLUDED.logo_url;
	`
	_, err := tx.Exec(ctx, query, teamID, name, tag, logoURL)
	return err
}
