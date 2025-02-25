package physics

import (
	"math"
)

// Проверяет пересечение отрезка с окружностью
func isIntersectingSegmentCircle(x1, y1, x2, y2, cx, cy, r float64) bool {
	// Вектор AB (отрезок)
	dx, dy := x2-x1, y2-y1

	// Вектор AC (из начала отрезка к центру окружности)
	fx, fy := x1-cx, y1-cy

	// Коэффициенты квадратного уравнения для нахождения пересечения
	a := dx*dx + dy*dy
	b := 2 * (fx*dx + fy*dy)
	c := fx*fx + fy*fy - r*r

	// Если обе точки внутри окружности, отрезок внутри окружности
	if math.Sqrt((x1-cx)*(x1-cx)+(y1-cy)*(y1-cy)) < r &&
		math.Sqrt((x2-cx)*(x2-cx)+(y2-cy)*(y2-cy)) < r {
		return true
	}

	// Дискриминант квадратного уравнения
	d := b*b - 4*a*c
	if d < 0 {
		return false // Нет пересечений
	}

	// Вычисляем точки пересечения
	t1 := (-b - math.Sqrt(d)) / (2 * a)
	t2 := (-b + math.Sqrt(d)) / (2 * a)

	// Проверяем, попадают ли точки пересечения в границы отрезка (0 ≤ t ≤ 1)
	return (t1 >= 0 && t1 <= 1) || (t2 >= 0 && t2 <= 1)
}
