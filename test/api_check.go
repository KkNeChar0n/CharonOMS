package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "http://localhost:5001"

type Response struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func main() {
	fmt.Println("========================================")
	fmt.Println("财务管理模块 API 测试")
	fmt.Println("========================================\n")

	// 测试1: 获取收款列表（无认证）
	fmt.Println("测试1: 获取收款列表（无认证）")
	fmt.Println("----------------------------------------")
	resp, err := http.Get(baseURL + "/api/payment-collections")
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n\n", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应: %s\n\n", string(body))
		resp.Body.Close()
	}

	// 测试2: 获取分账明细列表（无认证）
	fmt.Println("测试2: 获取分账明细列表（无认证）")
	fmt.Println("----------------------------------------")
	resp, err = http.Get(baseURL + "/api/separate-accounts")
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n\n", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应: %s\n\n", string(body))
		resp.Body.Close()
	}

	// 测试3: 尝试创建收款（无认证）
	fmt.Println("测试3: 创建收款（无认证）")
	fmt.Println("----------------------------------------")
	createData := map[string]interface{}{
		"order_id":         1,
		"student_id":       1,
		"payment_scenario": 1,
		"payment_method":   0,
		"payment_amount":   100.00,
		"payer":            "测试用户",
		"payee_entity":     0,
		"merchant_order":   "TEST001",
	}
	jsonData, _ := json.Marshal(createData)
	resp, err = http.Post(baseURL+"/api/payment-collections", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n\n", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应: %s\n\n", string(body))
		resp.Body.Close()
	}

	// 测试4: 检查路由是否存在
	fmt.Println("测试4: 检查路由端点")
	fmt.Println("----------------------------------------")
	endpoints := []string{
		"/api/payment-collections",
		"/api/separate-accounts",
		"/api/orders",
		"/api/students",
	}

	for _, endpoint := range endpoints {
		resp, err := http.Get(baseURL + endpoint)
		if err != nil {
			fmt.Printf("❌ %s - 请求失败: %v\n", endpoint, err)
			continue
		}
		fmt.Printf("✅ %s - 状态码: %d\n", endpoint, resp.StatusCode)
		resp.Body.Close()
	}

	fmt.Println("\n========================================")
	fmt.Println("测试完成")
	fmt.Println("========================================")
}
