package docs

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"utils/request"
)

type EsBaseState struct {
	Id          string      `json:"_id"`
	Index       string      `json:"_index"`
	Type        string      `json:"_type"`
	SeqNo       int         `json:"_seq_no"`
	Version     int         `json:"_version"`
	PrimaryTerm int         `json:"_primary_term"`
	Error       interface{} `json:"error"`
	Status      int         `json:"status"`
}

type ShardsState struct {
	Total      int16 `json:"total"`
	Successful int16 `json:"successful"`
	Failed     int16 `json:"failed"`
}

type InsertRespReturn struct {
	EsBaseState
	Result string      `json:"result"`
	Shards ShardsState `json:"_shards"`
}

type QueryRespBody struct {
	EsBaseState
	Found  bool        `json:"found"`
	Source interface{} `json:"_source"` // 继承的结构体替换该字段内容
}

const timeout = time.Second * 3

var headers = map[string]string{"Content-Type": "Application/Json"}

func Insert(schema, addr string, port int, indexName, id string, data interface{}) *InsertRespReturn {
	var r = InsertRespReturn{}
	url := fmt.Sprintf("%s://%s:%d/%s/_doc", schema, addr, port, indexName)
	if id != "" {
		url = fmt.Sprintf("%s/%s", url, id)
	}
	resp := request.Post(url, data, headers, nil, timeout)
	if err := resp.Json(&r); err != nil {
		r.Error = err.Error()
		r.Status = 500
	}
	if resp.Err() != nil {
		r.Error = resp.Err().Error()
		r.Status = resp.Code()
	}
	return &r
}

func Update(schema, addr string, port int, indexName, id string, data interface{}) *InsertRespReturn {
	var r = InsertRespReturn{}
	url := fmt.Sprintf("%s://%s:%d/%s/_doc/%s/_update", schema, addr, port, indexName, id)
	data = map[string]interface{}{"doc": data}
	resp := request.Post(url, data, headers, nil, timeout)
	if err := resp.Json(&r); err != nil {
		r.Error = err.Error()
		r.Status = 500
	}
	if resp.Err() != nil {
		r.Error = resp.Err().Error()
		r.Status = resp.Code()
	}
	return &r
}

func Delete(schema, addr string, port int, indexName, id string) *InsertRespReturn {
	var r = InsertRespReturn{}
	url := fmt.Sprintf("%s://%s:%d/%s/_doc/%s", schema, addr, port, indexName, id)
	resp := request.Delete(url, nil, headers, nil, timeout)
	if err := resp.Json(&r); err != nil {
		r.Error = err.Error()
		r.Status = 500
	}
	if resp.Err() != nil {
		r.Error = resp.Err().Error()
		r.Status = resp.Code()
	}
	return &r
}

func Get(schema, addr string, port int, indexName, id string, result interface{}) error {
	url := fmt.Sprintf("%s://%s:%d/%s/_doc/%s", schema, addr, port, indexName, id)
	resp := request.Get(url, nil, headers, nil, timeout)
	if err := resp.Json(&result); err != nil {
		return err
	}
	if resp.Err() != nil {
		return resp.Err()
	}
	return nil
}

func IsExist(schema, addr string, port int, indexName, id string) (bool, error) {
	url := fmt.Sprintf("%s://%s:%d/%s/_doc/%s", schema, addr, port, indexName, id)
	resp := request.Head(url, nil, headers, nil, timeout)
	if resp.Err() != nil {
		return false, resp.Err()
	}
	if resp.Code() != 200 {
		return false, nil
	}
	return true, nil
}

func Search(schema, addr string, port int, indexName, dsl string, result interface{}) error {
	var D map[string]interface{}
	url := fmt.Sprintf("%s://%s:%d/%s/_doc/_search", schema, addr, port, indexName)
	if err := json.Unmarshal([]byte(dsl), &D); err != nil {
		return errors.New("Dsl解析错误, " + err.Error())
	}
	resp := request.Post(url, D, headers, nil, timeout)
	if err := resp.Json(&result); err != nil {
		return err
	}
	if resp.Err() != nil {
		return resp.Err()
	}
	return nil
}
