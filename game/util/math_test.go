package util

import (
	"math"
	"testing"
)

func TestV2Norm(t *testing.T) {
	x, y := V2Norm(3, 4)
	if x != 0.6 || y != 0.8 {
		t.Errorf("Expected (0.6, 0.8), got (%v, %v)", x, y)
	}
}

func TestV2Len(t *testing.T) {
	length := V2Len(3, 4)
	if length != 5 {
		t.Errorf("Expected 5, got %v", length)
	}
}

func TestV2Angle(t *testing.T) {
	angle := V2Angle(0, 1)
	expected := math.Pi / 2
	if angle != expected {
		t.Errorf("Expected %v, got %v", expected, angle)
	}
}

func TestVec2Rotate(t *testing.T) {
	v := Vec2{X: 1, Y: 0}
	rotated := v.Rotate(math.Pi / 2)
	expected := Vec2{X: 0, Y: 1}
	if rotated != expected {
		t.Errorf("Expected %v, got %v", expected, rotated)
	}
}

func TestVec2Mul(t *testing.T) {
	v := Vec2{X: 1, Y: 2}
	scaled := v.Mul(2)
	expected := Vec2{X: 2, Y: 4}
	if scaled != expected {
		t.Errorf("Expected %v, got %v", expected, scaled)
	}
}

func TestVec2Add(t *testing.T) {
	v1 := Vec2{X: 1, Y: 2}
	v2 := Vec2{X: 3, Y: 4}
	sum := v1.Add(v2)
	expected := Vec2{X: 4, Y: 6}
	if sum != expected {
		t.Errorf("Expected %v, got %v", expected, sum)
	}
}

func TestLerpAngle(t *testing.T) {
	a := 0.0
	b := math.Pi
	tVal := 0.5
	result := LerpAngle(a, b, tVal)
	expected := math.Pi / 2
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestShortestAngleDirection(t *testing.T) {
	tests := []struct {
		a, b, expected float64
	}{
		{0, math.Pi, math.Pi},
		{math.Pi, 0, -math.Pi},
		{0, math.Pi / 2, math.Pi / 2},
		{math.Pi / 2, 0, -math.Pi / 2},
		{0, -math.Pi / 2, -math.Pi / 2},
		{-math.Pi / 2, 0, math.Pi / 2},
		{math.Pi / 4, -math.Pi / 4, -math.Pi / 2},
		{math.Pi / 4, math.Pi / 4, 0},
		{math.Pi / 4, 3 * math.Pi / 4, math.Pi / 2},
		{3 * math.Pi / 2, -math.Pi / 2, 0},
	}

	for _, tt := range tests {
		result := ShortestAngleDirection(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("For angles %v and %v, expected %v, got %v", tt.a, tt.b, tt.expected, result)
		}
	}
}
