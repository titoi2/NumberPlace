/*


*/
package main

import (
	"fmt"
	"github.com/visualfc/go-ui/ui"
	"image/color"
	"strconv"
)

const (
	CELL_WIDTH = 24

	quiz_01 = "095060700020000060003000050360409000000000003080007026908053210100900800042000000"
)

const (
	CELL_KIND_NORMAL = iota
	CELL_KIND_FIXED
	CELL_KIND_ERROR
)

const LBL_SELECTED = "Selected :"

type CellInfo struct {
	kind int
	num  int
}

var placeCells [9][9]CellInfo // [Y][X]

var cursorx, cursory int

type blockPosTable struct {
	y [3]int
	x [3]int
}

var blockPosTables [9]blockPosTable
var selectedNum int

func main() {
	placeCells[0][0].kind = 0
	placeCells[0][0].num = 1

	cursorx = -1
	cursory = -1
	for i := 0; i < 9; i++ {
		for j := 0; j < 3; j++ {
			blockPosTables[i].y[j] = j + (i/3)*3
			blockPosTables[i].x[j] = j + (i%3)*3
		}
	}

	qix := 0
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if quiz_01[qix] == '0' {
				placeCells[y][x].kind = CELL_KIND_NORMAL
				placeCells[y][x].num = 0
			} else {
				placeCells[y][x].kind = CELL_KIND_FIXED
				placeCells[y][x].num = int(quiz_01[qix]) - '0'
			}
			qix++
		}
	}
	ui.Main(ui_main)
}

func ui_main() {
	w := ui.NewWidget()
	defer w.Close()
	w.SetWindowTitle("Number Place")

	place := ui.NewWidget()
	place.SetMinimumSize(ui.Sz(CELL_WIDTH*10, CELL_WIDTH*10))

	img := ui.NewImageWithSize(CELL_WIDTH*9, CELL_WIDTH*9)
	defer img.Close()

	imgPainter := ui.NewPainterWithImage(img)
	imgPainter.InitFrom(place)
	defer imgPainter.Close()

	place.OnPaintEvent(func(ev *ui.PaintEvent) {
		drawPlace(ev, place)
	})

	vbox := ui.NewVBoxLayout()

	hbox0 := ui.NewHBoxLayout()
	hbox := ui.NewHBoxLayout()
	hboxSelNum := ui.NewHBoxLayout()
	hbox7 := ui.NewHBoxLayout()
	hbox4 := ui.NewHBoxLayout()
	hbox1 := ui.NewHBoxLayout()

	vbox.AddWidget(place)
	vbox.AddLayout(hboxSelNum)
	vbox.AddLayout(hbox7)
	vbox.AddLayout(hbox4)
	vbox.AddLayout(hbox1)
	vbox.AddLayout(hbox0)
	vbox.AddLayout(hbox)

	lblSelectedNum := ui.NewLabelWithText(LBL_SELECTED)
	hboxSelNum.AddWidget(lblSelectedNum)

	Btn1 := ui.NewButtonWithText("1")
	Btn2 := ui.NewButtonWithText("2")
	Btn3 := ui.NewButtonWithText("3")
	Btn4 := ui.NewButtonWithText("4")
	Btn5 := ui.NewButtonWithText("5")
	Btn6 := ui.NewButtonWithText("6")
	Btn7 := ui.NewButtonWithText("7")
	Btn8 := ui.NewButtonWithText("8")
	Btn9 := ui.NewButtonWithText("9")
	BtnCLR := ui.NewButtonWithText("CLR")
	BtnReset := ui.NewButtonWithText("Reset")
	exitBtn := ui.NewButtonWithText("Exit")
	hbox7.AddWidget(Btn7)
	hbox7.AddWidget(Btn8)
	hbox7.AddWidget(Btn9)
	hbox7.AddStretch(0)
	hbox4.AddWidget(Btn4)
	hbox4.AddWidget(Btn5)
	hbox4.AddWidget(Btn6)
	hbox4.AddStretch(0)
	hbox1.AddWidget(Btn1)
	hbox1.AddWidget(Btn2)
	hbox1.AddWidget(Btn3)
	hbox1.AddStretch(0)

	hbox.AddSpacing(10)

	hbox0.AddWidget(BtnCLR)
	hbox0.AddStretch(0)
	hbox.AddStretch(0)
	hbox.AddWidget(BtnReset)
	hbox.AddWidget(exitBtn)

	selectdNumUpdate := func() {
		s := ""
		if selectedNum > 0 {
			s = strconv.Itoa(selectedNum)
		}
		lblSelectedNum.SetText(LBL_SELECTED + s)
	}

	place.OnMousePressEvent(func(ev *ui.MouseEvent) {
		pos := ev.Pos()
		fmt.Println("pos x", pos.X)
		fmt.Println("pos y", pos.Y)
		var ix int = pos.X / CELL_WIDTH
		var iy int = pos.Y / CELL_WIDTH
		if ix < 9 && iy < 9 {
			cursorx = ix
			cursory = iy
			if placeCells[iy][ix].kind != CELL_KIND_FIXED && placeCells[iy][ix].num == 0 {
				placeCells[iy][ix].num = selectedNum
			} else {
				selectedNum = placeCells[iy][ix].num
			}
			placeCheck()
			place.Update()
			selectdNumUpdate()
		}
	})

	numBtnFunc := func(n int) {
		if cursory == -1 {
			return
		}
		if placeCells[cursory][cursorx].kind != CELL_KIND_FIXED {
			placeCells[cursory][cursorx].num = n
			selectedNum = n
			placeCheck()
			place.Update()
			selectdNumUpdate()
		}
	}

	Btn1.OnClicked(func() {
		numBtnFunc(1)
	})
	Btn2.OnClicked(func() {
		numBtnFunc(2)
	})
	Btn3.OnClicked(func() {
		numBtnFunc(3)
	})
	Btn4.OnClicked(func() {
		numBtnFunc(4)
	})
	Btn5.OnClicked(func() {
		numBtnFunc(5)
	})
	Btn6.OnClicked(func() {
		numBtnFunc(6)
	})
	Btn7.OnClicked(func() {
		numBtnFunc(7)
	})
	Btn8.OnClicked(func() {
		numBtnFunc(8)
	})
	Btn9.OnClicked(func() {
		numBtnFunc(9)
	})

	BtnCLR.OnClicked(func() {
		numBtnFunc(0)
	})

	BtnReset.OnClicked(func() {
		selectedNum = 0
		selectdNumUpdate()
		for iy := 0; iy < 9; iy++ {
			for ix := 0; ix < 9; ix++ {
				if placeCells[iy][ix].kind != CELL_KIND_FIXED {
					placeCells[iy][ix].kind = CELL_KIND_NORMAL
					placeCells[iy][ix].num = 0
				}
			}
		}
		place.Update()
	})

	exitBtn.OnClicked(func() {
		w.Close()
	})

	w.SetLayout(vbox)
	w.Show()

	w.OnCloseEvent(func(e *ui.CloseEvent) {
	})

	ui.Run()
}

func placeCheck() {
	// ERROR CLEAR
	for iy := 0; iy < 9; iy++ {
		for ix := 0; ix < 9; ix++ {
			if placeCells[iy][ix].kind != CELL_KIND_FIXED {
				placeCells[iy][ix].kind = CELL_KIND_NORMAL
			}
		}
	}
	errorFound := false
	unfinished := false

	// COLUMN CHECK
	for iy := 0; iy < 9; iy++ {
		nums := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		for ix := 0; ix < 9; ix++ {
			nums[placeCells[iy][ix].num]++
		}
		for ix := 0; ix < 9; ix++ {
			if placeCells[iy][ix].kind != CELL_KIND_FIXED {
				if placeCells[iy][ix].num > 0 && nums[placeCells[iy][ix].num] > 1 {
					placeCells[iy][ix].kind = CELL_KIND_ERROR
					errorFound = true
					fmt.Printf("COLUMN ERROR y[%d] x[%d] n[%d]\n", iy, ix, placeCells[iy][ix].num)
				}
				if placeCells[iy][ix].num == 0 {
					// white cell found
					unfinished = true
				}
			}
		}
	}

	// LINE CHECK
	for ix := 0; ix < 9; ix++ {
		nums := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		for iy := 0; iy < 9; iy++ {
			nums[placeCells[iy][ix].num]++
		}
		for iy := 0; iy < 9; iy++ {
			if placeCells[iy][ix].num > 0 && nums[placeCells[iy][ix].num] > 1 {
				if placeCells[iy][ix].kind != CELL_KIND_FIXED {
					placeCells[iy][ix].kind = CELL_KIND_ERROR
					errorFound = true
					fmt.Printf("LINE ERROR y[%d] x[%d] n[%d]\n", iy, ix, placeCells[iy][ix].num)
				}
			}
		}
	}

	// BLOCK CHECK
	for bi := 0; bi < 9; bi++ {
		nums := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		for iy := 0; iy < 3; iy++ {
			for ix := 0; ix < 3; ix++ {
				y := blockPosTables[bi].y[iy]
				x := blockPosTables[bi].x[ix]

				nums[placeCells[y][x].num]++
			}
		}
		for iy := 0; iy < 3; iy++ {
			for ix := 0; ix < 3; ix++ {
				y := blockPosTables[bi].y[iy]
				x := blockPosTables[bi].x[ix]
				if placeCells[y][x].num > 0 && nums[placeCells[y][x].num] > 1 {
					if placeCells[y][x].kind != CELL_KIND_FIXED {
						placeCells[y][x].kind = CELL_KIND_ERROR
						errorFound = true
						fmt.Printf("BLOCK[%d] ERROR y[%d] x[%d] n[%d]\n", bi, y, x, placeCells[y][x].num)
					}
				}
			}
		}
	}

	if !errorFound && !unfinished {
		// GAME CLEAR
		fmt.Println("GAME CLEAR")
	}
}

func drawPlace(e *ui.PaintEvent, w *ui.Widget) {
	p := ui.NewPainter()
	defer p.Close()
	p.Begin(w)

	pen := ui.NewPen()
	pen.SetColor(color.RGBA{0, 0, 0, 0})

	ft := p.Font()
	ft.SetPixelSize(CELL_WIDTH * 8 / 10)
	p.SetFont(ft)
	for i := 0; i <= len(placeCells); i++ {
		if i%3 == 0 {
			pen.SetWidth(3)
		} else {
			pen.SetWidth(1)
		}
		p.SetPen(pen)
		p.DrawLine(ui.Pt(0, i*CELL_WIDTH), ui.Pt(9*CELL_WIDTH, i*CELL_WIDTH))
		p.DrawLine(ui.Pt(i*CELL_WIDTH, 0), ui.Pt(i*CELL_WIDTH, 9*CELL_WIDTH))
	}

	var pt ui.Point
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if placeCells[y][x].num != 0 {
				pt.X = x*CELL_WIDTH + 4 // TODO display position adjustment
				pt.Y = (y+1)*CELL_WIDTH - 4
				if placeCells[y][x].kind == CELL_KIND_ERROR {
					pen.SetColor(color.RGBA{255, 0, 0, 0})
				} else if placeCells[y][x].kind == CELL_KIND_FIXED {
					pen.SetColor(color.RGBA{0, 0, 0, 0})
				} else {
					pen.SetColor(color.RGBA{0, 0, 255, 0})
				}
				p.SetPen(pen)
				p.DrawText(pt, fmt.Sprintf("%c", placeCells[y][x].num+'0'))
			}
		}
	}
	if cursorx != -1 {
		var rc ui.Rect
		pen.SetColor(color.RGBA{200, 0, 0, 0})
		pen.SetWidth(3)

		p.SetPen(pen)
		rc.X = cursorx * CELL_WIDTH
		rc.Y = cursory * CELL_WIDTH
		rc.Width = CELL_WIDTH
		rc.Height = CELL_WIDTH
		p.DrawRect(rc)
	}
	p.End()
}
