package scripts

// Функция проверки пересечения отрезка с прямоугольником
func isIntersectingSegmentRectangle(x1, y1, x2, y2, xmin, ymin, xmax, ymax float64) bool {
	// Функция проверки, находится ли точка внутри прямоугольника
	isInside := func(x, y float64) bool {
		return x >= xmin && x <= xmax && y >= ymin && y <= ymax
	}

	// Если одна из точек отрезка внутри прямоугольника, то пересечение есть
	if isInside(x1, y1) || isInside(x2, y2) {
		return true
	}

	// Функция проверки пересечения двух отрезков
	doIntersect := func(ax, ay, bx, by, cx, cy, dx, dy float64) bool {
		cross := func(px, py, qx, qy, rx, ry float64) float64 {
			return (qx-px)*(ry-py) - (qy-py)*(rx-px)
		}
		s1 := cross(ax, ay, bx, by, cx, cy)
		s2 := cross(ax, ay, bx, by, dx, dy)
		s3 := cross(cx, cy, dx, dy, ax, ay)
		s4 := cross(cx, cy, dx, dy, bx, by)
		return (s1*s2 < 0 && s3*s4 < 0)
	}

	// Проверяем пересечение отрезка с каждой стороной прямоугольника
	return doIntersect(x1, y1, x2, y2, xmin, ymin, xmax, ymin) ||
		doIntersect(x1, y1, x2, y2, xmax, ymin, xmax, ymax) ||
		doIntersect(x1, y1, x2, y2, xmax, ymax, xmin, ymax) ||
		doIntersect(x1, y1, x2, y2, xmin, ymax, xmin, ymin)
}
