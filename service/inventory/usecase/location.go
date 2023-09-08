package usecase

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"golang.org/x/sync/semaphore"

	"go-poc/service/inventory/model"
	"go-poc/service/inventory/repository/port"
	"go-poc/utils"
	"go-poc/utils/log"
)

const (
	Cachelocation = "cache_location"
)

type Location interface {
	Upsert(ctx context.Context, inputs []model.LocationInput) (outputs []model.LocationOutput, err error)
	Delete(filter model.LocationFilter) error
	FindByID(ID uuid.UUID) (*model.Location, error)
	FindByFilter(filter model.LocationFilter) ([]*model.Location, error)
	FindPage(filter model.LocationFilter, page, limit int64) (utils.Pagination, error)
}

type service struct {
	main  port.MainRepository
	cache port.CacheRepository
}

func NewLocation(
	main port.MainRepository,
	cache port.CacheRepository,
) Location {
	return &service{
		main:  main,
		cache: cache,
	}
}

func (s *service) Upsert(ctx context.Context, inputs []model.LocationInput) (outputs []model.LocationOutput, err error) {
	var t = func(repoRegistry port.MainRepository) (interface{}, error) {
		locationRepository := repoRegistry.Location()
		ids := []uuid.UUID{}

		for _, input := range inputs {
			ids = append(ids, input.ID)
		}

		locations := []*model.Location{}
		if len(ids) > 0 {
			filter := model.LocationFilter{
				IDs: ids,
			}

			locations, err = locationRepository.FindByFilter(filter, true)
			if err != nil {
				return nil, stacktrace.Propagate(err, "find location by filter error")
			}
		}

		locationMap := make(map[uuid.UUID]model.Location)
		for _, locationData := range locations {
			locationMap[locationData.ID] = *locationData
		}

		upsertLocationWorker := 5
		if os.Getenv("UPDATE_LOCATION_WORKER") != "" {
			upsertLocationWorkerEnv, err := strconv.Atoi(os.Getenv("UPDATE_LOCATION_WORKER"))
			if err == nil {
				upsertLocationWorker = upsertLocationWorkerEnv
			}
		}

		outputChan := make(chan model.LocationOutput, len(inputs))
		workerSemaphore := semaphore.NewWeighted(int64(upsertLocationWorker))
		for _, inputData := range inputs {
			err := workerSemaphore.Acquire(ctx, 1)
			if err != nil {
				log.WithContext(ctx).Error(stacktrace.Propagate(err, "acquire worker error"))
				continue
			}

			go func(inputDataInWorker model.LocationInput) {
				defer workerSemaphore.Release(1)
				if locationData, exist := locationMap[inputDataInWorker.ID]; exist {
					locationData.Update(inputDataInWorker)
					err := locationRepository.Update(&locationData)
					if err != nil {
						output := model.LocationOutput{
							ID:      locationData.ID,
							Code:    locationData.Code,
							Message: stacktrace.RootCause(err).Error(),
						}

						outputChan <- output
						return
					}
					go s.cache.Location().Set(&locationData)
				} else {
					locationData := model.NewLocation(inputDataInWorker)
					err := locationRepository.Create(locationData)
					if err != nil {
						output := model.LocationOutput{
							ID:      locationData.ID,
							Code:    locationData.Code,
							Message: stacktrace.RootCause(err).Error(),
						}

						outputChan <- output
						return
					}
					go s.cache.Location().Set(locationData)
				}
			}(inputData)
		}

		if err := workerSemaphore.Acquire(ctx, int64(upsertLocationWorker)); err != nil {
			return nil, stacktrace.Propagate(err, "acquire worker error")
		}

		close(outputChan)
		for outputData := range outputChan {
			outputs = append(outputs, outputData)
		}

		if len(outputs) > 0 {
			return outputs, errors.New("internal server error")
		}

		return nil, nil
	}

	var out interface{}
	out, err = s.main.DoInTransaction(t)
	if err != nil {
		if out != nil {
			res := out.([]model.LocationOutput)
			return res, err
		}

		return nil, err
	}

	return nil, nil
}

func (s *service) Delete(filter model.LocationFilter) error {
	locationRepository := s.main.Location()

	if err := locationRepository.Delete(filter); err != nil {
		return stacktrace.Propagate(err, "delete location error")
	}

	return nil
}

func (s *service) FindByID(id uuid.UUID) (*model.Location, error) {
	locationDataCache, err := s.cache.Location().Get(id)
	if err == nil {
		return locationDataCache, nil
	}

	locationRepository := s.main.Location()
	locationData, err := locationRepository.FindByID(id)
	if err != nil {
		return nil, stacktrace.Propagate(err, "find location by id error")
	}

	go s.cache.Location().Set(locationData)

	return locationData, nil
}

func (s *service) FindByFilter(filter model.LocationFilter) ([]*model.Location, error) {
	locationRepository := s.main.Location()
	results, err := locationRepository.FindByFilter(filter, false)
	if err != nil {
		return []*model.Location{}, stacktrace.Propagate(err, "find location by filter error")
	}

	return results, nil
}

func (s *service) FindPage(filter model.LocationFilter, page, limit int64) (utils.Pagination, error) {
	locationRepository := s.main.Location()
	paginateEmpty := utils.PaginateEmpty()

	data, err := locationRepository.FindPage(filter, utils.GetOffset(page, limit), limit)
	if err != nil {
		return paginateEmpty, stacktrace.Propagate(err, "find location page error")
	}

	total, err := locationRepository.FindTotalByFilter(filter)
	if err != nil {
		return paginateEmpty, stacktrace.Propagate(err, "find total location by filter error")
	}

	return utils.PaginatePageLimit(data, total, page, limit), nil
}
