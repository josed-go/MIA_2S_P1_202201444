package manejadorDisco

import (
	"backend/estructuras"
	"backend/utilidades"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var letras = make(map[string]byte) // Mapa para almacenar la letra asignada a cada disco

func getLetra(path string) byte {
	if letra, exists := letras[path]; exists {
		return letra
	}
	// Si el disco no tiene una letra asignada, se le asigna la siguiente disponible
	newLetter := 'A' + byte(len(letras))
	letras[path] = newLetter
	return newLetter
}

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

	if fit != "b" && fit != "f" && fit != "w" {
		fmt.Println("Error: El parametro fit debe ser b - f - w")
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

func Mount(name string, path string, linea string) {
	fmt.Println("======INICIO MOUNT======")
	if path == "" {
		fmt.Println("Error: La ruta es obligatoria.")
		return
	}

	if name == "" {
		fmt.Println("Error: El nombre es obligatorio.")
		return
	}

	file, err := utilidades.OpenFile(path)
	if err != nil {
		utilidades.AgregarRespuesta("Error en linea " + linea + " : No se encontro la ruta:" + path)
		return
	}

	var mbrTemp estructuras.MBR
	if err := utilidades.ReadObject(file, &mbrTemp, 0); err != nil {
		return
	}

	mountedCount := 0
	diskLetter := getLetra(path) // Obtener la letra correspondiente al disco actual

	for i := 0; i < 4; i++ {
		if mbrTemp.Partitions[i].Status[0] == '1' {
			mountedCount++
		}
	}

	existe := false

	for i := 0; i < 4; i++ {
		partitionName := string(mbrTemp.Partitions[i].Name[:])
		partitionName = strings.TrimRight(partitionName, "\x00")
		if partitionName == name {
			if strings.Contains(string(mbrTemp.Partitions[i].Type[:]), "e") || strings.Contains(string(mbrTemp.Partitions[i].Type[:]), "l") {
				fmt.Println("Nombre:" + string(mbrTemp.Partitions[i].Name[:]))
				fmt.Println("Tipo:" + string(mbrTemp.Partitions[i].Type[:]))
				fmt.Println("Error: No se puede montar una partición extendida o lógica.")
				utilidades.AgregarRespuesta("Error en linea " + linea + " : No se puede montar una partición extendida o lógica.")
				return
			}

			if mbrTemp.Partitions[i].Status[0] == '1' {
				fmt.Println("Error: La particion con nombre: " + name + " ya se encuentra montada.")
				utilidades.AgregarRespuesta("Error en linea " + linea + " : La particion con nombre: " + name + " ya se encuentra montada.")
				return
			}

			mountedCount++ // Incrementar el número de partición montada
			mbrTemp.Partitions[i].Status[0] = '1'
			mbrTemp.Partitions[i].Correlative = int32(mountedCount)
			// Generar ID con 44 + número de partición + letra de disco
			mbrTemp.Partitions[i].Id = [4]byte{'4', '4', byte(mountedCount + '0'), diskLetter}
			fmt.Println("Partición montada exitosamente en la ruta: " + path)
			utilidades.AgregarRespuesta("Partición con nombre:" + name + " montada exitosamente en la ruta: " + path)
			existe = true
			break
		}
	}

	if !existe {
		fmt.Println("Error: No existe ninguna partición con ese nombre en el disco con ruta:" + path)
		utilidades.AgregarRespuesta("Error en linea " + linea + " : No existe ninguna partición con ese nombre en el disco con ruta:" + path)
		return
	}

	if err := utilidades.WriteObject(file, mbrTemp, 0); err != nil {
		return
	}

	var TempMBR2 estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}
	fmt.Println("******** MBR ********")
	estructuras.PrintMBR(TempMBR2)
	fmt.Println("******** MBR ********")
	defer file.Close()

	fmt.Println("======FIN MOUNT======")

}
