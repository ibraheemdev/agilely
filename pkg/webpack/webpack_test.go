package webpack

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ibraheemdev/poller/test"
)

func TestPreloadAssets(t *testing.T) {
	t.Parallel()
	test.MoveToRoot()
	manifest := PreloadAssets("test/manifest.json")
	want := []string{"static/js/runtime-main.a0de9167.js", "static/js/2.19dfafba.chunk.js", "static/css/main.5e844b86.chunk.css", "static/js/main.f096cf97.chunk.js"}
	if !reflect.DeepEqual(manifest.Entrypoints, want) {
		t.Errorf("expected %s, manifest: %s", want, manifest.Entrypoints)
	}

	got := "/static/css/main.5e844b86.chunk.css"
	if manifest.Assets["main.css"] != got {
		t.Errorf("expected %s, manifest: %s", got, manifest.Assets["main.css"])
	}

	for a := range manifest.Assets {
		if strings.HasSuffix(a, ".map") {
			t.Error("did not expect .map file")
		}
	}
}
