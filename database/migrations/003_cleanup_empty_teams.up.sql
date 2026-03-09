-- Delete matches_metadata for matches that reference teams with empty names
DELETE FROM matches_metadata
WHERE match_id IN (
  SELECT match_id FROM matches
  WHERE radiant_team_id IN (SELECT team_id FROM teams WHERE name = '')
     OR dire_team_id IN (SELECT team_id FROM teams WHERE name = '')
);

-- Delete series_match entries for series that reference teams with empty names
DELETE FROM series_match
WHERE series_id IN (
  SELECT series_id FROM series
  WHERE team_a_id IN (SELECT team_id FROM teams WHERE name = '')
     OR team_b_id IN (SELECT team_id FROM teams WHERE name = '')
);

-- Delete matches that reference teams with empty names
DELETE FROM matches
WHERE radiant_team_id IN (SELECT team_id FROM teams WHERE name = '')
   OR dire_team_id IN (SELECT team_id FROM teams WHERE name = '');

-- Delete series that reference teams with empty names
DELETE FROM series
WHERE team_a_id IN (SELECT team_id FROM teams WHERE name = '')
   OR team_b_id IN (SELECT team_id FROM teams WHERE name = '');

-- Finally delete teams with empty names
DELETE FROM teams WHERE name = '';
