package main

import "fmt"

type Station struct {
	StationCd    int
	StationGCd   int
	LineCd       int
	LineName     string
	StationName  string
	IsShinkansen bool
	Esort        int
	Lat          float64
	Lng          float64
	PrevCd       *int
	NextCd       *int
}
type Stations []*Station

func (s Stations) Reverse() Stations {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (s *Station) GroupKey() string {
	//	StationGCdとStationNameを組み合わせたものをグループコードにする
	return fmt.Sprintf("%d%s", s.StationGCd, s.StationName)
}

type NodeState int

const (
	None NodeState = iota
	Open
	Closed
)

type StationNode struct {
	station *Station
	edges   StationNodes
	state   NodeState
	cost    float64 // 親ノードから現在の駅までかかるコスト
	hcost   float64 // 現在の駅から到着駅までの推定コスト
	parent  *StationNode
}
type StationNodes []*StationNode

func (s *StationNode) Score() float64 {
	return s.cost + s.hcost
}

func (s StationNodes) Delete(i int) StationNodes {
	s = append(s[:i], s[i+1:]...)
	n := make(StationNodes, len(s))
	copy(n, s)
	return n
}
