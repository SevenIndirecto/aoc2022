package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type file struct {
	name string
	size int
}

type dir struct {
	name string
	files map[string]file
	dirs map[string]*dir
	parent *dir
	cachedSize int
}

func newDir(name string, parent *dir) dir {
	return dir{
		name: name,
		files: make(map[string]file, 0),
		dirs: make(map[string]*dir, 0),
		parent: parent,
		cachedSize: -1,
	}
}

func (d *dir) toString(depth int) string {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += " "
	}

	repr := indent + "- " + d.name + " (dir)\n"
	for _, subDir := range d.dirs {
		repr += subDir.toString(depth+1)
	}

	for _, f := range d.files {
		fileSize := strconv.Itoa(f.size)
		repr += indent + " - (" + f.name + ", size=" + fileSize + ")\n"
	}
	return repr
}

func (d *dir) size() int {
	if d.cachedSize >= 0 {
		return d.cachedSize
	}

	size := 0
	for _, f := range d.files {
		size += f.size
	}

	for _, subDir := range d.dirs {
		size += subDir.size()
	}
	d.cachedSize = size
	return size
}

func parseLog(lines []string) dir {
	root := newDir("/", nil)
	currentDir := &root

	for _, l := range lines {
		tokens := strings.Split(l, " ")

		if string(l[0]) == "$" {
			if len(tokens) == 2 {
				// ls command
				continue
			}

			// cd command
			if tokens[2] == ".." {
				currentDir = currentDir.parent
			} else if tokens[2] == "/" {
				currentDir = &root
			} else {
				dirName := tokens[2]
				_, exists := currentDir.dirs[dirName]
				if !exists {
					c := newDir(dirName, currentDir)
					currentDir.dirs[dirName] = &c
				}
				currentDir = currentDir.dirs[dirName]
			}
		} else {
			// list item output
			if tokens[0] == "dir" {
				// Ignore empty directories for now since they're size 0
				continue
			}
			filename := tokens[1]
			_, exists := currentDir.files[filename]
			if !exists {
				size, _ := strconv.Atoi(tokens[0])
				currentDir.files[filename] = file{filename, size}
			}
		}
	}

	return root
}

func findSmallDirs(d *dir, smallDirs []*dir) []*dir {
	// Process children
	for _, subDir := range d.dirs {
		smallDirs = findSmallDirs(subDir, smallDirs)
	}

	if d.size() <= 1e5 {
		smallDirs = append(smallDirs, d)
	}
	return smallDirs
}

func PartOne(lines []string) int {
	root := parseLog(lines)
	//fmt.Println(root.toString(0))
	smallDirs := findSmallDirs(&root, make([]*dir, 0))
	sum := 0
	for _, d := range smallDirs {
		sum += d.size() // Ok since it's cached
	}
	return sum
}

func PartTwo(lines []string) int {
	root := parseLog(lines)
	totalDiskSpace := 70000000
	requiredDiskSpace := 30000000
	unusedSpace := totalDiskSpace - root.size()
	spaceToFree := requiredDiskSpace - unusedSpace

	smallestDir := getSmallestDirToDelete(&root, &root, spaceToFree)
	return smallestDir.size()
}

func getSmallestDirToDelete(d *dir, candidate *dir, spaceToFree int) *dir {
	for _, subDir := range d.dirs {
		childCandidate := getSmallestDirToDelete(subDir, candidate, spaceToFree)
		if childCandidate.size() < candidate.size() {
			candidate = childCandidate
		}
	}

	if d.size() >= spaceToFree && d.size() < candidate.size() {
		candidate = d
	}
	return candidate
}

func LoadLines(path string) ([]string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	txt := string(dat)
	lines := strings.Split(txt, "\n")
	return lines[:len(lines)-1], nil
}

func main() {
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
