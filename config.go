package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	viper.OnConfigChange(func(e fsnotify.Event) {
		for _, p := range ps {
			read(p)
		}
	})
	AddConfig(&Log)
}

// 添加配置
func AddConfig(p ...interface{}) {
	for _, p := range p {
		read(p)
	}
	ps = append(ps, p...)
}

// 读取配置，补充缺少的配置
func read(p interface{}) {
	isWrite := false
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)
		key := genKey(v.Type().Name(), typeField.Name)
		switch typeField.Type.Kind() {
		case reflect.String:
			if !viper.InConfig(cc2sc(v.Type().Name())) || !viper.Sub(cc2sc(v.Type().Name())).InConfig(cc2sc(typeField.Name)) {
				isWrite = true
				viper.SetDefault(key, typeField.Tag.Get("default"))
			}
			if strings.Contains(strings.ToLower(typeField.Name), "dir") {
				createDir(viper.GetString(key))
			}
			field.SetString(viper.GetString(key))
		case reflect.Int:
			if !viper.InConfig(cc2sc(v.Type().Name())) || !viper.Sub(cc2sc(v.Type().Name())).InConfig(cc2sc(typeField.Name)) {
				isWrite = true
				i, _ := strconv.Atoi(typeField.Tag.Get("default"))
				viper.SetDefault(key, i)
			}
			field.SetInt(viper.GetInt64(key))
		case reflect.Float64:
			if !viper.InConfig(cc2sc(v.Type().Name())) || !viper.Sub(cc2sc(v.Type().Name())).InConfig(cc2sc(typeField.Name)) {
				isWrite = true
				f, _ := strconv.ParseFloat(typeField.Tag.Get("default"), 64)
				viper.SetDefault(key, f)
			}
			field.SetFloat(viper.GetFloat64(key))
		case reflect.Bool:
			if !viper.InConfig(cc2sc(v.Type().Name())) || !viper.Sub(cc2sc(v.Type().Name())).InConfig(cc2sc(typeField.Name)) {
				isWrite = true
				b, _ := strconv.ParseBool(typeField.Tag.Get("default"))
				viper.SetDefault(key, b)
			}
			field.SetBool(viper.GetBool(key))
		case reflect.Slice:
			switch reflect.New(field.Type().Elem()).Elem().Type().Kind() {
			case reflect.String:
				if !viper.InConfig(cc2sc(v.Type().Name())) || !viper.Sub(cc2sc(v.Type().Name())).InConfig(cc2sc(typeField.Name)) {
					isWrite = true
					ss := strings.Split(typeField.Tag.Get("default"), ",")
					viper.SetDefault(key, ss)
				}
				field.Set(reflect.ValueOf(viper.GetStringSlice(key)))
			case reflect.Int:
				if !viper.InConfig(cc2sc(v.Type().Name())) || !viper.Sub(cc2sc(v.Type().Name())).InConfig(cc2sc(typeField.Name)) {
					isWrite = true
					ss := strings.Split(typeField.Tag.Get("default"), ",")
					is := make([]int, len(ss))
					for i, s := range ss {
						is[i], _ = strconv.Atoi(s)
					}
					viper.SetDefault(key, is)
				}
				field.Set(reflect.ValueOf(viper.GetIntSlice(key)))
			}
		}
	}
	if isWrite {
	rewrite:
		i := 1
		err := viper.WriteConfig()
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
			i *= 2
			goto rewrite
		}
	}
}

// 创建文件夹
func createDir(path string) {
	path, _ = filepath.Abs(path)
	if _, err := os.Stat(path); err != nil {
		os.Mkdir(path, os.ModePerm)
	}
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
