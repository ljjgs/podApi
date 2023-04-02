package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/ljjgs/pod/proto/pod"
	from "github.com/ljjgs/podApi/plugin/form"
	podApi "github.com/ljjgs/podApi/proto/podApi/proto"
	"net/http"
	"strconv"
)

type PodApi struct {
	PodService pod.PodService
}

func (e *PodApi) DeletePodById(ctx context.Context, request *podApi.Request, response *podApi.Response) error {
	fmt.Println("接受到 podApi.DeletePodById 的请求")
	if _, ok := request.Get["pod_id"]; !ok {
		return errors.New("参数异常")
	}
	//获取要删除的ID
	podIdString := request.Get["pod_id"].Values[0]
	podId, err := strconv.ParseInt(podIdString, 10, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	//删除指定服务

	id, err := e.PodService.DeletePod(ctx, &pod.PodId{
		Id: podId,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	response.StatusCode = 200
	b, _ := json.Marshal(id)
	response.Body = string(b)
	return nil
}

func (e *PodApi) FindPodById(ctx context.Context, req *podApi.Request, res *podApi.Response) error {
	fmt.Println("接受到 podApi.FindPodById 的请求")
	if _, ok := req.Get["pod_id"]; !ok {
		res.StatusCode = 500
		return errors.New("参数异常")
	}
	//获取podid 参数
	podIdString := req.Get["pod_id"].Values[0]
	podId, err := strconv.ParseInt(podIdString, 10, 64)
	if err != nil {
		return err
	}
	//获取pod相关信息
	podInfo, err := e.PodService.FindPodByID(ctx, &pod.PodId{
		Id: podId,
	})
	if err != nil {
		return err
	}
	//json 返回pod信息
	res.StatusCode = 200
	b, _ := json.Marshal(podInfo)
	res.Body = string(b)
	fmt.Println(res.Body)
	return nil
}
func (e *PodApi) AddPod(ctx context.Context, req *podApi.Request, res *podApi.Response) error {
	fmt.Println("接受到 podApi.AddPod 的请求")
	addPodInfo := &pod.PodInfo{}
	//处理 port
	dataSlice, ok := req.Post["pod_port"]
	if ok {
		//特殊处理
		podSlice := []*pod.PodPort{}
		for _, v := range dataSlice.Values {
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				log.Error(err)
			}
			port := &pod.PodPort{
				ContainerPort: int32(i),
				Protocol:      "TCP",
			}
			podSlice = append(podSlice, port)
		}
		addPodInfo.PodPort = podSlice
	}
	//form类型转化到结构体中
	from.FromToPodStruct(req.Post, addPodInfo)

	response, err := e.PodService.AddPod(ctx, addPodInfo)
	if err != nil {
		log.Error(err)
		return err
	}

	res.StatusCode = 200
	b, _ := json.Marshal(response)
	res.Body = string(b)
	fmt.Println(res.Body)
	return nil
}
func (e *PodApi) UpdatePod(ctx context.Context, req *podApi.Request, res *podApi.Response) error {
	fmt.Println("收到增加pod的请求")
	addPodInfo := &pod.PodInfo{}
	//处理port
	if dataslice, pair := req.Post["pod_port"]; pair {
		podslice := []*pod.PodPort{}
		for _, value := range dataslice.Values {
			i, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}
			port := &pod.PodPort{
				ContainerPort: int32(i),
				Protocol:      "TCP",
			}
			podslice = append(podslice, port)
		}
		addPodInfo.PodPort = podslice
	}
	from.FromToPodStruct(req.Post, addPodInfo)
	addPod, err := e.PodService.UpdatePod(ctx, addPodInfo)
	if err != nil {
		return err
	}
	res.StatusCode = http.StatusOK
	marshal, err := json.Marshal(addPod)
	req.Body = string(marshal)
	return nil
}

func (e *PodApi) Call(ctx context.Context, req *podApi.Request, res *podApi.Response) error {
	fmt.Println("接受到 podApi.Call 的请求")
	allPod, err := e.PodService.FindAllPod(ctx, &pod.FindAll{})
	if err != nil {
		log.Error(err)
		return err
	}
	res.StatusCode = 200
	b, _ := json.Marshal(allPod)
	res.Body = string(b)
	return nil
}
