package pcalc

import (
	"fmt"
)

type Expression struct {
	children     []*Expression
	call         func(*Expression, *map[string]*Expression) *Expression
	value        func(*Expression) number
	description  func(*Expression) string
	listUnknowns func(*Expression) []string
}

func addExpressions(a, b *Expression) *Expression {
	return newExpressionWithValueFuncDescFormAndArgs(add, "%s+%s", a, b)
}

//Shell Methods
func (exp *Expression) mCall(m *map[string]*Expression) *Expression {
	return exp.call(exp, m)
}

func (exp *Expression) mValue() number {
	return exp.value(exp)
}

func (exp *Expression) mDescription() string {
	return exp.description(exp)
}

func (exp *Expression) mListUnknowns() []string {
	return exp.listUnknowns(exp)
}

//newExpression
func newExpressionWithConstant(con number) *Expression {
	e := new(Expression)
	e.children = nil
	e.call = func(*Expression, *map[string]*Expression) *Expression {
		return e
	}
	e.value = func(*Expression) number {
		return con
	}
	e.description = func(*Expression) string {
		return fmt.Sprint(con.description())
	}
	e.listUnknowns = func(*Expression) []string {
		return make([]string, 0)
	}
	return e
}

func newExpressionWithUnknown(unk string) *Expression {
	e := new(Expression)
	e.children = nil
	e.call = func(ne *Expression, m *map[string]*Expression) *Expression {
		for key, value := range *m {
			if key == unk {
				if value != nil {
					return value.mCall(m)
				} else {
					return e
				}
			}
		}
		return e
	}
	e.value = func(*Expression) number {
		return nil
	}
	e.description = func(*Expression) string {
		return unk
	}
	e.listUnknowns = func(*Expression) []string {
		return []string{unk}
	}
	return e
}

func newExpressionWithValueFuncDescFormAndArgs(valueFunc func(...number) number, descFormat string, args ...*Expression) *Expression {
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
	e.call = func(exp *Expression, m *map[string]*Expression) *Expression {
		ne := new(Expression)
		ne.children = make([]*Expression, len(exp.children))
		for i, v := range exp.children {
			ne.children[i] = v.mCall(m)
		}
		ne.call = exp.call
		ne.value = exp.value
		ne.description = exp.description
		ne.listUnknowns = exp.listUnknowns
		return ne
	}
	if valueFunc == nil {
		e.value = func(exp *Expression) number {
			if val := exp.children[0].mValue(); val != nil{
				return val
			} else {
				return nil
			}
		}
	} else {
		e.value = func(exp *Expression) number {
			cvals := make([]number, len(exp.children))
			for i, v := range exp.children {
				cvals[i] = v.mValue()
				if cvals[i] == nil {
					return nil
				}
			}
			return valueFunc(cvals...)
		}
	}
	e.description = func(exp *Expression) string {
		cdesc := make([]interface{}, len(exp.children))
		for i, v := range exp.children {
			cdesc[i] = v.mDescription()
		}
		return fmt.Sprintf(descFormat, cdesc...)
	}
	e.listUnknowns = func(exp *Expression) []string {
		unks := make([]string, 0)
		for _, c := range exp.children {
			cu := c.mListUnknowns()
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
