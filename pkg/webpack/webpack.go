package webpack

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Manifest :
type Manifest struct {
	Assets      map[string]string `json:"files"`
	Entrypoints []string          `json:"entrypoints"`
}

// PreloadAssets : Generates link tag for each webpack entrypoint in manifest.json
func PreloadAssets(file string) Manifest {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	w := Manifest{}
	err = json.Unmarshal(f, &w)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	for k := range w.Assets {
		if strings.HasSuffix(k, ".map") {
			delete(w.Assets, k)
		}
	}
	return w
}
