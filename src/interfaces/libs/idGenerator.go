package libs

type IdGeneratorInterface interface {
	Generate(prefix string) string
}
