package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page struct {
	Name    string
	Content tview.Primitive
}

func demoTree() tview.Primitive {
	var box = tview.NewTreeView()
	var root = tview.NewTreeNode("G:").SetSelectable(true).SetReference("G:")
	box.SetRoot(root).SetCurrentNode(root)

	var list func(parent *tview.TreeNode)
	list = func(parent *tview.TreeNode) {
		var p = parent.GetReference().(string)
		entries, _ := os.ReadDir(p)
		for _, entry := range entries {
			var node = tview.NewTreeNode(entry.Name())
			node.SetReference(filepath.Join(p, entry.Name()))
			if entry.IsDir() {
				node.SetColor(tcell.ColorBlue)
				node.SetSelectable(true)
				node.SetSelectedFunc(func() {
					if len(node.GetChildren()) > 0 {
						node.ClearChildren()
					} else {
						list(node)
					}
				})
			} else {
				node.SetSelectable(false)
			}
			parent.AddChild(node)
		}
	}
	list(root)
	return box
}

func demoTable() tview.Primitive {
	var box = tview.NewTable()
	box.SetBorders(true)
	var titles []string = []string{"province", "city", "town"}
	for i, v := range titles {
		box.SetCell(0, i, tview.NewTableCell(v))
	}
	for r := 1; r < 5; r++ {
		for c := 0; c < 3; c++ {
			box.SetCellSimple(r, c, fmt.Sprintf("%d_%d", r, c))
		}
	}
	return box
}

func demoGrid() tview.Primitive {
	var box = tview.NewGrid()
	box.SetTitle("box")
	box.SetBorders(true)
	box.SetColumns(8, 0, 9)
	box.SetRows(1, 0, 1)
	box.AddItem(
		tview.NewTextView().SetText("header").SetTextAlign(tview.AlignCenter),
		0, 0,
		1, 3,
		0, 0,
		false,
	)
	box.AddItem(
		tview.NewTextView().SetText("left bar").SetTextAlign(tview.AlignCenter),
		1, 0,
		1, 1,
		0, 0,
		false,
	)
	box.AddItem(
		tview.NewTextView().SetText("content").SetTextAlign(tview.AlignCenter),
		1, 1,
		1, 1,
		0, 0,
		false,
	)
	box.AddItem(
		tview.NewTextView().SetText("right bar").SetTextAlign(tview.AlignCenter),
		1, 2,
		1, 1,
		0, 0,
		false,
	)
	box.AddItem(
		tview.NewTextView().SetText("X").SetTextAlign(tview.AlignCenter),
		2, 0,
		1, 1,
		0, 0,
		false,
	)
	box.AddItem(
		tview.NewTextView().SetText("footer").SetTextAlign(tview.AlignCenter),
		2, 1,
		1, 2,
		0, 0,
		false,
	)
	return box
}

func demoModal() tview.Primitive {
	var box = tview.NewModal()
	box.SetTitle("tips")
	box.SetText("this is a tip")
	box.SetBorder(true)
	box.AddButtons([]string{"yes", "no", "cancel"})
	box.SetBlurFunc(func() {})
	box.SetDoneFunc(func(buttonIndex int, buttonLabel string) {

	})
	return box
}

func demoFrom() tview.Primitive {
	var box = tview.NewForm().
		AddInputField("username", "", 0, func(textToCheck string, lastChar rune) bool { return true }, func(text string) {}).
		AddPasswordField("password", "", 0, '*', func(text string) {}).
		AddDropDown("sex", []string{"male", "female", "other"}, 0, func(option string, optionIndex int) {}).
		AddTextArea("desc", "", 0, 0, 1000, func(text string) {}).
		AddCheckbox("privacy", false, func(checked bool) {}).
		AddButton("cancel", func() {
		}).
		AddButton("confirm", func() {
		})
	box.SetTitle("Login")
	return box
}

func demoList() tview.Primitive {
	var box = tview.NewList()
	box.SetTitle("hometown")
	box.AddItem("hb", "hubei", 'a', nil)
	box.AddItem("gd", "guangdong", 'b', nil)
	box.AddItem("bj", "beijing", 'c', nil)
	return box
}

func main() {
	var app = tview.NewApplication().EnableMouse(true)
	var title = "press PgUp/PgDn to switch"
	var index = 0

	var body = tview.NewPages()
	var pages []Page = []Page{
		{"form", demoFrom()},
		{"list", demoList()},
		{"modal", demoModal()},
		{"grid", demoGrid()},
		{"table", demoTable()},
		{"tree", demoTree()},
	}
	for _, page := range pages {
		body.AddPage(page.Name, page.Content, true, true)
	}
	body.SetTitle(title)
	body.SetBorder(true)
	body.SetChangedFunc(func() {
		body.SetTitle(fmt.Sprintf(" %s (%s) ", title, pages[index].Name))
	})

	body.SwitchToPage(pages[0].Name)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyPgDn {
			if index < len(pages)-1 {
				index = index + 1
				body.SwitchToPage(pages[index].Name)
			}
			return nil
		} else if event.Key() == tcell.KeyPgUp {
			if index > 0 {
				index = index - 1
				body.SwitchToPage(pages[index].Name)
			}
			return nil
		}
		return event
	})
	app.SetRoot(body, true).Run()
}
