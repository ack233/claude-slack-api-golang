package middlewares

import (
	"errors"
	"net/http"
	"regexp"
	"slackapi/pkgs/zlog"
	"time"

	"slackapi/pkgs/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CustomClaims struct {
	Token string `json:"token"`
	jwt.StandardClaims
}

// 一些变量

var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token")
)

func returnExpiredTokenErr(ctx *gin.Context, token string) {
	zlog.SugLog.Warn("token 过期", zap.Any("data", map[string]interface{}{
		"url":         ctx.Request.URL,
		"params":      ctx.Params,
		"tokenString": token,
	}))
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"error_code": 1,
		"message":    "授权已过期",
		"data":       map[string]interface{}{},
	})
	ctx.Abort()
}

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")

		if tokenString == "" {
			tokenString = ctx.Query("Authorization")
			if tokenString == "" {

				zlog.SugLog.Warn("请求未携带token，无权限访问", zap.Any("data", map[string]interface{}{
					"url":    ctx.Request.URL,
					"params": ctx.Params,
				}))
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error_code": 1,
					"message":    "请求未携带token，无权限访问",
					"data":       map[string]interface{}{},
				})
				ctx.Abort()
				return
			}
		}
		zlog.SugLog.Debug(tokenString)

		// parseToken 解析token包含的信息
		parse_relustr, err := ParseToken(tokenString)
		if err != nil {
			if err == TokenExpired {
				returnExpiredTokenErr(ctx, tokenString)
				return
			}
			zlog.SugLog.Error("token 错误", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error_code": 1,
				"message":    err.Error(),
				"data":       map[string]interface{}{},
			})
			ctx.Abort()
			return
		}
		ctx.Set("token", parse_relustr.Token)

	}
}

// CreateToken 生成一个token
func CreateToken(str string) (string, error) {
	claims := CustomClaims{
		Token: str,
		//5分钟后过期,
		//StandardClaims: jwt.StandardClaims{
		//	ExpiresAt: time.Now().Add(time.Hour * 5).Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SlackConfig.ENCRYPTION_KEY))
}

// 解析Tokne
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 需要从tokenString移除 Bearer，jwt-go这个包为了避免冗余，不会帮我们处理
	re, _ := regexp.Compile(`(?i)Bearer `)
	tokenString = re.ReplaceAllString(tokenString, "")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SlackConfig.ENCRYPTION_KEY), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.SlackConfig.ENCRYPTION_KEY, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
		return CreateToken(claims.Token)
	}
	return "", TokenInvalid
}
