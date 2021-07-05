package main

// Programa obtenido de https://appliedgo.net/pdf/

// En este caso se encuentran 2 archivos csv para probar como seria con distinta cantidad de columnas
// casos a tomar en cuenta como el tamaño del string en la cantidad de columnas, etc

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func main() {
	// aqui primero se carga el archivo CSV
	// la funcion path() solo es un pequeño intermedio para obtener el path de os.Args
	data, err := loadCSV(path())
	if err != nil {
		log.Printf("Failed loading PDF report: %s\n", err)
		return
	}
	var col int

	if len(data) != 0 {
		col = len(data[0])
	}

	// una vez tenemos la informacion se crea un documento pdf
	pdf := newReport() // creacion del la instancia pdf en la que se trabajara
	bytes, docType := imageToBytes("stats.png")
	insertImageFromBytes(pdf, bytes, docType)
	// pdf = image(pdf)                // le agregamos una imagen estetica al documento
	pdf = header(pdf, data[0], col) // se le añade el header a la tabla, los nombres de los datos que tiene
	pdf = table(pdf, data[1:], col) // se le añaden todos los datos a la tabla

	if pdf.Err() {
		log.Printf("Failed creating PDF report: %s\n", pdf.Error())
		return
	}

	err = savePDF(pdf, "reporte")
	if err != nil {
		log.Printf("Cannot save PDF: %s|n", err)
		return
	}
}

func loadCSV(path string) ([][]string, error) {
	// en este punto se obtiene toda la informacion del CSV,
	// se abre el archivo, se lee, y se regresan las filas del archivo leido

	var rows [][]string
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Cannot open '%s': %s\n", path, err.Error())
		f.Close()
		return rows, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	rows, err = r.ReadAll()
	if err != nil {
		log.Printf("Cannot read CSV data: %s", err.Error())
		return rows, err
	}
	return rows, err
}

func path() string {
	if len(os.Args) < 2 {
		return "cities.csv" // aqui es donde se da el nombre del archivo que se usara si no se entrega algun argumento
	}
	return os.Args[1]
}

func imageToBytes(path string) ([]byte, string) {

	// Esta funcion toma la imagen de local y la convierte a un []byte para su uso posterior
	// El motivo de esto es con fines practicos, ya que los servicios podrian mandar la imagen de esta forma
	// esto es para emular como tendrias que manejar la imagen si te la dan en bytes y no se tiene el archivo en local

	file, err := os.Open(path)

	bytes := make([]byte, 5)

	if err != nil {
		log.Printf("An error has occurred, Error: %s", err)
		file.Close()
		return bytes, ""
	}

	defer file.Close()
	nombre := file.Name()
	partes := strings.Split(nombre, ".")
	docType := ""
	if len(partes) > 1 {
		docType = partes[(len(partes) - 1)]
	}

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes = make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		log.Printf("Error: %s", err)
		return bytes, ""
	}

	return bytes, docType
}
