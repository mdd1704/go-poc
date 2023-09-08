package port

type InTransaction func(repoRegistry MainRepository) (interface{}, error)

type MainRepository interface {
	Location() LocationMainRepository
	Sourcing() SourcingMainRepository
	DoInTransaction(txFunc InTransaction) (out interface{}, err error)
}

type CacheRepository interface {
	Location() LocationCacheRepository
	Sourcing() SourcingCacheRepository
}
