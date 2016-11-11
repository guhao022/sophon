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

	// 获取星球基本信息
	Info(id string) *Star

	// 获取编号
	StarID() string

	// 行星或恒星
	StarStyle() int

	// 名称
	StarName() string

	// 星球主要形成物质
	ChemicalSubstance() map[float32]*Chemical
}

func StarText(code int) string {
	return styleText[code]
}

type star struct {
	// 类型
	style int

	// 唯一编号
	id string

	// 名称
	name string

	// 主要元素
	chem map[float32]*Chemical
}

func (s *star) Discovery(name string, style int) {
	s.style = StarFixedStar
	s.name = "地球"
}

func (s *star) Info(id string) *star {
	return s
}

func (s *star) StarID() string {
	return s.id
}

func (s *star) StarStyle() int {
	return s.style
}

func (s *star) StarName() string {
	return s.name
}

func (s *star) ChemicalSubstance() map[float32]*Chemical {
	return s.chem
}
