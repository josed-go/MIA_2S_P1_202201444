package analizador

import (
	"backend/comandos"
	"backend/manejadorDisco"
	"backend/sistema"
	"backend/usuarios"
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

	for scanner.Scan() {
		cont += 1
		line := scanner.Text()

		if strings.Contains(line, "#") {
			line = strings.Split(line, "#")[0]
		}

		line = strings.TrimSpace(line)

		if len(line) > 0 {
			fmt.Println("Comando leído:", line)
			command, params := getCommandAndParams(line)
			AnalyzeCommand(command, params, strconv.Itoa(cont))
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error leyendo el input:", err)
	}
}

func AnalyzeCommand(command string, params string, linea string) {

	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params, linea)
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params, linea)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params, linea)
	} else if strings.Contains(command, "mount") {
		fn_mount(params, linea)
	} else if strings.Contains(command, "cat") {
		fn_cat(params, linea)
	} else if strings.Contains(command, "mkfs") {
		fn_mkfs(params, linea)
	} else if strings.Contains(command, "rep") {
		fn_rep(params, linea)
	} else if strings.Contains(command, "login") {
		fn_login(params, linea)
	} else if strings.Contains(command, "logout") {
		usuarios.Logout()
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
	fit := fs.String("fit", "wf", "Ajuste")
	name := fs.String("name", "", "Nombre")
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(os.Args[1:])

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := strings.ToLower(match[1])  // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
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
		*fit = "wf"
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

	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: El parametro fit debe ser bf - ff - wf aqui")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El parametro fit debe ser bf - ff - wf")
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

/*func fn_mount(params string, linea string) {
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
}*/

func fn_mount(params string, linea string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre de la partición")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		fs.Set(flagName, flagValue)
	}

	if *path == "" || *name == "" {
		fmt.Println("Error: Path es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Path es obligatorio")
		return
	}

	if *name == "" {
		fmt.Println("Error: Name es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Name es obligatorio")
		return
	}

	lowercaseName := strings.ToLower(*name)
	manejadorDisco.Mount(*path, lowercaseName, linea)
}

func fn_cat(params string, linea string) {
	//fs := flag.NewFlagSet("cat", flag.ExitOnError)

	// Usaremos un mapa para almacenar los archivos
	files := make(map[int]string)

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]                   // match[1]: Captura y guarda el nombre del flag (por ejemplo, "file1", "file2", etc.)
		flagValue := strings.ToLower(match[2]) //strings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas

		flagValue = strings.Trim(flagValue, "\"")

		// Si el flagName empieza con "file" y es seguido por un número
		if strings.HasPrefix(flagName, "file") {
			// Extraer el número después de "file"
			fileNumber, err := strconv.Atoi(strings.TrimPrefix(flagName, "file"))
			if err != nil {
				fmt.Println("Error: Nombre de archivo inválido")
				utilidades.AgregarRespuesta("Error en linea " + linea + " : Nombre de archivo inválido")
				return
			}

			if flagValue == "" {
				fmt.Println("Error: parametro -file" + string(fileNumber) + " no contiene ninguna ruta")
				utilidades.AgregarRespuesta("Error en linea " + linea + " : parametro -file" + string(fileNumber) + " no contiene ninguna ruta")
			}

			files[fileNumber] = flagValue
		} else {
			fmt.Println("Error: Flag not found")
		}
	}

	// Convertir el mapa a un slice ordenado
	var orderedFiles []string
	for i := 1; i <= len(files); i++ {
		if file, exists := files[i]; exists {
			orderedFiles = append(orderedFiles, file)
		} else {
			fmt.Println("Error: Falta un archivo en la secuencia")
			utilidades.AgregarRespuesta("Error en linea " + linea + " : Falta un archivo en la secuencia")
			return
		}
	}

	if len(orderedFiles) == 0 {
		fmt.Println("Error: No se encontraron archivos")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : No se encontraron archivos")
		return
	}

	// Llamar a la función para manejar los archivos en orden
	comandos.Cat(orderedFiles, linea)
}

func fn_mkfs(input string, linea string) {
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "full", "Tipo")
	fs_ := fs.String("fs", "full", "Fs")

	// Parse the input string, not os.Args
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Verifica que se hayan establecido todas las flags necesarias
	if *id == "" {
		fmt.Println("Error: id es un parámetro obligatorio.")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : id es un parámetro obligatorio.")
		return
	}

	if *type_ == "" {
		/*fmt.Println("Error: type es un parámetro obligatorio.")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : type es un parámetro obligatorio.")
		return*/
		*type_ = "full"
	}

	// Llamar a la función
	sistema.Mkfs(*id, *type_, *fs_)
}

func fn_rep(params string, linea string) {

	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre")
	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {

		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "path", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}

	}

	if *id == "" {
		fmt.Println("Error: Id es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Id es obligatorio")

		return
	}

	if *path == "" {
		fmt.Println("Error: Path es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Path es obligatorio")
		return
	}

	if *name == "" {
		fmt.Println("Error: Name es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : Name es obligatorio")
		return
	}

	comandos.Reportes(*id, *path, *name, linea)

}

func fn_login(input string, linea string) {
	// Definir las flags
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "Id")

	// Parsearlas
	fs.Parse(os.Args[1:])

	// Match de flags en el input
	matches := re.FindAllStringSubmatch(input, -1)

	// Procesar el input
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	if *user == "" {
		fmt.Println("Error: El usuario es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El usuario es obligatorio")
		return
	}

	if *pass == "" {
		fmt.Println("Error: La contraseña es obligatoria")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : La contraseña es obligatoria")
		return
	}

	if *id == "" {
		fmt.Println("Error: El id es obligatorio")
		utilidades.AgregarRespuesta("Error en linea " + linea + " : El id es obligatorio")
		return
	}

	usuarios.Login(*user, *pass, *id)

}
