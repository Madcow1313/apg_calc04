package controller

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	Drawer "webCalc/drawer"
	Model "webCalc/model"
	Stack "webCalc/stack"
)

type Controller struct {
	expression,
	expressionBack,
	postfixExpression,
	lastResult,
	currentDir string
	priorities             map[string]int
	history                map[int64]string
	currentHistoryPosition int64
	logFile                *os.File
	xMax, xMin             float64
}

func (c *Controller) HandleMessage(message string) {
	if message == " unary_minus " {
		c.expression += message
		c.expressionBack += "-"
	} else if message == " unary_plus " {
		c.expression += message
		c.expressionBack += "+"
	} else if message == " sin " || message == " cos " || message == " tan " ||
		message == " asin " || message == " acos " || message == " atan " ||
		message == " ln " || message == " log " || message == " sqrt " || message == " mod " {
		c.expression += message + " ( "
		c.expressionBack += strings.TrimSuffix(message, " ") + "( "
	} else if message == "clear" {
		c.expression = ""
		c.expressionBack = ""
	} else {
		c.expression += message
		c.expressionBack += message
	}
}

func (c *Controller) HandleHistory(message string) {
	if message == "last" {
		history := c.history
		lh := int64(len(history))
		if lh != 0 {
			c.currentHistoryPosition = lh
			c.expressionBack = strings.ReplaceAll(history[lh], "unary_plus", "+")
			c.expressionBack = strings.ReplaceAll(history[lh], "unary_minus", "-")
			c.expression = history[lh]
		}
	} else if message == "prev" {
		history := c.history
		lh := int64(len(history))
		if lh != 0 && c.currentHistoryPosition != 0 {
			c.currentHistoryPosition -= 1
			c.expressionBack = strings.ReplaceAll(history[c.currentHistoryPosition], "unary_plus", "+")
			c.expressionBack = strings.ReplaceAll(history[c.currentHistoryPosition], "unary_minus", "-")
			c.expression = history[c.currentHistoryPosition]
		}
	} else if message == "next" {
		history := c.history
		lh := int64(len(history))
		if lh != 0 && c.currentHistoryPosition != int64(len(history)) {
			c.currentHistoryPosition += 1
			c.expressionBack = strings.ReplaceAll(history[c.currentHistoryPosition], "unary_plus", "+")
			c.expressionBack = strings.ReplaceAll(history[c.currentHistoryPosition], "unary_minus", "-")
			c.expression = history[c.currentHistoryPosition]
		}
	} else if message == "history_clear" {
		c.logFile.Truncate(0)
		for k := range c.history {
			delete(c.history, k)
		}
	}
}

func openLogFile(string) (*os.File, error) { return nil, nil }

func (c *Controller) Init() {
	c.fillPriorities()
	logFile, err := openLogFile(c.currentDir)
	if err == nil {
		c.logFile = logFile
		c.history = make(map[int64]string)
		c.readHistory()
	}
}

func (c *Controller) readHistory() {
	buf, err := os.ReadFile(c.logFile.Name())
	if err == nil {
		currentHistorySlice := strings.Split(string(buf), "\n")
		for _, str := range currentHistorySlice {
			firstHalf, secondHalf, _ := strings.Cut(str, " ")
			if len(firstHalf) != 0 {
				number, _ := strconv.ParseInt(firstHalf, 10, 64)
				c.history[number] = secondHalf
			}
		}
	}
}

func (c *Controller) writeLog(expression string) {
	if c.lastResult != "error" {
		var err error
		result := strconv.FormatFloat(float64(len(c.history)+1), 'G', 30, 64) + " " + expression + "\n"
		_, err = c.logFile.Write([]byte(result))
		if err != nil {
			fmt.Println(err)
		} else {
			c.history[int64(len(c.history))+1] = expression
		}
	}
}

func (c *Controller) invokeModel(expression string) {
	var m Model.Model
	c.postfixExpression, _ = c.infixToPostfix(strings.Fields(c.expression))
	m.PostfixExpression = c.postfixExpression
	if !m.StartComputeRPN() {
		c.lastResult = "error"
	} else {
		c.lastResult = strconv.FormatFloat(m.Result, 'f', -1, 64)
	}
	c.writeLog(expression)
}

func (c *Controller) invokeGraphic(xMin, xMax, yMin, yMax float64) (string, error) {
	var d Drawer.Drawer

	d.XMax, d.XMin, d.YMin, d.YMax = xMax, xMin, yMin, yMax
	d.PostfixExpression = c.postfixExpression
	d.CurrentDir = c.currentDir
	fileName, err := d.Draw()

	return fileName, err
}

func (c *Controller) fillPriorities() {
	c.priorities = make(map[string]int)
	c.priorities["sin"] = 5
	c.priorities["cos"] = 5
	c.priorities["tan"] = 5
	c.priorities["asin"] = 5
	c.priorities["acos"] = 5
	c.priorities["atan"] = 5
	c.priorities["ln"] = 5
	c.priorities["log"] = 5

	c.priorities["^"] = 4
	c.priorities["sqrt"] = 4
	c.priorities["unary_minus"] = 4
	c.priorities["unary_plus"] = 4

	c.priorities["*"] = 3
	c.priorities["/"] = 3

	c.priorities["+"] = 2
	c.priorities["-"] = 2

	c.priorities["mod"] = 1
}

func (c *Controller) getPriority(s string) int {
	result, exist := c.priorities[s]
	if !exist {
		return -1
	}
	return result
}

func (c *Controller) infixToPostfix(s []string) (string, error) {
	var postfixString string
	var stack Stack.Stack[string]
	for _, value := range s {
		prior := c.getPriority(value)
		if prior < 0 && value != "(" && value != ")" {
			postfixString = postfixString + " " + value
		} else if value == "(" {
			stack.Push(value)
		} else if value == ")" {
			for {
				str, _ := stack.Top()
				if stack.IsEmpty() || str == "(" {
					break
				}
				postfixString = postfixString + " " + str
				stack.Pop()
			}
			stack.Pop()
		} else {
			for {
				str, _ := stack.Top()
				if stack.IsEmpty() || !(prior <= c.getPriority(str)) {
					break
				}
				postfixString = postfixString + " " + str
				stack.Pop()
			}
			stack.Push(value)
		}
	}
	for {
		if stack.IsEmpty() {
			break
		}
		str, _ := stack.Pop()
		postfixString = postfixString + " " + str
	}
	return postfixString, nil
}
