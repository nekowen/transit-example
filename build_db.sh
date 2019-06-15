#!/bin/bash
if [ $# -ne 1 ];then
  cat <<_EOS_
Usage
  $0 csvVersion(e.g. 20190405)
_EOS_
exit -1
fi

if [ ! -x $(which sqlite3) ];then
  echo "sqlite3 is not installed"
  exit -1
fi

cd $(dirname $0)
date=$1
dbFileName=stations.db

sqlite3 -separator , ${dbFileName} ".import join${date}.csv joins"
sqlite3 -separator , ${dbFileName} ".import line${date}free.csv lines"
sqlite3 -separator , ${dbFileName} ".import station${date}free.csv stations"

sqlite3 ${dbFileName} "CREATE INDEX joins_station_cd1Index ON joins(station_cd1)"
sqlite3 ${dbFileName} "CREATE INDEX joins_station_cd2Index ON joins(station_cd2)"
sqlite3 ${dbFileName} "CREATE INDEX lines_line_cdIndex ON lines(line_cd)"