package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var ps []interface{}

type log struct {
	Level  string `default:"info"`
	MaxAge int    `default:"30"`
	Dir    string `default:"log"`
}

var Log log

var configPath = pflag.StringP("configPath", "c", ".", "配置文件路径")

func init() {
	pflag.Parse()
	viper.SetConfigName("config")
	viper.AddConfigPath(*configPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		viper.WriteConfigAs(*configPath + "/config.yml")
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(_ fsnotify.Event) {
		for _, p := range ps {
			go read(p)
		}
	})
	AddConfig(&Log)
}

// 添加配置
func AddConfig(p ...interface{}) {
	for _, p := range p {
		isSet(p)
		read(p)
	}
	ps = append(ps, p...)
}

// 配置是否完整，不完整补充写入默认值
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
		switch typeField.Type.Kind() {
		case reflect.String:
			viper.SetDefault(genKey(v.Type().Name(), typeField.Name), typeField.Tag.Get("default"))
		case reflect.Int:
			i, _ := strconv.Atoi(typeField.Tag.Get("default"))
			viper.SetDefault(genKey(v.Type().Name(), typeField.Name), i)
		case reflect.Float64:
			f, _ := strconv.ParseFloat(typeField.Tag.Get("default"), 64)
			viper.SetDefault(genKey(v.Type().Name(), typeField.Name), f)
		case reflect.Bool:
			b, _ := strconv.ParseBool(typeField.Tag.Get("default"))
			viper.SetDefault(genKey(v.Type().Name(), typeField.Name), b)
		case reflect.Slice:
			switch reflect.New(field.Type().Elem()).Elem().Type().Kind() {
			case reflect.String:
				ss := strings.Split(typeField.Tag.Get("default"), ",")
				viper.SetDefault(genKey(v.Type().Name(), typeField.Name), ss)
			case reflect.Int:
				ss := strings.Split(typeField.Tag.Get("default"), ",")
				is := make([]int, len(ss))
				for i, s := range ss {
					is[i], _ = strconv.Atoi(s)
				}
				viper.SetDefault(genKey(v.Type().Name(), typeField.Name), is)
			}
		}
	}
}

// 读取配置
func read(p interface{}) {
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)
		switch typeField.Type.Kind() {
		case reflect.String:
			if strings.Contains(strings.ToLower(typeField.Name), "dir") {
				field.SetString(createDir(viper.GetString(genKey(v.Type().Name(), typeField.Name))))
			} else {
				field.SetString(viper.GetString(genKey(v.Type().Name(), typeField.Name)))
			}
		case reflect.Int:
			field.SetInt(viper.GetInt64(genKey(v.Type().Name(), typeField.Name)))
		case reflect.Float64:
			field.SetFloat(viper.GetFloat64(genKey(v.Type().Name(), typeField.Name)))
		case reflect.Bool:
			field.SetBool(viper.GetBool(genKey(v.Type().Name(), typeField.Name)))
		case reflect.Slice:
			switch reflect.New(field.Type().Elem()).Elem().Type().Kind() {
			case reflect.String:
				field.Set(reflect.ValueOf(viper.GetStringSlice(genKey(v.Type().Name(), typeField.Name))))
			case reflect.Int:
				field.Set(reflect.ValueOf(viper.GetIntSlice(genKey(v.Type().Name(), typeField.Name))))
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

// 生成 key
func genKey(ss ...string) string {
	key := ""
	for _, s := range ss {
		key += "." + cc2sc(s)
	}
	return key[1:]
}

// 驼峰(CamelCase)转蛇形(snake_case)
func cc2sc(s string) string {
	var out []rune
	for i, r := range s {
		if i == 0 {
			out = append(out, unicode.ToLower(r))
		} else {
			if unicode.IsUpper(r) {
				out = append(out, '_')
			}
			out = append(out, unicode.ToLower(r))
		}
	}
	return string(out)
}
