package logexp

import "strings"

// 元表达式
type ExpressionMeta struct {
	Type       ExpressionType `json:"type"`
	IsNegative bool           `json:"is_negative"` // 是否取非
	Keyword    string         `json:"keyword"`     // 关键词
}

func (e *ExpressionMeta) GetIsNegative() bool {
	return e.IsNegative
}

func (e *ExpressionMeta) ReverseIsNegative() {
	e.IsNegative = !e.IsNegative
}

func (e *ExpressionMeta) GetType() ExpressionType {
	return e.Type
}

func (e *ExpressionMeta) GetExps() []IExpression {
	return []IExpression{}
}

func (e *ExpressionMeta) Match(text string) bool {
	res := false
	if strings.Contains(text, e.Keyword) {
		res = true
	}
	if e.IsNegative {
		res = !res
	}
	return res
}

func NewExpressionMeta(exp []rune, isNegative bool) (IExpression, *CstError) {
	mapNotInclude := map[rune]struct{}{
		'|': {}, '&': {}, '!': {}, '(': {}, ')': {},
	}
	for _, c := range exp {
		if _, ok := mapNotInclude[c]; ok {
			return nil, newCstError(ErrCodeInvalidExpression, "invalid meta expression: %v", string(exp))
		}
	}
	expMeta := ExpressionMeta{
		Type:       ExpressionType_Meta,
		IsNegative: isNegative,
		Keyword:    string(exp),
	}
	return &expMeta, nil
}
