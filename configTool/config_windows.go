package configTool

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/spf13/viper"
)

var wd = pflag.String("WorkingDirectory", ".", "工作目录")

func init() {
	pflag.Parse()
	viper.SetConfigName("config")
	viper.AddConfigPath(*wd)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfigAs(filepath.Join(*wd, "config.toml"))
		} else {
			fmt.Println(err)
		}
	}
}
