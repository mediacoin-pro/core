package goldb

import (
	"compress/flate"
	"os"

	"github.com/mediacoin-pro/core/common/bin"
)

type DumpOptions struct {
	Filter           *Query
	CompressionLevel int // 0-9 (0 - no compression; 9 - by default)
}

func (s *Storage) Dump(filepath string, options *DumpOptions) (err error) {
	op := DumpOptions{
		CompressionLevel: flate.DefaultCompression,
	}
	if options != nil {
		op = *options
	}

	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	fl, err := flate.NewWriter(file, op.CompressionLevel)
	if err != nil {
		return
	}
	defer fl.Close()
	w := bin.NewWriter(fl)

	q := op.Filter
	if q == nil {
		q = NewQuery(0)
		q.filter = nil // all keys
	}

	const SyncBatchSize = 32 * 1024 * 1024 // 32 MiB
	var nextSyncVol = int64(SyncBatchSize)
	err = s.Fetch(q, func(rec Record) error {
		w.WriteVar(rec.Key)
		w.WriteVar(rec.Value)
		if w.CntWritten > nextSyncVol {
			nextSyncVol += SyncBatchSize
			if err := file.Sync(); err != nil {
				return err
			}
		}
		return w.Error()
	})

	w.WriteVar([]byte{}) // EOF - null-key
	return
}

func (s *Storage) Restore(filepath string) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	fl := flate.NewReader(file)
	defer fl.Close()
	r := bin.NewReader(fl)

	s.mx.Lock()
	defer s.mx.Unlock()

	tr, err := s.db.OpenTransaction()
	if err != nil {
		return
	}
	var key, val []byte
	for i := 0; ; {
		if key, err = r.ReadBytes(); err != nil {
			return
		} else if len(key) == 0 { // EOF
			break
		}
		if val, err = r.ReadBytes(); err != nil {
			return
		}
		tr.Put(key, val, nil)

		if i++; i%10000 == 0 {
			if err = tr.Commit(); err != nil {
				return
			} else if tr, err = s.db.OpenTransaction(); err != nil {
				return
			}
		}
	}
	err = tr.Commit()
	return
}
