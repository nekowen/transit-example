package main

import (
	"database/sql"
	"log"
	"sort"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func buildData(f string) StationNodes {
	db, err := sql.Open("sqlite3", f)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(
		`SELECT station_cd, station_g_cd, line_cd, (SELECT line_name FROM lines WHERE lines.line_cd = stations.line_cd limit 1) as line_name, station_name, e_sort, lat, lon, (SELECT station_cd1 FROM joins WHERE station_cd2 = stations.station_cd AND joins.line_cd = stations.line_cd limit 1) as prev, (SELECT station_cd2 FROM joins WHERE station_cd1 = stations.station_cd AND joins.line_cd = stations.line_cd limit 1) as next FROM stations`,
	)

	defer rows.Close()

	_gcd := make(map[string][]int)
	_nodes := make(map[int]*StationNode)

	for rows.Next() {
		var stationCd int
		var stationGCd int
		var lineCd int
		var lineName string
		var stationName string
		var eSort int
		var lat float64
		var lng float64
		var prev *int
		var next *int

		if err := rows.Scan(&stationCd, &stationGCd, &lineCd, &lineName, &stationName, &eSort, &lat, &lng, &prev, &next); err != nil {
			log.Fatal("failed to scan row", err)
		}

		cdlen := len(strconv.Itoa(stationCd))
		isShinkansen := (cdlen == 6)

		s := &Station{stationCd, stationGCd, lineCd, lineName, stationName, isShinkansen, eSort, lat, lng, prev, next}
		_nodes[stationCd] = &StationNode{s, StationNodes{}, None, 0.0, 0.0, nil}
		grk := s.GroupKey()
		_gcd[grk] = append(_gcd[grk], stationCd)
	}

	for sid, node := range _nodes {
		var edges StationNodes
		grk := node.station.GroupKey()

		if node.station.PrevCd != nil {
			//	前の駅
			edges = append(edges, _nodes[*node.station.PrevCd])
		}
		if node.station.NextCd != nil {
			//	次の駅
			edges = append(edges, _nodes[*node.station.NextCd])
		}

		//	グループに属する駅を全て追加
		for _, gid := range _gcd[grk] {
			if sid == gid {
				continue
			}
			edges = append(edges, _nodes[gid])
		}

		node.edges = edges
	}

	i := 0
	nodes := make(StationNodes, len(_nodes))
	for _, val := range _nodes {
		nodes[i] = val
		i++
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].station.StationCd > nodes[j].station.StationCd
	})
	return nodes
}

func buildRoutes(end *StationNode) Stations {
	node := end
	routes := Stations{}
	prev := -1
	for {
		if prev != node.station.StationGCd {
			routes = append(routes, node.station)
		}
		prev = node.station.StationGCd
		node = node.parent
		if node == nil {
			break
		}
	}

	return routes.Reverse()
}
