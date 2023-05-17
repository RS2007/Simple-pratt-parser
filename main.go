package main

/*
  SPECS: Very minimal pratt parser: only parses single digit addition subtraction multiplication and division
*/
import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

type TokenType int

type TokenArray []Token

func (c TokenArray) Reverse() {
	s := 0
	e := len(c) - 1
	for e-s >= 1 {
		temp := c[s]
		c[s] = c[e]
		c[e] = temp
		s += 1
		e -= 1
	}
}

const (
	Integer TokenType = iota
	Operand
	Prefix
)

type Token interface {
	getTokenType() TokenType
	getExpressionValue() string
}

type Expression interface {
	getExpressionValue() string
}

type IntegerToken struct {
	value int32
}

type OperatorToken struct {
	literal string
}

type PrefixToken struct {
	literal string
}

func (i PrefixToken) getTokenType() TokenType {
	return Prefix
}

func (i IntegerToken) getTokenType() TokenType {
	return Integer
}

func (i OperatorToken) getTokenType() TokenType {
	return Operand
}

func (i IntegerToken) getExpressionValue() string {
	return strconv.Itoa(int(i.value))
}

func (i OperatorToken) getExpressionValue() string {
	return i.literal
}

func (i PrefixToken) getExpressionValue() string {
	return i.literal
}

type Lexer struct {
	tokens TokenArray
}

type InfixExpression struct {
	lhs Expression
	rhs Expression
	op  string
}

type PrefixExpression struct {
	op  string
	rhs Expression
}

func (i PrefixExpression) getExpressionValue() string {
	return i.op + strconv.Itoa(evalExpression(i.rhs))
}

func (i InfixExpression) getExpressionValue() string {
	returnString := "("
	// for _, expr := range i.rhs {
	// 	returnString += " ("
	// 	expr_infix, ok := (expr).(InfixExpression)
	// 	if !ok {
	// 		expr_int, ok := (expr).(IntegerToken)
	// 		if !ok {
	// 			panic("Undesired behaviour")
	// 		}
	// 		returnString += " " + strconv.Itoa(int(expr_int.value))
	// 	}
	// 	returnString += " " + expr_infix.getExpressionValue()
	// }

	return returnString
}

func New(input string) *Lexer {
	var buffer bytes.Buffer
	for _, char := range input {
		buffer.WriteByte(byte(char))
	}
	charArray := buffer.Bytes()
	tokenArray := make(TokenArray, 0)
	for _, c := range charArray {
		if c == ' ' || c == '\r' || c == '\t' || c == '\n' {
			continue
		} else if c >= '0' && c <= '9' {
			intValue, err := strconv.Atoi(string(c))
			if err != nil {
				panic("This should not have happened")
			}
			tokenArray = append(tokenArray, IntegerToken{
				value: int32(intValue),
			})
		} else if c == '+' || c == '-' || c == '*' || c == '/' {
			tokenArray = append(tokenArray, OperatorToken{
				literal: string(c),
			})
		}
	}
	tokenArray.Reverse()
	l := &Lexer{
		tokens: tokenArray,
	}
	return l
}

func (l *Lexer) next() Token {
	if len(l.tokens) == 0 {
		return nil
	}
	lastToken := l.tokens[len(l.tokens)-1]
	l.tokens = l.tokens[:len(l.tokens)-1]
	return lastToken
}

func (l *Lexer) peek() Token {
	if len(l.tokens) < 1 {
		return nil
	}
	return l.tokens[len(l.tokens)-1]
}

func parse(l *Lexer, min_bp int) Expression {
	operatorBindingPowerMap := map[string][]int{
		"+": {1, 2},
		"-": {1, 2},
		"*": {3, 4},
		"/": {3, 4},
	}

	prefixBindingPowerMap := map[string][]int{
		"+": {0, 5},
		"-": {0, 5},
	}
	var lhs Expression

	lhsExpr := l.next()
	switch lhsExpr.getTokenType() {
	case Integer:
		lhs, _ = lhsExpr.(IntegerToken)
		break
	case Operand:
		lhs_prefix_token, _ := lhsExpr.(OperatorToken)
		r_bp := prefixBindingPowerMap[lhs_prefix_token.literal][1]
		rhs := parse(l, r_bp)
		lhs = PrefixExpression{
			op:  lhs_prefix_token.literal,
			rhs: rhs,
		}
		break
	}
	for {
		if l.peek() == nil {
			break
		}
		op, ok := l.peek().(OperatorToken)
		if !ok {
			panic("Integer should be followed by an operand")
		}
		l_bp, r_bp := operatorBindingPowerMap[op.literal][0], operatorBindingPowerMap[op.literal][1]
		if l_bp < min_bp {
			break
		}
		l.next()

		rhs := parse(l, r_bp)
		lhs = InfixExpression{
			lhs: lhs,
			rhs: rhs,
			op:  op.literal,
		}
	}
	return lhs
}

func evalPrefix(e PrefixExpression) int {
	switch e.op {
	case "+":
		return (1) * evalExpression(e.rhs)
	case "-":
		return (-1) * evalExpression(e.rhs)
	}
	panic("Should not reach here")
}

func evalExpression(e Expression) int {
	operationFunctionMap := map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}
	int_value, ok := e.(IntegerToken)
	if ok {
		return int(int_value.value)
	}
	prefix_expr, ok := e.(PrefixExpression)

	if ok {
		return evalPrefix(prefix_expr)
	}
	expression, ok := e.(InfixExpression)
	if !ok {
		panic("Undesired behaviour")
	}
	return operationFunctionMap[expression.op](evalExpression(expression.lhs), evalExpression(expression.rhs))

}

func main() {
	var a string
	in := bufio.NewReader(os.Stdin)
	a, err := in.ReadString('\n')
	if err != nil {
		panic("Error reading string")
	}
	lexer := New(a)
	parsed := parse(lexer, 0)
	parsed_infix, ok := (parsed).(InfixExpression)
	if !ok {
		parsed_int, ok := (parsed).(IntegerToken)
		if !ok {
			parsed_prefix, ok := (parsed).(PrefixExpression)
			if !ok {
				panic("Undesired behaviour")
			}
			fmt.Println(evalExpression(parsed_prefix))
			return
		}
		fmt.Println(parsed_int.getExpressionValue())
		return
	}
	fmt.Println(evalExpression(parsed_infix))
}
