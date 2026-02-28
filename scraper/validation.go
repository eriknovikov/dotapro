package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

// Validator provides validation functions for scraper data
type Validator struct {
	// Configuration for validation rules
	MinMatchDuration int    // Minimum match duration in seconds
	MaxMatchDuration int    // Maximum match duration in seconds
	ValidTiers       []string // Valid league tiers
}

// NewValidator creates a new Validator with default rules
func NewValidator() *Validator {
	return &Validator{
		MinMatchDuration: 300,  // 5 minutes minimum
		MaxMatchDuration: 7200, // 2 hours maximum
		ValidTiers:       []string{"premium", "professional", "amateur", "unknown"},
	}
}

// ValidateMatch validates a Match struct before insertion
func (v *Validator) ValidateMatch(m Match) error {
	if m.MatchID <= 0 {
		return &ValidationError{
			Field:   "match_id",
			Value:   m.MatchID,
			Message: "must be positive",
		}
	}

	if m.Duration < v.MinMatchDuration || m.Duration > v.MaxMatchDuration {
		return &ValidationError{
			Field:   "duration",
			Value:   m.Duration,
			Message: fmt.Sprintf("must be between %d and %d seconds", v.MinMatchDuration, v.MaxMatchDuration),
		}
	}

	// Check if start time is reasonable (not in the future, not too old)
	if m.StartTime.After(time.Now().Add(5 * time.Minute)) {
		return &ValidationError{
			Field:   "start_time",
			Value:   m.StartTime,
			Message: "cannot be in the future",
		}
	}

	// Check if start time is not too old (more than 1 year)
	if m.StartTime.Before(time.Now().AddDate(-1, 0, 0)) {
		log.Warn().
			Int64("match_id", m.MatchID).
			Time("start_time", m.StartTime).
			Msg("match start time is very old")
	}

	// Validate team IDs if present
	if m.RadiantTeamID < 0 {
		return &ValidationError{
			Field:   "radiant_team_id",
			Value:   m.RadiantTeamID,
			Message: "cannot be negative",
		}
	}

	if m.DireTeamID < 0 {
		return &ValidationError{
			Field:   "dire_team_id",
			Value:   m.DireTeamID,
			Message: "cannot be negative",
		}
	}

	// Validate league ID if present
	if m.LeagueID < 0 {
		return &ValidationError{
			Field:   "league_id",
			Value:   m.LeagueID,
			Message: "cannot be negative",
		}
	}

	// Validate heroes arrays
	if len(m.RadiantHeroes) > 5 {
		return &ValidationError{
			Field:   "radiant_heroes",
			Value:   len(m.RadiantHeroes),
			Message: "cannot have more than 5 heroes",
		}
	}

	if len(m.DireHeroes) > 5 {
		return &ValidationError{
			Field:   "dire_heroes",
			Value:   len(m.DireHeroes),
			Message: "cannot have more than 5 heroes",
		}
	}

	// Validate players arrays
	if len(m.RadiantPlayers) > 5 {
		return &ValidationError{
			Field:   "radiant_players",
			Value:   len(m.RadiantPlayers),
			Message: "cannot have more than 5 players",
		}
	}

	if len(m.DirePlayers) > 5 {
		return &ValidationError{
			Field:   "dire_players",
			Value:   len(m.DirePlayers),
			Message: "cannot have more than 5 players",
		}
	}

	return nil
}

// ValidateMatchMetadata validates MatchMetadata before insertion
func (v *Validator) ValidateMatchMetadata(md MatchMetadata) error {
	if md.MatchID <= 0 {
		return &ValidationError{
			Field:   "match_id",
			Value:   md.MatchID,
			Message: "must be positive",
		}
	}

	// Validate scores
	if md.RadiantScore < 0 || md.RadiantScore > 100 {
		return &ValidationError{
			Field:   "radiant_score",
			Value:   md.RadiantScore,
			Message: "must be between 0 and 100",
		}
	}

	if md.DireScore < 0 || md.DireScore > 100 {
		return &ValidationError{
			Field:   "dire_score",
			Value:   md.DireScore,
			Message: "must be between 0 and 100",
		}
	}

	// Validate version
	if md.Version < 0 {
		return &ValidationError{
			Field:   "version",
			Value:   md.Version,
			Message: "cannot be negative",
		}
	}

	// Validate captains if present
	if md.RadiantCaptain != nil && *md.RadiantCaptain <= 0 {
		return &ValidationError{
			Field:   "radiant_captain",
			Value:   *md.RadiantCaptain,
			Message: "must be positive",
		}
	}

	if md.DireCaptain != nil && *md.DireCaptain <= 0 {
		return &ValidationError{
			Field:   "dire_captain",
			Value:   *md.DireCaptain,
			Message: "must be positive",
		}
	}

	return nil
}

// ValidateLeague validates a League struct before insertion
func (v *Validator) ValidateLeague(l League) error {
	if l.LeagueID <= 0 {
		return &ValidationError{
			Field:   "league_id",
			Value:   l.LeagueID,
			Message: "must be positive",
		}
	}

	if l.Name == "" {
		return &ValidationError{
			Field:   "name",
			Value:   l.Name,
			Message: "cannot be empty",
		}
	}

	// Validate tier if present
	if l.Tier != "" {
		valid := false
		for _, t := range v.ValidTiers {
			if l.Tier == t {
				valid = true
				break
			}
		}
		if !valid {
			log.Warn().
				Int64("league_id", l.LeagueID).
				Str("tier", l.Tier).
				Msg("unusual league tier")
		}
	}

	return nil
}

// ValidateTeam validates a Team struct before insertion
func (v *Validator) ValidateTeam(t Team) error {
	if t.TeamID <= 0 {
		return &ValidationError{
			Field:   "team_id",
			Value:   t.TeamID,
			Message: "must be positive",
		}
	}

	if t.Name == "" {
		return &ValidationError{
			Field:   "name",
			Value:   t.Name,
			Message: "cannot be empty",
		}
	}

	return nil
}

// ValidatePlayer validates a Player struct before insertion
func (v *Validator) ValidatePlayer(p Player) error {
	if p.PlayerID <= 0 {
		return &ValidationError{
			Field:   "player_id",
			Value:   p.PlayerID,
			Message: "must be positive",
		}
	}

	if p.Name == "" {
		return &ValidationError{
			Field:   "name",
			Value:   p.Name,
			Message: "cannot be empty",
		}
	}

	return nil
}

// ValidateSeries validates a Series struct before insertion
func (v *Validator) ValidateSeries(s Series) error {
	if s.SeriesID <= 0 {
		return &ValidationError{
			Field:   "series_id",
			Value:   s.SeriesID,
			Message: "must be positive",
		}
	}

	// Check if start time is reasonable
	if s.StartTime.After(time.Now().Add(5 * time.Minute)) {
		return &ValidationError{
			Field:   "start_time",
			Value:   s.StartTime,
			Message: "cannot be in the future",
		}
	}

	// Validate scores
	if s.TeamAScore < 0 || s.TeamAScore > 10 {
		return &ValidationError{
			Field:   "team_a_score",
			Value:   s.TeamAScore,
			Message: "must be between 0 and 10",
		}
	}

	if s.TeamBScore < 0 || s.TeamBScore > 10 {
		return &ValidationError{
			Field:   "team_b_score",
			Value:   s.TeamBScore,
			Message: "must be between 0 and 10",
		}
	}

	return nil
}

// ValidateODMatch validates an ODMatch before processing
func (v *Validator) ValidateODMatch(m ODMatch) error {
	if m.MatchID <= 0 {
		return &ValidationError{
			Field:   "match_id",
			Value:   m.MatchID,
			Message: "must be positive",
		}
	}

	if m.Duration < v.MinMatchDuration || m.Duration > v.MaxMatchDuration {
		return &ValidationError{
			Field:   "duration",
			Value:   m.Duration,
			Message: fmt.Sprintf("must be between %d and %d seconds", v.MinMatchDuration, v.MaxMatchDuration),
		}
	}

	// Validate league
	if m.League.ID < 0 {
		return &ValidationError{
			Field:   "league.id",
			Value:   m.League.ID,
			Message: "cannot be negative",
		}
	}

	// Validate teams
	if m.RadiantTeam.ID < 0 {
		return &ValidationError{
			Field:   "radiant_team.id",
			Value:   m.RadiantTeam.ID,
			Message: "cannot be negative",
		}
	}

	if m.DireTeam.ID < 0 {
		return &ValidationError{
			Field:   "dire_team.id",
			Value:   m.DireTeam.ID,
			Message: "cannot be negative",
		}
	}

	// Validate scores
	if m.RadiantTeam.Score < 0 || m.RadiantTeam.Score > 100 {
		return &ValidationError{
			Field:   "radiant_team.score",
			Value:   m.RadiantTeam.Score,
			Message: "must be between 0 and 100",
		}
	}

	if m.DireTeam.Score < 0 || m.DireTeam.Score > 100 {
		return &ValidationError{
			Field:   "dire_team.score",
			Value:   m.DireTeam.Score,
			Message: "must be between 0 and 100",
		}
	}

	return nil
}

// ValidateBatch validates a batch of matches and returns valid/invalid counts
func (v *Validator) ValidateBatch(matches []Match, metadata []MatchMetadata) (validCount int, invalidCount int, errors []error) {
	if len(matches) != len(metadata) {
		errors = append(errors, fmt.Errorf("matches and metadata slices have different lengths"))
		return 0, 0, errors
	}

	for i := range matches {
		if err := v.ValidateMatch(matches[i]); err != nil {
			invalidCount++
			errors = append(errors, fmt.Errorf("match %d: %w", matches[i].MatchID, err))
			continue
		}

		if err := v.ValidateMatchMetadata(metadata[i]); err != nil {
			invalidCount++
			errors = append(errors, fmt.Errorf("match metadata %d: %w", metadata[i].MatchID, err))
			continue
		}

		validCount++
	}

	return validCount, invalidCount, errors
}
