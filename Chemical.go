package sophon

type ChemicalNum int

const (
	// 元素周期表
	H ChemicalNum = 1001 + iota
	He
	Li
	Be
	B
	C
	N
	O
	F
	Ne
	Na
	Mg
	Al
	Si
	P
	S
	Cl
	Ar
	K
	Ca
	Sc
	Ti
	V
	Cr
	Mn
	Fe
	Co
	Ni
	Cu
	Zn
	Ga
	Ge
	As
	Se
	Br
	Kr
	Rb
	Sr
	Y
	Zr
	Nb
	Mo
	Tc
	Ru
	Rh
	Pd
	Ag
	Cd
	In
	Sn
	Sb
	Te
	I
	Xe
	Cs
	Ba
	La
	Ce
	Pr
	Nd
	Pm
	Sm
	En
	Gd
	Tb
	Dy
	Ho
	Er
	Tm
	Yb
	Lu
	Hf
	Ta
	W
	Re
	Os
	Ir
	Pt
	Au
	Hg
	Tl
	Pb
	Bi
	Po
	At
	Rn
	Fr
	Ra
	Ac
	Th
	Pa
	U
	Np
	Pu
	Am
	Cm
	Bk
	Cf
	Es
	Fm
	Md
	No
	Lr
	Rf
	Db
	Sg
	Bh
	Hs
	Mt
	Ds
	Rg
	Cn
	Nh
	Fl
	Mc
	Lv
	Ts
	Og
)

// 元素周期数关联名称
var chemText = map[ChemicalNum]string{
	H:  "氢",
	He: "氦",
	Li: "锂",
	Be: "铍",
	B:  "硼",
	C:  "碳",
	N:  "氮",
	O:  "氧",
	F:  "氟",
	Ne: "氖",
	Na: "钠",
	Mg: "镁",
	Al: "铝",
	Si: "硅",
	P:  "磷",
	S:  "硫",
	Cl: "氯",
	Ar: "氩",
	K:  "钾",
	Ca: "钙",
	Sc: "钪",
	Ti: "钛",
	V:  "钡",
	Cr: "铬",
	Mn: "锰",
	Fe: "铁",
	Co: "钴",
	Ni: "镍",
	Cu: "铜",
	Zn: "锌",
	Ga: "镓",
	Ge: "锗",
	As: "砷",
	Se: "硒",
	Br: "溴",
	Kr: "氪",
	Rb: "铷",
	Sr: "锶",
	Y:  "钇",
	Zr: "钴",
	Nb: "铌",
	Mo: "钼",
	Tc: "锝",
	Ru: "钌",
	Rh: "铑",
	Pd: "钯",
	Ag: "银",
	Cd: "镉",
	In: "铟",
	Sn: "锡",
	Sb: "锑",
	Te: "碲",
	I:  "碘",
	Xe: "氙",
	Cs: "铯",
	Ba: "钡",
	La: "镧",
	Ce: "铈",
	Pr: "镨",
	Nd: "钕",
	Pm: "钷",
	Sm: "钐",
	En: "铕",
	Gd: "钆",
	Tb: "铽",
	Dy: "镝",
	Ho: "钬",
	Er: "铒",
	Tm: "铥",
	Yb: "镱",
	Lu: "镥",
	Hf: "铪",
	Ta: "钽",
	W:  "钨",
	Re: "铼",
	Os: "锇",
	Ir: "铱",
	Pt: "铂",
	Au: "金",
	Hg: "汞",
	Tl: "铊",
	Pb: "铅",
	Bi: "铋",
	Po: "钋",
	At: "砹",
	Rn: "氡",
	Fr: "钫",
	Ra: "镭",
	Ac: "锕",
	Th: "钍",
	Pa: "镤",
	U:  "铀",
	Np: "镎",
	Pu: "钚",
	Am: "镅",
	Cm: "锔",
	Bk: "锫",
	Cf: "锎",
	Es: "镶",
	Fm: "镄",
	Md: "钔",
	No: "锘",
	Lr: "铹",
	Rf: "钅卢",
	Db: "钅杜",
	Sg: "钅喜",
	Bh: "钅波",
	Hs: "钅黑",
	Mt: "钅麦",
	Ds: "鐽",
	Rg: "錀",
	Cn: "鎶",
	Nh: "鈤（推测）",
	Fl: "钅夫",
	Mc: "镆（推测）",
	Lv: "钅立",
	Ts: "钿（推测）",
	Og: "奥气（推测）",
}

func ChemText(code int) string {
	return chemText[ChemicalNum(code)]
}

// 生成化学物质的函数
type Contain func() map[int]ChemicalNum

// 基本物质
type Chemical interface {
	Generate(f Contain, name string)  // 生成基本物质
	ID() string                       // 获取基本物质的编号
	Name() string                     // 获取基本物质的名称
	Composition() map[int]ChemicalNum // 基本物质的化学成分
}

// 默认实现化学物质
type chem struct {
	id       string              // 物质编号
	name     string              // 物质名称
	comprise map[int]ChemicalNum // 物质化学成分
}

func New() Chemical {
	c := new(chem)

	c.id = "1"

	return c
}

// 生成基本物质
func (c *chem) Generate(f Contain, name string) {
	c.comprise = f()
	c.name = name
}

// 获取元素符号
func (c *chem) ID() string {
	return c.id
}

// 获取元素名称
func (c *chem) Name() string {
	return c.name
}

// 获取元素成分
func (c *chem) Composition() map[int]ChemicalNum {
	return c.comprise
}
