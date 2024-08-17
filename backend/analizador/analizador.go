package analizador

import (
	"backend/manejadorDisco"
	"backend/utilidades"
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

//input := "mkdisk -size=3000 -unit=K -fit=BF -path=/home/bang/Disks/disk1.bin"

/*
parts[0] es "mkdisk"
*/

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		parametrosmin := strings.ToLower(params)
		return command, parametrosmin
	}
	return "", input

	/*Después de procesar la entrada:
	command será "mkdisk".
	params será "-size=3000 -unit=K -fit=BF -path=/home/bang/Disks/disk1.bin".*/
}

func Analyze(entrada string) {
	utilidades.LimpiarConsola()
	cont := 0
	scanner := bufio.NewScanner(strings.NewReader(entrada))
	// Leer el string línea por línea
	for scanner.Scan() {
		cont += 1
		line := scanner.Text() // Obtener la línea actual
		fmt.Println("Comando leído:", line)

		if len(line) > 0 {
			command, params := getCommandAndParams(line)

			AnalyzeCommand(command, params, strconv.Itoa(cont))
		}

	}

	// Manejo de errores durante la lectura
	if err := scanner.Err(); err != nil {
		fmt.Println("Error leyendo el input:", err)
	}
}

func AnalyzeCommand(command string, params string, linea string) {

	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params, linea)
		//} else if strings.Contains(command, "rep") {
		//fmt.Print("COMANDO REP")
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params, linea)

	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params, linea)
	} else if strings.Contains(command, "mount") {
		fn_mount(params, linea)
	} else {
		fmt.Println("Error: Commando invalido o no encontrado")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Commando invalido o no encontrado")
	}

}

func fn_mkdisk(params string, linea string) {
	// Definir flag
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamano")
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:])

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]                   // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := strings.ToLower(match[2]) //trings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Validaciones
	if *size <= 0 {
		fmt.Println("Error: El tamano debe ser mayor a 0")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El tamano debe ser mayor a 0")
		return
	}

	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: El parametro fit debe ser bf - ff - wf")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro fit debe ser bf - ff - wf")
		return
	}

	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: El parametro unit debe ser k - m")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro unit debe ser k - m")
		return
	}

	if *path == "" {
		fmt.Println("Error: El parametro path es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro path es obligatorio")
		return
	}

	// LLamamos a la funcion
	manejadorDisco.Mkdisk(*size, *fit, *unit, *path)
}

func fn_rmdisk(params string, linea string) {
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:])

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]                   // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := strings.ToLower(match[2]) //trings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	if *path == "" {
		fmt.Println("Error: El parametro path es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro path es obligatorio")
		return
	}

	// LLamamos a la funcion
	manejadorDisco.Rmdisk(*path, linea)
}

func fn_fdisk(params string, linea string) {
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamano")
	unit := fs.String("unit", "k", "Unidad")
	typ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "", "Ajuste")
	name := fs.String("name", "", "Nombre")
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:])

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]                   // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := strings.ToLower(match[2]) //trings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path", "type", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	if *fit == "" {
		*fit = "w"
	}

	// Validaciones
	if *size <= 0 {
		fmt.Println("Error: El tamano debe ser mayor a 0")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El tamano debe ser mayor a 0")
		return
	}

	if *unit != "k" && *unit != "m" && *unit != "b" {
		fmt.Println("Error: El parametro unit debe ser k - m")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro unit debe ser b - k - m")
		return
	}

	if *typ != "p" && *typ != "e" && *typ != "l" {
		fmt.Println("Error: El parametro type debe ser p - e - l")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro type debe ser p - e - l")
		return
	}

	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error: El parametro fit debe ser b - f - w")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro fit debe ser b - f - w")
		return
	}

	if *name == "" {
		fmt.Println("Error: El parametro name es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro name es obligatorio")
		return
	}

	if *path == "" {
		fmt.Println("Error: El parametro path es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro path es obligatorio")
		return
	}

	manejadorDisco.Fdisk(*size, *fit, *unit, *path, *typ, *name, linea)
}

func fn_mount(params string, linea string) {
	// Definir flag
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	name := fs.String("name", "", "Nombre")
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:])

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]                   // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := strings.ToLower(match[2]) //trings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "name", "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	if *path == "" {
		fmt.Println("Error: El parametro path es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro path es obligatorio")
		return
	}

	if *name == "" {
		fmt.Println("Error: El parametro name es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro name es obligatorio")
		return
	}

	// LLamamos a la funcion
	manejadorDisco.Mount(*name, *path, linea)
}
