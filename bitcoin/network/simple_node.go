package network

type SimpleNode struct {
	host      string
	port      int
	isTestnet bool
	isLogging bool
}

func NewSimpleNode(host string, port int, isTestnet bool, isLogging bool) *SimpleNode {
	return &SimpleNode{host, port, isTestnet, isLogging}
}
