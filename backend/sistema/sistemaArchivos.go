package sistema

import (
	"backend/estructuras"
	"backend/manejadorDisco"
	"backend/utilidades"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func Mkfs(id string, type_ string, fs_ string) {
	fmt.Println("======INICIO MKFS======")
	fmt.Println("Id:", id)
	fmt.Println("Type:", type_)
	fmt.Println("Fs:", fs_)

	currentTime := time.Now()

	// Buscar la partición montada por ID
	var mountedPartition manejadorDisco.ParticionMontada
	var partitionFound bool

	for _, partitions := range manejadorDisco.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.ID == id {
				mountedPartition = partition
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Particion no encontrada")
		utilidades.AgregarRespuesta("Particion con ID: " + id + " no encontrada")
		return
	}

	if mountedPartition.Status != '1' { // Verifica si la partición está montada
		fmt.Println("La particion aun no esta montada")
		utilidades.AgregarRespuesta("La particion con ID: " + id + " aun no esta montada")
		return
	}

	// Abrir archivo binario
	file, err := utilidades.OpenFile(mountedPartition.Path)
	if err != nil {
		return
	}

	var TempMBR estructuras.MBR
	// Leer objeto desde archivo binario
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Imprimir objeto
	estructuras.PrintMBR(TempMBR)

	fmt.Println("-------------")

	var index int = -1
	// Iterar sobre las particiones para encontrar la que tiene el nombre correspondiente
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				index = i
				break
			}
		}
	}

	if index != -1 {
		estructuras.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Particion no encontrada (2)")
		utilidades.AgregarRespuesta("Particion no encontrada (2)")
		return
	}

	numerador := int32(TempMBR.Partitions[index].Size - int32(binary.Size(estructuras.Superblock{})))
	denominador_base := int32(4 + int32(binary.Size(estructuras.Inode{})) + 3*int32(binary.Size(estructuras.Fileblock{})))
	var temp int32 = 0
	if type_ == "full" {
		temp = 0
	} else {
		fmt.Print("Error por el momento solo está disponible FULL.")
		utilidades.AgregarRespuesta("Error por el momento solo está disponible FULL.")
	}
	denominador := denominador_base + temp
	n := int32(numerador / denominador)

	fmt.Println("INODOS:", n)

	// Crear el Superblock con todos los campos calculados
	var newSuperblock estructuras.Superblock
	newSuperblock.S_filesystem_type = 2 // EXT2
	newSuperblock.S_inodes_count = n
	newSuperblock.S_blocks_count = 3 * n
	newSuperblock.S_free_blocks_count = 3*n - 2
	newSuperblock.S_free_inodes_count = n - 2
	copy(newSuperblock.S_mtime[:], currentTime.Format("2006-01-02 15:04:05"))
	copy(newSuperblock.S_umtime[:], currentTime.Format("2006-01-02 15:04:05"))
	newSuperblock.S_mnt_count = 1
	newSuperblock.S_magic = 0xEF53
	newSuperblock.S_inode_size = int32(binary.Size(estructuras.Inode{}))
	newSuperblock.S_block_size = int32(binary.Size(estructuras.Fileblock{}))

	// Calcula las posiciones de inicio
	newSuperblock.S_bm_inode_start = TempMBR.Partitions[index].Start + int32(binary.Size(estructuras.Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*newSuperblock.S_inode_size

	if type_ == "full" {
		create_ext2(n, TempMBR.Partitions[index], newSuperblock, currentTime.Format("2006-01-02 15:04:05"), file)
	} else {
		fmt.Println("EXT3 no está soportado.")
		utilidades.AgregarRespuesta("EXT3 no está soportado.")
	}

	// Cerrar archivo binario
	defer file.Close()

	fmt.Println("======FIN MKFS======")
}

func create_ext2(n int32, partition estructuras.Partition, newSuperblock estructuras.Superblock, date string, file *os.File) {
	fmt.Println("======Start CREATE EXT2======")
	fmt.Println("INODOS:", n)

	// Imprimir Superblock inicial
	estructuras.PrintSuperblock(newSuperblock)
	fmt.Println("Date:", date)

	// Escribe los bitmaps de inodos y bloques en el archivo
	for i := int32(0); i < n; i++ {
		if err := utilidades.WriteObject(file, byte(0), int64(newSuperblock.S_bm_inode_start+i)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	for i := int32(0); i < 3*n; i++ {
		if err := utilidades.WriteObject(file, byte(0), int64(newSuperblock.S_bm_block_start+i)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	// Inicializa inodos y bloques con valores predeterminados
	if err := initInodesAndBlocks(n, newSuperblock, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Crea la carpeta raíz y el archivo users.txt
	if err := createRootAndUsersFile(newSuperblock, date, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Escribe el superbloque actualizado al archivo
	if err := utilidades.WriteObject(file, newSuperblock, int64(partition.Start)); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Marca los primeros inodos y bloques como usados
	if err := markUsedInodesAndBlocks(newSuperblock, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Leer e imprimir los inodos después de formatear
	fmt.Println("====== Imprimiendo Inodos ======")
	for i := int32(0); i < n; i++ {
		var inode estructuras.Inode
		offset := int64(newSuperblock.S_inode_start + i*int32(binary.Size(estructuras.Inode{})))
		if err := utilidades.ReadObject(file, &inode, offset); err != nil {
			fmt.Println("Error al leer inodo: ", err)
			return
		}
		estructuras.PrintInode(inode)
	}

	// Leer e imprimir los Folderblocks y Fileblocks después de formatear
	fmt.Println("====== Imprimiendo Folderblocks y Fileblocks ======")

	// Imprimir Folderblocks
	for i := int32(0); i < 1; i++ {
		var folderblock estructuras.Folderblock
		offset := int64(newSuperblock.S_block_start + i*int32(binary.Size(estructuras.Folderblock{})))
		if err := utilidades.ReadObject(file, &folderblock, offset); err != nil {
			fmt.Println("Error al leer Folderblock: ", err)
			return
		}
		estructuras.PrintFolderblock(folderblock)
	}

	// Imprimir Fileblocks
	for i := int32(0); i < 1; i++ {
		var fileblock estructuras.Fileblock
		offset := int64(newSuperblock.S_block_start + int32(binary.Size(estructuras.Folderblock{})) + i*int32(binary.Size(estructuras.Fileblock{})))
		if err := utilidades.ReadObject(file, &fileblock, offset); err != nil {
			fmt.Println("Error al leer Fileblock: ", err)
			return
		}
		estructuras.PrintFileblock(fileblock)
	}

	// Imprimir el Superblock final
	estructuras.PrintSuperblock(newSuperblock)

	fmt.Println("======End CREATE EXT2======")
}

// Función auxiliar para inicializar inodos y bloques
func initInodesAndBlocks(n int32, newSuperblock estructuras.Superblock, file *os.File) error {
	var newInode estructuras.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		if err := utilidades.WriteObject(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(estructuras.Inode{})))); err != nil {
			return err
		}
	}

	var newFileblock estructuras.Fileblock
	for i := int32(0); i < 3*n; i++ {
		if err := utilidades.WriteObject(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
			return err
		}
	}

	return nil
}

// Función auxiliar para crear la carpeta raíz y el archivo users.txt
func createRootAndUsersFile(newSuperblock estructuras.Superblock, date string, file *os.File) error {
	var Inode0, Inode1 estructuras.Inode
	initInode(&Inode0, date)
	initInode(&Inode1, date)

	Inode0.I_block[0] = 0
	Inode1.I_block[0] = 1

	// Asignar el tamaño real del contenido
	data := "1,G,root\n1,U,root,root,123\n"
	actualSize := int32(len(data))
	Inode1.I_size = actualSize // Esto ahora refleja el tamaño real del contenido

	var Fileblock1 estructuras.Fileblock
	copy(Fileblock1.B_content[:], data) // Copia segura de datos a Fileblock

	var Folderblock0 estructuras.Folderblock
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[2].B_inodo = 1
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")

	// Escribir los inodos y bloques en las posiciones correctas
	if err := utilidades.WriteObject(file, Inode0, int64(newSuperblock.S_inode_start)); err != nil {
		return err
	}
	if err := utilidades.WriteObject(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(estructuras.Inode{})))); err != nil {
		return err
	}
	if err := utilidades.WriteObject(file, Folderblock0, int64(newSuperblock.S_block_start)); err != nil {
		return err
	}
	if err := utilidades.WriteObject(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(estructuras.Folderblock{})))); err != nil {
		return err
	}

	return nil
}

// Función auxiliar para inicializar un inodo
func initInode(inode *estructuras.Inode, date string) {
	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_size = 0
	copy(inode.I_atime[:], date)
	copy(inode.I_ctime[:], date)
	copy(inode.I_mtime[:], date)
	copy(inode.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		inode.I_block[i] = -1
	}
}

// Función auxiliar para marcar los inodos y bloques usados
func markUsedInodesAndBlocks(newSuperblock estructuras.Superblock, file *os.File) error {
	if err := utilidades.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start)); err != nil {
		return err
	}
	if err := utilidades.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start+1)); err != nil {
		return err
	}
	if err := utilidades.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start)); err != nil {
		return err
	}
	if err := utilidades.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start+1)); err != nil {
		return err
	}
	return nil
}
