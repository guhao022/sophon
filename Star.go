package sophon

const (
	StarFixedStar = 101 + iota
	StarPlanet
)

var styleText = map[int]string{
	StarFixedStar:		"恒星",
	StarPlanet:			  "行星",
}

// 星球接口
// 星球包括行星和恒星，恒星质量大且不宜居住，行星围绕恒星且有几率产生物种
type Star interface {
	// 发现星球
	Discovery(name string, style int)

	// 行星或恒星
	StarStyle() string

	// 星球主要形成物质
	ChemicalSubstance() map[string]map[float32]*Chemical
}

func StarText(code int) string {
	return styleText[code]
}
