package main

import (
	"log"
	"math"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nekowen/transit-example/utils"
)

func findRoute(nodes StationNodes, origin *StationNode, dest *StationNode) Stations {
	if origin == nil || dest == nil {
		log.Fatal("origin or destination params are nil")
	}

	origin.state = Open
	origin.cost = 0
	origin.hcost = estimateTime(origin.station, dest.station, 60) / 2

	openNodes := StationNodes{origin}

	for {
		if len(openNodes) == 0 {
			// 経路が見つからなかった
			break
		}

		minscore := math.MaxFloat64
		mincost := math.MaxFloat64
		nextIndex := -1
		f := false
		for i, node := range openNodes {
			if node == dest {
				f = true
				break
			}
			score := node.Score()
			if score > minscore || score == minscore && node.cost >= mincost {
				continue
			}

			minscore = score
			mincost = node.cost
			nextIndex = i
		}

		if f {
			return buildRoutes(dest)
		}

		p := openNodes[nextIndex]
		openNodes = openNodes.Delete(nextIndex)

		for _, edge := range p.edges {
			if edge.state == None {
				addCost := 0.0
				if p.station.GroupKey() == edge.station.GroupKey() && p.station.LineCd != edge.station.LineCd {
					//	乗り換え可能な駅。乗り換えには1分かかるとする
					addCost += 60.0
				}
				speed := 60.0
				if edge.station.IsShinkansen {
					//	新幹線は早い
					speed = 100.0
				}

				if p.parent != nil {
					edge.cost = p.parent.cost + estimateTime(p.station, edge.station, speed) + addCost
				} else {
					//	開始駅構内の乗り換えコストは考えない
					edge.cost = estimateTime(p.station, edge.station, speed)
				}
				edge.hcost = estimateTime(edge.station, dest.station, speed)
				edge.state = Open
				edge.parent = p

				openNodes = append(openNodes, edge)
			}
		}
		p.state = Closed
	}

	return Stations{}
}

func estimateTime(from *Station, to *Station, speed float64) float64 {
	return utils.Distance(from.Lat, from.Lng, to.Lat, to.Lng) / 60 / 1000 * 60 * 60
}

func findStationByName(a StationNodes, s string) *StationNode {
	for _, st := range a {
		if st.station.StationName == s {
			return st
		}
	}
	return nil
}

func findStationByID(a StationNodes, n int) *StationNode {
	for _, st := range a {
		if st.station.StationCd == n {
			return st
		}
	}
	return nil
}
