package config

import (
	"log"
)

// InitLogger configures the standard logger. Call early in startup if needed.
func InitLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
