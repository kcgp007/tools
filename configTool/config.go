package configTool

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/viper"
)

type log struct {
	Level  string        `default:"info"`
	MaxAge time.Duration `default:"30"`
}

var Log log

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfigAs("./config.toml")
		} else {
			fmt.Println(err)
		}
	}
	Add(Log)
}

// 添加配置
func Add(ps ...interface{}) {
	for _, p := range ps {
		config(p, nil)
	}
	viper.WriteConfig()
}

// 加载默认配置及配置文件数据
func config(p interface{}, ss []string) {
	v := reflect.ValueOf(p).Elem()
	if ss == nil {
		ss = append(ss, v.Type().Name())
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)
		ss = append(ss, typeField.Name)
		switch typeField.Type.Kind() {
		case reflect.Struct:
			config(field, ss)
		case reflect.String:
			viper.SetDefault(genKey(ss...), typeField.Tag.Get("default"))
			if strings.Contains(strings.ToLower(typeField.Name), "dir") {
				createDir(viper.GetString(genKey(ss...)))
			}
			field.SetString(viper.GetString(genKey(ss...)))
		case reflect.Int:
			i, _ := strconv.Atoi(typeField.Tag.Get("default"))
			viper.SetDefault(genKey(ss...), i)
			field.SetInt(viper.GetInt64(genKey(ss...)))
		case reflect.Float64:
			f, _ := strconv.ParseFloat(typeField.Tag.Get("default"), 64)
			viper.SetDefault(genKey(ss...), f)
			field.SetFloat(viper.GetFloat64(genKey(ss...)))
		case reflect.Bool:
			b, _ := strconv.ParseBool(typeField.Tag.Get("default"))
			viper.SetDefault(genKey(ss...), b)
			field.SetBool(viper.GetBool(genKey(ss...)))
		case reflect.Slice:
			switch reflect.New(field.Type().Elem()).Elem().Type().Kind() {
			case reflect.String:
				ss := strings.Split(typeField.Tag.Get("default"), ",")
				viper.SetDefault(genKey(ss...), ss)
				field.Set(reflect.ValueOf(viper.GetStringSlice(genKey(ss...))))
			case reflect.Int:
				ss := strings.Split(typeField.Tag.Get("default"), ",")
				is := make([]int, len(ss))
				for i, s := range ss {
					is[i], _ = strconv.Atoi(s)
				}
				viper.SetDefault(genKey(ss...), is)
				field.Set(reflect.ValueOf(viper.GetIntSlice(genKey(ss...))))
			}
		default:
			switch typeField.Type.String() {
			case "time.Duration":
				i, _ := strconv.Atoi(typeField.Tag.Get("default"))
				viper.SetDefault(genKey(ss...), i)
				field.Set(reflect.ValueOf(viper.GetDuration(genKey(ss...))))
			}
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
