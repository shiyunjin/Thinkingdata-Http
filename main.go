package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/ThinkingDataAnalytics/go-sdk/v2/src/thinkingdata"
	"github.com/spf13/cast"
	"pkg.moe/pkg/logger"
)

type BaseData struct {
	AccountId  interface{} `json:"pid,omitempty"`
	DistinctId string      `json:"#distinct_id,omitempty"`
}

func main() {
	config := thinkingdata.TDLogConsumerConfig{
		Directory: "./data",
	}
	// 初始化 logConsumer
	consumer, _ := thinkingdata.NewLogConsumerWithConfig(config)
	// 创建 te 对象
	te := thinkingdata.New(consumer)

	server := http.Server{Addr: ":4477"}
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean(r.URL.Path)
		trimmedPath := strings.TrimPrefix(cleanPath, "/")

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Get().Error("io.ReadAll: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{}
		if err := json.Unmarshal(payload, &data); err != nil {
			logger.Get().Error("json.Unmarshal: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		baseData := BaseData{}
		if err := json.Unmarshal(payload, &baseData); err != nil {
			logger.Get().Error("json.Unmarshal: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch trimmedPath {
		case "UserSet":
			if err := te.UserSet(cast.ToString(baseData.AccountId), baseData.DistinctId, data); err != nil {
				logger.Get().Error("te.UserSet: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case "UserSetOnce":
			if err := te.UserSetOnce(cast.ToString(baseData.AccountId), baseData.DistinctId, data); err != nil {
				logger.Get().Error("te.UserSetOnce: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case "UserAdd":
			if err := te.UserAdd(cast.ToString(baseData.AccountId), baseData.DistinctId, data); err != nil {
				logger.Get().Error("te.UserAdd: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case "UserDelete":
			if err := te.UserDelete(cast.ToString(baseData.AccountId), baseData.DistinctId); err != nil {
				logger.Get().Error("te.UserDelete: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		default:
			if err := te.Track(cast.ToString(baseData.AccountId), baseData.DistinctId, trimmedPath, data); err != nil {
				logger.Get().Error("te.Track: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})

	fmt.Println("Server is running at :4477")
	err := server.ListenAndServe()
	if err != nil {
		logger.Get().Error("server.ListenAndServe: %v", err)
	}
}
