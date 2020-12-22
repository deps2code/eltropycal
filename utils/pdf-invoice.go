package utils

import (
	"strconv"
	"time"

	"github.com/eltropycal/models/response"
	"github.com/jung-kurt/gofpdf"
)

func GenerateInvoicePdf(order response.OrderResponse, pdfFileName string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	setPDFFooter(pdf)
	setPDFHeader(pdf)

	pdf.AddPage()

	pdf.SetFont("Times", "B", 15)
	pdf.Ln(20)
	pdf.CellFormat(0, 10, "Date:  "+time.Now().Format("2006-01-02"), "", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.CellFormat(0, 10, "Customer", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, "Restaurant", "", 0, "R", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(0, 10, "Name: "+order.User.Name, "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, order.RestaurantAddress.AddressText, "", 0, "R", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(0, 10, "Address: "+order.DeliveryAddress.AddressText, "", 0, "L", false, 0, "")

	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(80, 10, "Items", "TB", 0, "L", false, 0, "")
	pdf.CellFormat(33, 10, "Quantity", "TB", 0, "L", false, 0, "")
	pdf.CellFormat(33, 10, "Rate", "TB", 0, "L", false, 0, "")
	pdf.CellFormat(33, 10, "Total Rate", "TB", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Times", "", 12)
	var totalBillAmount float64
	for _, item := range order.Items {
		totalBillAmount += item.Price * float64(item.Quantity)
		pdf.CellFormat(80, 10, item.FoodName, "TB", 0, "L", false, 0, "")
		pdf.CellFormat(33, 10, strconv.FormatInt(int64(item.Quantity), 10), "TB", 0, "L", false, 0, "")
		pdf.CellFormat(33, 10, strconv.FormatFloat(item.Price, 'f', 2, 64), "TB", 0, "L", false, 0, "")
		pdf.CellFormat(33, 10, strconv.FormatFloat(item.Price*float64(item.Quantity), 'f', 2, 64), "TB", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}
	deliveryDistance := Distance(order.RestaurantAddress.Lat, order.RestaurantAddress.Lng, order.DeliveryAddress.Lat, order.DeliveryAddress.Lng)
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(80, 10, "Total(Rupees)", "TB", 0, "R", false, 0, "")
	pdf.CellFormat(33, 10, "", "TB", 0, "L", false, 0, "")
	pdf.CellFormat(33, 10, "", "TB", 0, "L", false, 0, "")
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(33, 10, strconv.FormatFloat(totalBillAmount, 'f', 2, 64), "TB", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(80, 10, "Tax(5%)", "TB", 0, "R", false, 0, "")
	pdf.CellFormat(33, 10, "", "TB", 0, "L", false, 0, "")
	pdf.CellFormat(33, 10, "", "TB", 0, "L", false, 0, "")
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(33, 10, strconv.FormatFloat(totalBillAmount*0.05, 'f', 2, 64), "TB", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(80, 10, "Delivery Charge", "TB", 0, "R", false, 0, "")
	pdf.CellFormat(33, 10, "", "TB", 0, "L", false, 0, "")
	pdf.CellFormat(33, 10, "", "TB", 0, "L", false, 0, "")
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(33, 10, strconv.FormatFloat(deliveryDistance/1000, 'f', 2, 64), "TB", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(80, 10, "Total paid amount", "TB", 0, "R", false, 0, "")
	pdf.CellFormat(33, 10, "", "TB", 0, "L", false, 0, "")
	pdf.CellFormat(33, 10, "", "TB", 0, "L", false, 0, "")
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(33, 10, strconv.FormatFloat(totalBillAmount+(totalBillAmount*0.05)+(deliveryDistance/1000), 'f', 2, 64), "TB", 0, "L", false, 0, "")

	err := savePDF(pdf, pdfFileName)
	if err != nil {
		return err
	}
	return nil
}

func setPDFFooter(pdf *gofpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetFillColor(18, 84, 78)
		pdf.SetDrawColor(18, 84, 78)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetY(-20)
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 7, "Powered by:",
			"1", 0, "R", true, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(0, 7, "Eltropycal",
			"1", 0, "R", true, 0, "")
	})
}

func setPDFHeader(pdf *gofpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		pdf.SetFillColor(18, 84, 78)
		pdf.SetDrawColor(18, 84, 78)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Arial", "", 20)
		pdf.CellFormat(0, 15, "Invoice",
			"", 0, "C", true, 0, "")
		pdf.Ln(20)
	})
}

func savePDF(pdf *gofpdf.Fpdf, filename string) error {
	return pdf.OutputFileAndClose(filename + ".pdf")
}
