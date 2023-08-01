package web

import (
	"context"
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	loginStatusFailure = iota
	loginStatusSuccess
)

type login struct {
	Username   string `form:"username" json:"username" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
	VerifyCode string `form:"verify_code" json:"verify_code" binding:"required"`
}

type loginResponse struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

var identityKey = "id"

func jwtGinMiddleware(secretKey string) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:             "User",
		Key:               []byte(secretKey),
		Timeout:           time.Hour,
		MaxRefresh:        24 * time.Hour,
		IdentityKey:       identityKey,
		SendAuthorization: true,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*types.User); ok {
				return jwt.MapClaims{
					identityKey: v.UUID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &types.User{
				UUID: claims[identityKey].(string),
			}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": loginResponse{
					Token:  token,
					Expire: expire.Format(time.RFC3339),
				},
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
			})
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginParams login
			loginParams.Username = c.Query("username")
			loginParams.VerifyCode = c.Query("verify_code")
			Signature := c.Query("sign")
			Address := c.Query("address")
			if loginParams.Username == "" {
				return "", jwt.ErrMissingLoginValues
			}
			if loginParams.VerifyCode == "" && loginParams.Password == "" && Signature == "" {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginParams.Username
			var user interface{}
			if Signature != "" {
				user, _ = loginBySignature(c, userID, Address, Signature)
			}
			return user, nil
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(200, gin.H{
				"code":    code,
				"msg":     message,
				"success": false,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		//TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenLookup: "header: JwtAuthorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		RefreshResponse: func(c *gin.Context, code int, token string, t time.Time) {
			c.Next()
		},
	})
}

func loginBySignature(c *gin.Context, username, address, msg string) (interface{}, error) {
	verifyCode, err := GetVerifyCode(c.Request.Context(), username+"C")
	if err != nil {
		return nil, NewErrorCode(InvalidParams, c)
	}
	if verifyCode == "" {
		return nil, NewErrorCode(Unknown, c)
	}
	if address == "" {
		address = username
	}
	publicKey, err := VerifyMessage(verifyCode, msg)
	address = strings.ToUpper(address)
	publicKey = strings.ToUpper(publicKey)
	if publicKey != address {
		return nil, NewErrorCode(Unknown, c)
	}
	return &types.User{UUID: "", UserName: username, Role: 0}, nil
}

func AuthRequired(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, e := authMiddleware.GetClaimsFromJWT(ctx)
		if e == nil {
			fmt.Println("==============now token time left=====================")
			fmt.Println(ctx.Query("user_id"))
			token, _ := ctx.Get("JWT_TOKEN")
			fmt.Println(token)
			fmt.Println(int64(claims["exp"].(float64)) - authMiddleware.TimeFunc().Unix())
			if int64(claims["exp"].(float64)-authMiddleware.Timeout.Seconds()/2) < authMiddleware.TimeFunc().Unix() {
				tokenString, _, e := authMiddleware.RefreshToken(ctx)
				if e == nil {
					ctx.Header("new-token", tokenString)
				}
			}
		}
		ctx.Next()
	}
}

func VerifyMessage(message string, signedMessage string) (string, error) {
	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
	hash := crypto.Keccak256Hash(hashedMessage)
	// Get the bytes of the signed message
	decodedMessage := hexutil.MustDecode(signedMessage)
	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}
	// Recover a public key from the signed message
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if sigPublicKeyECDSA == nil {
		log.Errorf("Could not get a public get from the message signature")
	}
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}

func GetVerifyCode(ctx context.Context, key string) (string, error) {
	verifyCode := ""
	return verifyCode, nil
}
