package logexp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExpression(t *testing.T) {
	type Case struct {
		Exp          string // 表达式
		Valid        bool   // 表达式是否合法
		CompiledJson string // 表达式编译后的结果的json形式
	}
	testCases := []Case{
		{
			Exp:          "hello",
			Valid:        true,
			CompiledJson: `{"type":0,"is_negative":false,"keyword":"hello"}`,
		},
		{
			Exp:          "hello|hi",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"hi"}]}`,
		},
		{
			Exp:          "hello&hi",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"hi"}]}`,
		},
		{
			Exp:          "hello&hi|wow",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":2,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "hello|hi&wow",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":2,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hi"},{"type":0,"is_negative":false,"keyword":"wow"}]}]}`,
		},
		{
			Exp:          "hello|(hi&wow)",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":2,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hi"},{"type":0,"is_negative":false,"keyword":"wow"}]}]}`,
		},
		{
			Exp:          "(hello|hi)&wow",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "((hello|we)|hi)&wow",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"we"},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "((hello&we)|hi)&wow",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":false,"expressions":[{"type":2,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"we"}]},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "((hello&we)|hi)|wow",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":2,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"we"}]},{"type":0,"is_negative":false,"keyword":"hi"},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "((hello&!we)|hi)&wow",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":false,"expressions":[{"type":2,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":true,"keyword":"we"}]},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "(!(hello&!we)|hi)&wow",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":false,"expressions":[{"type":2,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":true,"keyword":"we"}]},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "!(!(hello&!we)|hi)&wow",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":true,"expressions":[{"type":2,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":true,"keyword":"we"}]},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "!(!(!(hello&!we)|hi)&wow)",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":true,"expressions":[{"type":1,"is_negative":true,"expressions":[{"type":2,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":true,"keyword":"we"}]},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "!!(!(!(hello&!we)|hi)&wow)",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":true,"expressions":[{"type":2,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":true,"keyword":"we"}]},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "!!(!(!(hello&!!we)|hi)&wow)",
			Valid:        true,
			CompiledJson: `{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":true,"expressions":[{"type":2,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"we"}]},{"type":0,"is_negative":false,"keyword":"hi"}]},{"type":0,"is_negative":false,"keyword":"wow"}]}`,
		},
		{
			Exp:          "!!(!(!(hello&!!we&中国)|hi|深圳)&wow|空 格)",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":2,"is_negative":false,"expressions":[{"type":1,"is_negative":true,"expressions":[{"type":2,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"we"},{"type":0,"is_negative":false,"keyword":"中国"}]},{"type":0,"is_negative":false,"keyword":"hi"},{"type":0,"is_negative":false,"keyword":"深圳"}]},{"type":0,"is_negative":false,"keyword":"wow"}]},{"type":0,"is_negative":false,"keyword":"空 格"}]}`,
		},

		{
			Exp:   "hello&(hi)wow",
			Valid: false,
		},
		{
			Exp:   "hello&(hi|)wow",
			Valid: false,
		},
		{
			Exp:   "(hello&(hi)&wow",
			Valid: false,
		},

		{
			Exp:          "((hello|(a|b))|aa)|hiwow",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"hello"},{"type":0,"is_negative":false,"keyword":"a"},{"type":0,"is_negative":false,"keyword":"b"},{"type":0,"is_negative":false,"keyword":"aa"},{"type":0,"is_negative":false,"keyword":"hiwow"}]}`,
		},
		{
			Exp:          "(a|b)",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":0,"is_negative":false,"keyword":"a"},{"type":0,"is_negative":false,"keyword":"b"}]}`,
		},
		{
			Exp:          "!(a|b)",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"a"},{"type":0,"is_negative":false,"keyword":"b"}]}`,
		},
		{
			Exp:          "!(a|(b|c))",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"a"},{"type":0,"is_negative":false,"keyword":"b"},{"type":0,"is_negative":false,"keyword":"c"}]}`,
		},
		{
			Exp:          "!(a|!(b|c))",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"a"},{"type":1,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"b"},{"type":0,"is_negative":false,"keyword":"c"}]}]}`,
		},
		{
			Exp:          "(!a|!(b|c))",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":false,"expressions":[{"type":0,"is_negative":true,"keyword":"a"},{"type":1,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"b"},{"type":0,"is_negative":false,"keyword":"c"}]}]}`,
		},
		{
			Exp:          "!(a|!(b&c))",
			Valid:        true,
			CompiledJson: `{"type":1,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"a"},{"type":2,"is_negative":true,"expressions":[{"type":0,"is_negative":false,"keyword":"b"},{"type":0,"is_negative":false,"keyword":"c"}]}]}`,
		},
	}
	for idx, cas := range testCases {
		expression, cerr := Compile(cas.Exp)
		if cas.Valid {
			assert.Equal(t, (*CstError)(nil), cerr, fmt.Sprintf("case %v: %v  err:%v", idx, cas.Exp, cerr))
			assert.Equal(t, cas.CompiledJson, expression.ToJson(), fmt.Sprintf("case %v: %v", idx, cas.Exp))
		} else {
			if !assert.NotEqual(t, (*CstError)(nil), cerr, fmt.Sprintf("case %v: %v", idx, cas.Exp)) {
				fmt.Printf("compiled: %v", expression.ToJson())
			}
		}
	}
}

func TestMatch(t *testing.T) {
	type Case struct {
		Exp   string
		Text  string
		Match bool
	}
	testCases := []Case{
		{
			Exp:   "hello|hi",
			Text:  "hello world",
			Match: true,
		},
		{
			Exp:   "hello|hi",
			Text:  "helllo world",
			Match: false,
		},
		{
			Exp:   "(!(hello&!we)|hi)&wow",
			Text:  "hello world",
			Match: false,
		},
		{
			Exp:   "(!(hello&!we)|hi)&wow",
			Text:  "we hello world wow",
			Match: true,
		},
	}
	for idx, cas := range testCases {
		expression, cerr := Compile(cas.Exp)
		if cerr != nil {
			t.Error(cerr)
		} else {
			assert.Equal(t, cas.Match, expression.Match(cas.Text), fmt.Sprintf("case %v: %v", idx, cas.Exp))
		}
	}
}

func BenchmarkNew(b *testing.B) {
	for idx := 0; idx < b.N; idx++ {
		exp, _ := Compile("hello|hi|we")
		// s1 := "hello"
		// s2 := "hi"
		// s3 := "we"
		for i := 0; i <= 10000; i++ {
			exp.Match("wewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewehello")
			// strings.Contains("wewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewehello", s1)
			// strings.Contains("wewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewehello", s2)
			// strings.Contains("wewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewewehello", s3)
			// NewExpression("(!(hello&!we)|hi)&wow")
			// NewExpression("hellohellohellohellohellohellohellohellohellohellohellohello|hihellohellohellohellohellohellohellohellohellohellohello")
		}
	}
}
