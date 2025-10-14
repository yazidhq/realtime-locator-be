package utils

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type FilterOptions struct {
    Field string
    Op    string
    Value string
}

func BuildDynamicFilters(queryParams map[string][]string, allowedOps []string) []FilterOptions {
	allowedOpMap := make(map[string]bool)
	for _, o := range allowedOps {
		allowedOpMap[strings.ToLower(o)] = true
	}

	var filters []FilterOptions

	for key, values := range queryParams {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			field := key[7 : len(key)-1]

			value := ""
			if len(values) > 0 {
				value = values[0]
			}

			opParam := fmt.Sprintf("op[%s]", field)
			op := "="
			if opValues, ok := queryParams[opParam]; ok && len(opValues) > 0 {
				tmpOp := strings.ToLower(opValues[0])
				if allowedOpMap[tmpOp] {
					op = tmpOp
				}
			}

			filters = append(filters, FilterOptions{
				Field: field,
				Op:    op,
				Value: value,
			})
		}
	}

	return filters
}

func BuildCQLFilter(filters []FilterOptions) string {
	var parts []string
	for _, f := range filters {
		if _, err := strconv.Atoi(f.Value); err == nil {
			parts = append(parts, fmt.Sprintf("%s %s %s", f.Field, f.Op, f.Value))
		} else {
			if strings.ToLower(f.Op) == "like" {
				parts = append(parts, fmt.Sprintf("%s ILIKE '%%%s%%'", f.Field, f.Value))
			} else {
				parts = append(parts, fmt.Sprintf("%s %s '%s'", f.Field, f.Op, f.Value))
			}
		}
	}
	return strings.Join(parts, " AND ")
}


func ApplyDynamicFilters(db *gorm.DB, filters []FilterOptions) *gorm.DB {
    for _, f := range filters {
        switch f.Op {
        case "like":
            db = db.Where(fmt.Sprintf("%s LIKE ?", f.Field), "%"+f.Value+"%")
        default:
            db = db.Where(fmt.Sprintf("%s %s ?", f.Field, f.Op), f.Value)
        }
    }
    return db
}