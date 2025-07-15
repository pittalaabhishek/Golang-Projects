// package main

// import (
// 	"fmt"
// 	"strconv"
// 	"time"

// 	"github.com/jung-kurt/gofpdf"
// )

// type InvoiceItem struct {
// 	UnitName       string
// 	PricePerUnit   int
// 	UnitsPurchased int
// }

// func main() {
// 	data := []InvoiceItem{
// 		{
// 			UnitName:       "2x6 Lumber - 8'",
// 			PricePerUnit:   375,
// 			UnitsPurchased: 220,
// 		},
// 		{
// 			UnitName:       "Drywall Sheet",
// 			PricePerUnit:   822,
// 			UnitsPurchased: 50,
// 		},
// 		{
// 			UnitName:       "Paint",
// 			PricePerUnit:   1455,
// 			UnitsPurchased: 3,
// 		},
// 	}

// 	generateInvoice(data, "invoice.pdf")
// 	fmt.Println("Invoice generated: invoice.pdf")
// }

// func generateInvoice(items []InvoiceItem, filename string) {
// 	pdf := gofpdf.New("P", "mm", "A4", "")
// 	pdf.AddPage()

// 	pdf.SetFont("Arial", "B", 20)
// 	pdf.Cell(190, 10, "INVOICE")
// 	pdf.Ln(15)

// 	pdf.SetFont("Arial", "B", 12)
// 	pdf.Cell(95, 6, "BuildCorp Materials")
// 	pdf.Ln(6)
// 	pdf.SetFont("Arial", "", 10)
// 	pdf.Cell(95, 5, "123 Construction Ave")
// 	pdf.Ln(5)
// 	pdf.Cell(95, 5, "Builder City, BC 12345")
// 	pdf.Ln(5)
// 	pdf.Cell(95, 5, "Phone: (555) 123-4567")
// 	pdf.Ln(15)

// 	pdf.SetFont("Arial", "B", 10)
// 	pdf.Cell(40, 6, "Invoice #: INV-2024-001")
// 	pdf.Ln(6)
// 	pdf.Cell(40, 6, "Date: "+time.Now().Format("2006-01-02"))
// 	pdf.Ln(6)
// 	pdf.Cell(40, 6, "Due Date: "+time.Now().AddDate(0, 0, 30).Format("2006-01-02"))
// 	pdf.Ln(15)
	
// 	pdf.SetFont("Arial", "B", 12)
// 	pdf.Cell(40, 6, "Bill To:")
// 	pdf.Ln(8)
// 	pdf.SetFont("Arial", "", 10)
// 	pdf.Cell(40, 5, "John Contractor")
// 	pdf.Ln(5)
// 	pdf.Cell(40, 5, "456 Builder St")
// 	pdf.Ln(5)
// 	pdf.Cell(40, 5, "Construction Town, CT 67890")
// 	pdf.Ln(15)

// 	pdf.SetFont("Arial", "B", 10)
// 	pdf.SetFillColor(220, 220, 220)
// 	pdf.CellFormat(80, 8, "Description", "1", 0, "L", true, 0, "")
// 	pdf.CellFormat(30, 8, "Price Each", "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(30, 8, "Quantity", "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(30, 8, "Total", "1", 0, "C", true, 0, "")
// 	pdf.Ln(8)

// 	pdf.SetFont("Arial", "", 10)
// 	pdf.SetFillColor(255, 255, 255)
	
// 	var grandTotal int
// 	for _, item := range items {
// 		total := item.PricePerUnit * item.UnitsPurchased
// 		grandTotal += total

// 		pdf.CellFormat(80, 8, item.UnitName, "1", 0, "L", false, 0, "")
// 		pdf.CellFormat(30, 8, formatCurrency(item.PricePerUnit), "1", 0, "C", false, 0, "")
// 		pdf.CellFormat(30, 8, strconv.Itoa(item.UnitsPurchased), "1", 0, "C", false, 0, "")
// 		pdf.CellFormat(30, 8, formatCurrency(total), "1", 0, "C", false, 0, "")
// 		pdf.Ln(8)
// 	}

// 	pdf.Ln(5)
// 	pdf.SetFont("Arial", "B", 12)
// 	pdf.Cell(140, 10, "")
// 	pdf.Cell(20, 10, "TOTAL:")
// 	pdf.CellFormat(30, 10, formatCurrency(grandTotal), "1", 0, "C", false, 0, "")
// 	pdf.Ln(15)

// 	pdf.SetFont("Arial", "I", 10)
// 	pdf.Cell(190, 5, "Thank you for your business!")
// 	pdf.Ln(5)
// 	pdf.Cell(190, 5, "Payment due within 30 days")

// 	err := pdf.OutputFileAndClose(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func formatCurrency(cents int) string {
// 	dollars := float64(cents) / 100.0
// 	return fmt.Sprintf("$%.2f", dollars)
// }
package main

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	studentName := "Pittala Abhishek"
	
	generateStyledCertificate(studentName, "gophercises_certificate.pdf")
	fmt.Println("Certificate generated: gophercises_certificate.pdf")
}

func generateStyledCertificate(name, filename string) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFillColor(250, 250, 250)
	pdf.Rect(0, 0, 297, 210, "F")

	// Decorative border
	pdf.SetLineWidth(3)
	pdf.SetDrawColor(0, 50, 100)
	pdf.Rect(20, 20, 257, 170, "D")
	
	pdf.SetLineWidth(1)
	pdf.SetDrawColor(0, 100, 200)
	pdf.Rect(25, 25, 247, 160, "D")

	pdf.SetFillColor(0, 50, 100)
	pdf.Rect(40, 40, 217, 25, "F")
	
	pdf.SetFont("Arial", "B", 28)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(40, 48)
	pdf.CellFormat(217, 15, "CERTIFICATE OF COMPLETION", "", 0, "C", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 14)
	pdf.SetXY(40, 85)
	pdf.CellFormat(217, 10, "This certifies that", "", 0, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(0, 100, 0)
	pdf.SetXY(40, 105)
	pdf.CellFormat(217, 15, name, "", 0, "C", false, 0, "")
	pdf.Line(80, 122, 217, 122)

	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(40, 130)
	pdf.CellFormat(217, 10, "has successfully completed", "", 0, "C", false, 0, "")
	
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(0, 0, 150)
	pdf.SetXY(40, 145)
	pdf.CellFormat(217, 12, "GOPHERCISES - Go Programming Exercises", "", 0, "C", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(0, 0, 0)
	currentDate := time.Now().Format("July 15, 2025")
	
	pdf.SetXY(60, 175)
	pdf.CellFormat(60, 8, "Date: "+currentDate, "", 0, "L", false, 0, "")
	
	pdf.SetXY(180, 175)
	pdf.CellFormat(60, 8, "Instructor: Jon Calhoun", "", 0, "L", false, 0, "")

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		panic(err)
	}
}