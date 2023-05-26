package enum

type ElasticIndex int

const (
	ElasticIndex_Movies ElasticIndex = iota
)

func (index ElasticIndex) String() string {
	return []string{
		"movies",
	}[index]
}
