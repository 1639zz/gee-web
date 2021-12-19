package gee

import "strings"

/**前缀树的实现
/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，
/static/css/geektutu.css匹配到/static/*filepath，解析结果为{filepath: "css/geektutu.css"}。
*/

//树结构体的节点
type node struct {
	pattern  string  //待匹配路由 例如 /p/:lang
	part     string  //路由的一部分 例子：例如 :lang
	children []*node //子节点 [a,b,c]
	isWild   bool    //是否精确匹配 part是否含有: 或者* 为true
}

//寻找第一个节点
func (n *node) matchChild(part string) *node {
	//遍历循环 查找子节点
	for _, child := range n.children {
		//如果是子节点 并且为精确匹配，返回子节点
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//查询操作
func (n *node) matchChildren(part string) []*node {
	//将nodes指在第一个位置上
	nodes := make([]*node, 0)
	//循环遍历、判断
	for _, child := range n.children {
		if child.part == part || child.isWild {
			//得到子节点的所有node值
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入操作
//递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个
func (n *node) insert(pattern string, parts []string, height int) {
	//如果长度相等，直接设置待匹配路由为参数值，并返回
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	//将路由信息放在路由组内
	part := parts[height]
	//寻找第一个子节点
	child := n.matchChild(part)
	//如果子节点为空
	if child == nil {
		//判断是否为精确匹配
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//插入子节点的值
	child.insert(pattern, parts, height+1)
}

//查询操作
func (n *node) search(parts []string, height int) *node {
	//判断长度是否相等，并且字符是否含有*
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//如果待匹配路由为空，则匹配失败
		if n.pattern == "" {
			return nil
		}
		return n
	}
	//将路由信息放在路由组内
	part := parts[height]
	//寻找所有子节点
	children := n.matchChildren(part)
	//遍历循环
	for _, child := range children {
		//得到寻找结果
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
