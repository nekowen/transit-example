package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	originStationName = flag.String("o", "大宮", "出発駅名")
	destStationName   = flag.String("d", "東京", "到着駅名")
)

func main() {
	flag.Parse()

	nodes := buildData("./stations.db")
	origin := findStationByName(nodes, *originStationName)
	if origin == nil {
		log.Fatalf("出発駅が見つかりません: %s", *originStationName)
	}
	destination := findStationByName(nodes, *destStationName)
	if destination == nil {
		log.Fatalf("到着駅が見つかりません: %s", *destStationName)
	}
	routes := findRoute(nodes, origin, destination)

	rlen := len(routes)
	if rlen == 0 {
		log.Fatalln("ルートが見つかりませんでした")
	} else {
		prevl := -1
		for i, n := range routes {
			suffix := "→"
			if rlen-1 == i {
				suffix = "\n"
			}
			if prevl != n.LineCd {
				fmt.Printf("%s(%sへ乗り換え)%s", n.StationName, n.LineName, suffix)
			} else {
				fmt.Printf("%s%s", n.StationName, suffix)
			}
			prevl = n.LineCd
		}
	}
}
