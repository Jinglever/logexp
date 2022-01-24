package logexp

// “且”表达式
type ExpressionAnd struct {
	OrgExp     []rune         `json:"-"` // 原始表达式
	iter       int            `json:"-"` // 遍历原表达式的浮标
	Type       ExpressionType `json:"type"`
	IsNegative bool           `json:"is_negative"` // 是否取非
	Exps       []IExpression  `json:"expressions"` // “且”表达式的内部应该只有“或”表达式
}

func (e *ExpressionAnd) GetIsNegative() bool {
	return e.IsNegative
}

func (e *ExpressionAnd) ReverseIsNegative() {
	e.IsNegative = !e.IsNegative
}

func (e *ExpressionAnd) GetType() ExpressionType {
	return e.Type
}

func (e *ExpressionAnd) GetExps() []IExpression {
	return e.Exps
}

func (e *ExpressionAnd) Match(text string) bool {
	res := true
	for i := range e.Exps {
		// “且”表达式里，只要有一个条件不通过，就算不通过
		if !e.Exps[i].Match(text) {
			res = false
			break
		}
	}
	if e.IsNegative {
		res = !res
	}
	return res
}

/*
 * @Param exp: 表达式字符串
 * @Param isNegative: 是否取非
 * @Param mode: 该表达式尝试过编译模式（ExpressionType_Or | ExpressionType_And)
 */
func NewExpressionAnd(exp []rune, isNegative bool, mode int) (IExpression, *CstError) {
	// 表达式不能为空字符串
	if len(exp) == 0 {
		return nil, newCstError(ErrCodeInvalidExpression, "invalid expression: %v", string(exp))
	}
	// 不重复当作同样的表达式类型来编译，否则会死循环
	if mode|int(ExpressionType_And) == mode {
		// 程序跑到这里，说明肯定表达式有语法问题
		return nil, newCstError(ErrCodeInvalidExpression, "invalid expression: %v", string(exp))
	}
	mode |= int(ExpressionType_And)

	expAnd := ExpressionAnd{
		OrgExp:     []rune(exp),
		iter:       0,
		Type:       ExpressionType_And,
		IsNegative: isNegative,
		Exps:       make([]IExpression, 0),
	}

	for expAnd.iter < len(expAnd.OrgExp) {
		// 提取下一个子表达式
		subExp, cerr := expAnd.getNextSubExp()
		if cerr != nil {
			return nil, cerr
		}

		// 修整子表达式
		subExp.Trim()

		// 编译子表达式
		var exp IExpression
		if subExp.IsMeta {
			// 元表达式
			exp, cerr = NewExpressionMeta(subExp.Exp, subExp.IsNegative)
		} else {
			if len(subExp.Exp) == len(expAnd.OrgExp) {
				// 如果子表达式跟原表达式完全相同，那么mode要透传进去
				exp, cerr = NewExpressionOr(subExp.Exp, subExp.IsNegative, mode)
			} else {
				exp, cerr = NewExpressionOr(subExp.Exp, subExp.IsNegative, 0)
			}
		}
		if cerr != nil {
			return nil, cerr
		}

		// 装载子表达式
		if exp.GetType() == ExpressionType_And && !exp.GetIsNegative() {
			// 子表达式如果跟父表达式同类型，直接展开，减少层级
			expAnd.Exps = append(expAnd.Exps, exp.GetExps()...)
		} else {
			expAnd.Exps = append(expAnd.Exps, exp)
		}
	}

	// 如果只有一个子表达式，那么可以直接往上层提，减少不必要的层级
	if len(expAnd.Exps) == 1 {
		if expAnd.IsNegative {
			expAnd.Exps[0].ReverseIsNegative()
		}
		return expAnd.Exps[0], nil
	}

	return &expAnd, nil
}

// 把原始表达式当作或表达式，分解提取子表达式
func (e *ExpressionAnd) getNextSubExp() (*SubExp, *CstError) {
	subExp := SubExp{
		IsNegative:   false,
		Exp:          make([]rune, 0, len(e.OrgExp)),
		IsMeta:       true,
		BracketStack: make([]int, len(e.OrgExp)),
	}
	bsIdx := 0 // bracketStack的浮标
	for i := e.iter; i < len(e.OrgExp); i++ {
		c := e.OrgExp[i]
		switch c {
		case rune('('):
			subExp.BracketStack[bsIdx] = len(subExp.Exp) // 记录左括号在子表达式里的下标
			bsIdx++
			subExp.Exp = append(subExp.Exp, c)
		case rune(')'):
			bsIdx--
			if bsIdx < 0 { // 括号不匹配
				return nil, newCstError(ErrCodeInvalidExpression, "invalid expression: %v", string(e.OrgExp))
			}
			subExp.Exp = append(subExp.Exp, c)
		case rune('|'):
			subExp.Exp = append(subExp.Exp, c)
			subExp.IsMeta = false
		case rune('&'):
			if bsIdx > 0 { // 括号内的'&'符号，不能作为当前层级分割子表达式的标识
				subExp.Exp = append(subExp.Exp, c)
				subExp.IsMeta = false
			} else {
				e.iter = i + 1
				return &subExp, nil
			}
		default:
			subExp.Exp = append(subExp.Exp, c)
		}
	}
	if bsIdx > 0 { // 子表达式里，括号不配对
		return nil, newCstError(ErrCodeInvalidExpression, "invalid expression: %v", string(subExp.Exp))
	}
	e.iter = len(e.OrgExp)
	return &subExp, nil
}
