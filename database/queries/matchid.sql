SELECT
    m.match_id,
    m.start_time,
    m.duration,
    JSON_BUILD_OBJECT(
    	'id', m.dire_team_id,
    	'name',dire.name,
    	'score', md.dire_score,
        'captain', md.dire_captain
    ) AS dire_team,
    JSON_BUILD_OBJECT(
    	'id', m.radiant_team_id,
    	'name',radiant.name,
    	'score', md.radiant_score,
        'captain', md.radiant_captain
    ) AS radiant_team,
    m.radiant_win,
    md.series_id,
    json_build_object(
        'id', l.league_id,
        'name', l.name,
        'tier', l.tier
    ) AS league,
    --extra fields 
    md.picks_bans,
    md.players_data,
    md.radiant_gold_adv,
    md.radiant_xp_adv,
    m.patch,
    m.version
FROM matches m
LEFT JOIN matches_metadata md USING(match_id)
LEFT JOIN leagues l USING(league_id)
LEFT JOIN teams radiant ON (radiant.team_id = m.radiant_team_id)
LEFT JOIN teams dire ON (dire.team_id = m.dire_team_id)
WHERE match_id = $match_id;