package usuarios

import (
	"backend/estructuras"
	"backend/manejadorDisco"
	"backend/utilidades"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Grupo string
var Usuario string
var Password string

func SetUsuarioLogeado(grupo, user, pass string) {
	Grupo = grupo
	Usuario = user
	Password = pass
}

func Login(user string, pass string, id string) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	// Verificar si el usuario ya está logueado buscando en las particiones montadas
	mountedPartitions := manejadorDisco.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false

	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.ID == id && partition.Logeado { // Verifica si ya está logueado
				fmt.Println("Ya existe un usuario logueado!")
				utilidades.AgregarRespuesta("Ya existe un usuario logueado!")
				return
			}
			if partition.ID == id { // Encuentra la partición correcta
				filepath = partition.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No se encontró ninguna partición montada con el ID proporcionado")
		utilidades.AgregarRespuesta("Error: No se encontró ninguna partición montada con el ID proporcionado")
		return
	}

	// Abrir archivo binario
	file, err := utilidades.OpenFile(filepath)
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

	// Buscar el archivo de usuarios /users.txt -> retorna índice del Inodo
	indexInode := InitSearch("/users.txt", file, tempSuperblock)

	var crrInode estructuras.Inode
	// Leer el Inodo desde el archivo binario
	if err := utilidades.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(estructuras.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo:", err)
		return
	}

	// Leer datos del archivo
	data := GetInodeFileData(crrInode, file, tempSuperblock)

	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	// Iterar a través de las líneas para verificar las credenciales
	for _, line := range lines {
		words := strings.Split(line, ",")

		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], pass)) {
				SetUsuarioLogeado(words[2], words[3], words[4])
				login = true
				break
			}
		}
	}

	// Imprimir información del Inodo
	fmt.Println("Inode", crrInode.I_block)

	// Si las credenciales son correctas y marcamos como logueado
	if login {
		fmt.Println("Usuario logueado con exito")
		utilidades.AgregarRespuesta("Usuario logueado con exito")
		manejadorDisco.MarkPartitionAsLogeado(id) // Marcar la partición como logueada
	}

	fmt.Println("======End LOGIN======")
}

func InitSearch(path string, file *os.File, tempSuperblock estructuras.Superblock) int32 {
	fmt.Println("======Start BUSQUEDA INICIAL ======")
	fmt.Println("path:", path)
	// path = "/ruta/nueva"

	// split the path by /
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:]

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 estructuras.Inode
	// Read object from bin file
	if err := utilidades.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1
	}

	fmt.Println("======End BUSQUEDA INICIAL======")

	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

// stack
func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func SarchInodeByPath(StepsPath []string, Inode estructuras.Inode, file *os.File, tempSuperblock estructuras.Superblock) int32 {
	fmt.Println("======Start BUSQUEDA INODO POR PATH======")
	index := int32(0)
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {
				//CASO DIRECTO

				var crrFolderBlock estructuras.Folderblock
				// Read object from bin file
				if err := utilidades.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(estructuras.Folderblock{})))); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_content {
					// fmt.Println("Folder found======")
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					if strings.Contains(string(folder.B_name[:]), SearchedName) {

						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)
						if len(StepsPath) == 0 {
							fmt.Println("Folder found======")
							return folder.B_inodo
						} else {
							fmt.Println("NextInode======")
							var NextInode estructuras.Inode
							// Read object from bin file
							if err := utilidades.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(estructuras.Inode{})))); err != nil {
								return -1
							}
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End BUSQUEDA INODO POR PATH======")
	return 0
}

func GetInodeFileData(Inode estructuras.Inode, file *os.File, tempSuperblock estructuras.Superblock) string {
	fmt.Println("======Start CONTENIDO DEL BLOQUE======")
	index := int32(0)
	// define content as a string
	var content string

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			//Dentro de los directos
			if index < 13 {
				var crrFileBlock estructuras.Fileblock
				// Read object from bin file
				if err := utilidades.ReadObject(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_content[:])

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End CONTENIDO DEL BLOQUE======")
	return content
}

// MKUSER
func AppendToFileBlock(inode *estructuras.Inode, newData string, file *os.File, superblock estructuras.Superblock) error {
	// Leer el contenido existente del archivo utilizando la función GetInodeFileData
	existingData := GetInodeFileData(*inode, file, superblock)

	// Concatenar el nuevo contenido
	fullData := existingData + newData

	// Asegurarse de que el contenido no exceda el tamaño del bloque
	if len(fullData) > len(inode.I_block)*binary.Size(estructuras.Fileblock{}) {
		// Si el contenido excede, necesitas manejar bloques adicionales
		return fmt.Errorf("el tamaño del archivo excede la capacidad del bloque actual y no se ha implementado la creación de bloques adicionales")
	}

	// Escribir el contenido actualizado en el bloque existente
	var updatedFileBlock estructuras.Fileblock
	copy(updatedFileBlock.B_content[:], fullData)
	if err := utilidades.WriteObject(file, updatedFileBlock, int64(superblock.S_block_start+inode.I_block[0]*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
		return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
	}

	// Actualizar el tamaño del inodo
	inode.I_size = int32(len(fullData))
	if err := utilidades.WriteObject(file, *inode, int64(superblock.S_inode_start+inode.I_block[0]*int32(binary.Size(estructuras.Inode{})))); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}

func Logout() {
	fmt.Println("======Start LOGOUT======")

	// Obtener las particiones montadas
	mountedPartitions := manejadorDisco.GetMountedPartitions()

	var partitionFound bool

	// Buscar una partición que esté logueada
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.Logeado { // Si la partición está logueada
				manejadorDisco.MarkPartitionAsDeslogeado(partition.ID) // Marcar la partición como deslogueada
				fmt.Println("Sesión cerrada exitosamente")
				utilidades.AgregarRespuesta("Sesión cerrada exitosamente")
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No se encontró ninguna sesión activa")
		utilidades.AgregarRespuesta("Error: No se encontró ninguna sesión activa")
	}

	fmt.Println("======End LOGOUT======")
}

// Función para verificar si un usuario está logueado
func IsUserLoggedIn() bool {
	ParticionesMount := manejadorDisco.GetMountedPartitions()

	for _, partitions := range ParticionesMount {
		for _, partition := range partitions {
			// Verifica si alguna partición tiene un usuario logueado
			if partition.Logeado {
				return true
			}
		}
	}

	return false
}

/*func AddUser(username, password, group string) {
	fmt.Println("======Start ADD USER======")

	// Obtener las particiones montadas
	mountedPartitions := manejadorDisco.GetMountedPartitions()
	var filepath string
	var partitionFound bool

	// Buscar la partición que está logueada
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.Logeado { // Si la partición está logueada
				filepath = partition.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No hay ninguna sesión activa")
		utilidades.AgregarRespuesta("Error: No hay ninguna sesión activa")
		return
	}

	// Abrir el archivo binario
	file, err := utilidades.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Leer el MBR desde el archivo binario
	var TempMBR estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Buscar la partición correcta donde está logueado el usuario
	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 && TempMBR.Partitions[i].Status[0] == '1' {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Error: No se encontró una partición válida")
		return
	}

	// Leer el Superblock desde el archivo binario
	var tempSuperblock estructuras.Superblock
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Buscar el archivo de usuarios /users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock)

	var crrInode estructuras.Inode
	// Leer el Inodo del archivo de usuarios
	if err := utilidades.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(estructuras.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo:", err)
		return
	}

	// Obtener el contenido actual del archivo de usuarios
	data := GetInodeFileData(crrInode, file, tempSuperblock)

	// Dividir el contenido en líneas
	lines := strings.Split(data, "\n")

	// Determinar el siguiente número para el nuevo usuario
	var nextID int = 1
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) > 0 {
			id, err := strconv.Atoi(words[0])
			if err == nil && id >= nextID {
				nextID = id + 1
			}
		}
	}

	// Formatear la nueva entrada de usuario
	newUser := fmt.Sprintf("%d, U, %s, %s, %s\n", nextID, group, username, password)

	// Agregar el nuevo usuario al archivo
	err = AppendToFileBlock(&crrInode, newUser, file, tempSuperblock)
	if err != nil {
		fmt.Println("Error al agregar el usuario:", err)
		return
	}

	fmt.Println("Usuario agregado con éxito:", newUser)
	utilidades.AgregarRespuesta("Usuario agregado con éxito")
	fmt.Println("======End ADD USER======")
}*/

func AddUser(user string, pass string, grp string) {
	if !IsUserLoggedIn() {
		fmt.Println("Error: No hay un usuario logueado")
		utilidades.AgregarRespuesta("Error: No hay un usuario logueado")
		return
	}

	if Grupo != "root" {
		fmt.Println("Error: El usuario no tiene permiso de lectura (permiso 777 requerido)")
		utilidades.AgregarRespuesta("Error: El usuario no tiene permiso de lectura (permiso 777 requerido)")
		return
	}

	ParticionesMount := manejadorDisco.GetMountedPartitions()
	var filepath string
	var id string

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

	file, err := utilidades.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

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

	var tempSuperblock estructuras.Superblock
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el superblock:", err)
		return
	}

	indexInode := BuscarStart("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo /users.txt")
		return
	}

	var crrInode estructuras.Inode
	if err := utilidades.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(estructuras.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el inodo del archivo /users.txt")
		return
	}

	data := GetInodeFileData(crrInode, file, tempSuperblock)

	cleanedData := LImpiarNull(data)

	if strings.Contains(cleanedData, fmt.Sprintf("U,%s,", user)) {
		fmt.Println("Error: El usuario ya existe")
		utilidades.AgregarRespuesta("Error: El usuario ya existe")
		return
	}

	lastGroupID := 1
	lines := strings.Split(cleanedData, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "1,G,") {
			lastGroupID = 1
		} else if strings.HasPrefix(line, "2,G,") {
			lastGroupID = 2
		}
	}

	newUserData := fmt.Sprintf("%d,U,%s,%s, %s\n", lastGroupID, grp, user, pass)
	cleanedData += newUserData

	fmt.Println("Data:", cleanedData)

	if err := ActuaFileBlock(&crrInode, cleanedData, file, tempSuperblock); err != nil {
		fmt.Println("Error: No se pudo actualizar el archivo /users.txt:", err)
		return
	}

	fmt.Println("Usuario creado con éxito:", user)
	utilidades.AgregarRespuesta("Usuario creado con éxito " + user)
}

func LImpiarNull(data string) string {
	cleaneData := strings.TrimRight(data, "\x00")
	return cleaneData
}

func ActuaFileBlock(inode *estructuras.Inode, newData string, file *os.File, superblock estructuras.Superblock) error {
	fmt.Println("FullData:", newData)

	// Escribir el contenido actualizado en el bloque
	var updatedFileBlock estructuras.Fileblock
	copy(updatedFileBlock.B_content[:], newData)
	if err := utilidades.WriteObject(file, updatedFileBlock, int64(superblock.S_block_start+inode.I_block[0]*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
		return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
	}

	// Actualizar el tamaño del inodo
	inode.I_size = int32(len(newData))
	if err := utilidades.WriteObject(file, *inode, int64(superblock.S_inode_start+inode.I_block[0]*int32(binary.Size(estructuras.Inode{})))); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}

// Función para verificar si el usuario tiene permisos
func tienePermiso() bool {
	ParticionesMount := manejadorDisco.GetMountedPartitions()
	var filepath string
	var id string

	for _, partitions := range ParticionesMount {
		for _, partition := range partitions {
			// Verifica si alguna partición tiene un usuario logueado
			if partition.Logeado {
				filepath = partition.Path
				id = partition.ID
				break
			}
		}
	}

	file, err := utilidades.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return false
	}
	defer file.Close()

	var TempMBR estructuras.MBR

	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return false
	}

	var index int = -1

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 {
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) {
				if TempMBR.Partitions[i].Status[0] == '1' {
					index = i
				} else {
					return false
				}
				break
			}
		}
	}

	if index == -1 {
		return false
	}

	var tempSuperblock estructuras.Superblock
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		return false
	}

	indexInode := BuscarStart("/users.txt", file, tempSuperblock)

	var crrInode estructuras.Inode

	if err := utilidades.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(estructuras.Inode{})))); err != nil {
		return false
	}

	perm := string(crrInode.I_perm[:])
	return strings.Contains(perm, "777")
}

// Función modificada para buscar y leer Fileblocks en lugar de Folderblocks
func BuscarStart(path string, file *os.File, tempSuperblock estructuras.Superblock) int32 {
	TempStepsPath := strings.Split(path, "/")
	RutaPasada := TempStepsPath[1:]

	var Inode0 estructuras.Inode
	if err := utilidades.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1
	}

	return BuscarInodoRuta(RutaPasada, Inode0, file, tempSuperblock)
}

// Cambiado para manejar Fileblock en lugar de Folderblock
func BuscarInodoRuta(RutaPasada []string, Inode estructuras.Inode, file *os.File, tempSuperblock estructuras.Superblock) int32 {
	SearchedName := strings.Replace(pop(&RutaPasada), " ", "", -1)

	for _, block := range Inode.I_block {
		if block != -1 {
			if len(RutaPasada) == 0 { // Caso donde encontramos el archivo
				var fileblock estructuras.Fileblock
				if err := utilidades.ReadObject(file, &fileblock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
					return -1
				}

				estructuras.PrintFileblock(fileblock) // Imprime el contenido del Fileblock
				return 1
			} else {
				// En este caso seguimos buscando en los bloques de carpetas
				var crrFolderBlock estructuras.Folderblock
				if err := utilidades.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(estructuras.Folderblock{})))); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_content {
					if strings.Contains(string(folder.B_name[:]), SearchedName) {
						var NextInode estructuras.Inode
						if err := utilidades.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(estructuras.Inode{})))); err != nil {
							return -1
						}

						return BuscarInodoRuta(RutaPasada, NextInode, file, tempSuperblock)
					}
				}
			}
		}
	}

	return -1
}

func AddGroup(grupos string) {

	if !IsUserLoggedIn() {
		fmt.Println("Error: No hay un usuario logueado")
		utilidades.AgregarRespuesta("Error: No hay un usuario logueado")
		return
	}

	// Verificar si el usuario tiene permiso para escribir
	if Grupo != "root" {
		fmt.Println("Error: El usuario no tiene permiso de escritura")
		utilidades.AgregarRespuesta("Error: El usuario no tiene permiso de escritura")
		return
	}

	ParticionesMount := manejadorDisco.GetMountedPartitions()
	var filepath string
	var id string

	for _, particiones := range ParticionesMount {
		for _, particion := range particiones {
			if particion.Logeado {
				filepath = particion.Path
				id = particion.ID
				break
			}
		}
		if filepath != "" {
			break
		}
	}

	file, err := utilidades.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		utilidades.AgregarRespuesta("Error: No se pudo abrir el archivo" + filepath)
		return
	}
	defer file.Close()

	var TempMBR estructuras.MBR
	if err := utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Encontrar la partición activa
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
		utilidades.AgregarRespuesta("Error: No se encontró la partición")
		return
	}

	// Leer el Superbloque
	var tempSuperblock estructuras.Superblock
	if err := utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		utilidades.AgregarRespuesta("Error: No se pudo leer el Superblock")
		return
	}

	indexInode := BuscarStart("/user.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo usuarios.txt")
		utilidades.AgregarRespuesta("Error: No se encontró el archivo usuarios.txt")
		return
	}

	var crrInode estructuras.Inode
	if err := utilidades.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(estructuras.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inode del archivo usuarios.txt")
		return
	}

	var lastFileBlock estructuras.Fileblock
	var lastBlockIndex int32 = -1

	for _, block := range crrInode.I_block {
		if block != -1 {
			// Leer cada FileBlock
			var fileBlock estructuras.Fileblock
			if err := utilidades.ReadObject(file, &fileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
				fmt.Println("Error: No se pudo leer el Fileblock")
				return
			}
			lastFileBlock = fileBlock
			lastBlockIndex = block
		}
	}

	if lastBlockIndex == -1 {
		fmt.Println("Error: No se encontró un bloque válido en usuarios.txt")
		utilidades.AgregarRespuesta("Error: No se encontró un bloque válido en usuarios.txt")
		return
	}

	newGroupEntry := fmt.Sprintf("%d, G, %s\n", verNuevoGruPO(lastFileBlock), grupos)

	copy(lastFileBlock.B_content[len(string(lastFileBlock.B_content[:])):], []byte(newGroupEntry))

	if err := utilidades.WriteObject(file, lastFileBlock, int64(tempSuperblock.S_block_start+lastBlockIndex*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
		fmt.Println("Error: No se pudo escribir el bloque actualizado")
		utilidades.AgregarRespuesta("Error: No se pudo escribir el bloque actualizado")
		return
	}

	if err := AgregarGRupo(&lastFileBlock, grupos); err != nil {
		fmt.Println(err)
		utilidades.AgregarRespuesta(err.Error())
		return
	}

	if err := utilidades.WriteObject(file, lastFileBlock, int64(tempSuperblock.S_block_start+lastBlockIndex*int32(binary.Size(estructuras.Fileblock{})))); err != nil {
		fmt.Println("Error: No se pudo escribir el bloque actualizado")
		utilidades.AgregarRespuesta("Error: No se pudo escribir el bloque actualizado")
		return
	}

	fmt.Println("Grupo creado exitosamente")
	utilidades.AgregarRespuesta("Grupo creado exitosamente")
}

func verNuevoGruPO(fileBlock estructuras.Fileblock) int {

	content := string(fileBlock.B_content[:])
	lines := strings.Split(content, "\n")
	var lastID int

	for _, line := range lines {
		if strings.Contains(line, "G") {
			parts := strings.Split(line, ",")
			if len(parts) > 0 {
				id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
				if err == nil {
					lastID = id
				}
			}
		}
	}
	return lastID + 1
}

func AgregarGRupo(fileblock *estructuras.Fileblock, nombreDeGrupo string) error {

	content := strings.TrimRight(string(fileblock.B_content[:]), "\x00")

	lines := strings.Split(content, "\n")

	IdGrupo := 0
	for _, line := range lines {

		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		fields := strings.Split(trimmedLine, ",")

		if len(fields) >= 2 && strings.TrimSpace(fields[1]) == "G" {

			if len(fields) >= 3 && strings.TrimSpace(fields[2]) == nombreDeGrupo {
				utilidades.AgregarRespuesta("El grupo " + nombreDeGrupo + " ya existe en la partición")
				return fmt.Errorf("Error: El grupo '%s' ya existe en la particion", nombreDeGrupo)
			}

			// Obtener el ID del grupo
			id, err := strconv.Atoi(strings.TrimSpace(fields[0]))
			if err != nil {
				utilidades.AgregarRespuesta("Error al analizar el ID del grupo")
				return fmt.Errorf("Error al analizar el ID del grupo: %v", err)
			}

			if id > IdGrupo {
				IdGrupo = id
			}
		}
	}

	nuevoIDgrupo := IdGrupo + 1

	newGroupEntry := fmt.Sprintf("%d, G, %s", nuevoIDgrupo, nombreDeGrupo)

	content += "\n" + newGroupEntry

	copy(fileblock.B_content[:], []byte(content))

	return nil
}
