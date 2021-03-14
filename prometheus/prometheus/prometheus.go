package prometheus

import (
	"errors"
	"fmt"
	"time"
	"utils/request"
)

type QueryResponseBody struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// 查询指定范时间范围内的数据
func QueryRange(serverIp string, port int, promQL string, startTime time.Time, endTime time.Time,
	step int, timeout time.Duration) (*QueryResponseBody, error) {
	var err error
	url := fmt.Sprintf("http://%s:%d/api/v1/query_range", serverIp, port)
	queryResponseBody := QueryResponseBody{}
	params := map[string]string{
		"query": promQL,
		"start": string(startTime.Unix()),
		"end":   string(endTime.Unix()),
		"step":  string(step),
	}
	req := request.Get(url, nil, nil, params, timeout)
	if req.Err() != nil {
		return &QueryResponseBody{}, err
	}
	if req.Code() != 200 {
		switch req.Code() {
		case 404:
			err = errors.New("请求参数错误或数据丢失")
		case 402:
			err = errors.New("promQL错误，执行无效")
		case 503:
			err = errors.New("请求超时或者被中断")
		default:
			err = errors.New(fmt.Sprintf("prometheus响应错误，响应状态码是:%d", req.Code()))
		}
	}
	if err = req.Json(&queryResponseBody); err != nil {
		return &QueryResponseBody{}, err
	}
	return &queryResponseBody, nil
}

// 标记需要清除的数据，标记后不会立即执行，调用/api/v1/admin/tsdb/clean_tombstones后才从磁盘删除
func FlagDeleteData(serverIp string, port int, promQL string, timeout time.Duration) error {
	url := fmt.Sprintf("http://%s:%d/api/v1/admin/tsdb/delete_series", serverIp, port)
	parms := map[string]string{
		"match[]": promQL,
	}
	req := request.Post(url, nil, nil, parms, timeout)
	if req.Err() != nil {
		return req.Err()
	}
	if 299 < req.Code() && req.Code() < 200 {
		return errors.New(fmt.Sprintf("标记删除的数据错误，响应码:%d", req.Code()))
	}
	return nil
}

// 标记指定时间范围需要清除的数据，标记后不会立即执行，调用/api/v1/admin/tsdb/clean_tombstones后才从磁盘删除
func FlagDeleteDataRange(serverIp string, port int, promQL string,
	startTime time.Time, endTime time.Time, timeout time.Duration) error {
	url := fmt.Sprintf("http://%s:%d/api/v1/admin/tsdb/delete_series", serverIp, port)
	parms := map[string]string{
		"match[]": promQL,
		"start":   string(startTime.Unix()),
		"end":     string(endTime.Unix()),
	}
	req := request.Post(url, nil, nil, parms, timeout)
	if req.Err() != nil {
		return req.Err()
	}
	if 299 < req.Code() && req.Code() < 200 {
		return errors.New(fmt.Sprintf("标记删除的数据错误，响应码:%d", req.Code()))
	}
	return nil
}

// 从磁盘中清理标记需要删除的数据
func CleanData(serverIp string, port int, timeout time.Duration) error {
	url := fmt.Sprintf("http://%s:%d/api/v1/admin/tsdb/clean_tombstones", serverIp, port)
	req := request.Post(url, nil, nil, nil, timeout)
	if req.Err() != nil {
		return req.Err()
	}
	if 299 < req.Code() && req.Code() < 200 {
		return errors.New(fmt.Sprintf("标记删除的数据错误，响应码:%d", req.Code()))
	}
	return nil
}
