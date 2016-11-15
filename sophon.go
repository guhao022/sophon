package sophon

type Sophon struct {
	//
}

// 设置集群节点
var node *Node

// 初始化
func init()  {

	n, err := NewNode(1)

	if err != nil {
		CLog("[ERRO]", err)
	}

	node = n
}

