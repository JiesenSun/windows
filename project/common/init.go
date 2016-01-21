package common

import (
	"database/sql"
	"flag"

	"gopkg.in/mgo.v2"

	"project/common/config"
	"project/common/mongo"
	"project/common/mysql"
	"project/common/nsq"
	"project/common/redis"
	"project/common/syslog"
)

var (
	DEBUG                              = true
	g_ini_conf  config.ConfigContainer = nil
	g_conf_file                        = flag.String("conf_file", "../config/default.conf", "the program config filename")
)

func Config() config.ConfigContainer {
	return g_ini_conf
}

func syslogInit() {
	ini_conf := Config()
	cfg := syslog.NewConfig()
	cfg.ModelName = ini_conf.DefaultString("module_name", "unknow_module_name")
	cfg.LogLevel = ini_conf.DefaultInt("log_level", 0)

	cfg.ConsoleLog = ini_conf.DefaultBool("log_console_enable", false)
	cfg.ConsoleLogLevel = ini_conf.DefaultInt("log_console_level", 0)

	cfg.FileLog = ini_conf.DefaultBool("log_file_enable", false)
	cfg.FileName = ini_conf.DefaultString("log_file_name", "logfile")
	cfg.Directory = ini_conf.DefaultString("log_file_path", "../log")
	cfg.FileLogLevel = ini_conf.DefaultInt("log_file_level", 0)

	cfg.MysqlLog = ini_conf.DefaultBool("log_mysql_enable", false)
	cfg.MysqlLogLevel = ini_conf.DefaultInt("log_mysql_level", 0)
	cfg.Host = ini_conf.DefaultString("mysql_host", "127.0.0.1")
	cfg.Database = ini_conf.DefaultString("mysql_log_db", "logDB")
	cfg.User = ini_conf.DefaultString("mysql_user", "root")
	cfg.Password = ini_conf.DefaultString("mysql_pwd", "Youkang@0814")

	cfg.MongoLog = ini_conf.DefaultBool("log_mongo_enable", false)
	cfg.MongoLogLevel = ini_conf.DefaultInt("log_mongo_level", 0)
	cfg.MongoAddr = ini_conf.DefaultString("mongodb_address", "127.0.0.1:27017")

	cfg.NsqdLog = ini_conf.DefaultBool("log_nsq_client_enable", false)
	cfg.NsqdAddrs = ini_conf.DefaultString("nsq_address", "127.0.0.1:4150")
	cfg.NsqdLogLevel = ini_conf.DefaultInt("log_nsq_client_level", 0)

	cfg.NsqdSvrLog = ini_conf.DefaultBool("log_nsq_server_enable", false)
	if cfg.NsqdSvrLog == true {
		cfg.FileLog = false
		cfg.ConsoleLog = false
		cfg.MysqlLog = false
		cfg.MongoLog = false
		cfg.NsqdLog = false
	}

	cfg.NsqdSvrConfig.FileLog = ini_conf.DefaultBool("log_file_enable", false)
	cfg.NsqdSvrConfig.FileName = cfg.FileName
	cfg.NsqdSvrConfig.Directory = cfg.Directory
	cfg.NsqdSvrConfig.FileLogLevel = cfg.FileLogLevel

	cfg.NsqdSvrConfig.NsqdLog = false // 必须false 避免死循环日志 *g_log_nsq_client_enable
	cfg.NsqdSvrConfig.NsqdAddrs = cfg.NsqdAddrs
	cfg.NsqdSvrConfig.NsqdLogLevel = cfg.NsqdLogLevel

	cfg.NsqdSvrConfig.MysqlLog = ini_conf.DefaultBool("log_mysql_enable", false)
	cfg.NsqdSvrConfig.Database = cfg.Database
	cfg.NsqdSvrConfig.MysqlLogLevel = cfg.MysqlLogLevel
	cfg.NsqdSvrConfig.Host = cfg.Host
	cfg.NsqdSvrConfig.User = cfg.User
	cfg.NsqdSvrConfig.Password = cfg.Password

	cfg.NsqdSvrConfig.MongoLog = ini_conf.DefaultBool("log_mongo_enable", false)
	cfg.NsqdSvrConfig.MongoLogLevel = cfg.MongoLogLevel
	cfg.NsqdSvrConfig.MongoAddr = cfg.MongoAddr

	cfg.NsqdSvrConfig.ConsoleLog = ini_conf.DefaultBool("log_console_enable", false)
	cfg.NsqdSvrConfig.ConsoleLogLevel = cfg.ConsoleLogLevel
	syslog.SysLogInit(cfg)
}

func init() {
	flag.Parse()
	ini_conf, err := config.NewConfig("ini", *g_conf_file)
	if err != nil {
		ini_conf = config.NewFakeConfig()
		println(err.Error())
	}
	g_ini_conf = ini_conf

	syslogInit()
	NsqInit(g_ini_conf.DefaultString("nsq_address", "127.0.0.1:4150"))
	MongoInit(g_ini_conf.DefaultString("mongo_address", "127.0.0.1:27017"))
	RedisInit(g_ini_conf.DefaultString("redis_address", "127.0.0.1:6379"))
}

//=======================nsq common===========================
func NsqInit(addrs string) {
	nsq.Init(addrs)
}
func NsqPublish(topic string, body []byte) error {
	return nsq.Publish(topic, body)
}
func NsqConsumer(topic, channel string, handle nsq.Handler) (*nsq.ConsumerT, error) {
	return nsq.Consumer(topic, channel, handle)
}
func NsqConsumerGO(topic, channel string, goCount uint, handle nsq.Handler) (*nsq.ConsumerT, error) {
	return nsq.ConsumerGO(topic, channel, goCount, handle)
}
func NsqDeinit() {
	nsq.Deinit()
}

//========================mongo common====================
func MongoInit(mongodbAddr string) {
	mongo.Init(mongodbAddr, 2)
}
func MongoGet() *mgo.Session {
	return mongo.Get()
}
func MongoPut(sess *mgo.Session) {
	mongo.Put(sess)
}
func MongoCollection(db, table string) *mgo.Collection {
	return mongo.Collection(db, table)
}
func MongoDeinit() {
	mongo.Deinit()
}

//=========================mysql common========================
func MysqlInit(host, dbname, user, passwd string) {
	mysql.Init(host, dbname, user, passwd, 2)
}
func MysqlGet() *sql.DB {
	return mysql.Get()
}
func MysqlPut(sqlDB *sql.DB) {
	mysql.Put(sqlDB)
}
func MysqlExec(sqlStr string, args ...interface{}) error {
	return mysql.Exec(sqlStr, args...)
}
func MysqlExecRet(sqlStr string, args ...interface{}) (uint64, error) {
	return mysql.ExecRet(sqlStr, args...)
}
func MysqlQuery(sqlStr string, args ...interface{}) (*sql.Rows, error) {
	return mysql.Query(sqlStr, args...)
}
func MysqlDeinit() {
	mysql.Deinit()
}

//=====================redis common====================
func RedisInit(addrs string) {
	redis.Init(addrs)
}
func RedisDeinit() {
	redis.Deinit()
}
