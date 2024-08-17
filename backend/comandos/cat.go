package comandos

import (
	"backend/utilidades"
	"fmt"
	"io"
)

func Cat(files []string, linea string) {
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
}
