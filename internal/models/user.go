package models

type User struct {
	Id            string `bson:"_id,omitempty"`
	Email         string
	Name          string
	Age           int32
	Subscribers   []string
	Subscriptions []string
	Token         string
}

type SubscribeEvent struct {
	SubscriberId string
	ListenerId   string
}
