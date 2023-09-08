package port

type InTransaction func(repoRegistry MainRepository) (interface{}, error)

type MainRepository interface {
	Channel() ChannelMainRepository
	DoInTransaction(txFunc InTransaction) (out interface{}, err error)
}

type CacheRepository interface {
	Channel() ChannelCacheRepository
}
