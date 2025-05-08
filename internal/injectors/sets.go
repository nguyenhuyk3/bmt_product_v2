package injectors

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/internal/implementations/message_broker/writers"
	"bmt_product_service/internal/implementations/redis"
	"bmt_product_service/internal/implementations/upload"
	"bmt_product_service/internal/injectors/provider"

	"github.com/google/wire"
)

var uploadServiceSet = wire.NewSet(
	upload.NewUploadService,
)

var kafkaWriterSet = wire.NewSet(
	writers.NewKafkaWriter,
)

var dbSet = wire.NewSet(
	provider.ProvidePgxPool,
	writers.NewKafkaWriter,
	sqlc.NewStore,
)

var redisSet = wire.NewSet(
	redis.NewRedisClient,
)
