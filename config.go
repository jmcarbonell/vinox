package SAKnife

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Singleton para almacenar los datos clave-valor
type ConfigManager struct {
	datos    map[string]string
	filename string
	sync.RWMutex
}

var instancia *ConfigManager
var once sync.Once

func obtenerInstancia(afilename string) *ConfigManager {

	once.Do(func() {
		instancia = &ConfigManager{

			datos: make(map[string]string),
		}
		instancia.filename = afilename
		instancia.cargarDatos() // Reemplaza "tu_archivo.txt" con la ruta correcta
	})
	return instancia
}
func guardar(archivo string) error {
	instancia.RLock()
	defer instancia.RUnlock()

	file, err := os.Create(archivo)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for clave, valor := range instancia.datos {
		_, err := writer.WriteString(clave + "=" + valor + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

func (cm *ConfigManager) cargarDatos() {
	cm.Lock()
	defer cm.Unlock()

	archivoTexto, err := os.Open(cm.filename)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer archivoTexto.Close()

	lector := bufio.NewReader(archivoTexto)

	for {
		linea, err := lector.ReadString('\n')
		if err != nil {
			break
		}

		parts := strings.Split(linea, "=")
		if len(parts) == 2 {
			clave := strings.TrimSpace(parts[0])
			valor := strings.TrimSpace(parts[1])
			cm.datos[clave] = valor
		}
	}
}
func (cm *ConfigManager) Get(clave string) (string, bool) {
	cm.RLock()
	defer cm.RUnlock()

	valor, encontrado := cm.datos[clave]
	return valor, encontrado
}

func (cm *ConfigManager) Put(clave, avalue string) (string, bool) {
	cm.RLock()
	defer cm.RUnlock()

	valor, encontrado := cm.datos[clave]

	cm.datos[clave] = avalue
	return valor, encontrado
}

func (cm *ConfigManager) Has(clave string) bool {
	cm.RLock()
	defer cm.RUnlock()

	_, encontrado := cm.datos[clave]
	return encontrado
}

func (cm *ConfigManager) Flush() {
	cm.RLock()
	defer cm.RUnlock()
	guardar(cm.filename)
}
func (cm *ConfigManager) GetFilename() string {
	cm.RLock()
	defer cm.RUnlock()
	return cm.filename
}
func (cm *ConfigManager) ConvertirAJSON() (string, error) {
	instancia.RLock()
	defer cm.RUnlock()

	// Convertir el mapa de datos en JSON
	jsonBytes, err := json.MarshalIndent(cm.datos, "", "  ")
	if err != nil {
		return "", err
	}

	// Convertir los bytes JSON a una cadena
	jsonStr := string(jsonBytes)

	return jsonStr, nil
}

func ConfInst(afilename string) *ConfigManager {
	return obtenerInstancia(afilename)
}

func (this *ConfigManager) SaveToFile(filename string) error {

	this.RLock()
	defer this.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for clave, valor := range this.datos {
		_, err := writer.WriteString(clave + "=" + valor + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

func (this *ConfigManager) SaveJSONToFile(filename string) error {

	this.RLock()
	defer this.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	salida, err := this.ConvertirAJSON()
	if err != nil {
		return err
	}
	_, err = writer.WriteString(salida)
	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}
