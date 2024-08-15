package utilidades

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// Funcion para crear un archivo binario
func CreateFile(name string) error {
	//Se asegura que el archivo existe
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Crear archivo
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funcion para abrir un archivo binario ead/write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

func DeleteFile(name string, linea string) error {
	// Verifica si el archivo existe
	if _, err := os.Stat(name); os.IsNotExist(err) {
		AgregarRespuesta("Error en linea " + linea + " : El archivo " + name + " no existe")
		fmt.Printf("Error: El archivo %s no existe", name)
		return err
	}

	// Intenta eliminar el archivo
	err := os.Remove(name)
	if err != nil {
		fmt.Printf("Error eliminando archivo %s: %v", name, err)
		return err
	}

	fmt.Printf("Archivo %s eliminado correctamente\n", name)
	return nil
}

// Funcion para escribir un objecto en un archivo binario
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Funcion para leer un objeto de un archivo binario
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}
