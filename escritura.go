package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func newReport() *gofpdf.Fpdf {
	// la creacion del pdf es simple, y esta dado por:
	// .New(
	// 	"orientationStr",  esta puede ser "L" landscape (paisaje: horizontal) o  "P" portrait (retrato: vertical) "" = "P"
	//	"unitStr", esta es la unidad de medida en la que se pasaran los tamaños ("pt", "cm", "in", "mm") "" = "mm"
	//	"sizeStr", esto se refiere a el formato de la hoja del pdf ("A3", "A4", "A5", "Letter", "Legal", "Tabloid") "" = "A4"
	//	"fontDirStr", este valor es una direccion hacia un directorio para usar alguna font especifica descargada
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

	// NOTA: los valores pueden salirse de su celda si la info es muy grande,
	// algo que se puede hacer es crear una lista con tamaños para cada columna
	// en una columna de 5 hacer  []float64 {20, 30, 25, 20, 20, 50}
	// esto ya cuando sea necesario y se tenga conocimiento de como es la informacion que se recibira

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

func savePDF(pdf *gofpdf.Fpdf, name string) error {
	fileName := fmt.Sprintf("%s.pdf", name)
	return pdf.OutputFileAndClose(fileName)
}

func Image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {

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

func insertImageFromBytes(pdf *gofpdf.Fpdf, data []byte, docType string) {

	// Otra forma de insertar una imagen en pdf es usando los bytes de la imagen y no el archivo de la imagen como tal
	// esto es necesario cuando un servicio nos proporciona un slice de bytes en vez del archivo

	colocarImagen := true
	imageData := bytes.NewBuffer(data)
	if len(data) == 0 || docType == "" {
		log.Println("Error: No image data or docType miss")
		colocarImagen = false
	}

	if colocarImagen {
		// Se registra la imagen al pdf (solo se registra para su uso posterior)
		pdf.RegisterImageOptionsReader("Imagen", gofpdf.ImageOptions{ImageType: docType, ReadDpi: false}, imageData)

		// se inserta la imagen, la opcion Image toma de las imagenes registradas, buscando por el nombre que se le dio y agrega esa imagen a la pagina actual
		pdf.Image("Imagen", 225, 10, 25, 25, false, "", 0, "")

	}

}
