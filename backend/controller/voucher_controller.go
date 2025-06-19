package controller

import (
	"backend/models"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type VoucherController struct {
	DBClient *models.DBClient
}

func NewVoucherController(dbClient *models.DBClient) *VoucherController {
	return &VoucherController{DBClient: dbClient}
}

func (c *VoucherController) CreateVoucher(write http.ResponseWriter, read *http.Request) {
	var voucher models.Voucher

	err := json.NewDecoder(read.Body).Decode(&voucher)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	voucher.Status = 1

	err = c.DBClient.CreateVoucher(&voucher)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusCreated)
	json.NewEncoder(write).Encode(voucher)
}

func (c *VoucherController) UpdateVoucher(write http.ResponseWriter, read *http.Request) {
	idStr := mux.Vars(read)["id"]

	id, _ := strconv.Atoi(idStr)

	var voucher models.Voucher

	err := json.NewDecoder(read.Body).Decode(&voucher)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.DBClient.UpdateVoucher(&voucher, id)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(voucher)
}

func (c *VoucherController) GetVoucherById(write http.ResponseWriter, read *http.Request) {
	idStr := mux.Vars(read)["id"]

	id, _ := strconv.Atoi(idStr)

	user, err := c.DBClient.GetVoucherById(id)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusOK)
	json.NewEncoder(write).Encode(user)
}

func (c *VoucherController) DeleteVoucher(write http.ResponseWriter, read *http.Request) {
	idStr := mux.Vars(read)["id"]

	id, _ := strconv.Atoi(idStr)

	err := c.DBClient.DeleteVoucher(id)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	write.WriteHeader(http.StatusNoContent)
}

func (c *VoucherController) GetAllVoucher(write http.ResponseWriter, read *http.Request) {
	page, err := strconv.Atoi(read.URL.Query().Get("page"))
	if err != nil || page < 0 {
		page = 1
	}

	size, err := strconv.Atoi(read.URL.Query().Get("size"))
	if err != nil || size < 0 {
		size = 10
	}

	posts, totalPosts, err := c.DBClient.GetAllVoucher(page, size)
	if err != nil {
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalPosts + int64(size) - 1) / int64(size)

	response := struct {
		Page       int              `json:"page"`
		Size       int              `json:"size"`
		TotalItems int64            `json:"totalItems"`
		TotalPages int64            `json:"totalPages"`
		Data       []models.Voucher `json:"data"`
	}{
		Page:       page,
		Size:       size,
		TotalItems: totalPosts,
		TotalPages: totalPages,
		Data:       posts,
	}

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)

	err = json.NewEncoder(write).Encode(response)
	if err != nil {
		http.Error(write, "Error encoding response", http.StatusInternalServerError)
	}

}

func (c *VoucherController) UploadVoucherCSV(write http.ResponseWriter, read *http.Request) {
	read.ParseMultipartForm(10 << 20)

	file, _, err := read.FormFile("file")
	if err != nil {
		http.Error(write, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var successCount, failCount int
	var errorRows []int
	var vouchers []models.Voucher

	_, _ = reader.Read()

	i := 2
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) != 3 {
			failCount++
			errorRows = append(errorRows, i)
			i++
			continue
		}

		percent, err := strconv.Atoi(record[1])
		if err != nil || percent < 1 || percent > 100 {
			failCount++
			errorRows = append(errorRows, i)
			i++
			continue
		}

		vouchers = append(vouchers, models.Voucher{
			VoucherCode:     record[0],
			DiscountPercent: percent,
			ExpiryDate:      record[2],
			Status:          1,
		})
		successCount++
		i++
	}

	if len(vouchers) > 0 {
		c.DBClient.CreateBulkVouchers(&vouchers)
	}

	json.NewEncoder(write).Encode(map[string]interface{}{
		"success": successCount,
		"failed":  failCount,
		"errors":  errorRows,
	})
}

func (c *VoucherController) ExportCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=vouchers.csv")

	err := c.DBClient.ExportVouchersToCSV(w)
	if err != nil {
		http.Error(w, "Gagal export voucher", http.StatusInternalServerError)
	}
}
