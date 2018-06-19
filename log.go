package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"strings"
)

// Field names are converted to camel case by default, no need to add
// `db:"first_name"`, if you want to disable a filed add `db:"-"` tag.
type EventLogException struct {
	LogId          string `db:"logid"`
	Hostname       string `db:"hostname"`
	Logger         string `db:"logger"`
	Tcode          string `db:"tcode"`
	Tid            string `db:"tid"`
	ExceptionNames string `db:"exceptionnames"`
	ContextLogs    string `db:"contextlogs"`
	Timestamp      string `db:"timestamp"`
}

func findLog(logid string) (*EventLogException, error) {
	cluster := gocql.NewCluster(strings.Fields(*cassandraCluster)...)
	cluster.Keyspace = "blackcat"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	stmt, names := qb.Select("event_log_exception").Where(qb.Eq("logId")).ToCql()
	var p EventLogException

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"logId": logid})
	if err := q.GetRelease(&p); err != nil {
		if err == gocql.ErrNotFound {
			fmt.Println("Not Found")
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &p, nil
}
