package dao

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

type Dao struct {
	session *gocql.Session
}

type NetworkSpeedResult struct {
	Date     time.Time
	Upload   string
	Download string
	Ping     string
}

func NewDao(keySpace string, ips ...string) (*Dao, error){
	cluster := gocql.NewCluster(ips...)
	cluster.Keyspace = keySpace
	session, err := cluster.CreateSession()
	return &Dao{
		session: session,
	}, err
}

func (d *Dao) Insert(result NetworkSpeedResult) error {
	return d.session.Query(fmt.Sprintf("INSERT INTO NETWORK_SPEED (DATE, INTERVAL, DOWNLOAD, UPLOAD, PING) VALUES ('%s', '%s', '%s', '%s', '%s');",
		result.Date.Format("20060102"), result.Date.Format("2006-01-02 15:04:05"), result.Download, result.Upload, result.Ping)).Exec()
}