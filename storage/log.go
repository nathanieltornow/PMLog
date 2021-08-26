package storage

type Log interface {
	// Append stores record and returns the local-sequence-number lsn, which will be used to find the record
	// on an Commit request
	Append(record string) (lsn uint64, err error)

	// Commit commits the records which is stored on local-sequence-number lsn with the global-sequence-number gsn on
	// on the log of color. The record can be committed for multiple colors.
	Commit(lsn uint64, color uint32, gsn uint64) (err error)

	// Read returns the record stored on global-sequence-number gsn on the log of color
	Read(color uint32, gsn uint64) (record string, err error)

	// Trim deletes all records of the log of color before global-sequence-number gsn
	Trim(color uint32, gsn uint64) (err error)
}
