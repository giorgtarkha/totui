package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/giorgtarkha/totui/tui/common"
)

type CLITUI struct {
	cmdTree *common.Tree
}

func (c *CLITUI) Init() tea.Cmd {
	root := common.NewRCTreeNode(&common.RCTreeNodeParams{
		Label: "root",
		Children: []*common.RCTreeNode{
			common.NewRCTreeNode(&common.RCTreeNodeParams{
				Label: "child1",
				Children: []*common.RCTreeNode{
					common.NewRCTreeNode(&common.RCTreeNodeParams{
						Label: "child11",
					}),
					common.NewRCTreeNode(&common.RCTreeNodeParams{
						Label: "child12",
					}),
				},
			}),
			common.NewRCTreeNode(&common.RCTreeNodeParams{
				Label: "child2",
				Children: []*common.RCTreeNode{
					common.NewRCTreeNode(&common.RCTreeNodeParams{
						Label: "child21",
					}),
					common.NewRCTreeNode(&common.RCTreeNodeParams{
						Label: "child22",
						Children: []*common.RCTreeNode{
							common.NewRCTreeNode(&common.RCTreeNodeParams{
								Label: "child221",
							}),
							common.NewRCTreeNode(&common.RCTreeNodeParams{
								Label: "child222",
							}),
						},
					}),
					common.NewRCTreeNode(&common.RCTreeNodeParams{
						Label: "child23",
						Children: []*common.RCTreeNode{
							common.NewRCTreeNode(&common.RCTreeNodeParams{
								Label: "child231",
							}),
							common.NewRCTreeNode(&common.RCTreeNodeParams{
								Label: "child232",
							}),
							common.NewRCTreeNode(&common.RCTreeNodeParams{
								Label: "child233",
							}),
						},
					}),
					common.NewRCTreeNode(&common.RCTreeNodeParams{
						Label: "child24",
					}),
				},
			}),
			common.NewRCTreeNode(&common.RCTreeNodeParams{
				Label:    "child3",
				Children: []*common.RCTreeNode{},
			}),
			common.NewRCTreeNode(&common.RCTreeNodeParams{
				Label: "child4",
				Children: []*common.RCTreeNode{
					common.NewRCTreeNode(&common.RCTreeNodeParams{
						Label: "child41",
					}),
					common.NewRCTreeNode(&common.RCTreeNodeParams{
						Label: "child42",
					}),
				},
			}),
		},
	})
	c.cmdTree = root.ToTree()
	return nil
}

func (c *CLITUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		}
	}

	_, cmd := c.cmdTree.Update(msg)
	return c, cmd
}

func (c *CLITUI) View() string {
	return c.renderCmdTree()
}

func (c *CLITUI) renderCmdTree() string {
	width := c.cmdTree.GetMaxWidth() + 2
	headerStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(width)
	headerText := "Commands\n"

	treeStyle := lipgloss.NewStyle().
		Align(lipgloss.Left, lipgloss.Top).
		PaddingLeft(1)

	header := headerStyle.Render(headerText)
	tree := treeStyle.Render(c.cmdTree.View())

	combined := lipgloss.JoinVertical(lipgloss.Top, header, tree)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Left, lipgloss.Top).
		Width(width)

	return borderStyle.Render(combined)
}
