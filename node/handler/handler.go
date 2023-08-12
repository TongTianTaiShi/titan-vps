package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/types"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("handler")

type (
	// RemoteAddr client address
	RemoteAddr struct{}
	// user id (node id)
	ID struct{}
	// LoginType filecoin tron eth ...
	LoginType struct{}
)

// Handler represents an HTTP handler that also adds remote client address and node ID to the request context
type Handler struct {
	// handler *auth.Handler
	verify func(ctx context.Context, token string) (*types.JWTPayload, error)
	next   http.HandlerFunc
}

// GetRemoteAddr returns the remote address of the client
func GetRemoteAddr(ctx context.Context) string {
	v, ok := ctx.Value(RemoteAddr{}).(string)
	if !ok {
		return ""
	}
	return v
}

// GetID returns the ID of the client
func GetID(ctx context.Context) string {
	if !api.HasPerm(ctx, api.RoleDefault, api.RoleUser) {
		return ""
	}

	v, ok := ctx.Value(ID{}).(string)
	if !ok {
		return ""
	}

	return v
}

// GetLoginType returns the login type of the client
func GetLoginType(ctx context.Context) types.LoginType {
	if !api.HasPerm(ctx, api.RoleDefault, api.RoleUser) {
		return -1
	}

	v, ok := ctx.Value(LoginType{}).(types.LoginType)
	if !ok {
		return -1
	}

	return v
}

// New returns a new HTTP handler with the given auth handler and additional request context fields
func New(verify func(ctx context.Context, token string) (*types.JWTPayload, error), next http.HandlerFunc) http.Handler {
	return &Handler{verify, next}
}

// ServeHTTP serves an HTTP request with the added client remote address and node ID in the request context
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteAddr := getClientIP(r)
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, RemoteAddr{}, remoteAddr)

	token := r.Header.Get("Authorization")

	if token == "" {
		token = r.FormValue("token")
		if token != "" {
			token = "Bearer " + token
		}
	}

	if token != "" {
		if !strings.HasPrefix(token, "Bearer ") {
			log.Warn("missing Bearer prefix in auth header")
			w.WriteHeader(401)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		payload, err := h.verify(ctx, token)
		if err != nil {
			log.Warnf("JWT Verification failed (originating from %s): %s, token:%s", r.RemoteAddr, err, token)
			w.WriteHeader(401)
			return
		}

		ctx = context.WithValue(ctx, ID{}, payload.ID)
		ctx = context.WithValue(ctx, LoginType{}, payload.LoginType)
		ctx = api.WithPerm(ctx, payload.Allow)
	}

	h.next(w, r.WithContext(ctx))
}

func getClientIP(req *http.Request) string {
	clientIP := req.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(req.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}
	return ""
}
