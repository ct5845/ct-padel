package fileserver

import (
	"crypto/md5"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CachedFileServer struct {
	dir   string
	etags map[string]string
	mutex sync.RWMutex
}

type responseWriterWrapper struct {
	http.ResponseWriter
	etag      string
	hasETag   bool
	headerSet bool
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	if !rw.headerSet && rw.hasETag {
		rw.Header().Set("ETag", rw.etag)
		rw.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
		rw.headerSet = true
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

func NewCachedFileServer(dir string) *CachedFileServer {
	cfs := &CachedFileServer{
		dir:   dir,
		etags: make(map[string]string),
	}
	cfs.buildETags()
	return cfs
}

func (cfs *CachedFileServer) buildETags() {
	cfs.mutex.Lock()
	defer cfs.mutex.Unlock()

	err := filepath.Walk(cfs.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				slog.Error("Error opening file", "path", path, "error", err)
				return nil
			}
			defer file.Close()

			hash := md5.New()
			if _, err := io.Copy(hash, file); err != nil {
				slog.Error("Error hashing file", "path", path, "error", err)
				return nil
			}

			// Calculate relative path from dir and convert to URL path
			relPath, err := filepath.Rel(cfs.dir, path)
			if err != nil {
				return err
			}

			// Convert Windows paths to URL paths
			urlPath := strings.ReplaceAll(relPath, "\\", "/")
			etag := fmt.Sprintf(`"%x"`, hash.Sum(nil))
			cfs.etags[urlPath] = etag
		}
		return nil
	})

	if err != nil {
		slog.Error("Error building ETags", "error", err)
	}
}

func (cfs *CachedFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cfs.mutex.RLock()
	etag, hasETag := cfs.etags[r.URL.Path]
	cfs.mutex.RUnlock()


	// Create a wrapped ResponseWriter to intercept header writes
	wrappedWriter := &responseWriterWrapper{
		ResponseWriter: w,
		etag:          etag,
		hasETag:       hasETag,
		headerSet:     false,
	}

	if hasETag {
		// Handle If-None-Match header
		if match := r.Header.Get("If-None-Match"); match == etag {
			w.Header().Set("ETag", etag)
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// Serve the file using standard file server with our wrapped writer
	fileServer := http.FileServer(http.Dir(cfs.dir))
	fileServer.ServeHTTP(wrappedWriter, r)
}

func (cfs *CachedFileServer) RefreshETags() {
	cfs.buildETags()
}
