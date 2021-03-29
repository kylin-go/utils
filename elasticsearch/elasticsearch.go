package elasticsearch

import (
	"errors"
	"fmt"
	"time"
	"utils/request"
)

type ES struct {
	Schema     string
	ServerAddr string
	ServerPort int
}

type responseErrState struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

const timeout = time.Second * 3

var headers = map[string]string{"Content-Type": "Application/Json"}

func (e *ES) IndexCreate(indexName string, shardNum, replicaNum int) error {
	var response responseErrState
	addr := fmt.Sprintf("%s://%s:%d/%s", e.Schema, e.ServerAddr, e.ServerPort, indexName)
	body := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   shardNum,
			"number_of_replicas": replicaNum,
		},
	}
	resp := request.Put(addr, body, headers, nil, timeout)
	_ = resp.Json(&response)
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return resp.Err()
}
