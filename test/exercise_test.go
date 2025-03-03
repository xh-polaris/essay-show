package test

import (
	"fmt"
	"github.com/xh-polaris/essay-show/biz/application/dto/essay/show"
	"github.com/xh-polaris/essay-show/biz/application/service"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/exercise"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/log"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/user"
	"golang.org/x/net/context"
	"testing"
)

func Test(t *testing.T) {
	c := config.GetConfig()
	s := service.ExerciseService{
		ExerciseMapper: exercise.NewMongoMapper(c),
		LogMapper:      log.NewMongoMapper(c),
		UserMapper:     user.NewMongoMapper(c),
	}
	ctx := context.Context(context.Background())
	e, err := s.CreateExercise(ctx, &show.CreateExerciseReq{
		LogId: "675a8f5fc6c523a81d50484e",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(e)

}
