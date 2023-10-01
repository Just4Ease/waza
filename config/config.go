package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Secrets struct {
	Port        string      `json:"PORT"`
	Environment Environment `json:"ENVIRONMENT"`
}

var ss Secrets

const ServiceName = "waza"
const empty = ""

func init() {
	importPath := fmt.Sprintf("%s/config", strings.ReplaceAll(ServiceName, "-", "."))
	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err == nil {
		env := filepath.Join(p.Dir, "../.env")
		_ = godotenv.Load(env)
	}

	ss = Secrets{}

	if ss.Environment = Environment(os.Getenv("ENVIRONMENT")); ss.Environment.IsValid() != nil {
		log.Fatal("Error in environment variables: ", err)
	}

	if ss.Port = os.Getenv("PORT"); ss.Port == empty {
		ss.Port = "4000" // Defaults the port to 4000
	}

}

func GetSecrets() Secrets {
	return ss
}
