package mongoDb

import "grpc-gateway-project/internal/services"

type ConnectionMongo struct {
	MongoHost       string
	MongoPort       string
	MongoDBName     string
	MongoCollection string
}

func NewConfig() *ConnectionMongo {
	return &ConnectionMongo{
		MongoHost:       services.GetEnv("HOST_MONGO", "mongo"),
		MongoPort:       services.GetEnv("PORT_MONGO", "27017"),
		MongoDBName:     services.GetEnv("DBNAME_MONGO", "gateway"),
		MongoCollection: services.GetEnv("COLLECTION_MONGO", "users"),
	}
}
