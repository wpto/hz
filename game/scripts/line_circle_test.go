package scripts

import "testing"

// Тесты для функции
func TestIsIntersectingSegmentCircle(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2, cx, cy, r float64
		expected                  bool
	}{
		{1, 1, 4, 4, 3, 3, 2, true},  // Отрезок пересекает окружность
		{5, 5, 6, 6, 3, 3, 2, false}, // Отрезок полностью вне окружности
		{2, 2, 3, 3, 3, 3, 2, true},  // Отрезок внутри окружности
		{1, 1, 2, 2, 3, 3, 2, true},  // Отрезок касается окружности
		{0, 0, 1, 1, 3, 3, 1, false}, // Отрезок рядом, но не пересекает
	}

	for _, tt := range tests {
		result := isIntersectingSegmentCircle(tt.x1, tt.y1, tt.x2, tt.y2, tt.cx, tt.cy, tt.r)
		if result != tt.expected {
			t.Errorf("For segment (%.1f, %.1f) -> (%.1f, %.1f) and circle (%.1f, %.1f, %.1f), expected %v but got %v",
				tt.x1, tt.y1, tt.x2, tt.y2, tt.cx, tt.cy, tt.r, tt.expected, result)
		}
	}
}
