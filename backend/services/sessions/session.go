package sessions

import (
	"bitbucket.org/daforester/go-audience/services/sessions/mongo"
	"bitbucket.org/daforester/go-backend-repository/component/libs/strgen"
	"encoding/base32"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memcached"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"strconv"
)

type sessionMethod int

const (
	Cookie    sessionMethod = 0
	Redis     sessionMethod = 1
	Memcached sessionMethod = 2
	MongoDB   sessionMethod = 3
)

func Default(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}

func SetupSessionUsingViper(r *gin.Engine, sessionName string) error {
	method := viper.GetString("SESSION_METHOD")

	authKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_AUTHKEY"))
	if err != nil {
		return err
	}

	encKey, err := base32.StdEncoding.DecodeString(viper.GetString("SESSION_ENCKEY"))
	if err != nil {
		return err
	}

	switch method {
	default:
		fallthrough
	case "Cookie":
		return SetupSession(Cookie, r, sessionName, authKey, encKey)
	case "Redis":
		address := viper.GetString("SESSION_REDIS_ADDRESS")
		network := viper.GetString("SESSION_REDIS_NETWORK_TYPE")
		password := viper.GetString("SESSION_REDIS_PASSWORD")
		size := viper.GetString("SESSION_REDIS_MAX_CONNECTIONS")

		return SetupSession(Redis, r, sessionName, authKey, encKey, size, network, address, password)
	case "Memcached":
		address := viper.GetString("SESSION_MEMCACHED_ADDRESS")
		port := viper.GetString("SESSION_MEMCACHED_PORT")
		keyPrefix := viper.GetString("SESSION_MEMCACHED_ADDRESS")

		return SetupSession(Memcached, r, sessionName, authKey, encKey, address, port, keyPrefix)
	case "MongoDB":
		//address (127.0.0.1), port(27017), path (/), name, maxage (3600), ensureTTL (true)
		address := viper.GetString("SESSION_MONGODB_ADDRESS")
		port := viper.GetString("SESSION_MONGODB_PORT")
		path := viper.GetString("SESSION_MONGODB_PATH")
		name := viper.GetString("SESSION_MONGODB_NAME")
		maxage := viper.GetString("SESSION_MONGODB_MAXAGE")
		ensureTTL := viper.GetString("SESSION_MONGODB_TTL")

		return SetupSession(MongoDB, r, sessionName, authKey, encKey, address, port, path, name, maxage, ensureTTL)
	}
}

func SetupSession(method sessionMethod, r *gin.Engine, sessionName string, authKey []byte, encKey []byte, options ...string) error {
	var store sessions.Store
	var err error

	if len(authKey) == 0 {
		authKey, err = strgen.RandBytes(64)

		fmt.Println("Generated AuthKey")
		fmt.Println(base32.StdEncoding.EncodeToString(authKey))

		if err != nil {
			return err
		}
	}
	if len(encKey) == 0 {
		encKey, err = strgen.RandBytes(32)

		fmt.Println("Generated EncKey")
		fmt.Println(base32.StdEncoding.EncodeToString(encKey))

		if err != nil {
			return err
		}
	}

	switch method {
	default:
		fallthrough
	case Cookie:
		store = setupCookieSession(authKey, encKey, options...)

	case Redis:
		store, err = setupRedisSession(authKey, encKey, options...)

	case Memcached:
		store, err = setupMemcachedSession(authKey, encKey, options...)

	case MongoDB:
		store, err = setupMongoDBSession(authKey, encKey, options...)
	}

	if err != nil {
		return err
	}

	r.Use(sessions.Sessions(sessionName, store))

	return nil
}

func setupCookieSession(authKey []byte, encKey []byte, options ...string) cookie.Store {
	if len(encKey) > 0 {
		return cookie.NewStore(authKey, encKey)
	} else {
		return cookie.NewStore(authKey)
	}
}

func setupRedisSession(authKey []byte, encKey []byte, options ...string) (redis.Store, error) {
	if len(options) < 4 {
		return nil, errors.New("require parameters for redis connection: size (10), network (tcp), address (127.0.0.1), password")
	}

	size, err := strconv.ParseInt(options[0], 10, 16)

	if err != nil {
		return nil, err
	}

	network := options[1]
	address := options[2]
	password := options[3]

	return redis.NewStore(int(size), network, address, password, authKey, encKey)
}

func setupMemcachedSession(authKey []byte, encKey []byte, options ...string) (memcached.Store, error) {
	if len(options) < 3 {
		return nil, errors.New("require parameters for Memcached connection: address (127.0.0.1), port(27017), keyPrefix")
	}

	server := fmt.Sprintf("%s:%s", options[0], options[1])
	keyPrefix := options[2]

	return memcached.NewStore(memcache.New(server), keyPrefix, authKey, encKey), nil
}

func setupMongoDBSession(authKey []byte, encKey []byte, options ...string) (mongo.Store, error) {
	if len(options) < 6 {
		return nil, errors.New("require parameters for redis connection: address (127.0.0.1), port(27017), path (/), name, maxage (3600), ensureTTL (true)")
	}

	url := fmt.Sprintf("%s:%s%s", options[0], options[1], options[2])

	session, err := mgo.Dial(url)

	if err != nil {
		return nil, err
	}

	name := options[3]
	maxAge, err := strconv.ParseInt(options[4], 10, 64)

	if err != nil {
		return nil, err
	}

	ensureTTL, err := strconv.ParseBool(options[5])

	if err != nil {
		return nil, err
	}

	c := session.DB("").C(name)

	return mongo.NewStore(c, int(maxAge), ensureTTL, authKey, encKey), nil
}
