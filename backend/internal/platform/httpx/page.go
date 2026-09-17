package httpx

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Pagination bounds shared by every list endpoint's contract
// (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 4.1).
const (
	DefaultPageLimit = 25
	MaxPageLimit     = 100
	MaxSearchLength  = 64
)

// PageParams is the validated shape of the limit/offset/search query
// parameters a paginated list endpoint reads.
type PageParams struct {
	Limit  int32
	Offset int32
	Search string
}

// ParsePage reads limit/offset/search from the request's query string.
// limit defaults to DefaultPageLimit and is clamped to [1, MaxPageLimit];
// offset is clamped to >= 0; search is trimmed, capped at MaxSearchLength
// and has LIKE wildcards escaped so it's safe to interpolate into an
// `ILIKE '%'||search||'%'` clause — a bare "%" or "_" would otherwise widen
// the match arbitrarily, and one "%" alone would scan the whole table.
func ParsePage(c *gin.Context) PageParams {
	limit := int32(DefaultPageLimit)
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
		limit = int32(v)
	}
	if limit > MaxPageLimit {
		limit = MaxPageLimit
	}

	var offset int32
	if v, err := strconv.Atoi(c.Query("offset")); err == nil && v > 0 {
		offset = int32(v)
	}

	search := strings.TrimSpace(c.Query("search"))
	if len(search) > MaxSearchLength {
		search = search[:MaxSearchLength]
	}

	return PageParams{Limit: limit, Offset: offset, Search: EscapeLikeTerm(search)}
}

// EscapeLikeTerm neutralises ILIKE's own wildcards ('%', '_') and its
// escape character ('\') in user input bound for an ILIKE pattern.
func EscapeLikeTerm(term string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(term)
}

// Page wraps one page of items with the metadata the frontend's DataGrid
// server mode needs — never a bare array, so `total` can travel with it.
type Page[T any] struct {
	Items  []T   `json:"items"`
	Total  int64 `json:"total"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
