package Drawer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	Model "webCalc/model"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

type Drawer struct {
	Expression             string
	CurrentDir             string
	XMin, XMax, YMin, YMax float64
}

func (d *Drawer) initPlotterFunction(x float64) float64 {
	var m Model.Model
	X := strconv.FormatFloat(x, 'f', -1, 64)
	m.Expression = strings.ReplaceAll(d.Expression, "X", X)
	m.StartComputeRPN()
	if m.Result == math.Inf(1) || m.Result == math.Inf(-1) || math.IsNaN(m.Result) {
		return math.NaN()
	}
	return m.Result
}

func (d *Drawer) Draw() (string, error) {
	var xys plotter.XYs
	fmt.Println(d.XMax, d.XMin)
	for x := d.XMin; x <= d.XMax; x += 0.1 {
		y := d.initPlotterFunction(x)
		if math.IsNaN(y) || y > d.YMax || y < d.YMin {
			continue
		}
		xys = append(xys, plotter.XY{X: x, Y: y})
	}
	ps, _, _ := plotter.NewLinePoints(xys)
	p := plot.New()
	p.Title.Text = "Plot of function"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.X.Label.Position = draw.PosRight
	p.Y.Label.Position = draw.PosTop
	p.X.Min = d.XMin
	p.X.Max = d.XMax
	p.Y.Min = d.YMin
	p.Y.Max = d.YMax
	p.Add(ps)
	err := p.Save(720, 480, "test.png")
	if err != nil {
		fmt.Printf("%v", err)
		return "", err
	}
	return "test.png", nil
}
