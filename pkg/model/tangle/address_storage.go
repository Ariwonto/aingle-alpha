package tangle

import (
	"time"

	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/objectstorage"

	"github.com/Ariwonto/aingle-alpha/pkg/model/aingle"
	"github.com/Ariwonto/aingle-alpha/pkg/profile"
)

var addressesStorage *objectstorage.ObjectStorage

type CachedAddress struct {
	objectstorage.CachedObject
}

type CachedAddresses []*CachedAddress

func (cachedAddresses CachedAddresses) Release(force ...bool) {
	for _, cachedAddress := range cachedAddresses {
		cachedAddress.Release(force...)
	}
}

func (c *CachedAddress) GetAddress() *aingle.Address {
	return c.Get().(*aingle.Address)
}

func databaseKeyPrefixForAddress(address aingle.Hash) []byte {
	return address
}

func databaseKeyPrefixForAddressTransaction(address aingle.Hash, txHash aingle.Hash, isValue bool) []byte {
	var isValueByte byte
	if isValue {
		isValueByte = aingle.AddressTxIsValue
	}

	result := append(databaseKeyPrefixForAddress(address), isValueByte)
	return append(result, txHash...)
}

func addressFactory(key []byte) (objectstorage.StorableObject, int, error) {
	address := aingle.NewAddress(key[:49], key[50:99], key[49] == aingle.AddressTxIsValue)
	return address, 99, nil
}

func GetAddressesStorageSize() int {
	return addressesStorage.GetSize()
}

func configureAddressesStorage(store kvstore.KVStore, opts profile.CacheOpts) {

	addressesStorage = objectstorage.New(
		store.WithRealm([]byte{StorePrefixAddresses}),
		addressFactory,
		objectstorage.CacheTime(time.Duration(opts.CacheTimeMs)*time.Millisecond),
		objectstorage.PersistenceEnabled(true),
		objectstorage.PartitionKey(49, 1, 49),
		objectstorage.KeysOnly(true),
		objectstorage.StoreOnCreation(true),
		objectstorage.LeakDetectionEnabled(opts.LeakDetectionOptions.Enabled,
			objectstorage.LeakDetectionOptions{
				MaxConsumersPerObject: opts.LeakDetectionOptions.MaxConsumersPerObject,
				MaxConsumerHoldTime:   time.Duration(opts.LeakDetectionOptions.MaxConsumerHoldTimeSec) * time.Second,
			}),
	)
}

// address +-0
func GetTransactionHashesForAddress(address aingle.Hash, valueOnly bool, forceRelease bool, maxFind ...int) aingle.Hashes {

	searchPrefix := databaseKeyPrefixForAddress(address)
	if valueOnly {
		var isValueByte byte = aingle.AddressTxIsValue
		searchPrefix = append(searchPrefix, isValueByte)
	}

	var txHashes aingle.Hashes

	i := 0
	addressesStorage.ForEachKeyOnly(func(key []byte) bool {
		i++
		if (len(maxFind) > 0) && (i > maxFind[0]) {
			return false
		}

		txHashes = append(txHashes, key[50:99])
		return true
	}, false, searchPrefix)

	return txHashes
}

// AddressConsumer consumes the given address during looping through all addresses in the persistence layer.
type AddressConsumer func(address aingle.Hash, txHash aingle.Hash, isValue bool) bool

// ForEachAddress loops over all addresses.
func ForEachAddress(consumer AddressConsumer, skipCache bool) {
	addressesStorage.ForEachKeyOnly(func(key []byte) bool {
		return consumer(key[:49], key[50:99], key[49] == aingle.AddressTxIsValue)
	}, skipCache)
}

// address +1
func StoreAddress(address aingle.Hash, txHash aingle.Hash, isValue bool) *CachedAddress {
	addressObj := aingle.NewAddress(address, txHash, isValue)
	return &CachedAddress{CachedObject: addressesStorage.Store(addressObj)}
}

// address +-0
func DeleteAddress(address aingle.Hash, txHash aingle.Hash) {
	addressesStorage.Delete(databaseKeyPrefixForAddressTransaction(address, txHash, false))
	addressesStorage.Delete(databaseKeyPrefixForAddressTransaction(address, txHash, true))
}

func ShutdownAddressStorage() {
	addressesStorage.Shutdown()
}

func FlushAddressStorage() {
	addressesStorage.Flush()
}
