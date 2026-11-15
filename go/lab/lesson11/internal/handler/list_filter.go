package handler

import (
	"errors"
	"net/url"
	"strconv"

	"golab/lesson11/internal/domain"
)

// parseListFilter는 ?done=true&q=go&limit=10&offset=20 같은 쿼리 문자열을 ListFilter로 바꿉니다.
// 형식이 틀린 값은 에러로 돌려주고, 값의 범위 검사는 서비스가 맡습니다.
func parseListFilter(q url.Values) (domain.ListFilter, error) {
	filter := domain.ListFilter{Query: q.Get("q")}

	switch v := q.Get("done"); v {
	case "":
	case "true", "false":
		done := v == "true"
		filter.Done = &done
	default:
		return domain.ListFilter{}, errors.New("done은 true 또는 false여야 합니다")
	}

	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return domain.ListFilter{}, errors.New("limit은 정수여야 합니다")
		}
		filter.Limit = n
	}
	if v := q.Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return domain.ListFilter{}, errors.New("offset은 정수여야 합니다")
		}
		filter.Offset = n
	}
	return filter, nil
}
