package curve

type SpaceFilling interface {
	Dims() []int
	Len() int
	Curve(v []int) int
	Space(d int) []int
}
