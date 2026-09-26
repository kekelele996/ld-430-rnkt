package util

// DefaultPageSize is applied when clients omit page_size.
const DefaultPageSize int64 = 20

// MaxPageSize caps client-provided page sizes to avoid abusive large queries.
const MaxPageSize int64 = 100

// NormalizePage returns a page number >= 1.
func NormalizePage(page int64) int64 {
	if page < 1 {
		return 1
	}
	return page
}

// NormalizePageSize returns a page size between 1 and MaxPageSize.
func NormalizePageSize(size int64) int64 {
	if size < 1 {
		return DefaultPageSize
	}
	if size > MaxPageSize {
		return MaxPageSize
	}
	return size
}

// Offset calculates the skip value for a repository query.
func Offset(page, pageSize int64) int64 {
	page = NormalizePage(page)
	pageSize = NormalizePageSize(pageSize)
	return (page - 1) * pageSize
}
