package config_test

import (
	"os"
	"payment-process/internal/config"
	"testing"

	"github.com/tj/assert"
)

func TestLoadConfig(t *testing.T) {
	// JSON simulado para las pruebas
	mockConfig := `
	{
		"Database": {
			"Host": "localhost",
			"Database": "testdb",
			"Password": "password123",
			"Port": "5432",
			"User": "testuser"
		},
		"EmailSender": {
			"AccountId": "testaccount",
			"Token": "testtoken"
		}
	}
	`

	// Crear un archivo temporal para simular conf.json
	tempFile, err := os.CreateTemp("", "conf*.json")
	if err != nil {
		t.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Limpiar el archivo después de la prueba

	// Escribir el JSON simulado en el archivo temporal
	_, err = tempFile.Write([]byte(mockConfig))
	if err != nil {
		t.Fatalf("Error al escribir en el archivo temporal: %v", err)
	}

	// Cerrar el archivo para asegurarse de que se puede leer
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Error al cerrar archivo temporal: %v", err)
	}

	// Cambiar el archivo de configuración esperado a la ruta temporal
	originalFile := "./conf.json"
	defer os.Rename(originalFile, originalFile+"_backup") // Restaurar después de la prueba
	os.Rename(tempFile.Name(), originalFile)

	// Probar LoadConfig
	config, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig devolvió un error: %v", err)
	}

	assert.Equal(t, "testdb", config.Database.Database)


	
}
