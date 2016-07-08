package taggraph

import (
	"reflect"
	"strings"
)

//Tagger is an interface for interacting with a node of a TagGrapher
type Tagger interface {
	Children() []string
	Parents() []string
	PathsToAllAncestors() [][]string
	PathsToAllAncestorsAsString(delim string) []string
	PathsToAllDescendents() [][]string
	PathsToAllDescendentsAsString(delim string) []string
	Name() string
}

type tag struct {
	name       string
	childTags  tags
	parentTags tags
}

func (t *tag) Name() string {
	return t.name
}

func (t *tag) Children() []string {
	return flatTags(t.childTags)
}

func (t *tag) Parents() []string {
	return flatTags(t.parentTags)
}

func flatTags(t tags) []string {
	arr := []string{}
	for _, v := range t {
		arr = append(arr, v.name)
	}
	return arr
}

func (t *tag) PathsToAllAncestors() [][]string {
	return ancestors(t, [][]string{{t.name}})
}

func (t *tag) PathsToAllAncestorsAsString(delim string) []string {
	paths := []string{}
	for _, v := range t.PathsToAllAncestors() {
		paths = append(paths, strings.Join(v, delim))
	}
	return paths
}

func (t *tag) PathsToAllDescendents() [][]string {
	return descendents(t, [][]string{{t.name}})
}

func (t *tag) PathsToAllDescendentsAsString(delim string) []string {
	paths := []string{}
	for _, v := range t.PathsToAllDescendents() {
		paths = append(paths, strings.Join(v, delim))
	}
	return paths
}

func ancestors(t *tag, paths [][]string) [][]string {
	newPaths := [][]string{}
	for _, path := range paths {
		for _, v := range t.parentTags {
			pathCopy := cloneStringSlice(path)
			if pathCopy[0] == t.name {
				//break ring
				if !contain(v.name, pathCopy) {
					pathCopy = prependStringSlice(pathCopy, v.name)
				} else {
					return paths
				}
			}
			newPaths = append(newPaths, pathCopy)
			if v.parentTags.Len() > 0 {
				newPaths = ancestors(v, newPaths)
			}
		}
	}
	return newPaths
}

func descendents(t *tag, paths [][]string) [][]string {
	newPaths := [][]string{}
	for _, path := range paths {
		for _, v := range t.childTags {
			if contain(v.name, path) {
				continue
			}
			pathCopy := cloneStringSlice(path)

			if pathCopy[len(pathCopy)-1] == t.name {
				//break ring
				if !contain(v.name, pathCopy) {
					pathCopy = append(pathCopy, v.name)
					newPaths = append(newPaths, pathCopy)
				} else {
					return paths
				}
			} else {
				newPaths = append(newPaths, pathCopy)
				break
			}
			if v.childTags.Len() > 0 {
				newPaths = descendents(v, newPaths)
			}
		}
	}
	return newPaths
}

func cloneStringSlice(slice []string) []string {
	return append([]string(nil), slice...)
}

func prependStringSlice(slice []string, val string) []string {
	return append(slice[:0], append([]string{val}, slice[0:]...)...)
}

// is obj in target?
func contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}
