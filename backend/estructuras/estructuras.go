package estructuras

import "fmt"

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
