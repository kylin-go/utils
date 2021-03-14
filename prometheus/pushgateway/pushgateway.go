package pushgateway

import (
	"errors"
	"fmt"
	"time"
	"utils/request"
)

func push2Gateway(serverIp string, port int, jobName, instanceName, metricName string, labels map[string]string, val float64, timestamp bool) error {
	var data string
	url := fmt.Sprintf("http://%s:%d/metrics/job/%s/instance/%s", serverIp, port, jobName, instanceName)
	label := ""
	for k, v := range labels {
		label = fmt.Sprintf(`%s%s=%s`, label, k, v)
	}
	label = label[:len(label)-1]
	if timestamp {
		data = fmt.Sprintf("%s%s %f %d\n", metricName, label, val, time.Now().Unix())
	} else {
		data = fmt.Sprintf("%s%s %f\n", metricName, label, val)
	}
	req := request.Post(url, data, map[string]string{"Content-Type": "text/plain"}, nil, time.Second*3)
	if req.Err() != nil {
		return req.Err()
	}
	if 299 < req.Code() && req.Code() < 200 {
		return errors.New(fmt.Sprintf("向pushgateway推送数据错误，响应码:%d", req.Code()))
	}
	return nil
}

func deleteGateway(serverIp string, port int, jobName, instanceName string) error {
	var url string
	if instanceName == "" {
		url = fmt.Sprintf("http://%s:%d/metrics/job/%s", serverIp, port, jobName)
	} else {
		url = fmt.Sprintf("http://%s:%d/metrics/job/%s/instance/%s", serverIp, port, jobName, instanceName)
	}
	req := request.Delete(url, nil, nil, nil, time.Second*3)
	if req.Err() != nil {
		return req.Err()
	}
	if 299 < req.Code() && req.Code() < 200 {
		return errors.New(fmt.Sprintf("删除pushgateway的数据错误，响应码:%d", req.Code()))
	}
	return nil
}

// 向pushGateway推送数据，不带当前数据时间戳
func PushWithoutTimestamp(serverIp string, port int, jobName, instanceName, metricName string, labels map[string]string, val float64) error {
	return push2Gateway(serverIp, port, jobName, instanceName, metricName, labels, val, false)
}

// 向pushGateway推送数据，带当前数据时间戳
func PushWithTimestamp(serverIp string, port int, jobName, instanceName, metricName string, labels map[string]string, val float64) error {
	return push2Gateway(serverIp, port, jobName, instanceName, metricName, labels, val, true)
}

// 删除pushGateway所有job的数据
func DeleteJob(serverIp string, port int, jobName string) error {
	return deleteGateway(serverIp, port, jobName, "")
}

// 删除pushGateway所有instance的数据
func DeleteInstance(serverIp string, port int, jobName, instanceName string) error {
	return deleteGateway(serverIp, port, jobName, instanceName)
}
