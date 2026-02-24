--query to run against opendota.com/explorer
WITH recent_pro_matches AS (
    SELECT 
        m.radiant_team_id, 
        m.dire_team_id
    FROM matches m
    JOIN leagues l ON m.leagueid = l.leagueid
    WHERE m.start_time > (extract(epoch from now()) - 2592000) -- Last 30 days
    AND l.tier IN ('professional', 'premium') -- Only important leagues
    AND m.radiant_team_id IS NOT NULL -- Exclude pub matches/unregistered teams
    AND m.dire_team_id IS NOT NULL
)
SELECT 
    t.team_id as id,
    t.name,
    t.logo_url,
    count(*) as matches_played
FROM (
    -- Unpivot: Combine Radiant and Dire teams into one list
    SELECT radiant_team_id as team_id FROM recent_pro_matches
    UNION ALL
    SELECT dire_team_id as team_id FROM recent_pro_matches
) as participation
JOIN teams t ON participation.team_id = t.team_id
GROUP BY t.team_id, t.name, t.logo_url
ORDER BY matches_played DESC
LIMIT 10