package usecase

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"golang.org/x/sync/semaphore"

	"go-poc/service/saleschannel/model"
	"go-poc/service/saleschannel/repository/port"
	"go-poc/utils"
	"go-poc/utils/log"
)

const (
	CacheChannel = "cache_channel"
)

type Channel interface {
	Upsert(ctx context.Context, inputs []model.ChannelInput) (outputs []model.ChannelOutput, err error)
	UpsertBatchFetching(ctx context.Context, inputs []model.ChannelInput) (outputs []model.ChannelOutput, err error)
	UpsertWithTransaction(ctx context.Context, inputs []model.ChannelInput) (outputs []model.ChannelOutput, err error)
	UpsertWithLock(ctx context.Context, inputs []model.ChannelInput) (outputs []model.ChannelOutput, err error)
	Delete(filter model.ChannelFilter) error
	FindByID(ID uuid.UUID) (*model.Channel, error)
	FindByFilter(filter model.ChannelFilter) ([]*model.Channel, error)
	FindPage(filter model.ChannelFilter, page, limit int64) (utils.Pagination, error)
}

type service struct {
	main  port.MainRepository
	cache port.CacheRepository
}

func NewChannel(
	main port.MainRepository,
	cache port.CacheRepository,
) Channel {
	return &service{
		main:  main,
		cache: cache,
	}
}

func (s *service) Upsert(ctx context.Context, inputs []model.ChannelInput) (outputs []model.ChannelOutput, err error) {
	channelRepository := s.main.Channel()
	upsertChannelWorker := 5
	if os.Getenv("UPDATE_CHANNEL_WORKER") != "" {
		upsertChannelWorkerEnv, err := strconv.Atoi(os.Getenv("UPDATE_CHANNEL_WORKER"))
		if err == nil {
			upsertChannelWorker = upsertChannelWorkerEnv
		}
	}

	outputChan := make(chan model.ChannelOutput, len(inputs))
	workerSemaphore := semaphore.NewWeighted(int64(upsertChannelWorker))

	for _, inputData := range inputs {
		err := workerSemaphore.Acquire(ctx, 1)
		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "acquire worker error"))
			continue
		}

		go func(inputDataInWorker model.ChannelInput) {
			defer workerSemaphore.Release(1)
			channelData, err := channelRepository.FindByID(inputDataInWorker.ID)
			if err == nil {
				channelData.Update(inputDataInWorker)
				err := channelRepository.Update(channelData)
				if err != nil {
					output := model.ChannelOutput{
						ID:      channelData.ID,
						Code:    channelData.Code,
						Message: stacktrace.RootCause(err).Error(),
					}

					outputChan <- output
					return
				}
				err = s.cache.Channel().Set(channelData)
				if err != nil {
					log.WithContext(ctx).Error(stacktrace.Propagate(err, "cache error"))
				}
			} else if stacktrace.RootCause(err) == sql.ErrNoRows {
				channelData := model.NewChannel(inputDataInWorker)
				err := channelRepository.Create(channelData)
				if err != nil {
					output := model.ChannelOutput{
						ID:      channelData.ID,
						Code:    channelData.Code,
						Message: stacktrace.RootCause(err).Error(),
					}

					outputChan <- output
					return
				}
				err = s.cache.Channel().Set(channelData)
				if err != nil {
					log.WithContext(ctx).Error(stacktrace.Propagate(err, "cache error"))
				}
			} else {
				output := model.ChannelOutput{
					ID:      channelData.ID,
					Code:    channelData.Code,
					Message: stacktrace.RootCause(err).Error(),
				}

				outputChan <- output
				return
			}
		}(inputData)
	}

	if err := workerSemaphore.Acquire(ctx, int64(upsertChannelWorker)); err != nil {
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

func (s *service) UpsertBatchFetching(ctx context.Context, inputs []model.ChannelInput) (outputs []model.ChannelOutput, err error) {
	channelRepository := s.main.Channel()
	ids := []uuid.UUID{}

	for _, input := range inputs {
		ids = append(ids, input.ID)
	}

	upsertChannelWorker := 5
	if os.Getenv("UPDATE_CHANNEL_WORKER") != "" {
		upsertChannelWorkerEnv, err := strconv.Atoi(os.Getenv("UPDATE_CHANNEL_WORKER"))
		if err == nil {
			upsertChannelWorker = upsertChannelWorkerEnv
		}
	}

	outputChan := make(chan model.ChannelOutput, len(inputs))
	workerSemaphore := semaphore.NewWeighted(int64(upsertChannelWorker))

	channels := []*model.Channel{}
	if len(ids) > 0 {
		filter := model.ChannelFilter{
			IDs: ids,
		}

		channels, err = channelRepository.FindByFilter(filter, true)
		if err != nil {
			return nil, stacktrace.Propagate(err, "find channel by filter error")
		}
	}

	channelMap := make(map[uuid.UUID]model.Channel)
	for _, channelData := range channels {
		channelMap[channelData.ID] = *channelData
	}

	for _, inputData := range inputs {
		err := workerSemaphore.Acquire(ctx, 1)
		if err != nil {
			log.WithContext(ctx).Error(stacktrace.Propagate(err, "acquire worker error"))
			continue
		}

		go func(inputDataInWorker model.ChannelInput) {
			defer workerSemaphore.Release(1)
			if channelData, exist := channelMap[inputDataInWorker.ID]; exist {
				channelData.Update(inputDataInWorker)
				err := channelRepository.Update(&channelData)
				if err != nil {
					output := model.ChannelOutput{
						ID:      channelData.ID,
						Code:    channelData.Code,
						Message: stacktrace.RootCause(err).Error(),
					}

					outputChan <- output
					return
				}
				err = s.cache.Channel().Set(&channelData)
				if err != nil {
					log.WithContext(ctx).Error(stacktrace.Propagate(err, "cache error"))
				}
			} else {
				channelData := model.NewChannel(inputDataInWorker)
				err := channelRepository.Create(channelData)
				if err != nil {
					output := model.ChannelOutput{
						ID:      channelData.ID,
						Code:    channelData.Code,
						Message: stacktrace.RootCause(err).Error(),
					}

					outputChan <- output
					return
				}
				err = s.cache.Channel().Set(channelData)
				if err != nil {
					log.WithContext(ctx).Error(stacktrace.Propagate(err, "cache error"))
				}
			}
		}(inputData)
	}

	if err := workerSemaphore.Acquire(ctx, int64(upsertChannelWorker)); err != nil {
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

func (s *service) UpsertWithTransaction(ctx context.Context, inputs []model.ChannelInput) (outputs []model.ChannelOutput, err error) {
	var t = func(repoRegistry port.MainRepository) (interface{}, error) {
		channelRepository := repoRegistry.Channel()
		ids := []uuid.UUID{}

		for _, input := range inputs {
			ids = append(ids, input.ID)
		}

		channels := []*model.Channel{}
		if len(ids) > 0 {
			filter := model.ChannelFilter{
				IDs: ids,
			}

			channels, err = channelRepository.FindByFilter(filter, true)
			if err != nil {
				return nil, stacktrace.Propagate(err, "find channel by filter error")
			}
		}

		channelMap := make(map[uuid.UUID]model.Channel)
		for _, channelData := range channels {
			channelMap[channelData.ID] = *channelData
		}

		upsertChannelWorker := 5
		if os.Getenv("UPDATE_CHANNEL_WORKER") != "" {
			upsertChannelWorkerEnv, err := strconv.Atoi(os.Getenv("UPDATE_CHANNEL_WORKER"))
			if err == nil {
				upsertChannelWorker = upsertChannelWorkerEnv
			}
		}

		outputChan := make(chan model.ChannelOutput, len(inputs))
		workerSemaphore := semaphore.NewWeighted(int64(upsertChannelWorker))
		for _, inputData := range inputs {
			err := workerSemaphore.Acquire(ctx, 1)
			if err != nil {
				log.WithContext(ctx).Error(stacktrace.Propagate(err, "acquire worker error"))
				continue
			}

			go func(inputDataInWorker model.ChannelInput) {
				defer workerSemaphore.Release(1)
				if channelData, exist := channelMap[inputDataInWorker.ID]; exist {
					channelData.Update(inputDataInWorker)
					err := channelRepository.Update(&channelData)
					if err != nil {
						output := model.ChannelOutput{
							ID:      channelData.ID,
							Code:    channelData.Code,
							Message: stacktrace.RootCause(err).Error(),
						}

						outputChan <- output
						return
					}
					go s.cache.Channel().Set(&channelData)
				} else {
					channelData := model.NewChannel(inputDataInWorker)
					err := channelRepository.Create(channelData)
					if err != nil {
						output := model.ChannelOutput{
							ID:      channelData.ID,
							Code:    channelData.Code,
							Message: stacktrace.RootCause(err).Error(),
						}

						outputChan <- output
						return
					}
					go s.cache.Channel().Set(channelData)
				}
			}(inputData)
		}

		if err := workerSemaphore.Acquire(ctx, int64(upsertChannelWorker)); err != nil {
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
			res := out.([]model.ChannelOutput)
			return res, err
		}

		return nil, err
	}

	return nil, nil
}

func (s *service) UpsertWithLock(ctx context.Context, inputs []model.ChannelInput) (outputs []model.ChannelOutput, err error) {
	var t = func(repoRegistry port.MainRepository) (interface{}, error) {
		channelRepository := repoRegistry.Channel()
		ids := []uuid.UUID{}

		for _, input := range inputs {
			ids = append(ids, input.ID)
		}

		channels := []*model.Channel{}
		if len(ids) > 0 {
			filter := model.ChannelFilter{
				IDs: ids,
			}

			channels, err = channelRepository.FindByFilter(filter, true)
			if err != nil {
				return nil, stacktrace.Propagate(err, "find channel by filter error")
			}
		}

		channelMap := make(map[uuid.UUID]model.Channel)
		for _, channelData := range channels {
			channelMap[channelData.ID] = *channelData
		}

		upsertChannelWorker := 5
		if os.Getenv("UPDATE_CHANNEL_WORKER") != "" {
			upsertChannelWorkerEnv, err := strconv.Atoi(os.Getenv("UPDATE_CHANNEL_WORKER"))
			if err == nil {
				upsertChannelWorker = upsertChannelWorkerEnv
			}
		}

		outputChan := make(chan model.ChannelOutput, len(inputs))
		workerSemaphore := semaphore.NewWeighted(int64(upsertChannelWorker))
		for _, inputData := range inputs {
			err := workerSemaphore.Acquire(ctx, 1)
			if err != nil {
				log.WithContext(ctx).Error(stacktrace.Propagate(err, "acquire worker error"))
				continue
			}

			go func(inputDataInWorker model.ChannelInput) {
				defer workerSemaphore.Release(1)
				// wait for concurrent testing
				time.Sleep(5 * time.Second)
				if channelData, exist := channelMap[inputDataInWorker.ID]; exist {
					channelData.Update(inputDataInWorker)
					err := channelRepository.Update(&channelData)
					if err != nil {
						output := model.ChannelOutput{
							ID:      channelData.ID,
							Code:    channelData.Code,
							Message: stacktrace.RootCause(err).Error(),
						}

						outputChan <- output
						return
					}
					go s.cache.Channel().Set(&channelData)
				} else {
					channelData := model.NewChannel(inputDataInWorker)
					err := channelRepository.Create(channelData)
					if err != nil {
						output := model.ChannelOutput{
							ID:      channelData.ID,
							Code:    channelData.Code,
							Message: stacktrace.RootCause(err).Error(),
						}

						outputChan <- output
						return
					}
					go s.cache.Channel().Set(channelData)
				}
			}(inputData)
		}

		if err := workerSemaphore.Acquire(ctx, int64(upsertChannelWorker)); err != nil {
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
			res := out.([]model.ChannelOutput)
			return res, err
		}

		return nil, err
	}

	return nil, nil
}

func (s *service) Delete(filter model.ChannelFilter) error {
	channelRepository := s.main.Channel()

	if err := channelRepository.Delete(filter); err != nil {
		return stacktrace.Propagate(err, "delete channel error")
	}

	return nil
}

func (s *service) FindByID(id uuid.UUID) (*model.Channel, error) {
	channelDataCache, err := s.cache.Channel().Get(id)
	if err == nil {
		return channelDataCache, nil
	}

	channelRepository := s.main.Channel()
	channelData, err := channelRepository.FindByID(id)
	if err != nil {
		return nil, stacktrace.Propagate(err, "find channel by id error")
	}

	go s.cache.Channel().Set(channelData)

	return channelData, nil
}

func (s *service) FindByFilter(filter model.ChannelFilter) ([]*model.Channel, error) {
	channelRepository := s.main.Channel()
	results, err := channelRepository.FindByFilter(filter, false)
	if err != nil {
		return []*model.Channel{}, stacktrace.Propagate(err, "find channel by filter error")
	}

	return results, nil
}

func (s *service) FindPage(filter model.ChannelFilter, page, limit int64) (utils.Pagination, error) {
	channelRepository := s.main.Channel()
	paginateEmpty := utils.PaginateEmpty()

	data, err := channelRepository.FindPage(filter, utils.GetOffset(page, limit), limit)
	if err != nil {
		return paginateEmpty, stacktrace.Propagate(err, "find channel page error")
	}

	total, err := channelRepository.FindTotalByFilter(filter)
	if err != nil {
		return paginateEmpty, stacktrace.Propagate(err, "find total channel by filter error")
	}

	return utils.PaginatePageLimit(data, total, page, limit), nil
}
