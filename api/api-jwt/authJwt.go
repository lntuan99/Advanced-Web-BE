package api_jwt

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// MapClaims type that uses the map[string]interface{} for JSON decoding
// This is the default claims type if you don't supply one
type MapClaims map[string]interface{}

// GinJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to FAKLoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
type GinJWTMiddleware struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Duration that a api-jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on userID and
	// password. Must return true on success, false on failure. Required.
	// Option return user data, if so, user data will be stored in Claim Array.
	Authenticator func(c *gin.Context) (interface{}, error)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(data interface{}, c *gin.Context) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on api-jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(data interface{}) MapClaims

	// LibraryUser can define own Unauthorized func.
	Unauthorized func(*gin.Context, int, string)

	// LibraryUser can define own LoginResponse func.
	LoginResponse func(*gin.Context, string, time.Time, interface{})

	// LibraryUser can define own RefreshResponse func.
	RefreshResponse func(*gin.Context, int, string, time.Time)

	// Set the identity handler function
	IdentityHandler func(jwt.MapClaims) interface{}

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time

	// HTTP Status messages for when something in the JWT middleware fails.
	// Check errors (e) to determine the appropriate errors message.
	HTTPStatusMessageFunc func(e error, c *gin.Context) string

	// Private key file for asymmetric algorithms
	PrivKeyFile string

	// Public key file for asymmetric algorithms
	PubKeyFile string

	// Private key
	privKey *rsa.PrivateKey

	// Public key
	pubKey *rsa.PublicKey

	// Optionally return the token as a cookie
	SendCookie bool

	// Allow insecure cookies for development over http
	SecureCookie bool

	// Required authentication
	IsRequired bool
}

var (
	// ErrMissingRealm indicates Realm name is required
	ErrMissingRealm = errors.New("realm is missing")

	// ErrMissingSecretKey indicates Secret key is required
	ErrMissingSecretKey = errors.New("secret key is required")

	// ErrForbidden when HTTP api-status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrMissingAuthenticatorFunc indicates Authenticator is required
	ErrMissingAuthenticatorFunc = errors.New("ginJWTMiddleware.Authenticator func is undefined")

	// ErrMissingLoginValues indicates a user tried to authenticate without username or password
	ErrMissingEmailLoginValues = errors.New("missing Email or Password")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	ErrFailedAuthentication = errors.New("incorrect Username or Password")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New("token is expired")

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("api-api-jwt header is empty")

	// ErrInvalidAuthHeader indicates api-api-jwt header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("api-api-jwt header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cokie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("private key invalid")

	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = errors.New("public key invalid")
)

// Login form structure.
type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (mw *GinJWTMiddleware) readKeys() error {
	err := mw.privateKey()
	if err != nil {
		return err
	}
	err = mw.publicKey()
	if err != nil {
		return err
	}
	return nil
}

func (mw *GinJWTMiddleware) privateKey() error {
	keyData, err := ioutil.ReadFile(mw.PrivKeyFile)
	if err != nil {
		return ErrNoPrivKeyFile
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPrivKey
	}
	mw.privKey = key
	return nil
}

func (mw *GinJWTMiddleware) publicKey() error {
	keyData, err := ioutil.ReadFile(mw.PubKeyFile)
	if err != nil {
		return ErrNoPubKeyFile
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPubKey
	}
	mw.pubKey = key
	return nil
}

func (mw *GinJWTMiddleware) usingPublicKeyAlgo() bool {
	switch mw.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

// MiddlewareInit initialize api-jwt configs.
func (mw *GinJWTMiddleware) MiddlewareInit() error {
	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 {
		mw.TokenHeadName = "Bearer"
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(data interface{}, c *gin.Context) bool {
			return true
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		}
	}

	if mw.LoginResponse == nil {
		mw.LoginResponse = func(c *gin.Context, token string, expire time.Time, data interface{}) {
			c.JSON(http.StatusOK, gin.H{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.RefreshResponse == nil {
		mw.RefreshResponse = func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(claims jwt.MapClaims) interface{} {
			return claims["id"]
		}
	}

	if mw.HTTPStatusMessageFunc == nil {
		mw.HTTPStatusMessageFunc = func(e error, c *gin.Context) string {
			return e.Error()
		}
	}

	if mw.Realm == "" {
		return ErrMissingRealm
	}

	if mw.usingPublicKeyAlgo() {
		return mw.readKeys()
	}

	if mw.Key == nil {
		return ErrMissingSecretKey
	}
	return nil
}

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *GinJWTMiddleware) MiddlewareFuncUser() gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			mw.unauthorized(c, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(err, nil))
		}
	}

	return func(c *gin.Context) {
		mw.middlewareImplUser(c)
	}
}

func (mw *GinJWTMiddleware) MiddlewareFuncAdminUser() gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			mw.unauthorized(c, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(err, nil))
		}
	}

	return func(c *gin.Context) {
		mw.middlewareImplAdminUser(c)
	}
}

func (mw *GinJWTMiddleware) middlewareImplUser(c *gin.Context) {
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, c))
		return
	}

	if token != nil {
		claims := token.Claims.(jwt.MapClaims)
		id := mw.IdentityHandler(claims)

		// Validate user id information
		idInUint, convertedErr := strconv.Atoi(fmt.Sprintf("%v", id))
		if convertedErr != nil {
			mw.unauthorized(c, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, c))
			return
		}

		existedUser, isExpired, user := model.User{}.FindUserByID(uint(idInUint))
		if existedUser == false {
			mw.unauthorized(c, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, c))
			return
		}

		if isExpired {
			mw.unauthorized(c, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, c))
			return
		}

		if !mw.Authorizator(user, c) {
			mw.unauthorized(c, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, c))
			return
		}

		// Set & Return
		c.Set("JWT_PAYLOAD", claims)
		c.Set("userId", util.ToInt(id))
		c.Set("user", user)
	}

	c.Next()
}

func (mw *GinJWTMiddleware) middlewareImplAdminUser(c *gin.Context) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)
	if user.IsAdmin {
		c.Next()
	} else {
		base.ResponseError(c, base.CodePermissionDenied)
		c.Abort()
	}
}

func (mw *GinJWTMiddleware) SignedString(token *jwt.Token) (string, error) {
	var tokenString string
	var err error
	if mw.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(mw.privKey)
	} else {
		tokenString, err = token.SignedString(mw.Key)
	}
	return tokenString, err
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the GinJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {
	token, _ := mw.parseToken(c)
	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken, c))
		return
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	newClaims["id"] = claims["id"]
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.SignedString(newToken)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation, c))
		return
	}

	mw.RefreshResponse(c, http.StatusOK, tokenString, expire)
}

// TokenGenerator method that clients can use to get a api-jwt token.
func (mw *GinJWTMiddleware) TokenGenerator(userID string, data MapClaims) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	expire := mw.TimeFunc().UTC().Add(mw.Timeout)
	claims["id"] = userID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.SignedString(token)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expire, nil
}

func (mw *GinJWTMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		if mw.IsRequired {
			return "", ErrEmptyAuthHeader
		} else {
			return "", nil
		}
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		if mw.IsRequired {
			return "", ErrInvalidAuthHeader
		} else {
			return "", nil
		}
	}

	return parts[1], nil
}

func (mw *GinJWTMiddleware) jwtFromQuery(c *gin.Context, key string) (string, error) {
	token := c.Query(key)

	if token == "" {
		if mw.IsRequired {
			return "", ErrEmptyQueryToken
		} else {
			return "", nil
		}
	}

	return token, nil
}

func (mw *GinJWTMiddleware) jwtFromCookie(c *gin.Context, key string) (string, error) {
	cookie, _ := c.Cookie(key)

	if cookie == "" {
		if mw.IsRequired {
			return "", ErrEmptyCookieToken
		} else {
			return "", nil
		}
	}

	return cookie, nil
}

func (mw *GinJWTMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split(mw.TokenLookup, ":")
	switch parts[0] {
	case "header":
		token, err = mw.jwtFromHeader(c, parts[1])
	case "query":
		token, err = mw.jwtFromQuery(c, parts[1])
	case "cookie":
		token, err = mw.jwtFromCookie(c, parts[1])
	}

	if err != nil {
		return nil, err
	}

	if token == "" && !mw.IsRequired {
		return nil, nil
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		// save token string if vaild
		c.Set("JWT_TOKEN", token)
		if mw.usingPublicKeyAlgo() {
			return mw.pubKey, nil
		}
		return mw.Key, nil
	})
}

func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {
	if mw.Realm == "" {
		mw.Realm = "gin api-jwt"
	}

	//c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	c.Abort()

	mw.Unauthorized(c, code, message)
}
