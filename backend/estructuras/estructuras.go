package estructuras

import (
	"backend/utilidades"
	"fmt"
	"strings"
	"unicode"
)

//"fmt"

type MBR struct {
	MbrSize      int32    // 4 bytes //int32 va desde -2,147,483,648 hasta 2,147,483,647.
	CreationDate [19]byte // 19 bytes
	Signature    int32    // 4 bytes
	Fit          [1]byte  // 1 byte
	Partitions   [4]Partition
	//Extendida    [1]byte
}

// Total bytes: 4+19+4+1 = 28

type Partition struct {
	Status      [1]byte  // 11
	Type        [1]byte  // 1
	Fit         [1]byte  // 1
	Start       int32    // 4
	Size        int32    // 4
	Name        [16]byte // 16
	Correlative int32    // 4
	Id          [4]byte  // 4
}

type EBR struct {
	PartMount [1]byte
	PartFit   [1]byte
	PartStart int32
	PartSize  int32
	PartNext  int32
	PartName  [16]byte
}

func PrintMBR(data MBR) {
	fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d", string(data.CreationDate[:]), string(data.Fit[:]), data.MbrSize))
	for i := 0; i < 4; i++ {
		PrintPartition(data.Partitions[i])
	}
}

func ImprimirParticion(datos Partition) {
	fmt.Printf("Nombre: %s, Tipo: %s, Inicio: %d, Tamaño: %d, Estado: %s, ID: %s, Ajuste: %s, Correlativo: %d\n",
		string(datos.Name[:]), string(datos.Type[:]), datos.Start, datos.Size, string(datos.Status[:]),
		string(datos.Id[:]), string(datos.Fit[:]), datos.Correlative)
}

func PrintPartition(data Partition) {
	fmt.Println(fmt.Sprintf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s, Correlative: %d", string(data.Name[:]), string(data.Type[:]), data.Start, data.Size, string(data.Status[:]), string(data.Id[:]), data.Correlative))
}

func PrintEBR(data EBR) {
	fmt.Println(fmt.Sprintf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c",
		string(data.PartName[:]),
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount))
}

//func PrintMBR(data MRB) {
//fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d, Signature: %d",
//string(data.CreationDate[:]),
//string(data.Fit[:]),
//data.MbrSize,
//data.Signature))
//}

//Estructuras relacionadas a EXT2

type Superblock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [19]byte
	S_umtime            [19]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_fist_ino          int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

func PrintSuperblock(sb Superblock) {
	fmt.Println("====== Superblock ======")
	fmt.Printf("S_filesystem_type: %d\n", sb.S_filesystem_type)
	fmt.Printf("S_inodes_count: %d\n", sb.S_inodes_count)
	fmt.Printf("S_blocks_count: %d\n", sb.S_blocks_count)
	fmt.Printf("S_free_blocks_count: %d\n", sb.S_free_blocks_count)
	fmt.Printf("S_free_inodes_count: %d\n", sb.S_free_inodes_count)
	fmt.Printf("S_mtime: %s\n", string(sb.S_mtime[:]))
	fmt.Printf("S_umtime: %s\n", string(sb.S_umtime[:]))
	fmt.Printf("S_mnt_count: %d\n", sb.S_mnt_count)
	fmt.Printf("S_magic: 0x%X\n", sb.S_magic) // Usamos 0x%X para mostrarlo en formato hexadecimal
	fmt.Printf("S_inode_size: %d\n", sb.S_inode_size)
	fmt.Printf("S_block_size: %d\n", sb.S_block_size)
	fmt.Printf("S_fist_ino: %d\n", sb.S_fist_ino)
	fmt.Printf("S_first_blo: %d\n", sb.S_first_blo)
	fmt.Printf("S_bm_inode_start: %d\n", sb.S_bm_inode_start)
	fmt.Printf("S_bm_block_start: %d\n", sb.S_bm_block_start)
	fmt.Printf("S_inode_start: %d\n", sb.S_inode_start)
	fmt.Printf("S_block_start: %d\n", sb.S_block_start)
	fmt.Println("========================")
}

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

func PrintInode(inode Inode) {
	fmt.Println("====== Inode ======")
	fmt.Printf("I_uid: %d\n", inode.I_uid)
	fmt.Printf("I_gid: %d\n", inode.I_gid)
	fmt.Printf("I_size: %d\n", inode.I_size)
	fmt.Printf("I_atime: %s\n", string(inode.I_atime[:]))
	fmt.Printf("I_ctime: %s\n", string(inode.I_ctime[:]))
	fmt.Printf("I_mtime: %s\n", string(inode.I_mtime[:]))
	fmt.Printf("I_type: %s\n", string(inode.I_type[:]))
	fmt.Printf("I_perm: %s\n", string(inode.I_perm[:]))
	fmt.Printf("I_block: %v\n", inode.I_block)
	fmt.Println("===================")
}

type Folderblock struct {
	B_content [4]Content
}

func PrintFolderblock(folderblock Folderblock) {
	fmt.Println("====== Folderblock ======")
	for i, content := range folderblock.B_content {
		fmt.Printf("Content %d: Name: %s, Inodo: %d\n", i, string(content.B_name[:]), content.B_inodo)
	}
	fmt.Println("=========================")
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type Fileblock struct {
	B_content [64]byte
}

func PrintFileblock(fileblock Fileblock) {
	fmt.Println("====== Fileblock ======")
	fmt.Printf("B_content: %s\n", string(fileblock.B_content[:]))
	fmt.Println("=======================")
}

func AgregarFileBlockConsola(fileblock Fileblock) {

	content := strings.TrimRight(string(fileblock.B_content[:]), "\x00")

	lines := strings.Split(content, "\n")

	for _, line := range lines {

		printableLine := strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) || r == '\t' || r == '\n' {
				return r
			}
			return -1
		}, line)

		if trimmedLine := strings.TrimSpace(printableLine); trimmedLine != "" {
			utilidades.AgregarRespuesta(trimmedLine)
		}
	}

}

type Pointerblock struct {
	B_pointers [16]int32
}

func PrintPointerblock(pointerblock Pointerblock) {
	fmt.Println("====== Pointerblock ======")
	for i, pointer := range pointerblock.B_pointers {
		fmt.Printf("Pointer %d: %d\n", i, pointer)
	}
	fmt.Println("=========================")
}
