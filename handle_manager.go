package sftp

import (
	"strconv"
	"sync"
	"syscall"

	"github.com/spf13/afero"
)

// handleManager manages open file handles
type handleManager struct {
	handleCount int
	openFilesMu sync.RWMutex
	openFiles   map[string]afero.File
}

func newHandleManager() *handleManager {
	return &handleManager{
		openFiles: make(map[string]afero.File),
	}
}

func (hm *handleManager) GetHandle(handle string) (afero.File, bool) {
	hm.openFilesMu.RLock()
	defer hm.openFilesMu.RUnlock()
	f, ok := hm.openFiles[handle]
	return f, ok
}

func (hm *handleManager) NextHandle(f afero.File) string {
	hm.openFilesMu.Lock()
	defer hm.openFilesMu.Unlock()
	hm.handleCount++
	handle := strconv.Itoa(hm.handleCount)
	hm.openFiles[handle] = f
	return handle
}

func (hm *handleManager) CloseHandle(handle string) error {
	hm.openFilesMu.Lock()
	defer hm.openFilesMu.Unlock()
	if f, ok := hm.openFiles[handle]; ok {
		delete(hm.openFiles, handle)
		return f.Close()
	}
	return syscall.EBADF
}

func (hm *handleManager) OpenHandles() map[string]afero.File {
	hm.openFilesMu.RLock()
	defer hm.openFilesMu.RUnlock()
	handles := make(map[string]afero.File)
	for h, f := range hm.openFiles {
		handles[h] = f
	}
	return handles
}
