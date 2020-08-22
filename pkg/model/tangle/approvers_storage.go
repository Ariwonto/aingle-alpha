package tangle

import (
	"time"

	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/objectstorage"

	"github.com/Ariwonto/aingle-alpha/pkg/model/aingle"
	"github.com/Ariwonto/aingle-alpha/pkg/profile"
)

var approversStorage *objectstorage.ObjectStorage

type CachedApprover struct {
	objectstorage.CachedObject
}

type CachedAppprovers []*CachedApprover

func (cachedApprovers CachedAppprovers) Release(force ...bool) {
	for _, cachedApprover := range cachedApprovers {
		cachedApprover.Release(force...)
	}
}

func (c *CachedApprover) GetApprover() *aingle.Approver {
	return c.Get().(*aingle.Approver)
}

func approversFactory(key []byte) (objectstorage.StorableObject, int, error) {
	approver := aingle.NewApprover(key[:49], key[49:98])
	return approver, 98, nil
}

func GetApproversStorageSize() int {
	return approversStorage.GetSize()
}

func configureApproversStorage(store kvstore.KVStore, opts profile.CacheOpts) {

	approversStorage = objectstorage.New(
		store.WithRealm([]byte{StorePrefixApprovers}),
		approversFactory,
		objectstorage.CacheTime(time.Duration(opts.CacheTimeMs)*time.Millisecond),
		objectstorage.PersistenceEnabled(true),
		objectstorage.PartitionKey(49, 49),
		objectstorage.KeysOnly(true),
		objectstorage.StoreOnCreation(true),
		objectstorage.LeakDetectionEnabled(opts.LeakDetectionOptions.Enabled,
			objectstorage.LeakDetectionOptions{
				MaxConsumersPerObject: opts.LeakDetectionOptions.MaxConsumersPerObject,
				MaxConsumerHoldTime:   time.Duration(opts.LeakDetectionOptions.MaxConsumerHoldTimeSec) * time.Second,
			}),
	)
}

// approvers +-0
func GetApproverHashes(txHash aingle.Hash, maxFind ...int) aingle.Hashes {
	var approverHashes aingle.Hashes

	i := 0
	approversStorage.ForEachKeyOnly(func(key []byte) bool {
		i++
		if (len(maxFind) > 0) && (i > maxFind[0]) {
			return false
		}

		approverHashes = append(approverHashes, key[49:98])
		return true
	}, false, txHash)

	return approverHashes
}

// ApproverConsumer consumes the given approver during looping through all approvers in the persistence layer.
type ApproverConsumer func(txHash aingle.Hash, approverHash aingle.Hash) bool

// ForEachApprover loops over all approvers.
func ForEachApprover(consumer ApproverConsumer, skipCache bool) {
	approversStorage.ForEachKeyOnly(func(key []byte) bool {
		return consumer(key[:49], key[49:98])
	}, skipCache)
}

// approvers +1
func StoreApprover(txHash aingle.Hash, approverHash aingle.Hash) *CachedApprover {
	approver := aingle.NewApprover(txHash, approverHash)
	return &CachedApprover{CachedObject: approversStorage.Store(approver)}
}

// approvers +-0
func DeleteApprover(txHash aingle.Hash, approverHash aingle.Hash) {
	approver := aingle.NewApprover(txHash, approverHash)
	approversStorage.Delete(approver.ObjectStorageKey())
}

// approvers +-0
func DeleteApprovers(txHash aingle.Hash) {

	var keysToDelete [][]byte

	approversStorage.ForEachKeyOnly(func(key []byte) bool {
		keysToDelete = append(keysToDelete, key)
		return true
	}, false, txHash)

	for _, key := range keysToDelete {
		approversStorage.Delete(key)
	}
}

func ShutdownApproversStorage() {
	approversStorage.Shutdown()
}

func FlushApproversStorage() {
	approversStorage.Flush()
}
