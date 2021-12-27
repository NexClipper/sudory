package model

import "time"

type ReqCluster struct {
	Name string `json:"name"`
}

type Cluster struct {
	ID      string    `xorm:"pk autoincr 'id'"`
	Name    string    `xorm:"name"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func (m *Cluster) GetType() string {
	return "CLUSTER"
}
