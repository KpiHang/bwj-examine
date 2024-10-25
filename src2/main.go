package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Println("starting curl ...")
	for {

		// 定义要执行的 curl 命令和参数
		cmd := exec.Command("curl", "deviceshifu-plate-reader.deviceshifu.svc.cluster.local/get_measurement")

		// 获取命令的输出
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return
		}

		average, values := Average(string(output)) // 计算平均值
		log.Printf("Average is: %.2f from values: %v", average, values)

		time.Sleep(3 * time.Second)
	}
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
