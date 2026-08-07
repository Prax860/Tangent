package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prax860/tangent/internal/types"
)

// Resolver resolves a package name into its canonical identifier.
type Resolver interface {
	Resolve(pkg string) (ResolvedPackage, error)
}

// ResolvedPackage contains metadata about a resolved package.
type ResolvedPackage struct {
	Original string
	Target   string
	Registry string
	Cached   bool
}

// httpClient is shared by every resolver below so the timeout is configured
// in exactly one place.
var httpClient = &http.Client{Timeout: 5 * time.Second}

// ResolveForWorkspace chooses the correct resolver based on the workspace.
func ResolveForWorkspace(ws types.WorkspaceType, pkg string) ResolvedPackage {

	var resolver Resolver

	switch ws {

	case types.WorkspaceGo:
		resolver = GoResolver{}

	case types.WorkspaceNode:
		resolver = NodeResolver{}

	case types.WorkspacePython:
		resolver = PythonResolver{}

	default:
		return ResolvedPackage{
			Original: pkg,
			Target:   pkg,
			Registry: "unknown",
			Cached:   false,
		}
	}

	resolved, err := resolver.Resolve(pkg)

	if err != nil {
		return ResolvedPackage{
			Original: pkg,
			Target:   pkg,
			Registry: string(ws),
			Cached:   false,
		}
	}

	return resolved
}

// --- Node: npm registry search API -----------------------------------------

// NodeResolver resolves against the public npm registry search API.
type NodeResolver struct{}

func (NodeResolver) Resolve(pkg string) (ResolvedPackage, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/-/v1/search?text=%s&size=1", pkg)

	resp, err := httpClient.Get(url)
	if err != nil {
		return ResolvedPackage{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Objects []struct {
			Package struct {
				Name string `json:"name"`
			} `json:"package"`
		} `json:"objects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ResolvedPackage{}, err
	}
	if len(result.Objects) == 0 {
		return ResolvedPackage{}, fmt.Errorf("no npm package found for %q", pkg)
	}

	return ResolvedPackage{
		Original: pkg,
		Target:   result.Objects[0].Package.Name,
		Registry: "npm",
		Cached:   false,
	}, nil
}

// --- Python: PyPI JSON API ---------------------------------------------

// PythonResolver resolves against PyPI's per-project JSON endpoint.
// PyPI has no public search API, so this is an exact-name existence check —
// which is fine, since pip install requires an exact name too.
type PythonResolver struct{}

func (PythonResolver) Resolve(pkg string) (ResolvedPackage, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", pkg)

	resp, err := httpClient.Get(url)
	if err != nil {
		return ResolvedPackage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResolvedPackage{}, fmt.Errorf("no PyPI package found for %q", pkg)
	}

	var result struct {
		Info struct {
			Name string `json:"name"`
		} `json:"info"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ResolvedPackage{}, err
	}

	return ResolvedPackage{
		Original: pkg,
		Target:   result.Info.Name,
		Registry: "pypi",
		Cached:   false,
	}, nil
}

// --- Go: verified seed table + common-pattern probe -------------------------

// GoResolver resolves short names (e.g. "gin") into full module import
// paths (e.g. "github.com/gin-gonic/gin").
//
// There's no public search API for the Go module proxy — only exact-path
// existence checks. goSeeds is a small, easily-extended hint table; every
// hit is still verified against proxy.golang.org before being trusted.
// Unseeded names fall back to the most common Go repo pattern
// "github.com/<pkg>/<pkg>", also verified. If nothing verifies, resolution
// fails and ResolveForWorkspace passes the original input straight through.
type GoResolver struct{}

var goSeeds = map[string]string{
	"gin":     "github.com/gin-gonic/gin",
	"echo":    "github.com/labstack/echo/v4",
	"cobra":   "github.com/spf13/cobra",
	"gorm":    "gorm.io/gorm",
	"chi":     "github.com/go-chi/chi/v5",
	"fiber":   "github.com/gofiber/fiber/v2",
	"testify": "github.com/stretchr/testify",
	"zap":     "go.uber.org/zap",
	"viper":   "github.com/spf13/viper",
}

func (GoResolver) Resolve(pkg string) (ResolvedPackage, error) {
	candidates := []string{}
	if seeded, ok := goSeeds[pkg]; ok {
		candidates = append(candidates, seeded)
	}
	candidates = append(candidates, fmt.Sprintf("github.com/%s/%s", pkg, pkg))

	for _, path := range candidates {
		if moduleExists(path) {
			return ResolvedPackage{
				Original: pkg,
				Target:   path,
				Registry: "go-proxy",
				Cached:   false,
			}, nil
		}
	}

	return ResolvedPackage{}, fmt.Errorf("no Go module found for %q", pkg)
}

func moduleExists(importPath string) bool {
	url := fmt.Sprintf("https://proxy.golang.org/%s/@latest", importPath)
	resp, err := httpClient.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}