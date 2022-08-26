package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"howett.net/plist"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	prefpath := path.Join(home, "Library/Preferences")

	files, err := filepath.Glob(prefpath + "/*.plist")
	if err != nil {
		log.Fatal(err)
	}

	fileData := map[string][]plistItem{}

	for _, file := range files {
		q, err := readPlist(file)
		if err != nil {
			log.Println(err)
			continue
		}

		fileData[file] = q
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op == fsnotify.Create {
					if oldQ, ok := fileData[event.Name]; ok {
						newQ, err := readPlist(event.Name)
						if err != nil {
							log.Println(err)
							continue
						}

						aString := plistString(oldQ)
						bString := plistString(newQ)

						edits := myers.ComputeEdits(span.URIFromPath("a.txt"), aString, bString)

						log.Println(gotextdiff.ToUnified(
							event.Name,
							event.Name,
							aString,
							edits,
						))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			log.Println("failed to watch:", file)
		}
	}

	<-make(chan struct{})
}

func readPlist(file string) ([]plistItem, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", file, err)
	}

	d := plist.NewDecoder(f)
	var foo map[string]interface{}
	d.Decode(&foo)

	q := []plistItem{}
	iterateMaps(foo, "", &q)

	return q, nil
}

type plistItem struct {
	key   string
	vtype string
	value string
}

func (p plistItem) String() string {
	return fmt.Sprintf("%s : %s = %#v", p.key, p.vtype, p.value)
}

func plistString(q []plistItem) string {
	var b bytes.Buffer

	sort.Slice(q, func(i, j int) bool {
		return q[j].key < q[i].key
	})

	for _, p := range q {
		b.WriteString(p.String() + "\n")
	}

	return b.String()
}

func iterateMaps[K int | string](items map[string]interface{}, key K, p *[]plistItem) {
	for k, i := range items {
		handleContent(i, fmt.Sprintf("%s.'%s'", fmt.Sprint(key), k), p)
	}
}

func iterateSlices[K int | string](items []interface{}, key K, p *[]plistItem) {
	for k, i := range items {
		handleContent(i, fmt.Sprintf("%s.'%d'", fmt.Sprint(key), k), p)
	}
}

func handleContent(i interface{}, key string, p *[]plistItem) {
	switch v := i.(type) {
	case string:
		*p = append(*p, plistItem{key, fmt.Sprintf("%T", v), fmt.Sprintf("%#v", v)})
	case bool:
		*p = append(*p, plistItem{key, fmt.Sprintf("%T", v), fmt.Sprintf("%#v", v)})
	case []uint8:
		*p = append(*p, plistItem{key, fmt.Sprintf("%T", v), fmt.Sprintf("%#v", v)})
	case uint64:
		*p = append(*p, plistItem{key, fmt.Sprintf("%T", v), fmt.Sprintf("%#v", v)})
	case int64:
		*p = append(*p, plistItem{key, fmt.Sprintf("%T", v), fmt.Sprintf("%#v", v)})
	case float32:
		*p = append(*p, plistItem{key, fmt.Sprintf("%T", v), fmt.Sprintf("%#v", v)})
	case float64:
		*p = append(*p, plistItem{key, fmt.Sprintf("%T", v), fmt.Sprintf("%#v", v)})
	case time.Time:
		*p = append(*p, plistItem{key, fmt.Sprintf("%T", v), fmt.Sprintf("%#v", v)})
	case map[string]interface{}:
		iterateMaps(v, key, p)
	case []interface{}:
		iterateSlices(v, key, p)
	default:
		fmt.Println(key)
		fmt.Printf("I don't know about type %T!\n", i)
	}

}
