package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const baseURL = "http://localhost:5001"

var token string

type Response struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func main() {
	fmt.Println("========================================")
	fmt.Println("财务管理模块集成测试")
	fmt.Println("========================================\n")

	// 步骤1: 登录获取token（需要手动提供）
	fmt.Println("步骤1: 认证测试")
	fmt.Println("----------------------------------------")
	fmt.Println("⚠️  需要先登录系统获取JWT token")
	fmt.Println("请打开浏览器访问: http://localhost:5001")
	fmt.Println("登录后，从浏览器开发者工具的Network标签中获取Authorization头")
	fmt.Println("")
	fmt.Print("请输入JWT token（或按Enter跳过认证测试）: ")

	// 读取token（简化处理）
	// 在实际测试中，可以从环境变量或配置文件读取
	token = os.Getenv("TEST_JWT_TOKEN")

	if token == "" {
		fmt.Println("未提供token，仅测试接口可访问性")
		fmt.Println("")
	}

	// 步骤2: 测试API响应格式
	fmt.Println("步骤2: 测试API响应格式")
	fmt.Println("----------------------------------------")
	testAPIResponse("/api/payment-collections", "GET")
	testAPIResponse("/api/separate-accounts", "GET")
	fmt.Println("")

	// 如果有token，进行完整的功能测试
	if token != "" {
		runAuthenticatedTests()
	} else {
		fmt.Println("跳过需要认证的测试")
		fmt.Println("提示: 设置环境变量 TEST_JWT_TOKEN 可运行完整测试")
		fmt.Println("      例如: export TEST_JWT_TOKEN=\"your-token-here\"")
	}

	fmt.Println("\n========================================")
	fmt.Println("测试完成")
	fmt.Println("========================================")
}

func testAPIResponse(endpoint, method string) {
	fmt.Printf("测试: %s %s\n", method, endpoint)

	var resp *http.Response
	var err error

	if method == "GET" {
		resp, err = http.Get(baseURL + endpoint)
	}

	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response Response
	json.Unmarshal(body, &response)

	fmt.Printf("   状态码: %d\n", resp.StatusCode)
	fmt.Printf("   响应码: %d\n", response.Code)
	fmt.Printf("   消息: %s\n", response.Message)

	// 验证响应格式
	if response.Code == 401 {
		fmt.Println("   ✅ 需要认证（预期行为）")
	} else {
		fmt.Printf("   ⚠️  意外的响应: %s\n", string(body))
	}
	fmt.Println("")
}

func runAuthenticatedTests() {
	fmt.Println("\n步骤3: 功能测试（需要认证）")
	fmt.Println("========================================\n")

	// 测试3.1: 获取收款列表
	fmt.Println("测试3.1: 获取收款列表")
	fmt.Println("----------------------------------------")
	resp := makeAuthRequest("GET", "/api/payment-collections?page=1&page_size=10", nil)
	printResponse(resp)
	fmt.Println("")

	// 测试3.2: 创建收款
	fmt.Println("测试3.2: 创建收款")
	fmt.Println("----------------------------------------")
	createData := map[string]interface{}{
		"order_id":         1,
		"student_id":       1,
		"payment_scenario": 1,
		"payment_method":   0,
		"payment_amount":   100.00,
		"payer":            "自动化测试",
		"payee_entity":     0,
		"merchant_order":   fmt.Sprintf("TEST%d", time.Now().Unix()),
	}
	resp = makeAuthRequest("POST", "/api/payment-collections", createData)
	var createResp Response
	json.Unmarshal(resp, &createResp)
	printResponse(resp)

	paymentID := 0
	if createResp.Code == 0 && createResp.Data != nil {
		if id, ok := createResp.Data["id"].(float64); ok {
			paymentID = int(id)
			fmt.Printf("✅ 创建成功，收款ID: %d\n", paymentID)
		}
	}
	fmt.Println("")

	// 如果创建成功，继续测试后续流程
	if paymentID > 0 {
		// 测试3.3: 确认到账
		fmt.Println("测试3.3: 确认到账")
		fmt.Println("----------------------------------------")
		confirmURL := fmt.Sprintf("/api/payment-collections/%d/confirm", paymentID)
		resp = makeAuthRequest("PUT", confirmURL, nil)
		printResponse(resp)
		fmt.Println("")

		// 测试3.4: 查询分账明细
		fmt.Println("测试3.4: 查询分账明细")
		fmt.Println("----------------------------------------")
		separateURL := fmt.Sprintf("/api/separate-accounts?payment_id=%d", paymentID)
		resp = makeAuthRequest("GET", separateURL, nil)
		printResponse(resp)
		fmt.Println("")

		// 测试3.5: 再次获取收款列表，验证状态
		fmt.Println("测试3.5: 验证收款状态")
		fmt.Println("----------------------------------------")
		listURL := fmt.Sprintf("/api/payment-collections?id=%d", paymentID)
		resp = makeAuthRequest("GET", listURL, nil)
		printResponse(resp)
		fmt.Println("")
	}

	// 测试3.6: 获取分账明细列表
	fmt.Println("测试3.6: 获取分账明细列表")
	fmt.Println("----------------------------------------")
	resp = makeAuthRequest("GET", "/api/separate-accounts?page=1&page_size=10", nil)
	printResponse(resp)
	fmt.Println("")
}

func makeAuthRequest(method, endpoint string, data interface{}) []byte {
	client := &http.Client{}
	var req *http.Request
	var err error

	if data != nil {
		jsonData, _ := json.Marshal(data)
		req, err = http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest(method, baseURL+endpoint, nil)
	}

	if err != nil {
		fmt.Printf("❌ 创建请求失败: %v\n", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func printResponse(body []byte) {
	if body == nil {
		return
	}

	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("响应: %s\n", string(body))
		return
	}

	fmt.Printf("响应码: %d\n", response.Code)
	fmt.Printf("消息: %s\n", response.Message)
	if response.Data != nil {
		dataJSON, _ := json.MarshalIndent(response.Data, "", "  ")
		fmt.Printf("数据: %s\n", string(dataJSON))
	}
}
