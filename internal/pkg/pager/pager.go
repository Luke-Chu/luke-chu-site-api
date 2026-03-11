package pager

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func Normalize(page, pageSize int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}

	switch {
	case pageSize < 1:
		pageSize = DefaultPageSize
	case pageSize > MaxPageSize:
		pageSize = MaxPageSize
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
