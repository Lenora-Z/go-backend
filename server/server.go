//Created by Goland
//@User: lenora
//@Date: 2021/2/5
//@Time: 2:13 下午
package server

import (
	"context"
	"fmt"
	"github.com/Lenora-Z/go-backend/conf"
	"github.com/Lenora-Z/go-backend/docs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"strings"
	"time"
)

type Server interface {
	Run(configPath string) error
	Close() error
}

type defaultServer struct {
	name    string
	port    string
	conf    *conf.ServerConfig
	db      *gorm.DB
	mongoDb *mongo.Database
	engine  *gin.Engine
}

func NewServer(name string, port string) Server {
	s := new(defaultServer)
	s.name = name
	s.port = port
	return s
}

func (ds *defaultServer) Run(configPath string) error {
	//config
	if err := ds.config(configPath); err != nil {
		return fmt.Errorf("ds.config():%s", err.Error())
	}

	//db
	if err := ds.dbClient(); err != nil {
		return fmt.Errorf("ds.dbClient:%s", err.Error())
	}

	// router
	if err := ds.initRouter(); err != nil {
		return fmt.Errorf("rs.router(): %s", err.Error())
	}

	// otherInit
	if err := ds.init(); err != nil {
		return fmt.Errorf("rs.init(): %s", err.Error())
	}

	//port
	port := ds.port
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return ds.engine.Run(port)
}

func (ds *defaultServer) Close() error {
	if ds.db != nil {
		if err := ds.db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (ds *defaultServer) config(configPath string) error {
	ds.conf = new(conf.ServerConfig)
	err := configor.Load(ds.conf, configPath)
	if err != nil {
		return err
	}
	return nil
}

func (ds *defaultServer) dbClient() error {
	db, err := getDbConnection(ds.conf.DbConfig, 0, 0)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}
	ds.db = db
	return nil
}

func (rs *defaultServer) mongoClient() error {
	db, err := getMongoConnection(rs.conf.MongoConfig)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}
	//rs.mongoDB = db
	rs.mongoDb = db.Database(rs.conf.MongoConfig.DbName)

	return nil
}

func (ds *defaultServer) initRouter() error {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	ds.engine = r
	if ds.conf.Debug {
		swaggerDoc(r)
	}
	ds.routers()
	return nil
}

func (ds *defaultServer) init() error {
	return nil
}

func getDbConnection(conf *conf.DBConfig, maxIdle, maxOpen int) (*gorm.DB, error) {
	format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"
	dsn := fmt.Sprintf(format, conf.User, conf.Password, conf.Name, conf.Port, conf.DbName, conf.Charset, url.QueryEscape(conf.Loc))
	logrus.Infof("dsn=%s", dsn)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.LogMode(conf.LogMode)
	idle := conf.MaxIdle
	if maxIdle > 0 {
		idle = maxIdle
	}
	open := conf.MaxOpen
	if open > 0 {
		open = maxOpen
	}
	db.DB().SetMaxIdleConns(idle)
	db.DB().SetMaxOpenConns(open)
	return db, nil
}

// https://mongodb-documentation.readthedocs.io/en/latest/reference/connection-string.html#gsc.tab=0
func getMongoConnection(conf *conf.MongoConfig) (*mongo.Client, error) {
	format := "mongodb://%s:%s@%s:%s/%s?maxPoolSize=%s&minPoolSize=%s"
	uri := fmt.Sprintf(format, conf.User, conf.Password, conf.IP, conf.Port, conf.DbName, conf.MaxOpen, conf.MaxIdle)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func swaggerDoc(ctx *gin.Engine) {
	docs.SwaggerInfo.Title = AppTitle
	docs.SwaggerInfo.Description = AppName
	docs.SwaggerInfo.Version = AppVersion
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	ctx.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

var validate *validator.Validate

func validatorInstance() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	return validate
}
