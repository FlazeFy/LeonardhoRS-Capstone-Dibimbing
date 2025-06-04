package utils

import (
	"fmt"
	"pelita/entity"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDFAssetFindingReport(c []entity.AssetFindingReport, filename string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetTitle("PELITA", false)
	pdf.AddPage()

	// Set Header
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 102, 204)
	pdf.CellFormat(0, 12, "PELITA", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "Pemeliharaan Inventaris Asset", "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// Set Letterhead
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Audit - Asset Finding")
	pdf.Ln(8)

	// Set header
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(40, 9, "Asset", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 9, "Category", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 9, "Notes", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 9, "Find At", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 9, "Floor - Room", "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 9, "Maintenance PIC", "1", 1, "C", true, 0, "")

	// Set body
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	for _, dt := range c {
		pdf.CellFormat(40, 8, dt.AssetName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 8, dt.FindingCategory, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 8, dt.FindingNotes, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 8, dt.CreatedAt.Format("2006-01-02 15:04:05"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%s - %s", dt.Floor, dt.RoomName), "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 8, fmt.Sprintf("%s - %s", dt.Username, dt.Email), "1", 1, "L", false, 0, "")
	}

	return pdf.OutputFileAndClose(filename)
}
