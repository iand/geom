package geom

type Path2 struct {
	Points []Point2 // waypoints
	dirs   []Vec2
	dists  []float32
	length float32
}

func NewPath2(pts []Point2) *Path2 {
	p := &Path2{
		Points: pts,
		dirs:   make([]Vec2, len(pts)-1),
		dists:  make([]float32, len(pts)-1),
	}

	for i := 0; i < len(pts)-1; i++ {
		p.dirs[i] = pts[i+1].Sub(pts[i])
		p.dists[i] = p.dirs[i].Len()
		p.length += p.dists[i]
		p.dirs[i] = p.dirs[i].Normalize()
	}

	return p
}

func (p *Path2) PositionAlong(d float32) Ray2 {
	if d <= 0 {
		return Ray2{
			Origin:    p.Points[0],
			Direction: p.dirs[0],
		}
	} else if d >= 1.0 {
		return Ray2{
			Origin:    p.Points[len(p.Points)-1],
			Direction: p.dirs[len(p.dirs)-1],
		}
	}

	l := d * p.length
	for i := 0; i < len(p.dists); i++ {
		if l <= p.dists[i] {
			return Ray2{
				Origin:    p.Points[i].Add(p.dirs[i].Mul(l)),
				Direction: p.dirs[i],
			}
		}
		l -= p.dists[i]
	}

	return Ray2{
		Origin:    p.Points[len(p.Points)-1],
		Direction: p.dirs[len(p.dirs)-1],
	}
}
