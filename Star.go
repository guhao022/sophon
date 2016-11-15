package sophon

const (
	StarFixedStar = 101 + iota
	StarPlanet
)

var styleText = map[int]string{
	StarFixedStar: "恒星",
	StarPlanet:    "行星",
}

// 星球接口
// 星球包括行星和恒星，恒星质量大且不宜居住，行星围绕恒星且有几率产生物种
type Star interface {
	// 发现星球
	Discovery(name string, style int)

	// 获取编号
	StarID() int64

	// 行星或恒星
	StarStyle() int

	// 名称
	StarName() string
}

func StarText(code int) string {
	return styleText[code]
}

type star struct {
	// 类型
	Style int    `json:"style"`

	// 唯一编号
	ID int64     `json:"id"`

	// 名称
	Name string    `json:"name"`
}

func NewStar() Star {
	id := node.Generate()

	s := new(star)
	s.ID = id.Int64()

	return s

}

func (s *star) Discovery(name string, style int) {
	s.Style = style
	s.Name = name
}

func (s *star) StarID() int64 {
	return s.ID
}

func (s *star) StarStyle() int {
	return s.Style
}

func (s *star) StarName() string {
	return s.Name
}

