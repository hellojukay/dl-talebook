package internal

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/bits-and-blooms/bitset"
)

const NoBookToDownload = -1

var (
	ErrStartBookID       = errors.New("the start book id should start from 1")
	ErrStartAndEndBookID = errors.New("start book id should below the available book id")
	ErrStorageFile       = errors.New("couldn't create file for storing download process")
)

// storage is a bit based
type storage struct {
	progress *bitset.BitSet // progress is used for file storage.
	assigned *bitset.BitSet // the assign status, memory based.
	lock     *sync.Mutex    // lock is used for concurrent request.
	file     string         // The storage file path for download progress.
}

// NewStorage Create a storge for save the download progress.
func NewStorage(start, size int64, path string) (*storage, error) {
	if start < 1 {
		return nil, ErrStartBookID
	}
	if start > size {
		return nil, ErrStartAndEndBookID
	}

	var progress *bitset.BitSet
	var file *os.File
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	startIndex := func(set *bitset.BitSet) {
		for i := uint(0); i < uint(start-1); i++ {
			set.Set(i)
		}
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// Create storage related file
		if file, err = os.Create(path); err != nil {
			return nil, ErrStorageFile
		} else {
			// Enrich file content.
			progress = bitset.New(uint(size))
			startIndex(progress)

			if err := saveStorage(file, progress); err != nil {
				return nil, err
			}
		}
	} else {
		// Load storage from file.
		if file, err = os.OpenFile(path, os.O_RDWR, 0644); err != nil {
			return nil, err
		}
		if progress, err = loadStorage(file); err != nil {
			return nil, err
		}

		// Recalculate start index.
		startIndex(progress)
	}

	assigned := bitset.New(progress.Len())
	progress.Copy(assigned)

	return &storage{
		progress: progress,
		assigned: assigned,
		lock:     new(sync.Mutex),
		file:     path,
	}, nil
}

func saveStorage(file *os.File, progress *bitset.BitSet) (err error) {
	_, err = progress.WriteTo(file)
	return
}

func loadStorage(file *os.File) (*bitset.BitSet, error) {
	set := new(bitset.BitSet)
	if _, err := set.ReadFrom(file); err != nil {
		return nil, err
	}

	return set, nil
}

// AcquireBookID would find the book id from assign array.
func (storage *storage) AcquireBookID() int64 {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	for i := uint(0); i < storage.assigned.Len(); i++ {
		if !storage.assigned.Test(i) {
			storage.assigned.Set(i)
			return int64(i + 1)
		}
	}

	return NoBookToDownload
}

// SaveBookID would save the download progress.
func (storage *storage) SaveBookID(bookID int64) error {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	if bookID >= int64(storage.progress.Len()) {
		return errors.New(fmt.Sprintf("invalid book id: %d", bookID))
	}

	i := uint(bookID - 1)
	storage.assigned.Set(i)
	storage.progress.Set(i)

	if file, err := os.OpenFile(storage.file, os.O_RDWR, 0644); err != nil {
		return err
	} else {
		defer func() { _ = file.Close() }()
		return saveStorage(file, storage.progress)
	}
}

// Finished would tell the called whether all the books have downloaded.
func (storage *storage) Finished() bool {
	return storage.progress.Count() == storage.progress.Len()
}

// Size would return the book size.
func (storage *storage) Size() uint {
	return storage.progress.Len()
}
