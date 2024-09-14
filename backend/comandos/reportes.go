package comandos

import (
	"backend/estructuras"
	"backend/manejadorDisco"
	"backend/utilidades"
	"bytes"
	"fmt"
	"html"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Reportes(id string, path string, name string, linea string) {
	switch name {
	case "mbr":
		repMBR(id, path, linea)
	case "disk":
		repDisk(id, path, linea)
	case "inode":
		repInode(id, path)
	case "sb":
		reportSuperBloque(id, path)
	case "bm_inode":
		repbmInodes(id, path)
	case "bm_block":
		rep_bmBlock(id, path)
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
				finEbr := partition.Start
				contPartLogic := 1
				for {
					var ebr estructuras.EBR
					if err := utilidades.ReadObject(file, &ebr, int64(finEbr)); err != nil {
						fmt.Println("Error al leer EBR:", err)
						utilidades.AgregarRespuesta("Error en linea " + linea + " : Error al leer EBR")
						break
					}

					textoDot += "<tr><td colspan='2' bgcolor=\"#fdfd96\">EBR</td></tr>\n"
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Next</td><td bgcolor=\"#fdf9c4\">%d</td></tr>\n", ebr.PartNext)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Size</td><td bgcolor=\"#fdf9c4\">%d</td></tr>\n", ebr.PartSize)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Start</td><td bgcolor=\"#fdf9c4\">%d</td></tr>\n", ebr.PartStart)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Name</td><td bgcolor=\"#fdf9c4\">%s</td></tr>\n", strings.Trim(string(ebr.PartName[:]), "\x00"))

					textoDot += fmt.Sprintf("<tr><td colspan='2' bgcolor=\"#ff6961\">Particion logica %d</td></tr>\n", contPartLogic)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#ffe4e1\">Fit</td><td bgcolor=\"#ffe4e1\">%s</td></tr>\n", string(ebr.PartFit[:]))
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#ffe4e1\">Start</td><td bgcolor=\"#ffe4e1\">%d</td></tr>\n", ebr.PartStart)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#ffe4e1\">Size</td><td bgcolor=\"#ffe4e1\">%d</td></tr>\n", ebr.PartSize)
					textoDot += fmt.Sprintf("<tr><td bgcolor=\"#ffe4e1\">Name</td><td bgcolor=\"#ffe4e1\">%s</td></tr>\n", strings.Trim(string(ebr.PartName[:]), "\x00"))

					contPartLogic++

					if ebr.PartNext <= 0 {
						break
					}
					finEbr = ebr.PartNext
				}
			}
		}
	}

	textoDot += "</table>\n"
	textoDot += ">];\n"
	textoDot += "}\n"

	rutaDot := "/home/jd/temps/mbr.dot"
	err = os.WriteFile(rutaDot, []byte(textoDot), 0644)
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

	cmd := exec.Command("dot", "-Tjpg", rutaDot, "-o", path)
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

func repDisk(id string, path string, linea string) {
	fmt.Println("====== INICIO REP DISK ======")
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

	// Variables para calcular el porcentaje
	totalSize := float64(TempMBR.MbrSize)
	usedSize := 0.0

	nombreConExtension := filepath.Base(mountedPartition.Path)

	textoDot := "digraph G {\n"
	textoDot += "label=\"" + nombreConExtension + "\"\n"
	textoDot += "labelloc=\"t\"\n"
	textoDot += "subgraph cluster1 {\n"
	textoDot += "label=\"\"\n"
	textoDot += "disco [shape=none label=<\n"
	textoDot += "<TABLE border=\"0\" cellspacing=\"4\" cellpadding=\"5\" color=\"blue\" >\n"
	textoDot += "<TR>\n"
	// MBR es siempre al inicio
	//textoDot += fmt.Sprintf("<tr><td bgcolor=\"#eaf7fb\">MBR</td><td bgcolor=\"#eaf7fb\">%d</td><td bgcolor=\"#eaf7fb\">%.2f%%</td></tr>\n", 512, (512/totalSize)*100)
	textoDot += "<TD border=\"1\"  cellpadding=\"65\">MBR</TD>\n"
	// Analizando particiones
	for i, partition := range TempMBR.Partitions {
		fmt.Println("Particion ", i+1, ": ")
		fmt.Println("Status: ", partition.Status[0])
		if partition.Status[0] != 0 {
			partSize := float64(partition.Size)
			fmt.Println("ESte esle tamano de la particion ", i+1, " : ", partSize)
			usedSize += partSize
			//textoDot += fmt.Sprintf("<tr><td bgcolor=\"#d8f8e1\">Particion %d</td><td bgcolor=\"#d8f8e1\">%d</td><td bgcolor=\"#d8f8e1\">%.2f%%</td></tr>\n", i+1)

			if partition.Type[0] == 'e' || partition.Type[0] == 'E' {
				finEbr := partition.Start
				contPartLogic := 1
				extSize := float64(partition.Size)
				extSizeD := partSize
				for {
					var ebr estructuras.EBR
					if err := utilidades.ReadObject(file, &ebr, int64(finEbr)); err != nil {
						fmt.Println("Error al leer EBR:", err)
						utilidades.AgregarRespuesta("Error en linea " + linea + " : Error al leer EBR")
						break
					}

					extSize -= float64(ebr.PartSize)
					contPartLogic++

					if ebr.PartNext <= 0 {
						if extSize > 0 {
							contPartLogic++
						}
						break
					}
					finEbr = ebr.PartNext
				}

				textoDot += "<TD border=\"1\" widht=\"75\">\n"
				textoDot += "<TABLE border=\"0\"  cellspacing=\"4\" cellpadding=\"10\">\n"
				textoDot += "<TR>\n"
				textoDot += fmt.Sprintf("<TD border=\"1\" colspan=\"%d\" cellpadding=\"75\">Extendida</TD>\n", contPartLogic+1)
				textoDot += "</TR>\n"
				textoDot += "<TR>\n"
				finEbrd := partition.Start
				for {
					var ebr estructuras.EBR
					if err := utilidades.ReadObject(file, &ebr, int64(finEbrd)); err != nil {
						fmt.Println("Error al leer EBR:", err)
						utilidades.AgregarRespuesta("Error en linea " + linea + " : Error al leer EBR")
						break
					}

					textoDot += "<TD border=\"1\" height=\"185\">EBR</TD>\n"

					ebrSize := float64(ebr.PartSize)
					fmt.Println("ESte esle tamano de la particion logica que se le esta sumando", i+1, " : ", ebrSize)
					usedSize += ebrSize
					extSizeD -= ebrSize

					//textoDot += fmt.Sprintf("<tr><td bgcolor=\"#fdf9c4\">Particion logica</td><td bgcolor=\"#fdf9c4\">%d</td><td bgcolor=\"#fdf9c4\">%.2f%%</td></tr>\n", ebr.PartSize, (ebrSize/totalSize)*100)
					textoDot += fmt.Sprintf("<TD border=\"1\" cellpadding=\"%d\">Logica<br/>%.2f%% por ciento del Disco</TD>\n", int(math.Round((ebrSize/totalSize)*100)), (ebrSize/totalSize)*100)
					fmt.Println("ESte esle tamano de la particion logica ", i+1, " : ", ebrSize)
					if ebr.PartNext <= 0 {
						if extSizeD > 0 {
							fmt.Println("ESte esle tamano de libre", i+1, " : ", extSizeD)
							textoDot += fmt.Sprintf("<TD border=\"1\" cellpadding=\"%d\">Libre<br/>%.2f%% por ciento del Disco</TD>\n", int(math.Round((extSizeD/totalSize)*100)), (extSizeD/totalSize)*100)
						}
						break
					}
					finEbrd = ebr.PartNext
				}
				textoDot += "</TR>\n"
				textoDot += "</TABLE>\n"
				textoDot += "</TD>\n"
			} else {
				textoDot += fmt.Sprintf("<TD border=\"1\" cellpadding=\"%d\">Primaria<br/>%.2f%% por ciento del Disco</TD>\n", int(math.Round((partSize/totalSize)*100)), (partSize/totalSize)*100)
			}
		}
	}

	// Espacio libre restante
	freeSize := totalSize - usedSize
	freePercentage := 100.0

	for _, partition := range TempMBR.Partitions {
		if partition.Status[0] != 0 {
			partSize := float64(partition.Size)
			freePercentage -= (partSize / totalSize) * 100
		}
	}
	fmt.Println("Espacio libre (calculado como complemento):", freePercentage)
	fmt.Println("ESte esle tamano total del disco ", totalSize)
	fmt.Printf("ESte esle tamano de la particion usada %.2f\n", usedSize)
	fmt.Println("ESte esle tamano de la particion libre ", freeSize)
	//textoDot += fmt.Sprintf("<tr><td bgcolor=\"#eaf7fb\">Libre</td><td bgcolor=\"#eaf7fb\">%.2f</td><td bgcolor=\"#eaf7fb\">%.2f%%</td></tr>\n", freeSize, (freeSize/totalSize)*100)
	textoDot += fmt.Sprintf("<TD border=\"1\" cellpadding=\"%d\">Libre<br/>%.2f%% por ciento del Disco</TD>\n", int(math.Round(freePercentage)), freePercentage)
	textoDot += "</TR>\n"
	textoDot += "</TABLE>\n"
	textoDot += ">];\n"
	textoDot += "}\n"
	textoDot += "}\n"

	// Guardar el archivo .dot y generar la imagen
	rutaDot := "/home/jd/temps/diskusage.dot"
	err = os.WriteFile(rutaDot, []byte(textoDot), 0644)
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

	cmd := exec.Command("dot", "-Tjpg", rutaDot, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		utilidades.AgregarRespuesta("Error al ejecutar Graphviz")
		fmt.Println("Error al ejecutar Graphviz:", err)
		fmt.Println("Detalles del error:", stderr.String())
		return
	}

	utilidades.AgregarRespuesta("Reporte de uso de disco generado exitosamente en " + path)
	fmt.Println("Reporte de uso de disco generado exitosamente")
	fmt.Println("====== FIN REP DISK ======")
}

func reportSuperBloque(id string, path string) {
	fmt.Println("-----INICIO REPORT_Superbloque----")
	fmt.Println("ID:", id)
	fmt.Println("Type (output path):", path)

	// Verificar si el usuario ya está logueado buscando en las particiones montadas
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
		utilidades.AgregarRespuesta("Error en linea  Partición no encontrada")
		return
	}

	// Abrir archivo binario
	file, err := utilidades.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR estructuras.MBR
	// Leer el MBR desde el archivo binario
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Imprimir el MBR
	estructuras.PrintMBR(TempMBR)
	fmt.Println("-------------")

	var index int = -1
	// Iterar sobre las particiones del MBR para encontrar la correcta
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				fmt.Println("Partition found")
				if TempMBR.Partitions[i].Status[0] == '1' {
					fmt.Println("Partition is mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		estructuras.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock estructuras.Superblock
	// Leer el Superblock desde el archivo binario
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Escribir el contenido del DOT
	textoDot := "digraph G {"
	textoDot += "  node [shape=none];"
	textoDot += "  Superbloque [label=<"
	textoDot += "  <table border='0' cellborder='1' cellspacing='0' cellpadding='10'>"
	textoDot += "    <tr><td colspan='2' bgcolor=\"#84b6f4\"><b>Superbloque</b></td></tr>"
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_filesystem_type</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_filesystem_type)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_inodes_count</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_inodes_count)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_blocks_count</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_blocks_count)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_free_blocks_count</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_free_blocks_count)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_free_inodes_count</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_free_inodes_count)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_mtime</b></td><td bgcolor=\"#bee6f4\">%s</td></tr>\n", string(tempSuperblock.S_mtime[:]))
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_umtime</b></td><td bgcolor=\"#bee6f4\">%s</td></tr>\n", string(tempSuperblock.S_umtime[:]))
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_mnt_count</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_mnt_count)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_magic</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_magic)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_inode_size</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_inode_size)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_block_size</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_block_size)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_fist_ino</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_fist_ino)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_first_blo</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_first_blo)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_bm_inode_start</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_bm_inode_start)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_bm_block_start</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_bm_block_start)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_inode_start</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_inode_start)
	textoDot += fmt.Sprintf("    <tr><td bgcolor=\"#bee6f4\"><b>S_block_start</b></td><td bgcolor=\"#bee6f4\">%d</td></tr>\n", tempSuperblock.S_block_start)
	textoDot += "  </table>"
	textoDot += ">];"
	textoDot += "}"

	rutaDot := "/home/jd/temps/superbloque.dot"
	err = os.WriteFile(rutaDot, []byte(textoDot), 0644)
	if err != nil {
		utilidades.AgregarRespuesta("Error al escribir el archivo DOT")
		fmt.Println("Error al escribir el archivo DOT:", err)
		return
	}

	// Verificar si el directorio de salida existe, si no, crearlo
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error al crear el directorio:", err)
			return
		}
	}

	// Generar la imagen en JPG usando Graphviz
	cmd := exec.Command("dot", "-Tjpg", rutaDot, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el archivo JPG:", err)
		fmt.Println("Detalles del error:", stderr.String())
		return
	}

	fmt.Println("Reporte del Superbloque generado correctamente en", path)
	utilidades.AgregarRespuesta("Reporte del Superbloque generado correctamente en " + path)
}

func repInode(id string, path string) {
	fmt.Println("====== INICIO REPORTE INODE ======")
	fmt.Println("ID:", id)
	fmt.Println("Path:", path)

	// Verificar si el usuario ya está logueado buscando en las particiones montadas
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
		utilidades.AgregarRespuesta("Error: Partición no encontrada")
		return
	}

	// Abrir archivo binario
	file, err := utilidades.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR estructuras.MBR
	// Leer el MBR desde el archivo binario
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	var index int = -1
	// Iterar sobre las particiones del MBR para encontrar la correcta
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				if TempMBR.Partitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Println("Partición no está montada")
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Println("Partición no encontrada")
		return
	}

	var tempSuperblock estructuras.Superblock
	// Leer el Superblock desde el archivo binario
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Escribir el contenido del DOT
	textoDot := "digraph G {"
	textoDot += "  node [shape=none];"
	textoDot += "  bgcolor=\"white\";"
	textoDot += "  edge [color=black];"

	// Leer los inodos y generar nodos para cada inodo utilizado
	for i := 0; i < int(tempSuperblock.S_inodes_count); i++ {
		var inode estructuras.Inode
		// Leer el inodo desde el archivo binario
		if err := utilidades.ReadObject(file, &inode, int64(tempSuperblock.S_inode_start)+int64(i)*int64(tempSuperblock.S_inode_size)); err != nil {
			fmt.Println("Error al leer el inodo:", err)
			continue
		}

		// Comprobar si el inodo está en uso
		if inode.I_size > 0 { // Suponiendo que un inodo con tamaño mayor a 0 está en uso
			textoDot += fmt.Sprintf(`  inode%d [label=<`, i)
			textoDot += `    <table border='0' cellborder='1' cellspacing='0' cellpadding='10'>`
			textoDot += fmt.Sprintf(`      <tr><td colspan='2' bgcolor='#77dd77'><b>Inode %d</b></td></tr>`, i)
			textoDot += fmt.Sprintf(`      <tr><td bgcolor='#d8f8e1'><b>UID</b></td><td bgcolor='#d8f8e1'>%d</td></tr>`, inode.I_uid)
			textoDot += fmt.Sprintf(`      <tr><td bgcolor='#d8f8e1'><b>GID</b></td><td bgcolor='#d8f8e1'>%d</td></tr>`, inode.I_gid)
			textoDot += fmt.Sprintf(`      <tr><td bgcolor='#d8f8e1'><b>Size</b></td><td bgcolor='#d8f8e1'>%d</td></tr>`, inode.I_size)
			textoDot += fmt.Sprintf(`      <tr><td bgcolor='#d8f8e1'><b>ATime</b></td><td bgcolor='#d8f8e1'>%s</td></tr>`, html.EscapeString(string(inode.I_atime[:])))
			textoDot += fmt.Sprintf(`      <tr><td bgcolor='#d8f8e1'><b>CTime</b></td><td bgcolor='#d8f8e1'>%s</td></tr>`, html.EscapeString(string(inode.I_ctime[:])))
			textoDot += fmt.Sprintf(`      <tr><td bgcolor='#d8f8e1'><b>MTime</b></td><td bgcolor='#d8f8e1'>%s</td></tr>`, html.EscapeString(string(inode.I_mtime[:])))
			textoDot += fmt.Sprintf(`      <tr><td bgcolor='#d8f8e1'><b>Blocks</b></td><td bgcolor='#d8f8e1'>%v</td></tr>`, inode.I_block)
			//fmt.Fprintf(dotFile, `      <tr><td bgcolor='#d8f8e1'><b>Type</b></td><td>%c</td></tr>`, inode.I_type[0]) // Acceso al primer byte
			textoDot += fmt.Sprintf(`      <tr><td bgcolor='#d8f8e1'><b>Perms</b></td><td bgcolor='#d8f8e1'>%v</td></tr>`, inode.I_perm)
			textoDot += `    </table>`
			textoDot += "  >];"
		}
	}

	// Cerrar el archivo DOT
	textoDot += "}"

	// Verificar si el directorio de salida existe, si no, crearlo
	rutaDot := "/home/jd/temps/inode.dot"
	err = os.WriteFile(rutaDot, []byte(textoDot), 0644)
	if err != nil {
		utilidades.AgregarRespuesta("Error al escribir el archivo DOT")
		fmt.Println("Error al escribir el archivo DOT:", err)
		return
	}

	// Verificar si el directorio de salida existe, si no, crearlo
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error al crear el directorio:", err)
			return
		}
	}

	// Generar la imagen en JPG usando Graphviz
	cmd := exec.Command("dot", "-Tjpg", rutaDot, "-o", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el archivo JPG:", err)
		fmt.Println("Detalles del error:", stderr.String())
		return
	}

	fmt.Println("Reporte de Inodes generado correctamente en", path)
	utilidades.AgregarRespuesta("Reporte de Inodes generado correctamente en " + path)
}

func repbmInodes(id string, path string) {
	fmt.Println("====== INICIO REPORTE BITMAP INODE ======")
	fmt.Println("ID:", id)
	fmt.Println("Path:", path)

	// Verificar si el directorio de salida existe, si no, crearlo
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error al crear el directorio:", err)
			return
		}
	}

	// Verificar si el usuario ya está logueado buscando en las particiones montadas
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
		utilidades.AgregarRespuesta("Error: Partición no encontrada")
		return
	}

	// Abrir archivo binario
	file, err := utilidades.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR estructuras.MBR
	// Leer el MBR desde el archivo binario
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	var index int = -1
	// Iterar sobre las particiones del MBR para encontrar la correcta
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				if TempMBR.Partitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Println("Partición no está montada")
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Println("Partición no encontrada")
		return
	}

	var tempSuperblock estructuras.Superblock
	// Leer el Superblock desde el archivo binario
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Leer el bitmap de inodos desde el archivo binario
	bitmapInode := make([]byte, tempSuperblock.S_inodes_count)
	if _, err := file.ReadAt(bitmapInode, int64(tempSuperblock.S_bm_inode_start)); err != nil {
		fmt.Println("Error: No se pudo leer el bitmap de inodos:", err)
		return
	}

	// Crear el archivo de salida para el reporte
	outputFile, err := os.Create(path)
	if err != nil {
		fmt.Println("Error al crear el archivo de reporte:", err)
		return
	}
	defer outputFile.Close()

	for i, bit := range bitmapInode {
		if i > 0 && i%20 == 0 {
			fmt.Fprintln(outputFile) // Nueva línea cada 20 bits
		}
		fmt.Fprintf(outputFile, "%d ", bit)
	}

	fmt.Println("Reporte de Bitmap de Inodos generado correctamente en", path)
	utilidades.AgregarRespuesta("Reporte de Bitmap de Inodos generado correctamente en " + path)
}

func rep_bmBlock(id string, path string) {
	fmt.Println("====== INICIO REPORTE BITMAP BLOQUES ======")
	fmt.Println("ID:", id)
	fmt.Println("Path:", path)

	// Verificar si el directorio de salida existe, si no, crearlo
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error al crear el directorio:", err)
			return
		}
	}

	// Verificar si el usuario ya está logueado buscando en las particiones montadas
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
		utilidades.AgregarRespuesta("Error: Partición no encontrada")
		return
	}

	// Abrir archivo binario
	file, err := utilidades.OpenFile(mountedPartition.Path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR estructuras.MBR
	// Leer el MBR desde el archivo binario
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	var index int = -1
	// Iterar sobre las particiones del MBR para encontrar la correcta
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				if TempMBR.Partitions[i].Status[0] == '1' {
					index = i
				} else {
					fmt.Println("Partición no está montada")
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Println("Partición no encontrada")
		return
	}

	var tempSuperblock estructuras.Superblock
	// Leer el Superblock desde el archivo binario
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Leer el bitmap de bloques desde el archivo binario
	bitmapBlock := make([]byte, tempSuperblock.S_blocks_count)
	if _, err := file.ReadAt(bitmapBlock, int64(tempSuperblock.S_bm_block_start)); err != nil {
		fmt.Println("Error: No se pudo leer el bitmap de bloques:", err)
		return
	}

	// Crear el archivo de salida para el reporte
	outputFile, err := os.Create(path)
	if err != nil {
		fmt.Println("Error al crear el archivo de reporte:", err)
		return
	}
	defer outputFile.Close()

	// Mostrar 20 bits por línea
	for i, bit := range bitmapBlock {
		if i > 0 && i%20 == 0 {
			fmt.Fprintln(outputFile) // Nueva línea cada 20 bits
		}
		fmt.Fprintf(outputFile, "%d ", bit)
	}

	fmt.Println("Reporte de Bitmap de Bloques generado correctamente en", path)
	utilidades.AgregarRespuesta("Reporte de Bitmap de Bloques generado correctamente en " + path)
}
