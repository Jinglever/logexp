package logexp

// “或”表达式
type ExpressionOr struct {
	OrgExp     []rune         `json:"-"` // 原始表达式
	iter       int            `json:"-"` // 遍历原表达式的浮标
	Type       ExpressionType `json:"type"`
	IsNegative bool           `json:"is_negative"` // 是否取非
	Exps       []IExpression  `json:"expressions"` // “或”表达式的内部应该只有“且”表达式
}

func (e *ExpressionOr) GetIsNegative() bool {
	return e.IsNegative
}

func (e *ExpressionOr) ReverseIsNegative() {
	e.IsNegative = !e.IsNegative
}

func (e *ExpressionOr) GetType() ExpressionType {
	return e.Type
}

func (e *ExpressionOr) GetExps() []IExpression {
	return e.Exps
}

func (e *ExpressionOr) Match(text string) bool {
	res := false
	for i := range e.Exps {
		// “或”表达式里，只要有一个条件通过，就算通过
		if e.Exps[i].Match(text) {
			res = true
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
 * @Param mode: 该表达式尝试过当作哪些表达式类型来编译，默认值0，尝试过的类型(ExpressionType_Or|ExpressionType_And)通过取或存储在此字段中
 */
func NewExpressionOr(exp []rune, isNegative bool, mode int) (IExpression, *CstError) {
	// 表达式不能为空字符串
	if len(exp) == 0 {
		return nil, newCstError(ErrCodeInvalidExpression, "invalid expression: %v", string(exp))
	}
	// 不重复当作同样的表达式类型来编译，否则会死循环
	if mode|int(ExpressionType_Or) == mode {
		// 程序跑到这里，说明肯定表达式有语法问题
		return nil, newCstError(ErrCodeInvalidExpression, "invalid expression: %v", string(exp))
	}
	mode |= int(ExpressionType_Or)

	expOr := ExpressionOr{
		OrgExp:     []rune(exp),
		iter:       0,
		Type:       ExpressionType_Or,
		IsNegative: isNegative,
		Exps:       make([]IExpression, 0),
	}

	for expOr.iter < len(expOr.OrgExp) {
		// 提取下一个子表达式
		subExp, cerr := expOr.getNextSubExp()
		if cerr != nil {
			return nil, cerr
		}

		// 修整子表达式
		isTrim := subExp.Trim()

		// 编译子表达式
		var exp IExpression
		if subExp.IsMeta {
			// 元表达式
			exp, cerr = NewExpressionMeta(subExp.Exp, subExp.IsNegative)

		} else if isTrim {
			// 如果子表达式经过修整发生了变化，那么我们就不能确定新的子表达式是否属于且表达式
			// 这个情况下，仍然要当或表达式来尝试编译
			exp, cerr = NewExpressionOr(subExp.Exp, subExp.IsNegative, 0)
		} else {
			if len(subExp.Exp) == len(expOr.OrgExp) {
				// 如果子表达式跟原表达式完全相同，那么mode要透传进去
				exp, cerr = NewExpressionAnd(subExp.Exp, subExp.IsNegative, mode)
			} else {
				exp, cerr = NewExpressionAnd(subExp.Exp, subExp.IsNegative, 0)
			}
		}
		if cerr != nil {
			return nil, cerr
		}

		// 装载子表达式
		if exp.GetType() == ExpressionType_Or && !exp.GetIsNegative() {
			// 子表达式如果跟父表达式同类型，直接展开，减少层级
			expOr.Exps = append(expOr.Exps, exp.GetExps()...)
		} else {
			expOr.Exps = append(expOr.Exps, exp)
		}
	}

	// 如果只有一个子表达式，那么可以直接往上层提，减少不必要的层级
	if len(expOr.Exps) == 1 {
		if expOr.IsNegative {
			expOr.Exps[0].ReverseIsNegative()
		}
		return expOr.Exps[0], nil
	}

	return &expOr, nil
}

// 把原始表达式当作或表达式，分解提取子表达式
func (e *ExpressionOr) getNextSubExp() (*SubExp, *CstError) {
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
		case rune('&'):
			subExp.Exp = append(subExp.Exp, c)
			subExp.IsMeta = false
		case rune('|'):
			if bsIdx > 0 { // 括号内的'|'符号，不能作为当前层级分割子表达式的标识
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
