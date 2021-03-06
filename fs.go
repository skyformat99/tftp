package tftp

import (
	"io"
	"net/http"
	"os"
)

type fileHandler struct {
	root http.FileSystem
}

// return a tftp.Handler that response file specified by request.
// argument fs could be http.Dir
func ReadonlyFileServer(fs http.FileSystem) Handler {
	return &fileHandler{root: fs}
}

func (h *fileHandler) ServeTFTPReadRequest(w io.WriteCloser, req *Request) error {
	f, err := h.root.Open(req.Filename)
	if err != nil {
		return toTFTPError(err)
	}
	stat, err := f.Stat()
	if err != nil {
		return toTFTPError(err)
	}
	if stat.IsDir() {
		return ErrFileNotFound
	}
	if _, err := io.Copy(w, f); err != nil {
		return err
	}
	w.Close()
	return nil
}

func (h *fileHandler) ServeTFTPWriteRequest(r io.Reader, req *Request) error {
	return ErrAccessViolation
}

func toTFTPError(err error) *tFTPError {
	if os.IsNotExist(err) {
		return ErrFileNotFound
	}
	if os.IsPermission(err) {
		return ErrAccessViolation
	}
	return ErrNotDefined
}
