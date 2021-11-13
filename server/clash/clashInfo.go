package clash

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"net/http"
	"os"
)

// GetClashInfo 读取clash的一些设置给前端
func GetClashInfo(c *gin.Context) {
	var rawConfig RawConfig
	config, err := os.ReadFile(mid.Dir + "/config/configClash.yaml")
	if err != nil {
		mid.Log.Error(err.Error())
	}
	toJSON, err := yaml.YAMLToJSON(config)
	if err != nil {
		mid.Log.Error(err.Error())
	}
	err = json.Unmarshal(toJSON, &rawConfig)
	if err != nil {
		mid.Log.Error(err.Error())
	}
	resultBody := mid.ResultBody{
		Code:    2000,
		Data:    rawConfig,
		Message: "查询成功",
	}
	c.JSON(http.StatusOK, resultBody)
}
