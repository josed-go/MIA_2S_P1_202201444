package utilidades

var respuesta string = ""

func AgregarRespuesta(res string) {
	if len(respuesta) == 0 {
		respuesta += res
	} else {
		respuesta += "\n" + res
	}
}

func ObtenerRespuestas() string {
	return respuesta
}

func LimpiarConsola() {
	respuesta = ""
}
