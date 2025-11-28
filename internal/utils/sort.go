package utils

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type SortOption struct {
    Field string
    Dir   string
}

func ParseSortParam(sortParam, defaultSort string) (field, dir string) {
    if sortParam == "" {
        sortParam = defaultSort
    }

    if sortParam == "" {
        return "created_at", "DESC"
    }

    dir = "ASC"
    field = sortParam
    if len(sortParam) > 0 && sortParam[0] == '-' {
        dir = "DESC"
        field = sortParam[1:]
    }

    return field, dir
}

func BuildDynamicSorts(queryParams map[string][]string, allowedFields []string) []SortOption {
    allowed := make(map[string]bool)
    for _, a := range allowedFields {
        allowed[a] = true
    }

    var sorts []SortOption

    for key, values := range queryParams {
        if strings.HasPrefix(key, "sort[") && strings.HasSuffix(key, "]") {
            field := key[5 : len(key)-1]
            if !allowed[field] {
                continue
            }
            dir := "ASC"
            if len(values) > 0 {
                v := strings.ToLower(values[0])
                if v == "desc" || v == "d" {
                    dir = "DESC"
                }
            }
            sorts = append(sorts, SortOption{Field: field, Dir: dir})
        }
    }

    if len(sorts) == 0 {
        if raw, ok := queryParams["sort"]; ok && len(raw) > 0 {
            tokens := strings.Split(raw[0], ",")
            for _, t := range tokens {
                t = strings.TrimSpace(t)
                if t == "" {
                    continue
                }
                dir := "ASC"
                field := t
                if field[0] == '-' {
                    dir = "DESC"
                    field = field[1:]
                }
                if allowed[field] {
                    sorts = append(sorts, SortOption{Field: field, Dir: dir})
                }
            }
        }
    }

    return sorts
}

func ApplyDynamicSort(db *gorm.DB, sorts []SortOption, defaultOrder string) *gorm.DB {
    if len(sorts) == 0 {
        if defaultOrder != "" {
            return db.Order(defaultOrder)
        }
        return db
    }

    for _, s := range sorts {
        dir := strings.ToUpper(s.Dir)
        if dir != "ASC" && dir != "DESC" {
            dir = "ASC"
        }
        db = db.Order(fmt.Sprintf("%s %s", s.Field, dir))
    }
    return db
}
