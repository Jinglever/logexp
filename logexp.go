/* 实现思路
 *    1、复杂表达式可以切分成子表达式，表达式不能为空字符串。
 *    2、表达式类型只有三种：或、且、元；这三种类型的表达式，均有一个属性：是否取非；元表达式内部不应该出现任何连接符。
 *    3、关于括号，括号必须配对；表达式最外层的括号应该被摘掉。
 *    4、逻辑非，只能出现在表达式最左侧，作用域是整个表达式
 *    5、如果表达式不能被解析成三种类型中的任意一种，那么这个表达式一定不符合语法
 */

package logexp

import "encoding/json"

type ExpressionType int32 // 表达式类型
const (
	ExpressionType_Meta ExpressionType = 0 // 元表达式（内部不包含'|'和'&'符号）
	ExpressionType_Or   ExpressionType = 1 // “或”表达式
	ExpressionType_And  ExpressionType = 2 // “且”表达式
)

type IExpression interface {
	/*
	 * 返回取非的标记
	 */
	GetIsNegative() bool
	/*
	 * 逆转取非的标记
	 */
	ReverseIsNegative()
	/*
	 * 返回表达式类型
	 */
	GetType() ExpressionType

	/*
	 * 返回所有子表达式
	 */
	GetExps() []IExpression

	/*
	 * 判断逻辑表达式是否匹配给定的文本
	 */
	Match(text string) bool
}

// 子表达式
type SubExp struct {
	IsNegative   bool   // 是否取非
	Exp          []rune // 表达式文本
	IsMeta       bool   // 是否元表达式 （不包含'|'和'&'符号)
	BracketStack []int  // 存储子表达式括号位置的栈，入栈的是左括号在表达式中的坐标
}

func (e *SubExp) Trim() bool {
	leftIdx := 0               // 修整后的表达式的最左侧字符在原表达式的位置
	rightIdx := len(e.Exp) - 1 // 修整后的表达式最右侧字符在原表达式的位置

	// 从左往右每次遇到'!'，isNegative都取反
	for leftIdx < len(e.Exp) && e.Exp[leftIdx] == '!' {
		e.IsNegative = !e.IsNegative
		leftIdx++
	}
	// 从左往右，对连续遇到的'('符号，如果跟表达式连续的右侧的')'匹配，那么认为该表达式被一对括号括起来了，可以把括号去掉
	// 这个判断可以借助bracketStack栈来完成，因为对于最外层的配对括号，它们的左括号的坐标一定顺序存储在bracketStack的前几个元素；阅读一下bracketStack的入栈出栈逻辑可以更好地理解
	for leftIdx < len(e.Exp) && e.Exp[leftIdx] == '(' && e.Exp[rightIdx] == ')' && e.BracketStack[len(e.Exp)-1-rightIdx] == leftIdx {
		leftIdx++
		rightIdx--
	}
	// 如果表达式没有被大括号整体括起来，而且它不是元表达式，那么'!'符号不能作用于整个表达式
	if leftIdx > 0 && rightIdx == len(e.Exp)-1 && !e.IsMeta {
		e.IsNegative = false
		leftIdx = 0
	}
	if leftIdx > 0 || rightIdx < len(e.Exp)-1 {
		e.Exp = e.Exp[leftIdx : rightIdx+1]
		e.BracketStack = make([]int, len(e.Exp)) // 重置
		return true
	} else {
		return false
	}
}

type LogExp struct {
	expression IExpression
}

func (e *LogExp) Match(text string) bool {
	return e.expression.Match(text)
}

func (e *LogExp) ToJson() string {
	buf, _ := json.Marshal(e.expression)
	return string(buf)
}

func Compile(exp string) (*LogExp, *CstError) {
	expression, cerr := NewExpressionOr([]rune(exp), false, 0)
	if cerr != nil {
		return nil, cerr
	}
	return &LogExp{expression: expression}, nil
}
