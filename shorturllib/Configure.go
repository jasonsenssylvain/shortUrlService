package shorturllib

import (
	"errors"
	"strconv"

	"github.com/ewangplay/config"
)

// 存储各个具体信息
type Configure struct {
	ConfigureMap map[string]string
}

//NewConfigure 解析fileName里的配置信息，返回配置信息
func NewConfigure(fileName string) (*Configure, error) {
	config := &Configure{}

	config.ConfigureMap = make(map[string]string)
	err := config.ParseConfigure(fileName)
	if err != nil {
		return nil, err
	}
	return config, nil
}

//ParseConfigure 读取配置文件信息
func (t *Configure) ParseConfigure(filename string) error {
	cfg, err := config.ReadDefault(filename)
	if err != nil {
		return err
	}

	t.loopConfigure("server", cfg)
	t.loopConfigure("service", cfg)
	t.loopConfigure("redis", cfg)

	return nil
}

// 通过循环section里的数据 存入到配置
func (t *Configure) loopConfigure(sectionName string, cfg *config.Config) error {
	if cfg.HasSection(sectionName) {
		section, err := cfg.SectionOptions(sectionName)
		if err == nil {
			for _, v := range section {
				options, err := cfg.String(sectionName, v)
				if err == nil {
					t.ConfigureMap[v] = options
				}
			}
		}
	}
	return nil
}

//GetPort 获取当前 服务的 端口
func (t *Configure) GetPort() (int, error) {
	portstr, ok := t.ConfigureMap["port"]
	if ok == false {
		return 9090, errors.New("no port set, use default")
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		return 9090, err
	}
	return port, nil
}

//GetRedisHost 获取redis主机地址
func (t *Configure) GetRedisHost() (string, error) {
	redishost, ok := t.ConfigureMap["redishost"]
	if ok == false {
		return "127.0.0.1", errors.New("no redis host, use default")
	}
	return redishost, nil
}

//GetRedisPort 获取 redis 的端口设置
func (t *Configure) GetRedisPort() (string, error) {
	redisport, ok := t.ConfigureMap["redisport"]
	if ok == false {
		return "6379", errors.New("no redis port, use default 6379")
	}
	return redisport, nil
}

//GetRedisStatus 判断是否需要开启
func (t *Configure) GetRedisStatus() bool {
	status, ok := t.ConfigureMap["status"]
	if ok == false {
		return true
	}
	if status == "true" {
		return true
	}
	return false
}

//GetHostInfo 短连接前缀域名
func (t *Configure) GetHostInfo() string {
	hostName, ok := t.ConfigureMap["hostname"]
	if ok == false {
		return "http://t.cn"
	}
	return hostName
}

//GetCounterType redis 还是 inner
func (t *Configure) GetCounterType() string {
	counterType, ok := t.ConfigureMap["counter"]
	if ok == false {
		return "inner"
	}
	return counterType
}
