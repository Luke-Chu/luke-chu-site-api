package pager

const (
	DefaultPage     = 1
	DefaultPageSize = 30
	MaxPageSize     = 60
)

func Normalize(page, pageSize int) (int, int) {
	return NormalizeWith(page, pageSize, DefaultPageSize, MaxPageSize)
}

func NormalizeWith(page, pageSize, defaultPageSize, maxPageSize int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}

	switch {
	case pageSize < 1:
		pageSize = defaultPageSize
	case pageSize > maxPageSize:
		pageSize = maxPageSize
	}

	return page, pageSize
}

func Offset(page, pageSize int) int {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	return (page - 1) * pageSize
}

func TotalPages(total int64, pageSize int) int {
	if total <= 0 {
		return 0
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	pages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		pages++
	}
	return int(pages)
}
