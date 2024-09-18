package comandos

import (
	"backend/estructuras"
	"backend/manejadorDisco"
	"backend/usuarios"
	"backend/utilidades"
	"encoding/binary"
	"fmt"
	"strings"
)

/*func Cat(files []string, linea string) {
	fmt.Println("======INICIO CAT======")
	for _, filePath := range files {
		// Intenta abrir el archivo
		file, err := utilidades.OpenFile(filePath)
		if err != nil {
			fmt.Printf("Error en línea %s: No se pudo abrir el archivo %s: %v\n", linea, filePath, err)
			utilidades.AgregarRespuesta("Error en linea " + linea + " : No se encontro la ruta:" + filePath)
			continue
		}
		defer file.Close()

		// Leer el contenido del archivo
		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error en línea %s: No se pudo leer el archivo %s: %v\n", linea, filePath, err)
			utilidades.AgregarRespuesta("Error en linea " + linea + " : No se pudo leer el archivo:" + filePath)
			continue
		}

		fmt.Println("Leyendo archivo: " + filePath)

		utilidades.AgregarRespuesta(string(content))
	}
	fmt.Println("======FIN CAT======")
}*/

func Cat(files []string) {
	// Check if a user is logged in
	if !usuarios.IsUserLoggedIn() {
		fmt.Println("Error: No hay un usuario logueado")
		utilidades.AgregarRespuesta("Error: No hay un usuario logueado")
		return
	}

	// Check if the user has permission
	if usuarios.Grupo == "root" {
		fmt.Println("Error: El usuario no tiene permiso de lectura (permiso 777 requerido)")
		utilidades.AgregarRespuesta("Error: El usuario no tiene permiso de lectura (permiso 777 requerido)")
		return
	}

	// Get the mounted partition information
	ParticionesMount := manejadorDisco.GetMountedPartitions()
	var filepath string
	var id string

	// Find the logged-in partition
	for _, partitions := range ParticionesMount {
		for _, partition := range partitions {
			if partition.Logeado {
				filepath = partition.Path
				id = partition.ID
				break
			}
		}
		if filepath != "" {
			break
		}
	}

	// Open the file
	file, err := utilidades.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Read the MBR
	var TempMBR estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Find the correct partition
	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 && strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
			if TempMBR.Partitions[i].Status[0] == '1' {
				index = i
				break
			}
		}
	}

	if index == -1 {
		fmt.Println("Error: No se encontró la partición")
		return
	}

	// Read the Superblock
	var tempSuperblock estructuras.Superblock
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Process each file in the input
	for _, filePath := range files {
		fmt.Printf("Imprimiendo el contenido de %s\n", filePath)

		indexInode := usuarios.BuscarStart(filePath, file, tempSuperblock)
		if indexInode == -1 {
			fmt.Printf("Error: No se pudo encontrar el archivo %s\n", filePath)
			continue
		}

		var crrInode estructuras.Inode
		if err := utilidades.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(estructuras.Inode{})))); err != nil {
			fmt.Printf("Error: No se pudo leer el Inode para %s\n", filePath)
			continue
		}

		// Read and print the content of each block in the file
		for _, block := range crrInode.I_block {
			if block != -1 {
				var fileblock estructuras.Fileblock
				if err := utilidades.ReadObject(file, &fileblock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
					fmt.Printf("Error: No se pudo leer el Fileblock para %s\n", filePath)
					continue
				}
				estructuras.AgregarFileBlockConsola(fileblock)
			}
		}

		fmt.Println("------FIN CAT------")
	}
}
