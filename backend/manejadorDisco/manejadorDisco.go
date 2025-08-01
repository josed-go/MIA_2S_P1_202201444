package manejadorDisco

import (
	"backend/estructuras"
	"backend/utilidades"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ParticionMontada struct {
	Path    string
	Name    string
	ID      string
	Status  byte // 0 -> No Montada 1 -> Montada
	Logeado bool
}

var particionesMontadas = make(map[string][]ParticionMontada)

//var letras = make(map[string]byte) // Mapa para almacenar la letra asignada a cada disco

func GetMountedPartitions() map[string][]ParticionMontada {
	return particionesMontadas
}

func MarkPartitionAsLogeado(id string) {
	for diskID, partitions := range particionesMontadas {
		for i, partition := range partitions {
			if partition.ID == id {
				particionesMontadas[diskID][i].Logeado = true
				fmt.Printf("Partición con ID %s marcada como logueada.\n", id)
				return
			}
		}
	}
	fmt.Printf("No se encontró la partición con ID %s para marcarla como logueada.\n", id)
}

func MarkPartitionAsDeslogeado(id string) {
	for _, partitions := range particionesMontadas {
		for i, partition := range partitions {
			if partition.ID == id {
				partitions[i].Logeado = false
				fmt.Println("Partición", id, "marcada como deslogueada.")
				utilidades.AgregarRespuesta("Se ha cerrado la sesion en la Partición " + id)
				return
			}
		}
	}
}

func PrintMountedPartitions() {
	fmt.Println("Particiones montadas:")

	if len(particionesMontadas) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}

	for diskID, partitions := range particionesMontadas {
		fmt.Printf("Disco ID: %s\n", diskID)
		for _, partition := range partitions {
			fmt.Printf(" - Partición Name: %s, ID: %s, Path: %s, Status: %c\n",
				partition.Name, partition.ID, partition.Path, partition.Status)
		}
	}
	fmt.Println("")
}

/*func getLetra(path string) byte {
	if letra, exists := letras[path]; exists {
		return letra
	}
	// Si el disco no tiene una letra asignada, se le asigna la siguiente disponible
	newLetter := 'A' + byte(len(letras))
	letras[path] = newLetter
	return newLetter
}*/

func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("======INICIO MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	fmt.Println("Path:", path)

	// Validar fit bf/ff/wf
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit debe ser bf, wf or ff")
		return
	}

	// Validar size > 0
	if size <= 0 {
		fmt.Println("Error: Size debe ser mayo a  0")
		return
	}

	// Validar unidar k - m
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades validas son k o m")
		return
	}

	// Create file
	err := utilidades.CreateFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	/*
		Si el usuario especifica unit = "k" (Kilobytes), el tamaño se multiplica por 1024 para convertirlo a bytes.
		Si el usuario especifica unit = "m" (Megabytes), el tamaño se multiplica por 1024 * 1024 para convertirlo a MEGA bytes.
	*/
	// Asignar tamanio
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// Open bin file
	file, err := utilidades.OpenFile(path)
	if err != nil {
		return
	}

	// Escribir los 0 en el archivo

	// create array of byte(0)
	for i := 0; i < size; i++ {
		err := utilidades.WriteObject(file, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// Crear MRB
	var newMRB estructuras.MBR
	newMRB.MbrSize = int32(size)
	newMRB.Signature = rand.Int31() // Numero random rand.Int31() genera solo números no negativos
	copy(newMRB.Fit[:], fit)

	// Obtener la fecha del sistema en formato YYYY-MM-DD
	currentTime := time.Now()
	formattedDate := currentTime.Format("02-01-2006 15:04:05")
	copy(newMRB.CreationDate[:], formattedDate)

	// Escribir el archivo
	if err := utilidades.WriteObject(file, newMRB, 0); err != nil {
		return
	}

	var TempMBR estructuras.MBR
	// Leer el archivo
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	//estructuras.PrintMBR(TempMBR)

	// Cerrar el archivo
	defer file.Close()

	fmt.Println("======FIN MKDISK======")

	utilidades.AgregarRespuesta("Disco creado correctamente, ruta: " + path)
}

func Rmdisk(path string, linea string) {
	fmt.Println("======INICIO RMDISK======")
	fmt.Println("Path:", path)

	// Create file
	err := utilidades.DeleteFile(path, linea)
	if err != nil {
		// Maneja el error si ocurre
		fmt.Println("Error:", err)
		return
	} else {
		// Confirmación de eliminación exitosa
		fmt.Println("Archivo eliminado exitosamente.")
		utilidades.AgregarRespuesta("Archivo con ruta: " + path + " eliminado correctamente")
	}

	fmt.Println("======FIN MKDISK======")
}

func Fdisk(size int, fit string, unit string, path string, typ string, name string, linea string) {
	fmt.Println("======INICIO FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Type:", typ)
	fmt.Println("Unit:", unit)
	fmt.Println("Path:", path)
	fmt.Println("Name:", name)

	// Validar size > 0
	if size <= 0 {
		fmt.Println("Error: Size debe ser mayo a  0")
		return
	}

	// Validar unidar k - m
	if unit != "k" && unit != "m" && unit != "b" {
		fmt.Println("Error: Las unidades validas son k o m o b")
		return
	}

	if typ != "p" && typ != "e" && typ != "l" {
		fmt.Println("Error: El parametro type debe ser p - e - l")
		return
	}

	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: El parametro fit debe ser bf - ff - wf")
		return
	}

	if name == "" {
		fmt.Println("Error: El parametro name es obligatorio")
		return
	}

	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	file, err := utilidades.OpenFile(path)
	if err != nil {
		utilidades.AgregarRespuesta("Error en linea " + linea + " : No se encontro la ruta:" + path)
		return
	}

	var mbrTemp estructuras.MBR

	if err := utilidades.ReadObject(file, &mbrTemp, 0); err != nil {
		utilidades.AgregarRespuesta("Ocurrio un error al acceder al disco en ruta:" + path)
		return
	}

	estructuras.PrintMBR(mbrTemp)

	fmt.Println("*****************")

	for i := 0; i < 4; i++ {
		partitionName := string(mbrTemp.Partitions[i].Name[:])
		partitionName = strings.TrimRight(partitionName, "\x00")
		fmt.Println("Nombre a ver:" + partitionName + " ---- Nombre validando:" + name)
		fmt.Println("Resultado: ", partitionName == name)
		if mbrTemp.Partitions[i].Size != 0 && partitionName == name {
			fmt.Println("Error en linea " + linea + " : Ya existe una particion llamada:" + name + " en la ruta:" + path)
			utilidades.AgregarRespuesta("Error en linea " + linea + " : Ya existe una particion llamada:" + name + " en la ruta:" + path)
			return
		}
	}

	var contP, contE, contT int

	var espacioUsado int32 = 0

	// Restar el espacio de MBR

	//espacioRestante = espacioRestante - 168

	for i := 0; i < 4; i++ {
		if mbrTemp.Partitions[i].Size != 0 {
			contT++
			espacioUsado += mbrTemp.Partitions[i].Size

			if mbrTemp.Partitions[i].Type[0] == 'p' {
				contP++
			} else if mbrTemp.Partitions[i].Type[0] == 'e' {
				contE++
			}
		}
	}

	if contT >= 4 && typ != "l" {
		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : No se pueden crear mas de 4 particiones")
		return
	}

	if typ == "e" && contE > 0 {
		fmt.Println("Error: Solo se permite una partición extendida por disco.")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Solo se permite una particion extendida por disco.")
		return
	}

	if typ == "l" && contE == 0 {
		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : No se puede crear una partición lógica sin una partición extendida.")
		return
	}

	if espacioUsado+int32(size) > mbrTemp.MbrSize {
		fmt.Println("Error: o hay suficiente espacio en el disco para crear esta partición.")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : No hay suficiente espacio en el disco para crear esta partición.")
		return
	}

	var posicion int32 = 0

	if contT == 0 {
		posicion = int32(binary.Size(mbrTemp))
	}

	if contT > 0 {
		posicion = mbrTemp.Partitions[contT-1].Start + mbrTemp.Partitions[contT-1].Size
	}

	for i := 0; i < 4; i++ {
		if mbrTemp.Partitions[i].Size == 0 {
			if typ == "p" || typ == "e" {
				mbrTemp.Partitions[i].Size = int32(size)
				mbrTemp.Partitions[i].Start = posicion
				copy(mbrTemp.Partitions[i].Name[:], name)
				copy(mbrTemp.Partitions[i].Fit[:], fit)
				copy(mbrTemp.Partitions[i].Status[:], "0")
				copy(mbrTemp.Partitions[i].Type[:], typ)
				mbrTemp.Partitions[i].Correlative = 0

				// CODIGO PARA LA EXTENDIDA Y LOGICAS
				if typ == "e" {
					// Inicializar el primer EBR en la partición extendida
					ebrStart := posicion // El primer EBR se coloca al inicio de la partición extendida
					ebr := estructuras.EBR{
						PartFit:   [1]byte{fit[0]},
						PartStart: ebrStart,
						PartSize:  0,
						PartNext:  -1,
					}
					copy(ebr.PartName[:], "")
					utilidades.WriteObject(file, ebr, int64(ebrStart))
				}

				break
			}
		}
	}

	if typ == "l" {
		var particionEx *estructuras.Partition
		for i := 0; i < 4; i++ {
			if mbrTemp.Partitions[i].Type[0] == 'e' {
				particionEx = &mbrTemp.Partitions[i]
				break
			}
		}

		if particionEx == nil {
			fmt.Println("Error: No se encontró una partición extendida para crear la partición lógica")
			return
		}

		// Encontrar el último EBR en la cadena
		ebrPos := particionEx.Start
		var lastEBR estructuras.EBR
		for {
			utilidades.ReadObject(file, &lastEBR, int64(ebrPos))

			if strings.Contains(string(lastEBR.PartName[:]), name) {
				fmt.Println("Error en linea " + linea + " : Ya existe una particion logica llamada:" + name + " en la ruta:" + path)
				utilidades.AgregarRespuesta("Error en linea " + linea + " : Ya existe una particion logica llamada:" + name + " en la ruta:" + path)
				return
			}

			if lastEBR.PartNext == -1 {
				break
			}
			ebrPos = lastEBR.PartNext
		}

		// Calcular la posición de inicio de la nueva partición lógica
		var newEBRPos int32
		if lastEBR.PartSize == 0 {
			// Es la primera partición lógica
			newEBRPos = ebrPos
		} else {
			// No es la primera partición lógica
			newEBRPos = lastEBR.PartStart + lastEBR.PartSize
		}

		// Verificar si hay espacio suficiente en la partición extendida
		if newEBRPos+int32(size)+int32(binary.Size(estructuras.EBR{})) > particionEx.Start+particionEx.Size {
			fmt.Println("Error: No hay suficiente espacio en la partición extendida para esta partición lógica")
			return
		}

		// Actualizar el EBR anterior
		if lastEBR.PartSize != 0 {
			lastEBR.PartNext = newEBRPos
			utilidades.WriteObject(file, lastEBR, int64(ebrPos))
		}

		fmt.Println("Imprimir el tamano del ebr")
		fmt.Println(int32(binary.Size(estructuras.EBR{})))
		// Crear y escribir el nuevo EBR
		newEBR := estructuras.EBR{
			PartFit:   [1]byte{fit[0]}, //[1]byte(fit[0]),
			PartStart: newEBRPos,       //+ int32(binary.Size(Structs.EBR{})),
			PartSize:  int32(size),
			PartNext:  -1,
		}
		copy(newEBR.PartName[:], name)
		utilidades.WriteObject(file, newEBR, int64(newEBRPos))

		fmt.Println("Partición lógica creada exitosamente")
		estructuras.PrintEBR(newEBR)
	}

	fmt.Println("------------------")
	fmt.Println("Tamaño del disco:", mbrTemp.MbrSize, "bytes")
	fmt.Println("Tamaño utilizado:", espacioUsado, "bytes")
	fmt.Println("Tamaño restante:", mbrTemp.MbrSize-espacioUsado, "bytes")
	fmt.Println("------------------")

	if err := utilidades.WriteObject(file, &mbrTemp, 0); err != nil {
		fmt.Println("Error: Could not write MBR to file")
		return
	}

	var TempMBR2 estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}

	estructuras.PrintMBR(TempMBR2)

	defer file.Close()

	fmt.Println("Partición con nombre : "+name+" creada con éxito en la ruta:", path)
	utilidades.AgregarRespuesta("Partición con nombre : " + name + " creada con éxito en la ruta: " + path)

	fmt.Println("======FIN FDISK======")
}

// Función para montar particiones
func Mount(path string, name string, linea string) {
	fmt.Println("======INICIO MOUNT======")
	file, err := utilidades.OpenFile(path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo en la ruta:", path)
		return
	}
	defer file.Close()

	var TempMBR estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR desde el archivo")
		return
	}

	fmt.Printf("Buscando partición con nombre: '%s'\n", name)

	partitionFound := false
	var partition estructuras.Partition
	var partitionIndex int

	nameBytes := [16]byte{}
	copy(nameBytes[:], []byte(name))

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Type[0] == 'p' && bytes.Equal(TempMBR.Partitions[i].Name[:], nameBytes[:]) {
			partition = TempMBR.Partitions[i]
			partitionIndex = i
			partitionFound = true
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: Partición no encontrada o no es una partición primaria")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Partición no encontrada o no es una partición primaria")
		return
	}

	if partition.Status[0] == '1' {
		fmt.Println("Error: La partición ya está montada")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : La partición ya está montada")
		return
	}

	diskID := generateDiskID(path)

	mountedPartitionsInDisk := particionesMontadas[diskID]
	var letter byte

	if len(mountedPartitionsInDisk) == 0 {

		if len(particionesMontadas) == 0 {
			letter = 'a'
		} else {
			lastDiskID := getLastDiskID()
			lastLetter := particionesMontadas[lastDiskID][0].ID[len(particionesMontadas[lastDiskID][0].ID)-1]
			letter = lastLetter + 1
		}
	} else {

		letter = mountedPartitionsInDisk[0].ID[len(mountedPartitionsInDisk[0].ID)-1]
	}

	carnet := "202201444" // Cambiar su carnet aquí
	lastTwoDigits := carnet[len(carnet)-2:]
	number := len(mountedPartitionsInDisk) + 1
	partitionID := fmt.Sprintf("%s%d%c", lastTwoDigits, number, letter)

	partition.Status[0] = '1'
	copy(partition.Id[:], partitionID)
	TempMBR.Partitions[partitionIndex] = partition
	particionesMontadas[diskID] = append(particionesMontadas[diskID], ParticionMontada{
		Path:   path,
		Name:   name,
		ID:     partitionID,
		Status: '1',
	})

	// Escribir el MBR actualizado al archivo
	if err := utilidades.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
		return
	}

	fmt.Printf("Partición montada con ID: %s\n", partitionID)
	utilidades.AgregarRespuesta("Partición " + name + " montada con ID: " + partitionID)

	fmt.Println("")
	// Imprimir el MBR actualizado
	fmt.Println("MBR actualizado:")
	estructuras.PrintMBR(TempMBR)
	fmt.Println("")

	PrintMountedPartitions()
	fmt.Println("======FIN MOUNT======")
}

// Función para obtener el ID del último disco montado
func getLastDiskID() string {
	var lastDiskID string
	for diskID := range particionesMontadas {
		lastDiskID = diskID
	}
	return lastDiskID
}

func generateDiskID(path string) string {
	return strings.ToLower(path)
}

func ShowPartitions(path string) {
	fmt.Println("======INICIO SHOW PARTITIONS======")
	utilidades.AgregarRespuesta("======INICIO SHOW PARTITIONS======")

	// Abre el archivo
	file, err := utilidades.OpenFile(path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo en la ruta:", path)
		return
	}
	defer file.Close()

	// Lee el MBR
	var TempMBR estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR desde el archivo")
		return
	}

	fmt.Println("Particiones montadas en disco:", path)
	utilidades.AgregarRespuesta("Particiones en disco: " + path)

	// Recorre las particiones
	for i := 0; i < 4; i++ {
		partition := TempMBR.Partitions[i]
		if partition.Size != 0 {
			if partition.Status[0] == '1' {
				// Convertir el nombre de partición y eliminar bytes nulos
				nombreParticion := strings.TrimRight(string(partition.Name[:]), "\x00")

				// Construye el string con la información de la partición
				valor := fmt.Sprintf(" - Partición Name: %s, ID: %s, Status: %c, Size: %d",
					nombreParticion, partition.Id, partition.Status[0], partition.Size)

				// Añade el resultado final
				utilidades.AgregarRespuesta(valor)
			}
		}
	}

	utilidades.AgregarRespuesta("======FIN SHOW PARTITIONS======")
	fmt.Println("======FIN SHOW PARTITIONS======")
}
