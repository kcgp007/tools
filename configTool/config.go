package configTool

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/kcgp007/tools/pathTool"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(pathTool.SmartWd())
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfigAs(filepath.Join(pathTool.SmartWd(), "config.toml"))
		} else {
			fmt.Println(err)
		}
	}
}

// 添加配置
func Add(ps ...interface{}) {
	for _, p := range ps {
		config(p, nil)
	}
	viper.WriteConfig()
}

// 加载默认配置及配置文件数据
func config(p interface{}, keys []string) {
	v := reflect.ValueOf(p).Elem()
	if keys == nil {
		keys = append(keys, v.Type().Name())
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)
		tmpKeys := append(keys, typeField.Name)
		switch typeField.Type.Kind() {
		case reflect.Struct:
			config(field, tmpKeys)
		case reflect.String:
			if strings.Contains(strings.ToLower(typeField.Name), "dir") {
				createDir(viper.GetString(genKey(tmpKeys...)))
			}
			viper.SetDefault(genKey(tmpKeys...), typeField.Tag.Get("default"))
			field.SetString(pathTool.SmartAbs(viper.GetString(genKey(tmpKeys...))))
		case reflect.Int:
			i, _ := strconv.Atoi(typeField.Tag.Get("default"))
			viper.SetDefault(genKey(tmpKeys...), i)
			field.SetInt(viper.GetInt64(genKey(tmpKeys...)))
		case reflect.Float64:
			f, _ := strconv.ParseFloat(typeField.Tag.Get("default"), 64)
			viper.SetDefault(genKey(tmpKeys...), f)
			field.SetFloat(viper.GetFloat64(genKey(tmpKeys...)))
		case reflect.Bool:
			b, _ := strconv.ParseBool(typeField.Tag.Get("default"))
			viper.SetDefault(genKey(tmpKeys...), b)
			field.SetBool(viper.GetBool(genKey(tmpKeys...)))
		case reflect.Slice:
			switch reflect.New(field.Type().Elem()).Elem().Type().Kind() {
			case reflect.String:
				ss := strings.Split(typeField.Tag.Get("default"), ",")
				viper.SetDefault(genKey(tmpKeys...), ss)
				field.Set(reflect.ValueOf(viper.GetStringSlice(genKey(tmpKeys...))))
			case reflect.Int:
				ss := strings.Split(typeField.Tag.Get("default"), ",")
				is := make([]int, len(ss))
				for i, s := range ss {
					is[i], _ = strconv.Atoi(s)
				}
				viper.SetDefault(genKey(tmpKeys...), is)
				field.Set(reflect.ValueOf(viper.GetIntSlice(genKey(tmpKeys...))))
			}
		default:
			switch typeField.Type.String() {
			case "time.Duration":
				i, _ := strconv.Atoi(typeField.Tag.Get("default"))
				viper.SetDefault(genKey(tmpKeys...), i)
				field.Set(reflect.ValueOf(viper.GetDuration(genKey(tmpKeys...))))
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
func genKey(keys ...string) string {
	key := ""
	for _, s := range keys {
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
