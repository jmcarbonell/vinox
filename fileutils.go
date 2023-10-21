package vinox

import (
	"path/filepath"
)

func ChangeExt(archivoOriginal, nuevaExtension string) string {
	// Obtener la ruta del archivo original
	rutaOriginal, nombreOriginal := filepath.Split(archivoOriginal)

	// Crear la nueva ruta con la extensión modificada
	nuevoNombre := nombreOriginal[:len(nombreOriginal)-len(filepath.Ext(nombreOriginal))] + nuevaExtension
	nuevaRuta := filepath.Join(rutaOriginal, nuevoNombre)

	return nuevaRuta
}
