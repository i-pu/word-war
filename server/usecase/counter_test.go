package usecase

import (
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/i-pu/word-war/server/domain/entity"
	"github.com/i-pu/word-war/server/domain/repository/mock"
	"github.com/i-pu/word-war/server/domain/service"
)

func TestCounterUsecase_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_repository.NewMockCounterRepository(ctrl)
	m.EXPECT().SetCounter(gomock.Any()).Return(nil)
	u := NewCounterUsecase(m, service.NewCounterService(m))
	c, err := u.Init(&entity.Counter{Value: 0})
	if err != nil {
		t.Errorf("%v", err)
	}
	if !(c.Value == 0) {
		t.Errorf("wow")
	}
}

func TestCounterUsecase_Incr(t *testing.T) {
	// m.EXPECT().IncrCounter().DoAndReturn(func() (int64, error) {
	// 	return 5, nil
	// })
}

func TestCounterUsecase_Get(t *testing.T) {
	// m.EXPECT().GetCounter().DoAndReturn(func() (int64, error) {
	// 	return 10, nil
	// })
}

