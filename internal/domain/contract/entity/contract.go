package entity

import "time"

// Contract 合同实体
type Contract struct {
	ID                   int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                 string    `gorm:"column:name;type:varchar(200);not null" json:"name"`
	StudentID            int       `gorm:"column:student_id;not null" json:"student_id"`
	Type                 int       `gorm:"column:type;not null" json:"type"`
	SignatureForm        int       `gorm:"column:signature_form;not null" json:"signature_form"`
	ContractAmount       float64   `gorm:"column:contract_amount;type:decimal(10,2);not null" json:"contract_amount"`
	Signatory            string    `gorm:"column:signatory;type:varchar(100)" json:"signatory"`
	InitiatingParty      string    `gorm:"column:initiating_party;type:varchar(100)" json:"initiating_party"`
	Initiator            string    `gorm:"column:initiator;type:varchar(50)" json:"initiator"`
	Status               int       `gorm:"column:status;default:0" json:"status"`
	PaymentStatus        int       `gorm:"column:payment_status;default:0" json:"payment_status"`
	TerminationAgreement string    `gorm:"column:termination_agreement;type:varchar(255)" json:"termination_agreement"`
	CreateTime           time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (Contract) TableName() string {
	return "contract"
}
