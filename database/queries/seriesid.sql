SELECT 
    series_id,
    JSON_BUILD_OBJECT(
        id, l.league_id,
        name, l.name
    ) as league,
    JSON_BUILD_OBJECT(
        id, a.team_id,
        name, a.name,
        score, team_a_score
    ) as team_a,
    JSON_BUILD_OBJECT(
        id, b.team_id,
        name, b.name,
        score, team_b_score
    ) as team_b,
    start_time
    --, matches : JSON with the list of matches for this series

FROM series s
LEFT JOIN leagues l USING (league_id)
LEFT JOIN teams a ON a.team_id = s.team_a_id
LEFT JOIN teams b ON b.team_id = s.team_b_id
WHERE s.series_id = $1