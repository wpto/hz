package physics

import "testing"

// func main() {
// 	x1, y1, x2, y2 := 1.0, 1.0, 4.0, 4.0
// 	xmin, ymin, xmax, ymax := 2.0, 2.0, 5.0, 5.0
// 	fmt.Println(isIntersectingSegmentRectangle(x1, y1, x2, y2, xmin, ymin, xmax, ymax)) // true
// }

// Тесты для функции
func TestIsIntersectingSegmentRectangle(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2, xmin, ymin, xmax, ymax float64
		expected                               bool
	}{
		{1, 1, 4, 4, 2, 2, 5, 5, true},  // Отрезок пересекает прямоугольник
		{6, 6, 7, 7, 2, 2, 5, 5, false}, // Отрезок полностью вне прямоугольника
		{3, 3, 4, 4, 2, 2, 5, 5, true},  // Отрезок внутри прямоугольника
		{0, 0, 2, 2, 2, 2, 5, 5, true},  // Отрезок касается угла прямоугольника
		{5, 1, 6, 2, 2, 2, 5, 5, false}, // Отрезок рядом, но не пересекает
		{4, 1, 6, 3, 2, 2, 5, 5, false}, // Отрезок пересекает угол прямоугольника
		{3.5, 1, 5.5, 3, 2, 2, 5, 5, true},
	}

	for _, tt := range tests {
		result := isIntersectingSegmentRectangle(tt.x1, tt.y1, tt.x2, tt.y2, tt.xmin, tt.ymin, tt.xmax, tt.ymax)
		if result != tt.expected {
			t.Errorf("For segment (%.1f, %.1f) -> (%.1f, %.1f) and rect (%.1f, %.1f, %.1f, %.1f), expected %v but got %v",
				tt.x1, tt.y1, tt.x2, tt.y2, tt.xmin, tt.ymin, tt.xmax, tt.ymax, tt.expected, result)
		}
	}
}
