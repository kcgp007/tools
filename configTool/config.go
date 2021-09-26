package configTool

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"tools/flagInit"
	"unicode"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var ps []interface{}

func init() {
	pflag.Parse()
	viper.SetConfigName("config")
	viper.AddConfigPath(*flagInit.ConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		if *flagInit.IsCompletion {
			return
		}
		fmt.Println(err)
		*flagInit.IsCompletion = true
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		for _, p := range ps {
			read(p)
		}
	})
}

// 添加配置
func AddConfig(p ...interface{}) {
	ps = append(ps, p...)
}

// 添加配置完成
func Done() {
	if *flagInit.IsCompletion {
		if err := viper.WriteConfigAs(*flagInit.ConfigPath + "/config.yml"); err != nil {
			fmt.Println(err)
		}
		for _, p := range ps {
			setDefault(p)
			read(p)
		}
		if err := viper.WriteConfig(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("已补全config文件，请重新启动程序，输入任意键退出...")
		fmt.Scanln()
		*flagInit.IsCompletion = true
		os.Exit(0)
	} else {
		for _, p := range ps {
			read(p)
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
		default:
			switch typeField.Type.String() {
			case "time.Duration":
				i, _ := strconv.Atoi(typeField.Tag.Get("default"))
				viper.SetDefault(genKey(v.Type().Name(), typeField.Name), i)
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
		key := genKey(v.Type().Name(), typeField.Name)
		switch typeField.Type.Kind() {
		case reflect.String:
			if strings.Contains(strings.ToLower(typeField.Name), "dir") {
				createDir(viper.GetString(key))
			}
			field.SetString(viper.GetString(key))
		case reflect.Int:
			field.SetInt(viper.GetInt64(key))
		case reflect.Float64:
			field.SetFloat(viper.GetFloat64(key))
		case reflect.Bool:
			field.SetBool(viper.GetBool(key))
		case reflect.Slice:
			switch reflect.New(field.Type().Elem()).Elem().Type().Kind() {
			case reflect.String:
				field.Set(reflect.ValueOf(viper.GetStringSlice(key)))
			case reflect.Int:
				field.Set(reflect.ValueOf(viper.GetIntSlice(key)))
			}
		default:
			switch typeField.Type.String() {
			case "time.Duration":
				field.Set(reflect.ValueOf(viper.GetDuration(key)))
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
