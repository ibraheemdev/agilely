package webpack

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ibraheemdev/poller/config"
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
	for k, a := range w.Assets {
		if strings.HasSuffix(k, ".map") {
			delete(w.Assets, k)
			break
		}
		w.Assets[k] = config.Config.Server.Static.BuildPath + a
	}
	return w
}
