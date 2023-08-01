package web

import (
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-explorer/core/errors"
	"net/http"
)

func CreateKeyPair(c *gin.Context) {
	regionID := c.Query("regionID")
	KeyPairName := c.Query("KeyPairName")
	k, s := getAccessKeys()
	keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, KeyPairName)
	if err != nil {
		log.Errorf("CreateKeyPair err: %s", err.Error())
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}
	c.JSON(http.StatusOK, respJSON(JsonObject{
		"keyInfo": keyInfo,
	}))
}

func getAccessKeys() (string, string) {
	return AliYunAccessKeyID, AliYunAccessKeySecret
}
