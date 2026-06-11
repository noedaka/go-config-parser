package parser

type Parser interface {
	ParseConfig(data []byte) (any, error)
}
