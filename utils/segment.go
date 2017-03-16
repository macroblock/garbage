package utils

// TSegment -
type TSegment struct {
	A int
	B int
}

// Length -
func (seg *TSegment) Length() int {
	return seg.B - seg.A
}

// IntersectSegment -
func IntersectSegment(seg1, seg2 TSegment) TSegment {
	return TSegment{A: Max(seg1.A, seg2.A), B: Min(seg1.B, seg2.B)}
}
