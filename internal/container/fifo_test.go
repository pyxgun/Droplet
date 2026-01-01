package container

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFifoCreateAndRemove_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.fifo")

	fifoHandler := &containerFifoHandler{}

	// create fifo
	if err := fifoHandler.createFifo(path); err != nil {
		t.Fatalf("createFifo failed: %v", err)
	}
	// verify fifo exists and is a FIFO
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("fifo not found: %v", err)
	}
	if info.Mode()&os.ModeNamedPipe == 0 {
		t.Fatalf("expected named pipe, got mode %v", info.Mode())
	}

	// remove fifo
	if err := fifoHandler.removeFifo(path); err != nil {
		t.Fatalf("removeFifo failed: %v", err)
	}
	// verify fifo not exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("fifo still exists after remove")
	}
}

func TestFifoCreateAndRemove_CreateError(t *testing.T) {
	t.Parallel()

	fifoHandler := &containerFifoHandler{}

	if err := fifoHandler.createFifo("/noexists/path"); err == nil {
		t.Fatalf("expected error for missing fifo, got nil")
	}
}

func TestFifoCreateAndRemove_RemoveError(t *testing.T) {
	t.Parallel()

	fifoHandler := &containerFifoHandler{}

	if err := fifoHandler.removeFifo("/noexists/path"); err == nil {
		t.Fatalf("expected error for missing fifo, got nil")
	}
}

func TestFifoReadAndWrite_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.fifo")

	fifoHandler := &containerFifoHandler{}

	// prepare fifo
	if err := fifoHandler.createFifo(path); err != nil {
		t.Fatalf("createFifo failed: %v", err)
	}

	errCh := make(chan error, 1)

	// run reader in goroutine
	go func() {
		errCh <- fifoHandler.readFifo(path)
	}()

	// give reader time to block on open()
	time.Sleep(50 * time.Millisecond)

	// writer writes signal
	if err := fifoHandler.writeFifo(path); err != nil {
		t.Fatalf("writeFifo failed: %v", err)
	}
	// wait for reader result (with timeout)
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("readFifo returned error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("readFifo did not complete (likely blocked)")
	}
}

func TestFifoReadAndWrite_ReadError(t *testing.T) {
	t.Parallel()

	fifoHandler := &containerFifoHandler{}

	if err := fifoHandler.readFifo("/noexists/path"); err == nil {
		t.Fatalf("expected error for missing fifo, got nil")
	}
}

func TestFifoReadAndWrite_WriteError(t *testing.T) {
	t.Parallel()

	fifoHandler := &containerFifoHandler{}

	if err := fifoHandler.writeFifo("/noexists/path"); err == nil {
		t.Fatalf("expected error for missing fifo, got nil")
	}
}
