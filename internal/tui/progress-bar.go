// a simple progress-bar
//
//	func main() {
//		fmt.Printf("simple progress-bar\n")
//		pbar := ProgressBar{size: 40, chunk: 2, BarChar: '#'}
//
//		for i := 0; i < 100; i++ {
//			if done := pbar.Next(); !done {
//				time.Sleep(time.Millisecond * 10)
//			}
//		}
//		pbar.Done()
//	}
package tui

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type ProgressBar struct {
	index    int
	size     int
	stepSize int
	BarChar  rune
	Out      string
}

func main() {
	fmt.Printf("simple progress bar\n")
	pbar := ProgressBar{size: 10, stepSize: 1, BarChar: '#'}

	for i := 0; i < 10; i++ {
		if i == 8 {
			pbar.MovePercent(100) // simulate that we are done
		}
		if done := pbar.MovePercent(10); !done {
			time.Sleep(time.Millisecond * 50)
		}
		fmt.Printf(pbar.Out)
	}
	fmt.Println("finished")
}

func (p *ProgressBar) Next() bool {
	return p.NextStep(1)
}

func (p *ProgressBar) NextStep(step int) bool {
	if step == 0 {
		panic("0 is invalid step argument")
	}
	p.index = p.index + p.stepSize*step
	if p.index > p.size {
		p.index = p.size
	}
	p.Out = fmt.Sprintf("[%s%s]\r", strings.Repeat(string(p.BarChar), p.index), strings.Repeat(" ", p.size-p.index))
	return p.index >= p.size
}

func (p *ProgressBar) MovePercent(percent int) bool {
	var i int = int(math.Ceil(float64(p.size / p.stepSize / 100.00)))
	return p.NextStep(i)
}
