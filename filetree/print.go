package filetree

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	SPACE_ONE       = " "
	SPACE_THREE     = "   "
	HORIZONTAL_LINE = "─"
	VERTICAL_LINE   = "│"
	T_PREFIX        = "├"
	END_PREFIX      = "└"
)

func PrintFileNodeTree(node *FileNode, prefix []string, depth, level int, human, sha256, imo bool) {
	for i, x := range node.Children {
		// first
		if i == 0 {
			prefix = append(prefix, T_PREFIX)
		}

		// last
		if i+1 == len(node.Children) {
			prefix = append(prefix[:len(prefix)-1], END_PREFIX)
		}

		if x.Type == TYPE_FILE {
			bin := []byte{}
			var e error
			if sha256 {
				bin, e = ioutil.ReadFile(x.AbsPath)
				if e != nil {
					panic(e)
				}
			}
			if sha256 && imo {
				fmt.Printf("%v%v %v %v %v %v\n", strings.Join(prefix, ""), strings.Repeat(HORIZONTAL_LINE, 2), x.Name, getSizeText(x.TotalSize, human), Sha256(bin), ImoHash(x.AbsPath))
			} else if sha256 {
				fmt.Printf("%v%v %v %v %v\n", strings.Join(prefix, ""), strings.Repeat(HORIZONTAL_LINE, 2), x.Name, getSizeText(x.TotalSize, human), Sha256(bin))
			} else if imo {
				fmt.Printf("%v%v %v %v %v\n", strings.Join(prefix, ""), strings.Repeat(HORIZONTAL_LINE, 2), x.Name, getSizeText(x.TotalSize, human), ImoHash(x.AbsPath))
			} else {
				fmt.Printf("%v%v %v %v\n", strings.Join(prefix, ""), strings.Repeat(HORIZONTAL_LINE, 2), x.Name, getSizeText(x.TotalSize, human))
			}
		} else {
			fmt.Printf("%v%v %v %v\n", strings.Join(prefix, ""), strings.Repeat(HORIZONTAL_LINE, 2), x.Name, getSizeText(x.TotalSize, human))
		}
		if x.Type == TYPE_DIR && depth < level {
			if x.Name[0] == '.' {
				continue
			}

			// last
			if i+1 == len(node.Children) {
				prefix = append(prefix[:len(prefix)-1], SPACE_ONE)
			} else {
				prefix = append(prefix[:len(prefix)-1], VERTICAL_LINE)
			}

			prefix = append(prefix, SPACE_THREE)
			PrintFileNodeTree(x, prefix, depth+1, level, human, sha256, imo)
			prefix = append(prefix[:len(prefix)-2], T_PREFIX)
		}
	}
}

func getSizeText(s int64, human bool) string {
	sizeTtext := fmt.Sprintf("%v", s)
	if human {
		sizeTtext = HumanSize(s, 1024)
	}
	return sizeTtext
}

func flat(node *FileNode) []*FileNode {
	list := []*FileNode{node}
	if node.Name[0] == '.' {
		return list
	}

	for _, x := range node.Children {
		if x.Type == TYPE_DIR {
			for _, y := range flat(x) {
				list = append(list, y)
			}
		} else {
			list = append(list, x)
		}
	}
	return list
}

func PrintFileNodeSimple(node *FileNode, human bool) {
	originChildren := node.Children
	node.Children = flat(node)
	node.Sort()
	defer func() { node.Children = originChildren }()

	for _, x := range node.Children {
		fmt.Printf("%v %v %v\n", string(x.Type[0]), x.RelPath, getSizeText(x.TotalSize, human))
	}
}

func HumanSize(s, factor int64) string {
	unit := "B"
	if s > factor {
		s /= factor
		unit = "KB"
	}
	if s > factor {
		s /= factor
		unit = "MB"
	}
	if s > factor {
		s /= factor
		unit = "GB"
	}
	if s > factor {
		s /= factor
		unit = "TB"
	}
	return fmt.Sprintf("%v%v", s, unit)
}
