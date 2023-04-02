package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/ljjgs/pod/common"
	go_micro_service_pod "github.com/ljjgs/pod/proto/pod"
	"github.com/ljjgs/podApi/handler"
	podApi "github.com/ljjgs/podApi/proto/podApi/proto"
	"github.com/opentracing/opentracing-go"
	"strconv"
)

var (
	serviceHost string = "podapi"
	servicePort        = "8081"
	consulHost         = "162.14.77.55"
	consulPort  int64  = 8500
	tracerHost         = "162.14.77.55"
	tracerPort         = 16686
	//hystrixPort           = 9002
	prometheusPort = "9090"
)

func main() {
	//1.注册中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			consulHost + ":" + strconv.FormatInt(consulPort, 10),
		}
	})

	//2.添加链路追踪
	t, io, err := common.NewTracer("go.micro.api.podApi", tracerHost+":"+strconv.Itoa(tracerPort))
	if err != nil {
		panic(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//3.添加熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	//4.添加日志
	//1）需要程序日志打入到日志文件中
	//2）在程序中添加filebeat.yml 文件
	//3) 启动filebeat，启动命令 ./filebeat -e -c filebeat.yml
	fmt.Println("日志统一记录在根目录 micro.log 文件中，请点击查看日志！")

	////6.启动熔断监听程序
	//go func() {
	//	//http://192.168.0.108:9092/turbine/turbine.stream
	//	//看板访问地址 http://127.0.0.1:9002/hystrix，url后面一定要带 /hystrix
	//	err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", strconv.Itoa(hystrixPort)), hystrixStreamHandler)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	//7.添加监控采集地址
	//common.PrometheusBoot(prometheusPort)

	//8.创建服务
	service := micro.NewService(
		micro.Server(
			server.NewServer(
				server.Address("podapi:8081"), // 指定要使用的端口地址
			),
		),
		micro.Name("go.micro.api.podApi"),
		micro.Version("latest"),
		//指定服务端口
		micro.Address(":"+servicePort),
		//添加注册中心，
		micro.Registry(consul),
		//添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
		//添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	service.Init()

	podService := go_micro_service_pod.NewPodService("go.micro.service.pod", service.Client())
	//注册控制器
	if err := podApi.RegisterPodApiHandler(service.Server(), &handler.PodApi{PodService: podService}); err != nil {
		panic(err)
	}
	// 启动服务
	if err := service.Run(); err != nil {
		panic(err)
	}
}
