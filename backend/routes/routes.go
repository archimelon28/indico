package routes

import (
	"backend/controller"
	"backend/middleware"
	"backend/models"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetupRouter(dbClient *models.DBClient) http.Handler {
	r := mux.NewRouter()

	voucherController := controller.NewVoucherController(dbClient)

	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")

	r.Handle("/voucher", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		voucherController.GetAllVoucher(w, r)
	}))).Methods("GET")
	r.Handle("/voucher", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		voucherController.CreateVoucher(w, r)
	}))).Methods("POST")
	r.Handle("/voucher/upload-csv", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		voucherController.UploadVoucherCSV(w, r)
	}))).Methods("POST")
	r.Handle("/voucher/export", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		voucherController.ExportCSV(w, r)
	}))).Methods("GET")
	r.Handle("/voucher/{id}", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		voucherController.UpdateVoucher(w, r)
	}))).Methods("PUT")
	r.Handle("/voucher/{id}", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		voucherController.GetVoucherById(w, r)
	}))).Methods("GET")
	r.Handle("/voucher/{id}", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		voucherController.DeleteVoucher(w, r)
	}))).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return c.Handler(r)
}
