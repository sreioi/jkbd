package ke4

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/sreioi/jkbd/core"
	"github.com/sreioi/jkbd/db"
	"github.com/sreioi/jkbd/env"
	"github.com/sreioi/jkbd/file"
	"github.com/sreioi/jkbd/request"
	"log"
	"strings"
	"sync"
)

// Ke4 科一题表结构
type Ke4 struct {
	Id         int    `json:"id" gorm:"UNIQUE,comment:问题id"`
	Question   string `json:"question" gorm:"type:varchar(1024);comment:问题"`
	OptionType int    `json:"option_type" gorm:"type:int(1);comment:(0->判断题1->单选题)"`
	OptionA    string `json:"optionA" gorm:"type:varchar(255);comment:选项a"`
	OptionB    string `json:"optionB" gorm:"type:varchar(255);comment:选项b"`
	OptionC    string `json:"optionC" gorm:"type:varchar(255);comment:选项c"`
	OptionD    string `json:"optionD" gorm:"type:varchar(255);comment:选项d"`
	OptionE    string `json:"optionE" gorm:"type:varchar(255);comment:选项e"`
	OptionF    string `json:"optionF" gorm:"type:varchar(255);comment:选项f"`
	OptionG    string `json:"optionG" gorm:"type:varchar(255);comment:选项g"`
	OptionH    string `json:"optionH" gorm:"type:varchar(255);comment:选项h"`
	MediaUrl   string `json:"media_url" gorm:"type:varchar(1024);comment:图片或视频url地址"`
	Answer     string `json:"answer" gorm:"type:varchar(5);comment:答案"`
	SoDesc     string `json:"so_desc" gorm:"type:varchar(255);comment:答案简要描述"`
	Desc       string `json:"desc" gorm:"type:text;comment:答案详情描述"`
	WrongRate  string `json:"wrongRate" gorm:"type:varchar(100);comment:错误率"`
}

// Ids 存放所有题目ID
var Ids []string

// consumersChannel 消费者队列 存放需要获取题目IDs
var consumersChannel = make(chan []string, env.Worker)

// wg WaitGroup
var wg sync.WaitGroup

func Pull() {
	if env.ID != "" {
		Ids = append(Ids, env.ID)
	} else {
		Ids = getIdsKe4()
	}

	// 开启worker个消费者去获取题目详情
	for i := 1; i <= env.Worker; i++ {
		go getQuestionByIds()
	}

	chunkSize := 30
	maxLen := len(Ids)
	// 将题目IDs放入消费者队列
	for i := 0; i < maxLen; i += chunkSize {
		wg.Add(1)
		end := i + chunkSize
		if end > maxLen {
			end = maxLen
		}
		consumersChannel <- Ids[i:end]
	}

	defer close(consumersChannel)
	// 等待所有消费者完成
	wg.Wait()
}

// getQuestionByIds 获取题目详情
func getQuestionByIds() {
	var url = env.API + "open/question/question-list.htm"

	for {
		idChunk, ok := <-consumersChannel
		if !ok {
			break
		}

		params := request.QueryParams{
			"_r":          core.GenerateRandomString("1"),
			"carType":     "car",
			"course":      "kemu3",
			"questionIds": strings.Join(idChunk, ","),
		}
		res := request.Get(url, params)
		if cast.ToInt(res["errorCode"]) != 0 {
			log.Printf("错误码：%s, 接口获取失败：%s", cast.ToString(res["errorCode"]), cast.ToString(res["message"]))
		}
		responseData := res["data"].([]interface{})
		transQuestion(responseData)

		wg.Done()
	}
}

// transQuestion 保存题目详情
func transQuestion(responseData []interface{}) {
	for _, results := range responseData {
		resultMap := results.(map[string]interface{})
		CpData := &Ke4{
			Id:         cast.ToInt(resultMap["questionId"]),
			Question:   cast.ToString(resultMap["question"]),
			OptionType: cast.ToInt(resultMap["optionType"]),
			OptionA:    cast.ToString(resultMap["optionA"]),
			OptionB:    cast.ToString(resultMap["optionB"]),
			OptionC:    cast.ToString(resultMap["optionC"]),
			OptionD:    cast.ToString(resultMap["optionD"]),
			OptionE:    cast.ToString(resultMap["optionE"]),
			OptionF:    cast.ToString(resultMap["optionF"]),
			OptionG:    cast.ToString(resultMap["optionG"]),
			OptionH:    cast.ToString(resultMap["optionH"]),
			MediaUrl:   cast.ToString(resultMap["mediaContent"]),
			Answer:     core.TransAnswer(cast.ToInt(resultMap["answer"])),
			SoDesc:     cast.ToString(resultMap["conciseExplain"]),
			Desc:       cast.ToString(resultMap["explain"]),
			WrongRate:  cast.ToString(resultMap["wrongRate"]),
		}

		db.DB.FirstOrCreate(&Ke4{}, CpData)
	}
}

// getIdsKe4 获取所有题目ID
func getIdsKe4() []string {
	fileName := "k4_question_ids.txt"
	if file.ExistFile(fileName) {
		if fileContent := file.ReadFile(fileName); fileContent != "" {
			return strings.Split(fileContent, ",")
		}
	}

	url := env.API + "open/exercise/sequence.htm"
	params := request.QueryParams{
		"_r":        core.GenerateRandomString("1"),
		"carStyle":  "xiaoche",
		"carType":   "car",
		"cityCode":  "110100",
		"course":    "kemu3",
		"kemuStyle": "kemu4",
	}
	res := request.Get(url, params)
	IDs := cast.ToStringSlice(res["data"])
	file.SaveFile(strings.Join(IDs, ","), fileName)
	return IDs
}

func CreateTableK4() {
	err := db.DB.AutoMigrate(&Ke4{})
	if err != nil {
		fmt.Println("创建k4表失败：", err)
		return
	}
	fmt.Println("创建k4表成功")
}
