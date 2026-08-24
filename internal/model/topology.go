package model

// CoilTopology 描述磁体线圈的拓扑（分段 + 邻接边，边视为无向）。
type CoilTopology struct {
	ID       string
	Name     string
	Segments []TopoSegment
	Edges    []TopoEdge
}

// TopoSegment 为线圈上的一个分段。
type TopoSegment struct {
	ID       string
	Label    string
	Position float64 // 沿磁体的位置（米）
}

// TopoEdge 为两个分段间的邻接边。
type TopoEdge struct {
	ID       string
	FromID   string
	ToID     string
	Distance float64 // 相邻分段间距离（米）
}

// Neighbors 返回与 segID 相邻的边（无向，双向展开）。
func (t *CoilTopology) Neighbors(segID string) []TopoEdge {
	var out []TopoEdge
	for _, e := range t.Edges {
		if e.FromID == segID {
			out = append(out, e)
		} else if e.ToID == segID {
			out = append(out, TopoEdge{ID: e.ID, FromID: e.ToID, ToID: e.FromID, Distance: e.Distance})
		}
	}
	return out
}

// PathDistance 返回两分段间沿拓扑的最短距离（米），不可达返回 -1。
func (t *CoilTopology) PathDistance(from, to string) float64 {
	if from == to {
		return 0
	}
	dist := map[string]float64{from: 0}
	queue := []string{from}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, e := range t.Neighbors(cur) {
			nd := dist[cur] + e.Distance
			if d, ok := dist[e.ToID]; !ok || nd < d {
				dist[e.ToID] = nd
				queue = append(queue, e.ToID)
			}
		}
	}
	if d, ok := dist[to]; ok {
		return d
	}
	return -1
}
