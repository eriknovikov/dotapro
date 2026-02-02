# API ENDPOINTS AND BEHAVIOUR

## MATCHES

- /matches/:id -> in depth info bout the match.
  info bout the match given: id, league, radiant, dire, etc

- /matches -> list of matches. filters:
  league L. team T. player P. hero H. Sort by = {newest,oldest}

## SERIES

/series -> the list of the SERIES played by team T, in league L. For each SERIES, display: id, team1, team2, score, league, date.

- /series/:id -> the info about each match played in the series
