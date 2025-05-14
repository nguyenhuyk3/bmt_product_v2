package initializations

import (
	"bmt_product_service/global"
	"fmt"
)

func Run() {
	loadConfigs()
	initPostgreSql()
	initRedis()
	initProductUploadReader()

	r := initRouter()

	r.Run(fmt.Sprintf("0.0.0.0:%s", global.Config.Server.ServerPort))
}
