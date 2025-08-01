package main

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
	"strconv"

	"github.com/fogleman/gg"
)

func drawPieChart(total, filled int) image.Image {
	const S = 128
	dc := gg.NewContext(S, S)
	dc.SetColor(color.White)
	dc.Clear()

	cx, cy := float64(S)/2, float64(S)/2
	radius := float64(S)/2 - 4

	angleStep := 2 * gg.Radians(180) / float64(total)
	startAngle := -gg.Radians(90)

	// Draw slices
	for i := 0; i < total; i++ {
		a0 := startAngle + angleStep*float64(i)
		a1 := a0 + angleStep

		if i < filled {
			dc.SetRGB(0.0, 0.4, 0.8) // blue
		} else {
			dc.SetRGB(0.8, 0.8, 0.8) // gray
		}
		dc.MoveTo(cx, cy)
		dc.LineTo(cx+radius*gg.Cos(a0), cy+radius*gg.Sin(a0))
		dc.Arc(cx, cy, radius, a0, a1)
		dc.ClosePath()
		dc.FillPreserve()
		dc.SetRGB(0.0, 0.0, 0.0)
		dc.Stroke()
	}
	return dc.Image()
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	total, _ := strconv.Atoi(r.URL.Query().Get("total"))
	if total < 1 {
		total = 12
	}
	filled, _ := strconv.Atoi(r.URL.Query().Get("filled"))
	if filled < 0 {
		filled = 0
	}
	if filled > total {
		filled = total
	}

	img := drawPieChart(total, filled)

	w.Header().Set("Content-Type", "image/png")
	err := gg.SavePNGWriter(w, img)
	if err != nil {
		http.Error(w, "Error generating image", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/image", imageHandler)
	fmt.Println("Listening on http://localhost:8080/image?total=12&filled=5")
	http.ListenAndServe(":8080", nil)
}
