package comandos

import (
	"backend/estructuras"
	"backend/manejadorDisco"
	"backend/utilidades"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Reportes(id string, path string, name string, linea string) {
	switch name {
	case "mbr":
		repMBR(id, path, linea)
	default:
		fmt.Println("Tipo de reporte no encontrado")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Tipo de reporte no encontrado")
	}
}

func repMBR(id string, path string, linea string) {
	fmt.Println("====== INICIO REP ======")
	fmt.Println("Id:", id)
	fmt.Println("Path:", path)

	var mountedPartition manejadorDisco.ParticionMontada
	var particionEncontrada bool

	for _, partitions := range manejadorDisco.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.ID == id {
				mountedPartition = partition
				particionEncontrada = true
				break
			}
		}
		if particionEncontrada {
			break
		}
	}

	if !particionEncontrada {
		fmt.Println("Partición no encontrada")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Partición no encontrada")
		return
	}

	file, err := utilidades.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Error al abrir el archivo")
		return
	}
	defer file.Close()

	var TempMBR estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Error al leer el MBR")
		return
	}

	textoDot := "digraph G {\n"
	textoDot += "node [shape=none];\n"
	textoDot += "tablaMBR [label=<\n"
	textoDot += "<table border='1' cellborder='1' cellspacing='0'>\n"
	textoDot += "<tr><td colspan='2' bgcolor=\"#84b6f4\">MBR</td></tr>\n"
	textoDot += fmt.Sprintf("<tr><td bgcolor=\"#eaf7fb\">Size</td><td bgcolor=\"#eaf7fb\">%d</td></tr>\n", TempMBR.MbrSize)
	textoDot += fmt.Sprintf("<tr><td bgcolor=\"#eaf7fb\">Signature</td><td bgcolor=\"#eaf7fb\">%d</td></tr>\n", TempMBR.Signature)
	textoDot += fmt.Sprintf("<tr><td bgcolor=\"#eaf7fb\">Fecha de creacion</td><td bgcolor=\"#eaf7fb\">%s</td></tr>\n", string(TempMBR.CreationDate[:]))

	for i, partition := range TempMBR.Partitions {
		if partition.Status[0] != 0 {

			textoDot += fmt.Sprintf("<tr><td colspan='2' bgcolor=\"#77dd77\">Particion %d</td></tr>\n", i+1)
			textoDot += fmt.Sprintf("<tr><td bgcolor=\"#d8f8e1\">Status</td><td bgcolor=\"#d8f8e1\">%d</td></tr>\n", partition.Status[0])
			textoDot += fmt.Sprintf("<tr><td bgcolor=\"#d8f8e1\">Type</td><td bgcolor=\"#d8f8e1\">%s</td></tr>\n", string(partition.Type[:]))
			textoDot += fmt.Sprintf("<tr><td bgcolor=\"#d8f8e1\">Fit</td><td bgcolor=\"#d8f8e1\">%s</td></tr>\n", string(partition.Fit[:]))
			textoDot += fmt.Sprintf("<tr><td bgcolor=\"#d8f8e1\">Start</td><td bgcolor=\"#d8f8e1\">%d</td></tr>\n", partition.Start)
			textoDot += fmt.Sprintf("<tr><td bgcolor=\"#d8f8e1\">Size</td><td bgcolor=\"#d8f8e1\">%d</td></tr>\n", partition.Size)
			textoDot += fmt.Sprintf("<tr><td bgcolor=\"#d8f8e1\">Name</td><td bgcolor=\"#d8f8e1\">%s</td></tr>\n", strings.Trim(string(partition.Name[:]), "\x00"))

			if partition.Type[0] == 'e' || partition.Type[0] == 'E' {
				ebrOffset := partition.Start
				logicalPartitionCount := 1
				for {
					var ebr estructuras.EBR
					if err := utilidades.ReadObject(file, &ebr, int64(ebrOffset)); err != nil {
						fmt.Println("Error al leer EBR:", err)
						utilidades.AgregarRespuesta("Error en linea " + linea + " : Error al leer EBR")
						break
					}

					textoDot += "<tr><td colspan='2' bgcolor=\"#fdfd96\">EBR</td></tr>\n"
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Next</td><td bgcolor=\"#fdf9c4\">%d</td></tr>\n", ebr.PartNext)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Size</td><td bgcolor=\"#fdf9c4\">%d</td></tr>\n", ebr.PartSize)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Start</td><td bgcolor=\"#fdf9c4\">%d</td></tr>\n", ebr.PartStart)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Name</td><td bgcolor=\"#fdf9c4\">%s</td></tr>\n", strings.Trim(string(ebr.PartName[:]), "\x00"))

					textoDot += fmt.Sprintf("<tr><td colspan='2' bgcolor=\"#ff6961\">Particion logica %d</td></tr>\n", logicalPartitionCount)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#ffe4e1\">Fit</td><td bgcolor=\"#ffe4e1\">%s</td></tr>\n", string(ebr.PartFit[:]))
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#ffe4e1\">Start</td><td bgcolor=\"#ffe4e1\">%d</td></tr>\n", ebr.PartStart)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#ffe4e1\">Size</td><td bgcolor=\"#ffe4e1\">%d</td></tr>\n", ebr.PartSize)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#ffe4e1\">Name</td><td bgcolor=\"#ffe4e1\">%s</td></tr>\n", strings.Trim(string(ebr.PartName[:]), "\x00"))

					logicalPartitionCount++

					if ebr.PartNext <= 0 {
						break
					}
					ebrOffset = ebr.PartNext
				}
			}
		}
	}

	textoDot += "</table>\n"
	textoDot += ">];\n"
	textoDot += "}\n"

	dotFilePath := "/home/jd/temps/mbr.dot"
	err = os.WriteFile(dotFilePath, []byte(textoDot), 0644)
	if err != nil {
		utilidades.AgregarRespuesta("Error al escribir el archivo DOT")
		fmt.Println("Error al escribir el archivo DOT:", err)
		return
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			utilidades.AgregarRespuesta("Error al crear el directorio")
			fmt.Println("Error al crear el directorio:", err)
			return
		}
	}

	cmd := exec.Command("dot", "-Tjpg", dotFilePath, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		utilidades.AgregarRespuesta("Error al ejecutar Graphviz")
		fmt.Println("Error al ejecutar Graphviz:", err)
		fmt.Println("Detalles del error:", stderr.String())
		return
	}

	utilidades.AgregarRespuesta("Reporte de MBR generado exitosamente en " + path)
	fmt.Println("Reporte de MBR generado exitosamente")
	fmt.Println("====== FIN REP ======")
}
