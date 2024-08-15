package estructuras

//"fmt"

type MRB struct {
	MbrSize      int32    // 4 bytes //int32 va desde -2,147,483,648 hasta 2,147,483,647.
	CreationDate [19]byte // 10 bytes
	Signature    int32    // 4 bytes
	Fit          [1]byte  // 1 byte
}

//func PrintMBR(data MRB) {
//fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d, Signature: %d",
//string(data.CreationDate[:]),
//string(data.Fit[:]),
//data.MbrSize,
//data.Signature))
//}
