package main

// programa obtenido de https://appliedgo.net/pdf/
// este diseño esta pensado solo para un pdf simple de un archivo especifico
// tiene que ser modificado depende de la cantidad de columnas

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	// aqui primero se carga el archivo CSV
	// la funcion path() solo es un pequeño intermedio para obtener el path de os.Args
	data := loadCSV(path())

	col := len(data[0])

	// una vez tenemos la informacion se crea un documento pdf
	pdf := newReport()
	pdf = image(pdf)
	pdf = header(pdf, data[0], col)
	pdf = table(pdf, data[1:], col)

	if pdf.Err() {
		log.Fatalf("Failed creating PDF report: %s\n", pdf.Error())
	}

	err := savePDF(pdf, "reporte")
	if err != nil {
		log.Fatalf("Cannot save PDF: %s|n", err)
	}
}

func loadCSV(path string) [][]string {
	// en este punto se obtiene toda la informacion del CSV,
	// se abre el archivo, se lee, y se regresan las filas del archivo leido
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", path, err.Error())
	}
	defer f.Close()
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatalln("Cannot read CSV data:", err.Error())
	}
	return rows
}

func path() string {
	if len(os.Args) < 2 {
		return "cities.csv" // aqui es donde se da el nombre del archivo que se usara si no se entrega algun argumento
	}
	return os.Args[1]
}

func newReport() *gofpdf.Fpdf {
	// la creacion del pdf es simple, y esta dado por:
	// .New(
	// 	"Orientacion",  esta puede ser "L" landscape (paisaje: horizontal) o  "P" portrait (retrato: vertical) "" = "P"
	//	"Unidad de medida", esta es la unidad de medida en la que se pasaran los tamaños ("pt", "cm", "in", "mm") "" = "mm"
	//	"Formato", esto se refiere a el formato de la hoja del pdf ("A3", "A4", "A5", "Letter", "Legal", "Tabloid") "" = "A4"
	//	"Font", este valor es una direccion hacia un directorio para usar alguna font especifica descargada
	// )
	// Los valores vacios tienen asignados una opcion default cada uno
	pdf := gofpdf.New("L", "mm", "Letter", "")

	pdf.AddPage()

	// set font te establece el tipo de letra que se usara, ya sea un tipo estandart o,
	// algun tipo agregado con addFont o AddFontFromReader, etc.

	// se construye:
	// .SetFont(
	//	"familyStr", Este parametro recibe el tipo de letra que se quiere usar
	// 	"styleStr", puede ser "B"(negritas), "I"(cursiva/inclinada), "U"(subrayado), "S"(tachado)
	// )
	pdf.SetFont("Times", "B", 28)

	//.Cell es la version simple de CellFormat,
	// solo se tienenn los tamaños, el texto y listo, no tiene bordes, ni aliniacion ni nada de eso
	pdf.Cell(40, 10, "Daily Report")

	pdf.Ln(12) // .Ln es un salto de linea, se le ingresa la distancia que se movera

	pdf.SetFont("Times", "", 20)
	pdf.Cell(40, 10, time.Now().Format("Mon Jan 2, 2006"))
	pdf.Ln(20)

	return pdf

}

func header(pdf *gofpdf.Fpdf, hdr []string, columnas int) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 16)
	pdf.SetFillColor(240, 240, 240)
	pageWidth, _ := pdf.GetPageSize()
	cellWidth := (pageWidth - 20) / float64(columnas)
	for _, str := range hdr {

		pdf.CellFormat(cellWidth, 7, str, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	return pdf
}

func table(pdf *gofpdf.Fpdf, tbl [][]string, columnas int) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 16)
	pdf.SetFillColor(255, 255, 255)

	// Aqui se toma los valores de la pagina para establecer que tan grande puede ser la celda a crear
	// se toma el tamaño del width, y se le restan 20 mm que seria aproximadamente los margenes

	pageWidth, _ := pdf.GetPageSize()
	cellWidth := (pageWidth - 20) / float64(columnas)

	for _, line := range tbl {

		for _, str := range line {

			pdf.CellFormat(cellWidth, 7, str, "1", 0, "C", false, 0, "")

		}
		pdf.Ln(-1)
	}
	return pdf
}

func image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {

	// para colocar una imagen, puedes usar .ImageOptions,
	// .ImageOptions usa siempre la pagina actual, los parametros que usa son
	// "ImageName", recibe un string con el nombre de la imagen a usar
	// "x", "y" son la posicion dentro de la pagina actual donde se colocara la imagen
	// "h", "w" son los tamaños que tomara la imagen
	// "flow", si es verdadero, el valor de "y" actual se desplaza y puede hacer un salto de pagina (verdadero  hace salto de linea, falso se conserva la linea actual)
	// "options", es un struct para opciones adicionales que se le pueden agregar

	pdf.ImageOptions("stats.png", 225, 10, 25, 25, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

func savePDF(pdf *gofpdf.Fpdf, name string) error {
	fileName := fmt.Sprintf("%s.pdf", name)
	return pdf.OutputFileAndClose(fileName)
}
