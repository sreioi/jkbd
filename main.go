package main

import (
	"github.com/gookit/color"
	flag "github.com/spf13/pflag"
	"github.com/sreioi/jkbd/db"
	"github.com/sreioi/jkbd/env"
	"github.com/sreioi/jkbd/ke1"
	"github.com/sreioi/jkbd/ke4"
	"os"
)

func init() {
	db.ConnDatabase() // 初始化数据库
}

func main() {
	flag.StringVarP(&env.Type, "type", "t", env.Type, "拉取数据类型[k1,k4]")
	flag.StringVarP(&env.ID, "id", "i", env.ID, "拉取问题ID")
	flag.IntVarP(&env.Worker, "worker", "w", env.Worker, "默认5个协程")
	flag.Parse()

	switch env.Type {
	case "k1":
		ke1.CreateTableK1()
		ke1.Pull()
	case "k4":
		ke4.CreateTableK4()
		ke4.Pull()
	default:
		color.Redln("拉取数据类型错误")
		os.Exit(1)
	}
}
