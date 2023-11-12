package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	Drawer "webCalc/drawer"
	Model "webCalc/model"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Expression,
	ExpressionBack,
	postfixExpression,
	LastResult,
	currentDir string
	history                map[int64]string
	currentHistoryPosition int64
	logFile                *os.File
	XMax, XMin             float64
}

func (Controller *Controller) GET(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		if Controller.LastResult == "error" {
			Controller.LastResult = ""
		}
		Controller.Expression = Controller.LastResult
		fmt.Println(Controller.LastResult)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"result": Controller.LastResult,
		})
	})
	router.GET("/help.html", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "help.html", nil)
	})
	router.GET("/graph_window.html", func(ctx *gin.Context) {
		_, err := Controller.InvokeGraphic(Controller.XMin, Controller.XMax, Controller.XMin, Controller.XMax)
		if err == nil {
			ctx.HTML(http.StatusOK, "graph_window.html", nil)
		}
	})
}

func (Controller *Controller) POST(router *gin.Engine) {
	router.POST("/", func(ctx *gin.Context) {
		q, _ := ctx.GetQuery("body")
		q = strings.Trim(q, "'")
		if q == " plus " {
			q = " + "
		} else if q == " divide " {
			q = " / "
		}
		if strings.HasPrefix(q, "x= ") {
			q = strings.TrimPrefix(q, "x= ")
			Controller.Expression = strings.ReplaceAll(Controller.Expression, "X", q)
		} else if strings.HasPrefix(q, "xy_min= ") {
			q = strings.TrimPrefix(q, "xy_min= ")
			Controller.XMin, _ = strconv.ParseFloat(q, 64)
		} else if strings.HasPrefix(q, "xy_max= ") {
			q = strings.TrimPrefix(q, "xy_max= ")
			Controller.XMax, _ = strconv.ParseFloat(q, 64)
		} else if q == "clear" {
			Controller.Expression = ""
		} else if q == "=" {
			Controller.InvokeModel()
		} else {
			Controller.HandleMessage(q)
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"result": Controller.LastResult,
		})
	})
}

func (c *Controller) HandleMessage(message string) {
	switch message {
	case "button_next":
		c.HandleHistory("next")
		break
	case "button_last":
		c.HandleHistory("last")
		break
	case "button_prev":
		c.HandleHistory("prev")
		break
	case "button_history_clear":
		c.HandleHistory("history_clear")
		break
	default:
		if message == " unary_minus " {
			c.Expression += message
			c.ExpressionBack += "-"
		} else if message == " unary_plus " {
			c.Expression += message
			c.ExpressionBack += "+"
		} else if message == " sin " || message == " cos " || message == " tan " ||
			message == " asin " || message == " acos " || message == " atan " ||
			message == " ln " || message == " log " || message == " sqrt " || message == " mod " {
			c.Expression += message
			c.ExpressionBack += strings.TrimSuffix(message, " ") + "( "
		} else if message == "clear" {
			c.Expression = ""
			c.ExpressionBack = ""
		} else {
			c.Expression += message
			c.ExpressionBack += message
		}
	}
}

func (c *Controller) HandleHistory(message string) {
	history := c.history
	lh := int64(len(history))
	c.ExpressionBack = strings.ReplaceAll(history[lh], "unary_plus", "+")
	c.ExpressionBack = strings.ReplaceAll(history[lh], "unary_minus", "-")
	if message == "last" {
		if lh != 0 {
			c.currentHistoryPosition = lh
			c.Expression = history[lh]
			c.LastResult = history[lh]
		}
	} else if message == "prev" {
		history := c.history
		lh := int64(len(history))
		if lh != 0 && c.currentHistoryPosition != 0 {
			c.currentHistoryPosition -= 1
			c.Expression = history[c.currentHistoryPosition]
			c.LastResult = history[c.currentHistoryPosition]
		}
	} else if message == "next" {
		history := c.history
		lh := int64(len(history))
		if lh != 0 && c.currentHistoryPosition != int64(len(history)) {
			c.currentHistoryPosition += 1
			c.Expression = history[c.currentHistoryPosition]
			c.LastResult = history[c.currentHistoryPosition]
		}
	} else if message == "history_clear" {
		c.logFile.Truncate(0)
		for k := range c.history {
			delete(c.history, k)
		}
	}
}
func openLogFile(currentDir string) (*os.File, error) {
	file, err := os.OpenFile(currentDir+"/calc_log", os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		file, err = os.Create("calc_log")
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

func (c *Controller) Init() {
	logFile, err := openLogFile(c.currentDir)
	if err == nil {
		c.logFile = logFile
		c.history = make(map[int64]string)
		c.readHistory()
	}
	router := gin.Default()

	router.LoadHTMLFiles("index.html", "help.html", "graph_window.html")
	router.StaticFile("./jsScript/index.js", "jsScript/index.js")
	router.StaticFile("./test.png", "test.png")

	c.GET(router)
	c.POST(router)
	router.Run()
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
	if c.LastResult != "error" {
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

func (c *Controller) InvokeModel() {
	var m Model.Model
	m.Expression = c.Expression
	m.FillPriorities()
	if !m.StartComputeRPN() {
		c.LastResult = "error"
	} else {
		c.LastResult = strconv.FormatFloat(m.Result, 'f', -1, 64)
	}
	c.writeLog(c.Expression)
}

func (c *Controller) InvokeGraphic(xMin, xMax, yMin, yMax float64) (string, error) {
	var d Drawer.Drawer

	d.XMax, d.XMin, d.YMin, d.YMax = xMax, xMin, yMin, yMax
	d.Expression = c.Expression
	d.CurrentDir = c.currentDir
	fileName, err := d.Draw()

	return fileName, err
}
