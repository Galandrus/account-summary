package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config estructura para manejar toda la configuración
type Config struct {
	Port     string
	MongoURI string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Default().Printf("Error loading .env file: %v", err)
	}

	config := &Config{
		Port:     getEnv("PORT"),
		MongoURI: getEnv("MONGO_URI"),
	}

	return config
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	log.Default().Fatalf("Environment variable %s is not set", key)
	return ""
}
