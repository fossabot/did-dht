// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package postgres

import (
	"context"
)

const failedRecordCount = `-- name: FailedRecordCount :one
SELECT count(*) AS exact_count FROM failed_records
`

func (q *Queries) FailedRecordCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, failedRecordCount)
	var exact_count int64
	err := row.Scan(&exact_count)
	return exact_count, err
}

const listFailedRecords = `-- name: ListFailedRecords :many
SELECT id, failure_count FROM failed_records
`

func (q *Queries) ListFailedRecords(ctx context.Context) ([]FailedRecord, error) {
	rows, err := q.db.Query(ctx, listFailedRecords)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FailedRecord
	for rows.Next() {
		var i FailedRecord
		if err := rows.Scan(&i.ID, &i.FailureCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRecords = `-- name: ListRecords :many
SELECT id, key, value, sig, seq FROM dht_records WHERE id > (SELECT id FROM dht_records WHERE dht_records.key = $1) ORDER BY id ASC LIMIT $2
`

type ListRecordsParams struct {
	Key   []byte
	Limit int32
}

func (q *Queries) ListRecords(ctx context.Context, arg ListRecordsParams) ([]DhtRecord, error) {
	rows, err := q.db.Query(ctx, listRecords, arg.Key, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DhtRecord
	for rows.Next() {
		var i DhtRecord
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Value,
			&i.Sig,
			&i.Seq,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRecordsFirstPage = `-- name: ListRecordsFirstPage :many
SELECT id, key, value, sig, seq FROM dht_records ORDER BY id ASC LIMIT $1
`

func (q *Queries) ListRecordsFirstPage(ctx context.Context, limit int32) ([]DhtRecord, error) {
	rows, err := q.db.Query(ctx, listRecordsFirstPage, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DhtRecord
	for rows.Next() {
		var i DhtRecord
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Value,
			&i.Sig,
			&i.Seq,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readRecord = `-- name: ReadRecord :one
SELECT id, key, value, sig, seq FROM dht_records WHERE key = $1 LIMIT 1
`

func (q *Queries) ReadRecord(ctx context.Context, key []byte) (DhtRecord, error) {
	row := q.db.QueryRow(ctx, readRecord, key)
	var i DhtRecord
	err := row.Scan(
		&i.ID,
		&i.Key,
		&i.Value,
		&i.Sig,
		&i.Seq,
	)
	return i, err
}

const recordCount = `-- name: RecordCount :one
SELECT count(*) AS exact_count FROM dht_records
`

func (q *Queries) RecordCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, recordCount)
	var exact_count int64
	err := row.Scan(&exact_count)
	return exact_count, err
}

const writeFailedRecord = `-- name: WriteFailedRecord :exec
INSERT INTO failed_records(id, failure_count)
VALUES($1, $2)
ON CONFLICT (id) DO UPDATE SET failure_count = failed_records.failure_count + 1
`

type WriteFailedRecordParams struct {
	ID           []byte
	FailureCount int32
}

func (q *Queries) WriteFailedRecord(ctx context.Context, arg WriteFailedRecordParams) error {
	_, err := q.db.Exec(ctx, writeFailedRecord, arg.ID, arg.FailureCount)
	return err
}

const writeRecord = `-- name: WriteRecord :exec
INSERT INTO dht_records(key, value, sig, seq) VALUES($1, $2, $3, $4)
`

type WriteRecordParams struct {
	Key   []byte
	Value []byte
	Sig   []byte
	Seq   int64
}

func (q *Queries) WriteRecord(ctx context.Context, arg WriteRecordParams) error {
	_, err := q.db.Exec(ctx, writeRecord,
		arg.Key,
		arg.Value,
		arg.Sig,
		arg.Seq,
	)
	return err
}
