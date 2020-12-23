package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrRuleParseError       = errors.New("rule parse error")
	ErrFieldTypeUnsupported = errors.New("field type unsupported")
	ErrUseForbidden         = errors.New("use forbidden")
	ErrNotAStruct           = errors.New("not a struct")

	ErrNotInList         = errors.New("not in list")
	ErrWrongType         = errors.New("wrong type")
	ErrLessThenMin       = errors.New("less then min")
	ErrMoreThenMax       = errors.New("more then max")
	ErrLengthNotMatched  = errors.New("length  not matched")
	ErrPatternNotMatched = errors.New("pattern not matched")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Field validation error %s: %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	return fmt.Sprintf("Field validation error(s) %d", len(v))
}

type Rule interface {
	Validate(v interface{}) error
}

// Правило "Входи в список"
// Из-за особенностей синтаксиса описания должно быть готово проверять и строки, и числа
// `validate:"in:admin,stuff"` или `validate:"in:200,404,500"`.
type RuleIn struct {
	strVals     []string
	intVals     []int
	stringsOnly bool
}

func NewRuleIn(ruleText string) *RuleIn {
	var ri RuleIn
	ri.strVals = strings.Split(ruleText, ",")
	ri.intVals = make([]int, len(ri.strVals))

	var err error
	for i, val := range ri.strVals {
		ri.intVals[i], err = strconv.Atoi(val)
		if err != nil {
			ri.stringsOnly = true
			break
		}
	}

	return &ri
}

func (ri *RuleIn) checkStrings(strVal string) bool {
	for _, val := range ri.strVals {
		if strVal == val {
			return true
		}
	}

	return false
}

func (ri *RuleIn) checkInts(iVal int) bool {
	for _, val := range ri.intVals {
		if iVal == val {
			return true
		}
	}

	return false
}

func (ri *RuleIn) Validate(v interface{}) error {
	sVal, isStr := v.(string)
	if isStr {
		if ri.checkStrings(sVal) {
			return nil
		}
	}

	if ri.stringsOnly {
		return ErrNotInList
	}

	iVal, isInt := v.(int64)
	if isInt {
		if ri.checkInts(int(iVal)) {
			return nil
		}
	}

	return ErrNotInList
}

// Правило "Проверить вложенную структуру"
// `validate:"nested"`.
type RuleNested struct {
}

func NewRuleNested(ruleText string) *RuleNested {
	return &RuleNested{}
}

func (r *RuleNested) Validate(v interface{}) error {
	return ErrUseForbidden
}

// Правило "Не меньше N"
// `validate:"min:18"`.
type RuleMin struct {
	min int
}

func NewRuleMin(ruleText string) *RuleMin {
	var err error
	var r RuleMin
	r.min, err = strconv.Atoi(ruleText)

	if err != nil {
		return nil
	}

	return &r
}

func (r *RuleMin) Validate(v interface{}) error {
	iVal, isInt := v.(int64)
	if !isInt {
		return ErrWrongType
	}

	if iVal < int64(r.min) {
		return ErrLessThenMin
	}

	return nil
}

// Правило "Не больше N"
// `validate:"max:50"`.
type RuleMax struct {
	max int
}

func NewRuleMax(ruleText string) *RuleMax {
	var err error
	var r RuleMax
	r.max, err = strconv.Atoi(ruleText)

	if err != nil {
		return nil
	}

	return &r
}

func (r *RuleMax) Validate(v interface{}) error {
	iVal, isInt := v.(int64)
	if !isInt {
		return ErrWrongType
	}

	if iVal > int64(r.max) {
		return ErrMoreThenMax
	}

	return nil
}

// Правило "Строка длиной ровно N"
// `validate:"len:11"`.
type RuleLen struct {
	len int
}

func NewRuleLen(ruleText string) *RuleLen {
	var err error
	var r RuleLen
	r.len, err = strconv.Atoi(ruleText)

	if err != nil {
		return nil
	}

	return &r
}

func (r *RuleLen) Validate(v interface{}) error {
	sVal, isStr := v.(string)
	if !isStr {
		return ErrWrongType
	}

	if len(sVal) != r.len {
		return ErrLengthNotMatched
	}

	return nil
}

// Правило "Строка, соответствующая регулярному выражению"
// `validate:"regexp:^\\w+@\\w+\\.\\w+$"`.
type RuleRegExp struct {
	re *regexp.Regexp
}

func NewRuleRegExp(ruleText string) *RuleRegExp {
	var ro RuleRegExp
	ro.re = regexp.MustCompile(ruleText)
	if ro.re == nil {
		return nil
	}

	return &ro
}

func (ro *RuleRegExp) Validate(v interface{}) error {
	sVal, isStr := v.(string)
	if !isStr {
		return ErrWrongType
	}

	if !ro.re.MatchString(sVal) {
		return ErrPatternNotMatched
	}

	return nil
}

// Правило "Выполняются все правила"
// `validate:"min:18|max:50"`.
type RuleGroup struct {
	rules []Rule
}

func NewRuleGroup(ruleText string) *RuleGroup {
	var ro RuleGroup
	rulesText := strings.Split(ruleText, "|")
	ro.rules = make([]Rule, len(rulesText))
	for i, text := range rulesText {
		ro.rules[i] = NewRule(text)
		if ro.rules[i] == nil {
			return nil
		}
	}
	return &ro
}

func (ro *RuleGroup) Validate(v interface{}) error {
	for _, rule := range ro.rules {
		err := rule.Validate(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// Создаём и инициализируем правило по его строковому виду.
func NewRule(ruleText string) Rule {
	rulesCount := strings.Count(ruleText, ":")
	switch {
	case rulesCount == 1:
		ruleParts := strings.Split(ruleText, ":")
		ruleText = ruleParts[1]
		switch ruleParts[0] {
		case "in":
			return NewRuleIn(ruleText)
		case "min":
			return NewRuleMin(ruleText)
		case "max":
			return NewRuleMax(ruleText)
		case "regexp":
			return NewRuleRegExp(ruleText)
		case "len":
			return NewRuleLen(ruleText)
		}

	case rulesCount-1 == strings.Count(ruleText, "|"):
		return NewRuleGroup(ruleText)
	case ruleText == "nested":
		return NewRuleNested(ruleText)
	}

	return nil
}

// Применяем правило к полю структуры.
func validateField(field reflect.Value, fieldRule Rule) error {
	switch field.Type().Kind() {
	case reflect.Int:
		return fieldRule.Validate(field.Int())
	case reflect.String:
		return fieldRule.Validate(field.String())
	case reflect.Slice:
		return validateSliceField(field, fieldRule)
	case reflect.Struct:
		return validateStruct(field.Interface(), ".")
	}

	return ErrFieldTypeUnsupported
}

// Применяем правило ко всем элементам среза.
func validateSliceField(field reflect.Value, fieldRule Rule) error {
	foundErrors := make(ValidationErrors, 0, field.Len())
	sliceEl := field.Type().Name() + "[%d]"
	for i := 0; i < field.Len(); i++ {
		err := validateField(field.Index(i), fieldRule)
		if err != nil {
			foundErrors = append(foundErrors, ValidationError{fmt.Sprintf(sliceEl, i), err})
		}
	}

	if len(foundErrors) == 0 {
		return nil
	}

	return foundErrors
}

func validateStruct(field interface{}, fieldOwner string) error {
	rValue := reflect.ValueOf(field)
	description := rValue.Type()

	if description.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	foundErrors := make(ValidationErrors, 0, description.NumField())

	for i := 0; i < description.NumField(); i++ {
		field := description.Field(i)
		rules, exist := field.Tag.Lookup("validate")
		if !exist {
			continue
		}

		rule := NewRule(rules)
		if rule == nil {
			foundErrors = append(foundErrors, ValidationError{fieldOwner + field.Name, ErrRuleParseError})
			continue
		}

		err := validateField(rValue.Field(i), rule)

		if err != nil {
			nested := ValidationErrors{}
			if errors.As(err, &nested) {
				for _, nestedErr := range nested {
					foundErrors = append(foundErrors, ValidationError{fieldOwner + field.Name + nestedErr.Field, nestedErr.Err})
				}
			} else {
				foundErrors = append(foundErrors, ValidationError{fieldOwner + field.Name, err})
			}
		}
	}

	if len(foundErrors) == 0 {
		return nil
	}

	return foundErrors
}

func Validate(v interface{}) error {
	return validateStruct(v, "")
}
