package storage

type Log interface {
	Append(color uint32, record string) (lsn uint64, err error)

	Commit(color uint32, lsn uint64, gsn uint64) (err error)

	Read(color uint32, gsn uint64) (record string, err error)

	Trim(color uint32, gsn uint64) (err error)
}
