package exercise

import (
	"encoding/json"
	"fmt"
	coze "github.com/coze-dev/coze-go"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/exercise"
	"github.com/xh-polaris/essay-show/biz/infrastructure/mapper/log"
	logx "github.com/xh-polaris/essay-show/biz/infrastructure/util/log"
	"golang.org/x/net/context"
	"strings"
	"sync"
	"time"
)

var instance *coze.CozeAPI
var once sync.Once

func getCozeCli() *coze.CozeAPI {
	once.Do(func() {
		auth := coze.NewTokenAuth(config.GetConfig().Coze.Key)
		api := coze.NewCozeAPI(auth, coze.WithBaseURL(coze.CnBaseURL))
		instance = &api
	})
	return instance
}

func GenerateExercise(ctx context.Context, grade int64, l *log.Log) (*exercise.Exercise, error) {
	m, err := parseLog(l)
	if err != nil {
		return nil, err
	}
	resp, err := generate(ctx, grade, m, l)
	if err != nil {
		return nil, err
	}
	e, err := parseExercise(resp)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// 将log的Response转换为Json格式
func parseLog(l *log.Log) (map[string]any, error) {
	m := make(map[string]any)
	err := json.Unmarshal([]byte(l.Response), &m)
	return m, err
}

// 将map形式的resp解析成exercise
func parseExercise(resp map[string]any) (*exercise.Exercise, error) {
	// 选择题数组
	cqs := make([]*exercise.ChoiceQuestion, 0)

	// 题目数组
	questions := resp["result"].([]any)
	for _, question := range questions {
		q := question.(map[string]any)
		cq := &exercise.ChoiceQuestion{Options: make([]*exercise.Option, 0)}
		for k, v := range q {
			switch k {
			case "question":
				cq.Question = v.(string)
			case "explaion":
				fallthrough
			case "explanation":
				cq.Explanation = v.(string)
			case "id":
				cq.Id = v.(string)
			default:
				detailQuestion := v.(map[string]any)
				opt := &exercise.Option{
					Option:  k,
					Content: detailQuestion["content"].(string),
					Score:   int64(detailQuestion["score"].(float64)),
				}
				cq.Options = append(cq.Options, opt)
			}
		}
		cqs = append(cqs, cq)
	}

	// 题目列表
	q := &exercise.Question{
		ChoiceQuestions: cqs,
	}
	// 作答记录
	records := make([]*exercise.Records, 0)
	h := &exercise.History{
		Records: records,
	}
	// 练习
	e := &exercise.Exercise{
		Question: q,
		History:  h,
		Like:     0,
		Status:   0,
	}
	return e, nil
}

func generate(ctx context.Context, grade int64, m map[string]any, l *log.Log) (map[string]any, error) {
	var retrieve *coze.RetrieveChatsResp
	cli, req := getCozeCli(), buildReq(grade, m, l)
	resp, err := cli.Chat.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	//// 最少需要四十五秒
	time.Sleep(45 * time.Second)
	// 60s timeout
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	chat := resp.Chat
	conversationID, chatID := chat.ConversationID, chat.ID
	for chat.Status == coze.ChatStatusInProgress {
		select {
		case <-timeout:
			if _, err = cli.Chat.Cancel(ctx, &coze.CancelChatsReq{
				ConversationID: conversationID,
				ChatID:         chatID,
			}); err != nil {
				logx.Error("generate exercise: timeout and cancel error %s", err)
			}
			logx.Info("generate exercise: timeout and cancel chat")
			return nil, consts.ErrExerciseTimeout
		case <-ticker.C:
			retrieve, err = cli.Chat.Retrieve(ctx, &coze.RetrieveChatsReq{
				ConversationID: conversationID,
				ChatID:         chatID,
			})
			if err != nil {
				continue
			}

			chat = retrieve.Chat
			if chat.Status == coze.ChatStatusCompleted {
				message, err := cli.Chat.Messages.List(ctx, &coze.ListChatsMessagesReq{
					ConversationID: conversationID,
					ChatID:         chat.ID,
				})
				if err != nil {
					logx.Error("generate exercise: list message error %s", err)
					return nil, consts.ErrExercise
				}
				result := make(map[string]any)
				plainResult := message.Messages[0].Content
				if err = json.Unmarshal([]byte(plainResult), &result); err != nil {
					logx.Error("generate exercise: list message error %s", err)
					return nil, consts.ErrExercise
				}
				return result, nil
			}
		}
	}

	logx.Error("generate exercise error", err)
	return nil, nil
}

func buildReq(grade int64, m map[string]any, l *log.Log) *coze.CreateChatsReq {
	// 作文正文
	var essay strings.Builder
	paragraphs := m["text"].([]any)
	for _, paragraph := range paragraphs {
		paragraph := paragraph.([]any)
		for _, sentence := range paragraph {
			essay.WriteString(sentence.(string))
		}
	}

	return &coze.CreateChatsReq{
		BotID:  config.GetConfig().Coze.BotId,
		UserID: "exercise",
		Messages: []*coze.Message{
			coze.BuildUserQuestionText(fmt.Sprintf("年级:%v,作文标题:%s\n正文:%s\n批改结果:%s\n", grade, m["title"], essay.String(), l.Response), nil),
		},
	}
}
