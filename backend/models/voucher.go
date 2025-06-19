package models

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Voucher struct {
	Id              int       `json:"id"`
	VoucherCode     string    `json:"voucher_code"`
	DiscountPercent int       `json:"discount_percent"`
	ExpiryDate      string    `json:"expiry_date"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status          int       `json:"status" gorm:"type:varchar(100);not null;"`
}

func (Voucher) TableName() string {
	return "voucher"
}

type DBClient struct {
	DB *gorm.DB
}

func NewDBClient(db *gorm.DB) *DBClient {
	return &DBClient{DB: db}
}

func (c *DBClient) CreateVoucher(voucher *Voucher) error {
	c.DB.Create(voucher)
	return nil
}

func (c *DBClient) CreateBulkVouchers(vouchers *[]Voucher) error {
	return c.DB.Create(vouchers).Error
}

func (c *DBClient) UpdateVoucher(voucher *Voucher, id int) error {
	var Exist Voucher
	result := c.DB.First(&Exist, id).Error
	if result != nil {
		return result
	}

	voucher.Id = id

	err := c.DB.Where("id = ?", id).Updates(voucher).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *DBClient) GetVoucherById(id int) (*Voucher, error) {
	var voucher Voucher
	rows := c.DB.Where("status = ?", 1).First(&voucher, id).Error
	if rows != nil {
		return nil, rows
	}

	return &voucher, nil

}

func (c *DBClient) DeleteVoucher(id int) error {
	var Exist Voucher
	result := c.DB.First(&Exist, id).Error
	if result != nil {
		return result
	}

	Exist.Status = 0
	err := c.DB.Save(&Exist).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *DBClient) GetAllVoucher(page, size int) ([]Voucher, int64, error) {
	var vouchers []Voucher

	offset := (page - 1) * size

	rows := c.DB.Limit(size).Offset(offset).Where("status = ?", 1).Find(&vouchers).Error
	if rows != nil {
		return nil, 0, rows
	}

	var totalVouchers int64

	rows = c.DB.Model(&Voucher{}).Count(&totalVouchers).Error
	if rows != nil {
		return nil, 0, rows
	}

	return vouchers, totalVouchers, nil
}

func (c *DBClient) ExportVouchersToCSV(w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"voucher_code", "discount_percent", "expiry_date"})

	var vouchers []Voucher
	if err := c.DB.Where("status = ?", 1).Find(&vouchers).Error; err != nil {
		return err
	}

	for _, v := range vouchers {
		parsedDate, err := time.Parse(time.RFC3339, v.ExpiryDate)
		if err != nil {
			return err
		}
		formattedDate := parsedDate.Format("02/01/2006")

		record := []string{
			v.VoucherCode,
			strconv.Itoa(v.DiscountPercent),
			formattedDate,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
