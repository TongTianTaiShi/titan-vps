package web

import (
	"time"

	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("web")

type Server struct {
	cfg    config.RouteCfg
	router *gin.Engine
}

var AliYunAccessKeyID string

var AliYunAccessKeySecret string

func NewServer(cfg *config.MallCfg) (*Server, error) {
	gin.SetMode(cfg.RouteCfg.Mode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))
	ConfigRouter(router, cfg.RouteCfg)
	s := &Server{
		cfg:    cfg.RouteCfg,
		router: router,
	}
	AliYunAccessKeyID = cfg.AliyunAccessKeyID
	AliYunAccessKeySecret = cfg.AliyunAccessKeySecret
	return s, nil
}

func ConfigRouter(router *gin.Engine, cfg config.RouteCfg) {
	apiV1 := router.Group("/api/v1")
	authMiddleware, err := jwtGinMiddleware(cfg.SecretKey)
	if err != nil {
		log.Fatalf("jwt auth middleware: %v", err)
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		log.Fatalf("authMiddleware.MiddlewareInit: %v", err)
	}
	// admin
	// storage
	apiV1.GET("/test", CreateKeyPair)
	apiV1.POST("/login", authMiddleware.LoginHandler)
	apiV1.POST("/logout", authMiddleware.LogoutHandler)
	apiV1.Use(authMiddleware.MiddlewareFunc())
	apiV1.Use(AuthRequired(authMiddleware))
	apiV1.GET("/refresh_token", authMiddleware.RefreshHandler)
}

func (s *Server) Run() {
	err := s.router.Run(s.cfg.ApiListen)
	if err != nil {
		log.Fatal(err)
	}
}
