package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// CopyAssetsToDist copies all files and directories from the specified assets directory
// to the dist/web directory.
func CopyAssetsToDist(assetsDir string) error {
	dstDir := filepath.Join("dist", assetsDir)

	return filepath.Walk(assetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path
		relPath, err := filepath.Rel(assetsDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// Create directory in destination
			return os.MkdirAll(targetPath, info.Mode())
		}

		// Open source file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Create destination file with same permissions
		dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// Copy content
		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

// UpdateWasmExec copies the wasm_exec.js file from GOROOT to dist/wasm_exec.js
// ensuring the JS glue code matches the Go version used to build the WASM binary.
func UpdateWasmExec() error {
	goroot := runtime.GOROOT()

	// Check possible locations for wasm_exec.js
	// Go 1.24+ moved it to lib/wasm, older versions had it in misc/wasm
	locations := []string{
		filepath.Join(goroot, "lib", "wasm", "wasm_exec.js"),
		filepath.Join(goroot, "misc", "wasm", "wasm_exec.js"),
	}

	var src string
	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			src = loc
			break
		}
	}

	if src == "" {
		return os.ErrNotExist
	}

	dst := filepath.Join("dist", "wasm_exec.js")

	log.Printf("Updating wasm_exec.js from %s to %s", src, dst)

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
