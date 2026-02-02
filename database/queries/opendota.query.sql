SELECT
    m.match_id,
    m.radiant_win,
    m.start_time AS start_time,
    m.duration,
    m.leagueid AS league_id,
    l.name AS league_name,
    l.tier AS league_tier,
    m.radiant_team_id,
    rt.name AS radiant_team_name,
    rt.tag AS radiant_team_tag,
    rt.logo_url AS radiant_team_logo_url,
    m.dire_team_id,
    dt.name AS dire_team_name,
    dt.tag AS dire_team_tag,
    dt.logo_url AS dire_team_logo_url,
    m.radiant_score,
    m.dire_score,
    m.radiant_captain,
    m.dire_captain,
    m.radiant_gold_adv,
    m.radiant_xp_adv,
    m.version,
    m.series_id,
    (
        SELECT JSON_AGG(
            JSON_BUILD_OBJECT(
                'hero_id', pm.hero_id,
                'account_id', pm.account_id,
                'player_slot', pm.player_slot,
                'kills', pm.kills,
                'deaths', pm.deaths,
                'assists', pm.assists,
                'gold_per_min', pm.gold_per_min,
                'xp_per_min', pm.xp_per_min,
                'last_hits', pm.last_hits,
                'denies', pm.denies,
                'level', pm.level,
                'item_0', pm.item_0,
                'item_1', pm.item_1,
                'item_2', pm.item_2,
                'item_3', pm.item_3,
                'item_4', pm.item_4,
                'item_5', pm.item_5,
                'item_neutral', pm.item_neutral,
                'backpack_0', pm.backpack_0,
                'backpack_1', pm.backpack_1,
                'backpack_2', pm.backpack_2,
                'net_worth', pm.net_worth,
                'isRadiant', (pm.player_slot < 128)
            )
            ORDER BY pm.player_slot
        )
        FROM player_matches pm
        WHERE pm.match_id = m.match_id
    ) AS players,
    (
        SELECT JSON_AGG(
            JSON_BUILD_OBJECT(
                'hero_id', pb.hero_id,
                'is_pick', pb.is_pick,
                'order', pb.ord,
                'team', pb.team
            )
            ORDER BY pb.ord
        )
        FROM picks_bans pb
        WHERE pb.match_id = m.match_id
    ) AS picks_bans
FROM matches m
LEFT JOIN leagues l ON m.leagueid = l.leagueid
LEFT JOIN teams rt ON m.radiant_team_id = rt.team_id
LEFT JOIN teams dt ON m.dire_team_id = dt.team_id
WHERE m.match_id IN (%s)
ORDER BY match_id ASC

