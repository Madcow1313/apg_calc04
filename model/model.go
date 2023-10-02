package Model

import (
	"math"
	"strconv"
	"strings"

	stack "webCalc/stack"
)

type Model struct {
	PostfixExpression string
	Result            float64
}

func cycleComputing(f func(float64, float64) float64, stack *stack.Stack[float64]) float64 {
	var Result float64
	Result, status := stack.Pop()
	if !status {
		return 0
	}
	val, _ := stack.Pop()
	Result = f(val, Result)
	return Result
}

func ComputeFunc(f func(float64) float64, stack *stack.Stack[float64]) float64 {
	var Result float64
	Result, status := stack.Pop()
	if !status {
		return 0
	}
	Result = f(Result)
	return Result
}

func Compute(oper string, stack *stack.Stack[float64]) float64 {
	var Result float64
	switch oper {
	case "+":
		return cycleComputing(func(f1, f2 float64) float64 { return f1 + f2 }, stack)
	case "-":
		return cycleComputing(func(f1, f2 float64) float64 { return f1 - f2 }, stack)
	case "*":
		return cycleComputing(func(f1, f2 float64) float64 { return f1 * f2 }, stack)
	case "/":
		return cycleComputing(func(f1, f2 float64) float64 { return f1 / f2 }, stack)
	case "^":
		return cycleComputing(func(f1, f2 float64) float64 { return math.Pow(f1, f2) }, stack)
	case "mod":
		return cycleComputing(func(f1, f2 float64) float64 { return math.Mod(f1, f2) }, stack)
	case "sin":
		return ComputeFunc(func(f1 float64) float64 { return math.Sin(f1) }, stack)
	case "cos":
		return ComputeFunc(func(f1 float64) float64 { return math.Cos(f1) }, stack)
	case "tan":
		return ComputeFunc(func(f1 float64) float64 { return math.Tan(f1) }, stack)
	case "asin":
		return ComputeFunc(func(f1 float64) float64 { return math.Asin(f1) }, stack)
	case "acos":
		return ComputeFunc(func(f1 float64) float64 { return math.Acos(f1) }, stack)
	case "atan":
		return ComputeFunc(func(f1 float64) float64 { return math.Atan(f1) }, stack)
	case "log":
		return ComputeFunc(func(f1 float64) float64 { return math.Log10(f1) }, stack)
	case "ln":
		return ComputeFunc(func(f1 float64) float64 { return math.Log(f1) }, stack)
	case "sqrt":
		return ComputeFunc(func(f1 float64) float64 { return math.Sqrt(f1) }, stack)
	case "unary_minus":
		return ComputeFunc(func(f1 float64) float64 { return -f1 }, stack)
	case "unary_plus":
		return ComputeFunc(func(f1 float64) float64 { return +f1 }, stack)
	}

	return Result
}

func getNumber(n string) (float64, error) {
	number, err := strconv.ParseFloat(n, 64)
	if err != nil {
		return number, err
	}
	return number, nil
}

func (m *Model) StartComputeRPN() bool {
	var Result float64
	var stack stack.Stack[float64]
	input := strings.Fields(m.PostfixExpression)
	for _, value := range input {
		n, err := getNumber(value)
		if err != nil {
			Result = Compute(value, &stack)
			stack.Push(Result)
		} else {
			stack.Push(n)
		}
	}
	Result, _ = stack.Pop()
	if !stack.IsEmpty() {
		return false
	}
	m.Result = Result
	return true
}
