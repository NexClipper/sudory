package model

import "time"

type ReqClient struct {
	AgetnID   string `json:"agent_id"`
	ClusterID uint64 `json:"cluster_id"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
}

type Client struct {
	ID        uint64    `xorm:"pk autoincr 'id'"`
	AgentID   string    `xorm:"agent_id"`
	ClusterID uint64    `xorm:"cluster_id"`
	Active    bool      `xorm:"active"`
	IP        string    `xorm:"ip"`
	Port      int       `xorm:"port"`
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
}

func (m *Client) GetType() string {
	return "CLIENT"
}
