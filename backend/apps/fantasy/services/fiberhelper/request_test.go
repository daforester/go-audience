package fiberhelper

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Request data must be copied, never pointed to. Controllers, helpers and middleware read
// requests only through this package; this test fails if one reaches into fiber.Ctx directly.
func TestNoDirectRequestAccess(t *testing.T) {
	forbidden := regexp.MustCompile(`\bc\.(Params|Query|Get|Body|Cookies|FormValue|Host|Hostname|BodyRaw|Bind)\(`)
	dirs := []string{"../../app", "../../middleware"}

	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for i, line := range strings.Split(string(data), "\n") {
				if forbidden.MatchString(line) {
					t.Errorf("%s:%d reads the request directly, use fiberhelper instead: %s", path, i+1, strings.TrimSpace(line))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
