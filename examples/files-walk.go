package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"utils"
)

type FileType int

const (
	TypeDir  FileType = 1
	TypeFile FileType = 2
)

type TreeNode struct {
	Path         string
	RelativePath string
	Children     []*TreeNode
	Parent       *TreeNode
	Type         FileType // "folder", "file"
}

func NewRootNode(path, relativePath string) *TreeNode {
	node := &TreeNode{
		Path:         "/", // not real..
		RelativePath: "/",
		Type:         TypeDir,
	}
	return node
}

func (node *TreeNode) AddChild(path, relativePath string, fileType FileType) *TreeNode {
	newChildNode := &TreeNode{
		Path:         "/", // not real..
		RelativePath: "/",
		Parent:       node,
		Type:         fileType,
	}
	node.Children = append(node.Children, newChildNode)
	return newChildNode
}

func Tree(basefolder, relativePath string, depth int) (tree *TreeNode, err error) {

	finalPath := filepath.Join(basefolder, relativePath)

	// ensure recurity. // 虽然这个被Deprecated了，但是我这里的需求应该没有bug，也没找到替代品，先用着。
	if !filepath.HasPrefix(finalPath, basefolder) {
		return nil, fmt.Errorf("permission denied, you don't have permission to access folder %s", relativePath)
	}

	fmt.Println("final path is : ", finalPath)

	separator := string(filepath.Separator)
	baseSlashCount := strings.Count(basefolder, separator)

	root := NewRootNode("/", "/")
	stack := utils.NewStack()
	lastDepth := 0

	var files []string
	err = filepath.Walk(finalPath, func(path string, info os.FileInfo, err error) error {
		// fmt.Println("...", path)
		// files = append(files, path)

		dep := strings.Count(path, separator) - baseSlashCount
		if dep > depth {
			return filepath.SkipDir
		}

		lastDepth = dep

		if err != nil {
			fmt.Printf("...erroccurde: %v", err)
			//处理文件读取异常
			return err // should be ?
		}

		if info.IsDir() {
			fmt.Println(".........", dep, " --> ", path)
			// root.AddChild()

			// 满足条件不用管
			return nil
			// 不满足条件
			return filepath.SkipDir
		} else {
			// file
		}

		return nil
	})
	// ! ..............................
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fmt.Println(file)
	}

	return nil, nil
}
