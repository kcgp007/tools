package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type log struct {
	Level  string `config:"level" default:"info"`
	MaxAge int    `config:"max_age" default:"30"`
	Dir    string `config:"dir" default:"log" isDir:""`
}

var Log log

var configPath = pflag.StringP("configPath", "c", ".", "配置文件路径")

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(*configPath)
	AddConfig(&Log)
}

// 添加配置
func AddConfig(p interface{}) {
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		setDefault(p)
		viper.WriteConfigAs(*configPath + "/config.yml")
	} else {
		isSet(p)
	}
	read(p)(*new(fsnotify.Event))
	viper.WatchConfig()
	viper.OnConfigChange(read(p))
}

// 配置是否完整，不完整重新写入默认值
func isSet(p interface{}) {
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		if !viper.IsSet(v.Type().Name() + "." + typeField.Tag.Get("config")) {
			setDefault(p)
			viper.WriteConfig()
			return
		}
	}
}

// 设置配置默认值
func setDefault(p interface{}) {
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)
		if !viper.IsSet(v.Type().Name() + "." + typeField.Tag.Get("config")) {

		}
		switch typeField.Type.Kind() {
		case reflect.String:
			viper.SetDefault(v.Type().Name()+"."+typeField.Tag.Get("config"), typeField.Tag.Get("default"))
		case reflect.Int:
			i, _ := strconv.Atoi(typeField.Tag.Get("default"))
			viper.SetDefault(v.Type().Name()+"."+typeField.Tag.Get("config"), i)
		case reflect.Float64:
			f, _ := strconv.ParseFloat(typeField.Tag.Get("default"), 64)
			viper.SetDefault(v.Type().Name()+"."+typeField.Tag.Get("config"), f)
		case reflect.Bool:
			b, _ := strconv.ParseBool(typeField.Tag.Get("default"))
			viper.SetDefault(v.Type().Name()+"."+typeField.Tag.Get("config"), b)
		case reflect.Slice:
			switch reflect.New(field.Type().Elem()).Elem().Type().Kind() {
			case reflect.String:
				ss := strings.Split(typeField.Tag.Get("default"), ",")
				viper.SetDefault(v.Type().Name()+"."+typeField.Tag.Get("config"), ss)
			case reflect.Int:
				ss := strings.Split(typeField.Tag.Get("default"), ",")
				is := make([]int, len(ss))
				for i, s := range ss {
					is[i], _ = strconv.Atoi(s)
				}
				viper.SetDefault(v.Type().Name()+"."+typeField.Tag.Get("config"), is)
			}
		}
	}
}

// 读取配置
func read(p interface{}) func(_ fsnotify.Event) {
	return func(_ fsnotify.Event) {
		v := reflect.ValueOf(p).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			typeField := v.Type().Field(i)
			switch typeField.Type.Kind() {
			case reflect.String:
				_, isDir := typeField.Tag.Lookup("isDir")
				if isDir {
					field.SetString(createDir(viper.GetString(v.Type().Name() + "." + typeField.Tag.Get("config"))))
				} else {
					field.SetString(viper.GetString(v.Type().Name() + "." + typeField.Tag.Get("config")))
				}
			case reflect.Int:
				field.SetInt(viper.GetInt64(v.Type().Name() + "." + typeField.Tag.Get("config")))
			case reflect.Float64:
				field.SetFloat(viper.GetFloat64(v.Type().Name() + "." + typeField.Tag.Get("config")))
			case reflect.Bool:
				field.SetBool(viper.GetBool(v.Type().Name() + "." + typeField.Tag.Get("config")))
			case reflect.Slice:
				switch reflect.New(field.Type().Elem()).Elem().Type().Kind() {
				case reflect.String:
					field.Set(reflect.ValueOf(viper.GetStringSlice(v.Type().Name() + "." + typeField.Tag.Get("config"))))
				case reflect.Int:
					field.Set(reflect.ValueOf(viper.GetIntSlice(v.Type().Name() + "." + typeField.Tag.Get("config"))))
				}
			}
		}
	}
}

// 创建文件夹
func createDir(path string) string {
	path, _ = filepath.Abs(path)
	if _, err := os.Stat(path); err != nil {
		os.Mkdir(path, os.ModePerm)
	}
	return path
}
