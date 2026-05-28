package dafi

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const (
	parameterPage      = "page"
	parameterLimit     = "limit"
	parameterSort      = "sort"
	parameterSelect    = "select"
	parameterRelations = "relations"
	defaultChaining    = And
)

type QueryParser struct {
	operators map[FilterOperator]struct{}
}

func NewQueryParser() *QueryParser {
	return &QueryParser{
		operators: map[FilterOperator]struct{}{
			Equal:          {},
			NotEqual:       {},
			Greater:        {},
			GreaterOrEqual: {},
			Less:           {},
			LessOrEqual:    {},
			Like:           {},
			In:             {},
			NotIn:          {},
			Contains:       {},
			NotContains:    {},
			Is:             {},
			IsNull:         {},
			IsNot:          {},
			IsNotNull:      {},
			Default:        {},
		},
	}
}

func (p *QueryParser) Parse(values url.Values) (Criteria, error) {
	criteria := Criteria{}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		values := values[key]
		if key == "datastar" {
			continue
		}

		if err := p.parseValues(key, values, &criteria); err != nil {
			return Criteria{}, err
		}
	}

	return criteria, nil
}

func (p *QueryParser) parseValues(key string, values []string, criteria *Criteria) error {
	for _, value := range values {
		if value == "" {
			continue
		}

		if key == parameterSelect {
			if err := p.parseSelect(value, criteria); err != nil {
				return err
			}

			continue
		}

		if key == parameterRelations {
			p.parseRelations(value, criteria)

			continue
		}

		parts := strings.SplitN(value, ":", 4)
		if len(parts) == 1 {
			continue
		}

		if err := p.parsePart(key, parts, criteria); err != nil {
			return err
		}
	}

	return nil
}

func (p *QueryParser) parsePart(key string, parts []string, criteria *Criteria) error {
	switch {
	case p.isPaginationPart(parts):
		return p.parsePagination(parts, &criteria.Pagination)
	case p.isSortPart(parts):
		criteria.Sorts = append(criteria.Sorts, p.parseSort(key, parts))
	default:
		filter, err := p.parseFilter(key, parts)
		if err != nil {
			return err
		}

		if filter.Module != "" {
			if criteria.FiltersByModule == nil {
				criteria.FiltersByModule = make(map[string]Filters)
			}

			if _, ok := criteria.FiltersByModule[filter.Module]; !ok {
				criteria.FiltersByModule[filter.Module] = Filters{}
			}
			criteria.FiltersByModule[filter.Module] = append(criteria.FiltersByModule[filter.Module], filter)

			if filter.OverridePreviousFilterChainingKey != "" && len(criteria.Filters) > 1 {
				criteria.FiltersByModule[filter.Module][len(criteria.FiltersByModule[filter.Module])-2].ChainingKey = filter.OverridePreviousFilterChainingKey
				criteria.FiltersByModule[filter.Module][len(criteria.FiltersByModule[filter.Module])-1].OverridePreviousFilterChainingKey = ""
			}

			return nil
		}

		criteria.Filters = append(criteria.Filters, filter)
		if filter.OverridePreviousFilterChainingKey != "" && len(criteria.Filters) > 1 {
			criteria.Filters[len(criteria.Filters)-2].ChainingKey = filter.OverridePreviousFilterChainingKey
			criteria.Filters[len(criteria.Filters)-1].OverridePreviousFilterChainingKey = ""
		}
	}

	return nil
}

func (p *QueryParser) isPaginationPart(parts []string) bool {
	return len(parts) == 2 && (parts[0] == parameterPage || parts[0] == parameterLimit)
}

func (p *QueryParser) isSortPart(parts []string) bool {
	return len(parts) == 2 && parts[0] == parameterSort
}

func (p *QueryParser) parsePagination(parts []string, pagination *Pagination) error {
	value, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("%w: %s: %w", ErrInvalidPaginationValue, parts[1], err)
	}

	switch parts[0] {
	case parameterPage:
		if value >= 0 {
			pagination.PageNumber = uint(value)
		}
	case parameterLimit:
		if value >= 0 {
			pagination.PageSize = uint(value)
		}
	}

	return nil
}

func (p *QueryParser) parseSort(field string, parts []string) Sort {
	return Sort{
		Field: SortField(field),
		Type:  SortType(strings.ToUpper(parts[1])),
	}
}

func (p *QueryParser) parseFilter(field string, parts []string) (Filter, error) {
	overridePreviousFilterChainingKey := FilterChainingKey("")
	if len(parts) == 4 {
		overridePreviousFilterChainingKey = FilterChainingKey(strings.ToUpper(parts[3]))
		parts = parts[1:]
	}
	if len(parts) == 3 && (strings.EqualFold(parts[0], string(And)) || strings.EqualFold(parts[0], string(Or))) {
		overridePreviousFilterChainingKey = FilterChainingKey(strings.ToUpper(parts[0]))
		parts = parts[1:]
	}

	operator := p.determineOperator(parts[0])
	chainingKey := p.determineChainingKey(parts)

	var value any = parts[1]
	if operator == In || operator == NotIn {
		value = strings.Split(parts[1], ",")
	}

	fieldSplit := strings.Split(field, ".")
	module := ""
	fieldKey := field
	if len(fieldSplit) == 2 {
		module = fieldSplit[0]
		fieldKey = fieldSplit[1]
	}

	return Filter{
		Module:                            module,
		Field:                             FilterField(fieldKey),
		Operator:                          operator,
		Value:                             value,
		ChainingKey:                       chainingKey,
		OverridePreviousFilterChainingKey: overridePreviousFilterChainingKey,
	}, nil
}

func (p *QueryParser) determineOperator(op string) FilterOperator {
	operator := FilterOperator(op)
	if _, ok := p.operators[operator]; !ok {
		return Equal
	}

	return operator
}

func (p *QueryParser) determineChainingKey(parts []string) FilterChainingKey {
	if len(parts) == 3 {
		return FilterChainingKey(strings.ToUpper(parts[2]))
	}

	return defaultChaining
}

func (p *QueryParser) parseSelect(value string, criteria *Criteria) error {
	if value == "*" {
		criteria.SelectColumns = nil

		return nil
	}

	fields := strings.Split(value, ",")
	selectColumns := make([]string, 0, len(fields))

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field != "" {
			selectColumns = append(selectColumns, field)
		}
	}

	criteria.SelectColumns = selectColumns

	return nil
}

func (p *QueryParser) parseRelations(value string, criteria *Criteria) {
	relations := strings.Split(value, ",")
	for _, rel := range relations {
		rel = strings.TrimSpace(rel)
		if rel != "" {
			criteria.Relations = append(criteria.Relations, rel)
		}
	}
}

func ValidateRelations(relations []string, allowedRelations map[string]struct{}) error {
	for _, rel := range relations {
		if _, ok := allowedRelations[rel]; !ok {
			return fmt.Errorf("%w: %s", ErrInvalidRelation, rel)
		}
	}

	return nil
}

func ValidateSelectFields(selectFields []string, validFields map[string]string) error {
	if len(selectFields) == 0 {
		return nil
	}

	for _, field := range selectFields {
		if _, ok := validFields[field]; !ok {
			return fmt.Errorf("%w: %s", ErrInvalidSelectField, field)
		}
	}

	return nil
}
