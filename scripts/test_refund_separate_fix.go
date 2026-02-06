package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// 测试退费分账明细修复
// 验证退费类分账是否按照售卖分账的分布正确生成
func main() {
	// 连接数据库
	dsn := "root:your_password@tcp(localhost:3306)/charonoms?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 测试场景2：需要冲回的情况
	fmt.Println("=== 测试场景2: 需要冲回 ===")
	testRefundWithChargeback(db)
}

func testRefundWithChargeback(db *sql.DB) {
	// 假设已经创建了测试数据并执行了退费审批通过
	// 这里验证分账明细的正确性

	orderID := 1 // 测试订单ID

	fmt.Println("\n1. 查询所有分账明细:")
	rows, err := db.Query(`
		SELECT id, childorders_id, payment_id, payment_type, separate_amount, type, parent_id
		FROM separate_account
		WHERE orders_id = ?
		ORDER BY id ASC
	`, orderID)
	if err != nil {
		log.Fatal("查询分账明细失败:", err)
	}
	defer rows.Close()

	separates := make(map[int]struct {
		ChildOrderID   int
		PaymentID      int
		PaymentType    int
		SeparateAmount float64
		Type           int
		ParentID       sql.NullInt64
	})

	fmt.Println("ID\t子订单ID\t收款ID\t类型\t金额\t\tParent")
	for rows.Next() {
		var id int
		var sep struct {
			ChildOrderID   int
			PaymentID      int
			PaymentType    int
			SeparateAmount float64
			Type           int
			ParentID       sql.NullInt64
		}
		err := rows.Scan(&id, &sep.ChildOrderID, &sep.PaymentID, &sep.PaymentType,
			&sep.SeparateAmount, &sep.Type, &sep.ParentID)
		if err != nil {
			log.Fatal("扫描数据失败:", err)
		}
		separates[id] = sep

		typeName := ""
		switch sep.Type {
		case 0:
			typeName = "售卖"
		case 1:
			typeName = "冲回"
		case 2:
			typeName = "退费"
		}

		parentStr := "NULL"
		if sep.ParentID.Valid {
			parentStr = fmt.Sprintf("%d", sep.ParentID.Int64)
		}

		fmt.Printf("%d\t%d\t\t%d\t%s\t%.2f\t\t%s\n",
			id, sep.ChildOrderID, sep.PaymentID, typeName, sep.SeparateAmount, parentStr)
	}

	fmt.Println("\n2. 验证退费类分账是否按照售卖分账分布生成:")

	// 查询退费子订单
	refundRows, err := db.Query(`
		SELECT childorder_id, refund_amount
		FROM refund_order_item
		WHERE refund_order_id = 1
		ORDER BY childorder_id ASC
	`)
	if err != nil {
		log.Fatal("查询退费子订单失败:", err)
	}
	defer refundRows.Close()

	for refundRows.Next() {
		var childOrderID int
		var refundAmount float64
		refundRows.Scan(&childOrderID, &refundAmount)

		fmt.Printf("\n子订单%d (退费金额: %.2f):\n", childOrderID, refundAmount)

		// 查询该子订单的售卖分账
		fmt.Println("  售卖分账:")
		sellRows, err := db.Query(`
			SELECT payment_id, separate_amount
			FROM separate_account
			WHERE childorders_id = ? AND type = 0
				AND id NOT IN (
					SELECT parent_id FROM separate_account
					WHERE childorders_id = ? AND parent_id IS NOT NULL
				)
			ORDER BY payment_id ASC
		`, childOrderID, childOrderID)
		if err != nil {
			log.Fatal("查询售卖分账失败:", err)
		}

		sellTotal := 0.0
		for sellRows.Next() {
			var paymentID int
			var amount float64
			sellRows.Scan(&paymentID, &amount)
			fmt.Printf("    收款%d: %.2f元\n", paymentID, amount)
			sellTotal += amount
		}
		sellRows.Close()
		fmt.Printf("  售卖分账总计: %.2f元\n", sellTotal)

		// 查询该子订单的退费分账
		fmt.Println("  退费分账:")
		refundSepRows, err := db.Query(`
			SELECT payment_id, separate_amount
			FROM separate_account
			WHERE childorders_id = ? AND type = 2
			ORDER BY payment_id ASC
		`, childOrderID)
		if err != nil {
			log.Fatal("查询退费分账失败:", err)
		}

		refundTotal := 0.0
		for refundSepRows.Next() {
			var paymentID int
			var amount float64
			refundSepRows.Scan(&paymentID, &amount)
			fmt.Printf("    收款%d: %.2f元\n", paymentID, amount)
			refundTotal += amount
		}
		refundSepRows.Close()
		fmt.Printf("  退费分账总计: %.2f元\n", refundTotal)

		// 验证
		if -refundTotal == refundAmount {
			fmt.Printf("  ✓ 退费分账总额正确 (%.2f == %.2f)\n", -refundTotal, refundAmount)
		} else {
			fmt.Printf("  ✗ 退费分账总额不正确 (%.2f != %.2f)\n", -refundTotal, refundAmount)
		}
	}

	fmt.Println("\n3. 验证子订单状态:")
	statusRows, err := db.Query(`
		SELECT id, amount_received, status
		FROM childorders
		WHERE parentsid = ?
		ORDER BY id ASC
	`, orderID)
	if err != nil {
		log.Fatal("查询子订单状态失败:", err)
	}
	defer statusRows.Close()

	fmt.Println("子订单ID\t实收金额\t状态")
	for statusRows.Next() {
		var id int
		var amountReceived float64
		var status int
		statusRows.Scan(&id, &amountReceived, &status)

		statusName := ""
		switch status {
		case 10:
			statusName = "未支付"
		case 20:
			statusName = "部分支付"
		case 30:
			statusName = "已支付"
		}

		fmt.Printf("%d\t\t%.2f\t\t%s\n", id, amountReceived, statusName)
	}
}
