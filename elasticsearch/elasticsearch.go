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

type EsQueryBody struct {
	Id          string      `json:"_id"`
	Index       string      `json:"_index"`
	Type        string      `json:"_type"`
	PrimaryTerm int         `json:"_primary_term"`
	SeqNo       int         `json:"_seq_no"`
	Version     int         `json:"_version"`
	Found       bool        `json:"found"`
	Source      interface{} `json:"_source"` // 继承的结构体替换该字段内容
}

type responseErrState struct {
	Error  interface{} `json:"error"`
	Status int         `json:"status"`
}

const timeout = time.Second * 3

var headers = map[string]string{"Content-Type": "Application/Json"}

func esErr(resp *request.Response) error {
	var response responseErrState
	_ = resp.Json(&response)
	if !(response.Error == "" || response.Error == nil) {
		return errors.New(resp.Text())
	}
	return resp.Err()
}

func (e *ES) IndexCreate(indexName string, shardNum, replicaNum int) error {
	addr := fmt.Sprintf("%s://%s:%d/%s", e.Schema, e.ServerAddr, e.ServerPort, indexName)
	body := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   shardNum,
			"number_of_replicas": replicaNum,
		},
	}
	resp := request.Put(addr, body, headers, nil, timeout)
	return esErr(resp)
}

func (e *ES) DocsInsertUpdate(indexName string, id interface{}, data interface{}) error {
	addr := fmt.Sprintf("%s://%s:%d/%s/_doc", e.Schema, e.ServerAddr, e.ServerPort, indexName)
	if !(id == "" || id == nil) {
		addr = fmt.Sprintf("%s/%v", addr, id)
	}
	resp := request.Post(addr, data, headers, nil, timeout)
	return esErr(resp)
}

func (e *ES) DocsDelete(indexName, id string) error {
	if id == "" {
		return errors.New("id参数不能为空")
	}
	addr := fmt.Sprintf("%s://%s:%d/%s/_doc/%s", e.Schema, e.ServerAddr, e.ServerPort, indexName, id)
	resp := request.Delete(addr, nil, headers, nil, timeout)
	return esErr(resp)
}

func (e *ES) DocsGet(indexName, id string, responseBody interface{}) error {
	if id == "" {
		return errors.New("id参数不能为空")
	}
	addr := fmt.Sprintf("%s://%s:%d/%s/_doc/%s", e.Schema, e.ServerAddr, e.ServerPort, indexName, id)
	resp := request.Get(addr, nil, headers, nil, timeout)
	if err := resp.Json(&responseBody); err != nil {
		return err
	}
	return esErr(resp)
}
