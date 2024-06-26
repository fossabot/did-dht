package postgres_test

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TBD54566975/did-dht/internal/did"
	"github.com/TBD54566975/did-dht/pkg/dht"
	"github.com/TBD54566975/did-dht/pkg/storage"
	"github.com/TBD54566975/did-dht/pkg/storage/db/postgres"
)

func getTestDB(t *testing.T) storage.Storage {
	uri := os.Getenv("TEST_DB")
	if uri == "" {
		t.SkipNow()
	}

	u, err := url.Parse(uri)
	require.NoError(t, err)
	if u.Scheme != "postgres" {
		t.SkipNow()
	}

	db, err := postgres.NewPostgres(uri)
	require.NoError(t, err)

	return db
}

func TestReadWrite(t *testing.T) {
	db := getTestDB(t)
	ctx := context.Background()

	beforeCnt, err := db.RecordCount(ctx)
	require.NoError(t, err)

	// create a did doc as a packet to store
	sk, doc, err := did.GenerateDIDDHT(did.CreateDIDDHTOpts{})
	require.NoError(t, err)
	require.NotEmpty(t, doc)

	packet, err := did.DHT(doc.ID).ToDNSPacket(*doc, nil, nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, packet)

	putMsg, err := dht.CreateDNSPublishRequest(sk, *packet)
	require.NoError(t, err)
	require.NotEmpty(t, putMsg)

	r := dht.RecordFromBEP44(putMsg)

	err = db.WriteRecord(ctx, r)
	require.NoError(t, err)

	r2, err := db.ReadRecord(ctx, r.ID())
	require.NoError(t, err)

	assert.Equal(t, r.Key, r2.Key)
	assert.Equal(t, r.Value, r2.Value)
	assert.Equal(t, r.Signature, r2.Signature)
	assert.Equal(t, r.SequenceNumber, r2.SequenceNumber)

	afterCnt, err := db.RecordCount(ctx)
	require.NoError(t, err)
	assert.Equal(t, beforeCnt+1, afterCnt)
}

func TestDBPagination(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()

	beforeCnt, err := db.RecordCount(ctx)
	require.NoError(t, err)

	preTestRecords, _, err := db.ListRecords(ctx, nil, 10)
	require.NoError(t, err)

	// store 10 records
	for i := 0; i < 10; i++ {
		// create a did doc as a packet to store
		sk, doc, err := did.GenerateDIDDHT(did.CreateDIDDHTOpts{})
		require.NoError(t, err)
		require.NotEmpty(t, doc)

		packet, err := did.DHT(doc.ID).ToDNSPacket(*doc, nil, nil, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, packet)

		putMsg, err := dht.CreateDNSPublishRequest(sk, *packet)
		require.NoError(t, err)
		require.NotEmpty(t, putMsg)

		// create record
		record := dht.RecordFromBEP44(putMsg)

		err = db.WriteRecord(ctx, record)
		assert.NoError(t, err)
	}

	// store 11th document
	// create a did doc as a packet to store
	sk, doc, err := did.GenerateDIDDHT(did.CreateDIDDHTOpts{})
	require.NoError(t, err)
	require.NotEmpty(t, doc)

	packet, err := did.DHT(doc.ID).ToDNSPacket(*doc, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, packet)

	putMsg, err := dht.CreateDNSPublishRequest(sk, *packet)
	require.NoError(t, err)
	require.NotEmpty(t, putMsg)

	// create eleventhRecord
	eleventhRecord := dht.RecordFromBEP44(putMsg)

	err = db.WriteRecord(ctx, eleventhRecord)
	assert.NoError(t, err)

	// read the first 10 back
	page, nextPageToken, err := db.ListRecords(ctx, nil, 10)
	assert.NoError(t, err)
	assert.Len(t, page, 10)

	page, nextPageToken, err = db.ListRecords(ctx, nextPageToken, 10+len(preTestRecords))
	assert.NoError(t, err)
	assert.Nil(t, nextPageToken)
	assert.Len(t, page, 1+len(preTestRecords))

	afterCnt, err := db.RecordCount(ctx)
	require.NoError(t, err)
	assert.Equal(t, beforeCnt+11, afterCnt)
}
