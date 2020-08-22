package tangle

import (
	"time"

	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/objectstorage"

	"github.com/Ariwonto/aingle-alpha/pkg/model/aingle"
	"github.com/Ariwonto/aingle-alpha/pkg/profile"
)

var tagsStorage *objectstorage.ObjectStorage

type CachedTag struct {
	objectstorage.CachedObject
}

type CachedTags []*CachedTag

// tag -1
func (cachedTags CachedTags) Release(force ...bool) {
	for _, cachedTag := range cachedTags {
		cachedTag.Release(force...)
	}
}

func (c *CachedTag) GetTag() *aingle.Tag {
	return c.Get().(*aingle.Tag)
}

func tagsFactory(key []byte) (objectstorage.StorableObject, int, error) {
	tag := aingle.NewTag(key[:17], key[17:66])
	return tag, 66, nil
}

func GetTagsStorageSize() int {
	return tagsStorage.GetSize()
}

func configureTagsStorage(store kvstore.KVStore, opts profile.CacheOpts) {

	tagsStorage = objectstorage.New(
		store.WithRealm([]byte{StorePrefixTags}),
		tagsFactory,
		objectstorage.CacheTime(time.Duration(opts.CacheTimeMs)*time.Millisecond),
		objectstorage.PersistenceEnabled(true),
		objectstorage.PartitionKey(17, 49),
		objectstorage.KeysOnly(true),
		objectstorage.StoreOnCreation(true),
		objectstorage.LeakDetectionEnabled(opts.LeakDetectionOptions.Enabled,
			objectstorage.LeakDetectionOptions{
				MaxConsumersPerObject: opts.LeakDetectionOptions.MaxConsumersPerObject,
				MaxConsumerHoldTime:   time.Duration(opts.LeakDetectionOptions.MaxConsumerHoldTimeSec) * time.Second,
			}),
	)
}

// tag +-0
func GetTagHashes(txTag aingle.Hash, forceRelease bool, maxFind ...int) aingle.Hashes {
	var tagHashes aingle.Hashes

	i := 0
	tagsStorage.ForEachKeyOnly(func(key []byte) bool {
		i++
		if (len(maxFind) > 0) && (i > maxFind[0]) {
			return false
		}

		tagHashes = append(tagHashes, aingle.Hash(key[17:66]))
		return true
	}, false, txTag)

	return tagHashes
}

// TagConsumer consumes the given tag during looping through all tags in the persistence layer.
type TagConsumer func(txTag aingle.Hash, txHash aingle.Hash) bool

// ForEachTag loops over all tags.
func ForEachTag(consumer TagConsumer, skipCache bool) {
	tagsStorage.ForEachKeyOnly(func(key []byte) bool {
		return consumer(key[:17], key[17:66])
	}, skipCache)
}

// tag +1
func StoreTag(txTag aingle.Hash, txHash aingle.Hash) *CachedTag {
	tag := aingle.NewTag(txTag[:17], txHash[:49])
	return &CachedTag{CachedObject: tagsStorage.Store(tag)}
}

// tag +-0
func DeleteTag(txTag aingle.Hash, txHash aingle.Hash) {
	tagsStorage.Delete(append(txTag[:17], txHash[:49]...))
}

func ShutdownTagsStorage() {
	tagsStorage.Shutdown()
}

func FlushTagsStorage() {
	tagsStorage.Flush()
}
