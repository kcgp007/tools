package configTool

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"

	"github.com/spf13/viper"
)

func init() {
	if runtime.GOOS == "windows" {
		ed, _ := os.Executable()
		wd, _ := os.Getwd()
		if ed != wd {
			os.Chdir(ed)
		}
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfigAs("config.toml")
		} else {
			fmt.Println(err)
		}
	}
}

// 加载默认配置及配置文件数据
func Add(c any) {
	AddWithKey(reflect.TypeOf(c).Elem().Name(), c)
}

// 加载默认配置及配置文件数据
func AddWithKey(key string, c any) {
	bs, _ := json.Marshal(c)
	m := make(map[string]any)
	json.Unmarshal(bs, &m)
	viper.SetDefault(key, m)
	viper.UnmarshalKey(key, &c)
	viper.WriteConfig()
}
