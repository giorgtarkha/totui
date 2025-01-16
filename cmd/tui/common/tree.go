package common

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TreeNodeSelectCallback func() error

type TreeNode struct {
	label              string
	expanded           bool
	leftSiblingOffset  int
	rightSiblingOffset int
	level              int
	onSelect           TreeNodeSelectCallback // TODO is the error needed at all
}

type TreeNodeParams struct {
	Label              string
	LeftSiblingOffset  int
	RightSiblingOffset int
	Level              int
	OnSelect           TreeNodeSelectCallback // TODO is the error needed at all
}

func NewTreeNode(p *TreeNodeParams) TreeNode {
	return TreeNode{
		label:              p.Label,
		leftSiblingOffset:  p.LeftSiblingOffset,
		rightSiblingOffset: p.RightSiblingOffset,
		level:              p.Level,
		expanded:           false,
		onSelect:           p.OnSelect,
	}
}

func (n *TreeNode) getBaseStyle() lipgloss.Style {
	return lipgloss.NewStyle().PaddingLeft(n.level * 2)
}

func (n *TreeNode) render(selected bool) string {
	prefix := "▸"
	if n.expanded {
		prefix = "▾"
	}
	if n.rightSiblingOffset <= 1 {
		prefix = " "
	}

	color := lipgloss.Color("205")
	if selected {
		color = lipgloss.Color("150")
	}
	style := n.getBaseStyle().Foreground(color)

	return style.Render(fmt.Sprintf("%s %s", prefix, n.label))
}

type Tree struct {
	nodes        []TreeNode
	currentIndex int
	maxWidth     int
}

type TreeParams struct {
	Nodes []TreeNode
}

func NewTree(p *TreeParams) *Tree {
	maxWidth := 0
	for _, n := range p.Nodes {
		width := lipgloss.Width(
			n.getBaseStyle().Render(fmt.Sprintf("  %s", n.label)),
		)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return &Tree{
		nodes:        p.Nodes,
		currentIndex: 0,
		maxWidth:     maxWidth,
	}
}

func (t *Tree) Init() tea.Cmd {
	return nil
}

func (t *Tree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if t.currentIndex > 0 {
				offset := t.nodes[t.currentIndex].leftSiblingOffset
				if offset == 0 || t.nodes[t.currentIndex-offset].expanded {
					offset = 1
				}
				if t.currentIndex-offset >= 0 {
					t.currentIndex -= offset
					// TODO How to handle error returned from callback, if needed at all
					if t.nodes[t.currentIndex].onSelect != nil {
						t.nodes[t.currentIndex].onSelect()
					}
				}
			}
		case "down":
			if t.currentIndex < len(t.nodes)-1 {
				offset := t.nodes[t.currentIndex].rightSiblingOffset
				if offset == 0 || t.nodes[t.currentIndex].expanded {
					offset = 1
				}
				if t.currentIndex+offset < len(t.nodes) {
					t.currentIndex += offset
					// TODO How to handle error returned from callback, if needed at all
					if t.nodes[t.currentIndex].onSelect != nil {
						t.nodes[t.currentIndex].onSelect()
					}
				}
			}
		case "enter":
			if t.nodes[t.currentIndex].rightSiblingOffset >= 1 {
				t.nodes[t.currentIndex].expanded = !t.nodes[t.currentIndex].expanded
			}
		}

	}

	return t, nil
}

func (t *Tree) View() string {
	sb := strings.Builder{}
	index := 0
	for index < len(t.nodes) {
		s := t.nodes[index].render(index == t.currentIndex)
		sb.WriteString(s + "\n")

		offset := 0
		if t.nodes[index].expanded {
			offset = 1
		} else {
			offset = t.nodes[index].rightSiblingOffset
		}
		if offset == 0 {
			offset = 1
		}
		index += offset
	}
	return sb.String()
}

func (t *Tree) GetMaxWidth() int {
	return t.maxWidth
}

type RCTreeNode struct {
	label    string
	children []*RCTreeNode
	onSelect TreeNodeSelectCallback
}

type RCTreeNodeParams struct {
	Label    string
	Children []*RCTreeNode
	OnSelect TreeNodeSelectCallback
}

func NewRCTreeNode(p *RCTreeNodeParams) *RCTreeNode {
	return &RCTreeNode{
		label:    p.Label,
		children: p.Children,
		onSelect: p.OnSelect,
	}
}

func (r *RCTreeNode) clearNils() {
	nonnil := []*RCTreeNode{}
	for _, n := range r.children {
		if n == nil {
			continue
		}
		n.clearNils()
		nonnil = append(nonnil, n)
	}
	r.children = nonnil
}

func (r *RCTreeNode) toTreeNode(
	leftSiblingOffset int,
	rightSiblingOffset int,
	level int,
) ([]TreeNode, int) {
	res := []TreeNode{
		TreeNode{
			label:              r.label,
			expanded:           false,
			leftSiblingOffset:  leftSiblingOffset,
			rightSiblingOffset: rightSiblingOffset,
			level:              level,
			onSelect:           r.onSelect,
		},
	}
	i := 0
	childCount := 0
	totalCount := 0
	var cres []TreeNode
	for i < len(r.children) {
		clo := 0
		if i > 0 {
			clo = childCount
		}
		cres, childCount = r.children[i].toTreeNode(clo, 0, level+1)
		totalCount += childCount
		cres[0].rightSiblingOffset = childCount
		res = append(res, cres...)
		i++
	}

	return res, totalCount + 1
}

func (r *RCTreeNode) ToTree() *Tree {
	r.clearNils()
	nodes, cc := r.toTreeNode(0, 0, 0)
	nodes[0].rightSiblingOffset = cc

	// TODO remove, for now this is good for debugging
	/*i := 0
	for i < len(nodes) {
		nodes[i].label = fmt.Sprintf("%s %d %d %d", nodes[i].label, nodes[i].leftSiblingOffset, nodes[i].rightSiblingOffset, nodes[i].level)
		i++
	}*/

	return NewTree(&TreeParams{
		Nodes: nodes,
	})
}
