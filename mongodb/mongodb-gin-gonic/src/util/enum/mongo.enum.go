package enum

type MongoCollection int

const (
	MongoCollection_Items MongoCollection = iota
)

func (index MongoCollection) String() string {
	return []string{
		"items",
	}[index]
}
