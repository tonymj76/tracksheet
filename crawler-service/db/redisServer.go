package db

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// RedisServer create a custom settings for redis
type RedisServerContext struct {
	URL         *URLModel
	ProductList []string
}

var Pool *work.WorkerPool

func NewRedisPool(maxActive, maxIdle int, wait bool, port string) *redis.Pool {
	return &redis.Pool{
		MaxActive: maxActive,
		MaxIdle:   maxIdle,
		Wait:      wait,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf(":%s", port))
		},
	}
}

func (r *RedisServerContext) SaveToDB(job *work.Job) error {
	err := DB.Create(r.URL).Error
	for _, value := range r.ProductList {
		link_model := URLProductLink{ProductURL: value, URLModelID: r.URL.ID, URLModel: *r.URL}
		err = DB.Create(&link_model).Error
	}
	return err
}

func (r *RedisServerContext) DeleteID(job *work.Job) error {
	var err error
	var productUrl URLProductLink
	if err := DB.Where("url_model_id = ?", r.URL.ID).Delete(&productUrl).Error; err != nil {
		return err
	}
	err = DB.Delete(&r.URL).Error
	return err
}
