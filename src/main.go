package main

import (
	"bytes"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Println("starting...")

	cli, cfg := Link_K8s()

	for {

		req := cli.CoreV1().RESTClient().Post().
			Namespace("default").
			Resource("pods").
			Name("nginx").
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: "nginx",
				Command:   []string{"curl", "deviceshifu-plate-reader.deviceshifu.svc.cluster.local/get_measurement"},
				Stdout:    true,
				Stderr:    true,
			}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
		if err != nil {
			panic(err)
		}

		var stdout, stderr bytes.Buffer
		err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
			Stdout: &stdout,
			Stderr: &stderr,
		})
		if err != nil {
			panic(err)
		}

		output := stdout.String()
		average, values := Average(output) // 计算平均值
		log.Printf("Average is: %.2f from values: %v", average, values)

		time.Sleep(3 * time.Second)
	}
}

func Link_K8s() (*kubernetes.Clientset, *rest.Config) {
	kubeconfig := "../.kube/config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panicf("failed to build config: %v", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("failed to create clientset: %v", err.Error())
	}
	return clientset, config
}

func Average(str string) (float64, []float64) {
	// 将字符串分割为子字符串切片
	strValues := strings.Fields(str)

	// 转换为浮点数切片
	var floatValues []float64
	for _, strVal := range strValues {
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			fmt.Printf("Error parsing '%s': %v\n", strVal, err)
			continue
		}
		floatValues = append(floatValues, floatVal)
	}

	// 计算平均值
	var sum float64
	for _, val := range floatValues {
		sum += val
	}
	average := sum / float64(len(floatValues))
	return average, floatValues
}
