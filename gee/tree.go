package gee

import "strings"

// user/login/name,user/login/id,user/login/a :*
type Node struct {
	Pattern  string  // 准确的 路径///user/login
	Part     string  //
	Children []*Node // 子节点
	IsWild   bool    // 是否为通配符
}

// 返回在n中查找part字段.,有返回该节点,没有返回空
func (n *Node) MatchChild(part string) *Node {
	for _, child := range n.Children {
		if child.Part == part || child.IsWild {
			return child
		}
	}
	return nil
}

// 查找返回与所有节点;查找n中part字段,返回节点数组(与MatchChild的区别:本函数返回所有匹配的节点而MatchChild返回第一个匹配的节点)
func (n *Node) MatchChildren(part string) []*Node {
	nodes := make([]*Node, 0)          //建立一个新节点
	for _, child := range n.Children { //遍历n中是否有panrt的相对或者绝对路径
		if child.Part == part || child.IsWild { //如果存在或者精确匹配
			nodes = append(nodes, child) //节点赋值
		}
	}
	return nodes //返回节点数组,没有就是空节点
}

// 如果当前处理的 URL 路径部分已经到了最后一部分（即 len(parts) == height），那么说明已经找到了该路由规则对应的叶子节点，将该叶子节点的模式设置为该规则对应的处理器函数 pattern，然后直接返回。

// 如果当前节点的子节点中有一个子节点的 part 值和当前处理的 URL 路径部分相等，或者是一个通配符（即 child.part == part || child.isWild），那么说明这个子节点可以用于匹配当前 URL 路径部分，就把 height 加一，递归调用 Insert 方法，将处理的 URL 路径部分的索引值向后移动一个位置。

// 如果找不到符合条件的子节点，说明当前节点的子节点中没有一个子节点能够匹配当前处理的 URL 路径部分，就创建一个新的子节点，将该子节点插入到当前节点的子节点列表中，然后递归调用 Insert 方法。

// 插入节点:pattern 为路由,路由的parts为字符串切片,,hight为树的深度.
// pattern:  user/login/*id
func (n *Node) Insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.Pattern = pattern
		return
	}

	part := parts[height]
	child := n.MatchChild(part)
	//如果没有路径创建新路径
	if child == nil {
		child = &Node{Part: part, IsWild: part[0] == ':' || part[0] == '*'}
		n.Children = append(n.Children, child)
	}
	child.Insert(pattern, parts, height+1)
}

// 如果当前处理的 URL 路径部分已经到了最后一部分（即 len(parts) == height），或者当前节点是一个以 * 开头的通配符节点，那么说明已经找到了匹配的叶子节点，如果该叶子节点的模式不为空，就返回该叶子节点，否则返回 nil。

// 如果当前节点的子节点中有一个子节点的 part 值和当前处理的 URL 路径部分相等，或者是一个通配符（即 child.part == part || child.isWild），那么说明这个子节点可以用于匹配当前 URL 路径部分，就把 height 加一，递归调用 Search 方法，将处理的 URL 路径部分的索引值向后移动一个位置。

// 如果找不到符合条件的子节点，说明当前节点的子节点中没有一个子节点能够匹配当前处理的 URL 路径部分，就返回 nil。
//
//pattern:user/login/*id(/name,/time) ==>parts{user,login,123}
func (n *Node) Search(parts []string, height int) *Node {
	if len(parts) == height || strings.HasPrefix(n.Part, "*") {
		if n.Pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.MatchChildren(part)

	for _, child := range children {
		result := child.Search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

//parts 是一个字符串切片，表示一个 URL 路径经过分割后的每一个部分。比如，对于一个 URL 路径 /users/:id/books/*title，经过分割后的 parts 内容应该是 []string{"users", ":id", "books", "*title"}。
