package database

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/trietphm/gruber/config"
	"github.com/trietphm/gruber/model/mredis"
)

// PgI is a interface for manipulating data in database postgresql
type RedisI interface {
	// Push driver location to redis
	PushDriverLocationGeo(driverID int, lat, lng float64) error

	// Remove driver location from redis geo
	RemoveDriverLocationGeo(driverID int) error

	// GetNearestDrivers Get near available driver near a geo location
	GetNearestDrivers(lat, lng, radius float64, limit int) ([]mredis.DriverLocation, error)
}

// Redis
type Redis struct {
	redis.Client
}

var _ RedisI = RedisI(Redis{})

func OpenRedisDB(conf config.Redis) (*Redis, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Password, // no password set
		DB:       0,             // use default DB
	})

	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Redis{*db}, nil
}

// PushDriverLocationGeo Push driver to geo data
func (db Redis) PushDriverLocationGeo(driverID int, lat, lng float64) error {
	location := redis.GeoLocation{
		Latitude:  lat,
		Longitude: lng,
		Name:      strconv.Itoa(driverID),
	}
	res := db.GeoAdd(mredis.KeyDriverGeo, &location)
	if res == nil {
		return errors.New("Can not execute GEOADD")
	}

	return res.Err()
}

// RemoveDriverLocationGeo remove a driver from redis geo
func (db Redis) RemoveDriverLocationGeo(driverID int) error {
	res := db.ZRem(mredis.KeyDriverGeo, driverID)
	if res == nil {
		return errors.New("Can not execute ZREM")
	}

	return res.Err()
}

// GetNearestDrivers get nearest driver in a radius via Redis GEORADIUS, unit is kilometer
func (db Redis) GetNearestDrivers(lat, lng, radius float64, limit int) (locations []mredis.DriverLocation, err error) {
	query := redis.GeoRadiusQuery{
		Radius: radius,
		Unit:   "km",
		Count:  limit,
		Sort:   "ASC",
	}
	res := db.GeoRadius(mredis.KeyDriverGeo, lat, lng, &query)
	if res == nil {
		err = errors.New("Can not execute GEORADIUS")
		return
	}

	return
}
