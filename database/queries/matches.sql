SELECT
    m.match_id,
    m.start_time,
    m.duration,
    JSON_BUILD_OBJECT(
    	'id', m.dire_team_id,
    	'name',dire.name,
    	'score', md.dire_score
    ) AS dire_team,
    JSON_BUILD_OBJECT(
    	'id', m.radiant_team_id,
    	'name',radiant.name,
    	'score', md.radiant_score
    ) AS radiant_team,
    m.radiant_win,
    md.series_id,
    JSON_BUILD_OBJECT(
        'id', l.league_id,
        'name', l.name
    ) AS league,
    m.radiant_heroes,
    m.dire_heroes
FROM matches m
LEFT JOIN matches_metadata md USING(match_id)
LEFT JOIN leagues l USING(league_id)
LEFT JOIN teams radiant ON (radiant.team_id = m.radiant_team_id)
LEFT JOIN teams dire ON (dire.team_id = m.dire_team_id)
WHERE 
    --filtering
    true
ORDER BY m.match_id DESC
LIMIT $1
OFFSET $2;



--FILTERS:

-- TEAM T
WHERE(m.radiant_team_id = T OR m.dire_team_id = T) 

-- HERO H
WHERE(m.radiant_heroes @> H OR m.dire_heroes @> H) 

--PLAYER P
WHERE (m.radiant_players @> P OR m.dire_players @> P)


--LEAGUE L 
WHERE m.league_id = L

-- SORT BY : [newest, oldest]




