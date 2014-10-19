package pcalc

import (
	"fmt"
)

type Expression struct {
	children            []*Expression
	privateCall         func(*Expression, *map[string]*Expression) *Expression
	privateValue        func(*Expression) number
	privateDescription  func(*Expression) string
	privateListUnknowns func(*Expression) []string
}

//Basic *Expression operator functions
func AddExpressions(a, b *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(add, "%s+%s", a, b)
}

func SubtractExpressions(a, b *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(subtract, "%s-%s", a, b)
}

func MultiplyExpressions(a, b *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(multiply, "%s*%s", a, b)
}

func DivideExpressions(a, b *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(divide, "%s/%s", a, b)
}

func RaiseExpressionToPower(a, b *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(raiseToPower, "%s^(%s)", a, b)
}

func DecimalLogarithmOfExpression(a *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(log, "log(%s)", a)
}

func NaturalLogarithmOfExpression(a *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(ln, "ln(%s)", a)
}

func ParenthesisEnclosedExpression(a *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(nil, "(%s)", a)
}

func SignInvertedExpression(a *Expression) *Expression {
	return NewExpressionWithValueFuncDescFormAndArgs(subtract, "(-%[2]s)", NewExpressionWithConstant(MakeFraction(0, 1)), a)
}

//Exported *Expression Methods
func (exp *Expression) Call(m *map[string]*Expression) *Expression {
	return exp.privateCall(exp, m)
}

func (exp *Expression) Value() number {
	return exp.privateValue(exp)
}

func (exp *Expression) Description() string {
	return exp.privateDescription(exp)
}

func (exp *Expression) ListUnknowns() []string {
	return exp.privateListUnknowns(exp)
}

//NewExpression
func NewExpressionWithConstant(con number) *Expression {
	e := new(Expression)
	e.children = nil
	e.privateCall = func(*Expression, *map[string]*Expression) *Expression {
		return e
	}
	e.privateValue = func(*Expression) number {
		return con
	}
	e.privateDescription = func(*Expression) string {
		return fmt.Sprint(con.Description())
	}
	e.privateListUnknowns = func(*Expression) []string {
		return make([]string, 0)
	}
	return e
}

func NewExpressionWithUnknown(unk string) *Expression {
	e := new(Expression)
	e.children = nil
	e.privateCall = func(ne *Expression, m *map[string]*Expression) *Expression {
		for key, value := range *m {
			if key == unk {
				if value != nil {
					return value.Call(m)
				} else {
					return e
				}
			}
		}
		return e
	}
	e.privateValue = func(*Expression) number {
		return nil
	}
	e.privateDescription = func(*Expression) string {
		return unk
	}
	e.privateListUnknowns = func(*Expression) []string {
		return []string{unk}
	}
	return e
}

func NewExpressionWithValueFuncDescFormAndArgs(valueFunc func(...number) number, descFormat string, args ...*Expression) *Expression {
	if len(args) < 1 {
		fmt.Println("Can not make combined expression without subexpression")
		return nil
	}
	for _, v := range args {
		if v == nil {
			fmt.Println("Can not make combined expression with subexpression == nil")
			return nil
		}
	}
	e := new(Expression)
	e.children = args
	e.privateCall = func(exp *Expression, m *map[string]*Expression) *Expression {
		ne := new(Expression)
		ne.children = make([]*Expression, len(exp.children))
		for i, v := range exp.children {
			ne.children[i] = v.Call(m)
		}
		ne.privateCall = exp.privateCall
		ne.privateValue = exp.privateValue
		ne.privateDescription = exp.privateDescription
		ne.privateListUnknowns = exp.privateListUnknowns
		return ne
	}
	if valueFunc == nil {
		e.privateValue = func(exp *Expression) number {
			if val := exp.children[0].Value(); val != nil {
				return val
			} else {
				return nil
			}
		}
	} else {
		e.privateValue = func(exp *Expression) number {
			cvals := make([]number, len(exp.children))
			for i, v := range exp.children {
				cvals[i] = v.Value()
				if cvals[i] == nil {
					return nil
				}
			}
			return valueFunc(cvals...)
		}
	}
	e.privateDescription = func(exp *Expression) string {
		cdesc := make([]interface{}, len(exp.children))
		for i, v := range exp.children {
			cdesc[i] = v.Description()
		}
		return fmt.Sprintf(descFormat, cdesc...)
	}
	e.privateListUnknowns = func(exp *Expression) []string {
		unks := make([]string, 0)
		for _, c := range exp.children {
			cu := c.ListUnknowns()
			for _, a := range cu {
				isdupe := false
				for _, b := range unks {
					if a == b {
						isdupe = true
						break
					}
				}
				if !isdupe {
					unks = append(unks, a)
				}
			}
		}
		return unks
	}
	return e
}
