package tangle

import (
	"encoding/binary"
	"time"

	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/objectstorage"

	"github.com/Ariwonto/aingle-alpha/pkg/model/aingle"
	"github.com/Ariwonto/aingle-alpha/pkg/model/milestone"
	"github.com/Ariwonto/aingle-alpha/pkg/profile"
)

var unconfirmedTxStorage *objectstorage.ObjectStorage

type CachedUnconfirmedTx struct {
	objectstorage.CachedObject
}

type CachedUnconfirmedTxs []*CachedUnconfirmedTx

func (cachedUnconfirmedTxs CachedUnconfirmedTxs) Release(force ...bool) {
	for _, cachedUnconfirmedTx := range cachedUnconfirmedTxs {
		cachedUnconfirmedTx.Release(force...)
	}
}

func (c *CachedUnconfirmedTx) GetUnconfirmedTx() *aingle.UnconfirmedTx {
	return c.Get().(*aingle.UnconfirmedTx)
}

func unconfirmedTxFactory(key []byte) (objectstorage.StorableObject, int, error) {

	unconfirmedTx := aingle.NewUnconfirmedTx(milestone.Index(binary.LittleEndian.Uint32(key[:4])), key[4:53])
	return unconfirmedTx, 53, nil
}

func GetUnconfirmedTxStorageSize() int {
	return unconfirmedTxStorage.GetSize()
}

func configureUnconfirmedTxStorage(store kvstore.KVStore, opts profile.CacheOpts) {

	unconfirmedTxStorage = objectstorage.New(
		store.WithRealm([]byte{StorePrefixUnconfirmedTransactions}),
		unconfirmedTxFactory,
		objectstorage.CacheTime(time.Duration(opts.CacheTimeMs)*time.Millisecond),
		objectstorage.PersistenceEnabled(true),
		objectstorage.PartitionKey(4, 49),
		objectstorage.KeysOnly(true),
		objectstorage.StoreOnCreation(true),
		objectstorage.LeakDetectionEnabled(opts.LeakDetectionOptions.Enabled,
			objectstorage.LeakDetectionOptions{
				MaxConsumersPerObject: opts.LeakDetectionOptions.MaxConsumersPerObject,
				MaxConsumerHoldTime:   time.Duration(opts.LeakDetectionOptions.MaxConsumerHoldTimeSec) * time.Second,
			}),
	)
}

// GetUnconfirmedTxHashes returns all hashes of unconfirmed transactions for that milestone.
func GetUnconfirmedTxHashes(msIndex milestone.Index, forceRelease bool) aingle.Hashes {

	var unconfirmedTxHashes aingle.Hashes

	key := make([]byte, 4)
	binary.LittleEndian.PutUint32(key, uint32(msIndex))

	unconfirmedTxStorage.ForEachKeyOnly(func(key []byte) bool {
		unconfirmedTxHashes = append(unconfirmedTxHashes, aingle.Hash(key[4:53]))
		return true
	}, false, key)

	return unconfirmedTxHashes
}

// UnconfirmedTxConsumer consumes the given unconfirmed transaction during looping through all unconfirmed transactions in the persistence layer.
type UnconfirmedTxConsumer func(msIndex milestone.Index, txHash aingle.Hash) bool

// ForEachUnconfirmedTx loops over all unconfirmed transactions.
func ForEachUnconfirmedTx(consumer UnconfirmedTxConsumer, skipCache bool) {
	unconfirmedTxStorage.ForEachKeyOnly(func(key []byte) bool {
		return consumer(milestone.Index(binary.LittleEndian.Uint32(key[:4])), key[4:53])
	}, skipCache)
}

// unconfirmedTx +1
func StoreUnconfirmedTx(msIndex milestone.Index, txHash aingle.Hash) *CachedUnconfirmedTx {
	unconfirmedTx := aingle.NewUnconfirmedTx(msIndex, txHash)
	return &CachedUnconfirmedTx{CachedObject: unconfirmedTxStorage.Store(unconfirmedTx)}
}

// DeleteUnconfirmedTxs deletes unconfirmed transaction entries.
func DeleteUnconfirmedTxs(msIndex milestone.Index) int {

	msIndexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(msIndexBytes, uint32(msIndex))

	var keysToDelete [][]byte

	unconfirmedTxStorage.ForEachKeyOnly(func(key []byte) bool {
		keysToDelete = append(keysToDelete, key)
		return true
	}, false, msIndexBytes)

	for _, key := range keysToDelete {
		unconfirmedTxStorage.Delete(key)
	}

	return len(keysToDelete)
}

func ShutdownUnconfirmedTxsStorage() {
	unconfirmedTxStorage.Shutdown()
}

func FlushUnconfirmedTxsStorage() {
	unconfirmedTxStorage.Flush()
}
