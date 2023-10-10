import (
	"http/net"
	"json"
	"os"
	"log"
)

func SetHeader() {
	value, ok = os.LookupEnv("JUPYTERHUB_USERNAME")
	if !ok {
		log.Fatal(
	}

	value, ok = os.LookupEnv("JUPYTERHUB_PASSWORD")
}
