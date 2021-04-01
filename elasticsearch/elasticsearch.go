package elasticsearch

import (
	"errors"
	"fmt"
	"time"
	"utils/elasticsearch/docs"
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

// 插入与覆盖文档
func (e *ES) DocsInsert(indexName, id string, data interface{}) *docs.InsertRespReturn {
	return docs.Insert(e.Schema, e.ServerAddr, e.ServerPort, indexName, id, data)
}

// 局部跟新文档
func (e *ES) DocsUpdate(indexName, id string, data interface{}) *docs.InsertRespReturn {
	return docs.Update(e.Schema, e.ServerAddr, e.ServerPort, indexName, id, data)
}

// 删除文档
func (e *ES) DocsDelete(indexName, id string) *docs.InsertRespReturn {
	return docs.Delete(e.Schema, e.ServerAddr, e.ServerPort, indexName, id)
}

// 按照ID获取文档
func (e *ES) DocsGet(indexName, id string, respData interface{}) error {
	return docs.Get(e.Schema, e.ServerAddr, e.ServerPort, indexName, id, respData)
}

// 根据文档ID查询是否存在
func (e *ES) DocsExist(indexName, id string) (bool, error) {
	return docs.IsExist(e.Schema, e.ServerAddr, e.ServerPort, indexName, id)
}

func (e *ES) DocsSearch(indexName, dsl string, result interface{}) error {
	return docs.Search(e.Schema, e.ServerAddr, e.ServerPort, indexName, dsl, result)
}
