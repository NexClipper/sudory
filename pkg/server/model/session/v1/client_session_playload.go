package v1

import "time"

type ClientSessionPlayload struct {
	Exp          time.Time `json:"exp,omitempty"`           //expiration_time
	Iat          time.Time `json:"iat,omitempty"`           //issued_at_time
	Uuid         string    `json:"uuid,omitempty"`          //token_uuid
	ClusterUuid  string    `json:"cluster-uuid,omitempty"`  //cluster_uuid
	ClientUuid   string    `json:"client-uuid,omitempty"`   //client_uuid
	PollInterval int       `json:"poll-interval,omitempty"` //config_poll_interval
	Loglevel     string    `json:"Loglevel,omitempty"`      //config_loglevel
}
