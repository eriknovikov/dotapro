package matches

import "github.com/uptrace/bun"

var EMPTY_MATCH_MODEL *Match

func (m *Model) buildQueryGetOne(id int64) *bun.SelectQuery {
	return m.DB.NewSelect().
		Model(EMPTY_MATCH_MODEL).
		Column("m.*").
		Where("match_id = ?", id)
}
