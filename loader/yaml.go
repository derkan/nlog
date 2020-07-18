// // Assimilated From: "https://github.com/kylelemons/go-gypsy" by "7 of 9".  v2.1.x 08cad36 on  Nov 5, 2011
// - changed parseNode for trimming spaces and '"'

// a simple configuration file, and as such does not support a lot of the more
// nuanced syntaxes allowed in full-fledged YAML.  YAML does not allow indent with
// tabs, and GYPSY does not ever consider a tab to be a space character.  It is
// recommended that your editor be configured to convert tabs to spaces when
// editing Gypsy config files.
//
// Gypsy understands the following to be a list:
//
//     - one
//     - two
//     - three
//
// This is parsed as a `yaml.List`, and can be retrieved from the
// `yaml.Node.List()` method.  In this case, each element of the `yaml.List` would
// be a `yaml.Scalar` whose value can be retrieved with the `yaml.Scalar.String()`
// method.
//
// Gypsy understands the following to be a mapping:
//
//     key:     value
//     foo:     bar
//     running: away
//
// A mapping is an unordered list of `key:value` pairs.  All whitespace after the
// colon is stripped from the value and is used for alignment purposes during
// export.  If the value is not a list or a map, everything after the first
// non-space character until the end of the line is used as the `yaml.Scalar`
// value.
//
// Gypsy allows arbitrary nesting of maps inside lists, lists inside of maps, and
// maps and/or lists nested inside of themselves.
//
// A map inside of a list:
//
//     - name: John Smith
//       age:  42
//     - name: Jane Smith
//       age:  45
//
// A list inside of a map:
//
//     schools:
//       - Meadow Glen
//       - Forest Creek
//       - Shady Grove
//     libraries:
//       - Joseph Hollingsworth Memorial
//       - Andrew Keriman Memorial
//
// A list of lists:
//
//     - - one
//       - two
//       - three
//     - - un
//       - deux
//       - trois
//     - - ichi
//       - ni
//       - san
//
// A map of maps:
//
//     google:
//       company: Google, Inc.
//       ticker:  GOOG
//       url:     http://google.com/
//     yahoo:
//       company: Yahoo, Inc.
//       ticker:  YHOO
//       url:     http://yahoo.com/
//
// In the case of a map of maps, all sub-keys must be on subsequent lines and
// indented equally.  It is allowable for the first key/value to be on the same
// line if there is more than one key/value pair, but this is not recommended.
//
// Values can also be expressed in long form (leading whitespace of the first line
// is removed from it and all subsequent lines).  In the normal (baz) case,
// newlines are treated as spaces, all indentation is removed.  In the folded case
// (bar), newlines are treated as spaces, except pairs of newlines (e.g. a blank
// line) are treated as a single newline, only the indentation level of the first
// line is removed, and newlines at the end of indented lines are preserved.  In
// the verbatim (foo) case, only the indent at the level of the first line is
// stripped.  The example:
//
//     foo: |
//       lorem ipsum dolor
//       sit amet
//     bar: >
//       lorem ipsum
//
//         dolor
//
//       sit amet
//     baz:
//       lorem ipsum
//        dolor sit amet

package loader

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Supporting types and constants

const (
	typUnknown = iota
	typSequence
	typMapping
	typScalar
)

var typNames = []string{
	"Unknown", "Sequence", "Mapping", "Scalar",
}

// A Node is a YAML Node which can be a Map, List or Scalar.
type Node interface {
	write(io.Writer, int, int)
}

// A Scalar is a YAML Scalar.
type Scalar string

// A List is a YAML Sequence of Nodes.
type List []Node

// A Map is a YAML Mapping which maps Strings to Nodes.
type Map map[string]Node

type NodeNotFound struct {
	Full string
	Spec string
}
type NodeTypeMismatch struct {
	Full     string
	Spec     string
	Token    string
	Node     Node
	Expected string
}

// A File represents the top-level YAML node found in a file.  It is intended
// for use as a configuration file.
type File struct {
	Root Node

	// TODO(kevlar): Add a cache?
}

// ReadFile reads a YAML configuration file from the given filename.
func ReadFile(filename string) (*File, error) {
	fin, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fin.Close()

	f := new(File)
	f.Root, err = Parse(fin)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ConfigReader reads a YAML configuration from a static string.  If an error is
// found, it will panic.  This is a utility function and is intended for use in
// initializers.
func Config(yamlconf string) (*File, error) {
	var err error
	buf := bytes.NewBufferString(yamlconf)

	f := new(File)
	f.Root, err = Parse(buf)
	if err != nil {
		return nil, err
	}

	return f, err
}

// ConfigFile reads a YAML configuration file from the given filename and
// panics if an error is found.  This is a utility function and is intended for
// use in initializers.
func ConfigFile(filename string) *File {
	f, err := ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return f
}

// Get retrieves a scalar from the file specified by a string of the same
// format as that expected by Child.  If the final node is not a Scalar, Get
// will return an error.
func (f *File) Get(def string, spec string, params ...interface{}) (string, error) {
	spec = fmt.Sprintf(spec, params...)
	node, err := Child(f.Root, spec)
	if err != nil {
		return def, err
	}

	if node == nil {
		return def, &NodeNotFound{
			Full: spec,
			Spec: spec,
		}
	}

	scalar, ok := node.(Scalar)
	if !ok {
		return def, &NodeTypeMismatch{
			Full:     spec,
			Spec:     spec,
			Token:    "$",
			Expected: "yaml.Scalar",
			Node:     node,
		}
	}
	return scalar.String(), nil
}

func (f *File) GetInt64(def int64, spec string, params ...interface{}) (int64, error) {
	spec = fmt.Sprintf(spec, params...)
	s, err := f.Get("", spec)
	if err != nil {
		return def, err
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil || s == "" {
		return def, err
	}

	return i, nil
}

func (f *File) GetInt(def int, spec string, params ...interface{}) (int, error) {
	spec = fmt.Sprintf(spec, params...)
	s, err := f.Get("", spec)
	if err != nil {
		return def, err
	}

	i, err := strconv.Atoi(s)
	if err != nil || s == "" {
		return def, err
	}

	return i, nil
}

func (f *File) GetBool(def bool, spec string, params ...interface{}) (bool, error) {
	spec = fmt.Sprintf(spec, params...)
	s, err := f.Get("", spec)
	if err != nil {
		return def, err
	}

	b, err := strconv.ParseBool(s)
	if err != nil || s == "" {
		return def, err
	}

	return b, nil
}

// Count retrieves a the number of elements in the specified list from the file
// using the same format as that expected by Child.  If the final node is not a
// List, Count will return an error.
func (f *File) Count(spec string, params ...interface{}) (int, error) {
	spec = fmt.Sprintf(spec, params...)
	node, err := Child(f.Root, spec)
	if err != nil {
		return -1, err
	}

	if node == nil {
		return -1, &NodeNotFound{
			Full: spec,
			Spec: spec,
		}
	}

	lst, ok := node.(List)
	if !ok {
		return -1, &NodeTypeMismatch{
			Full:     spec,
			Spec:     spec,
			Token:    "$",
			Expected: "yaml.List",
			Node:     node,
		}
	}
	return lst.Len(), nil
}

// Require retrieves a scalar from the file specified by a string of the same
// format as that expected by Child.  If the final node is not a Scalar, String
// will panic.  This is a convenience function for use in initializers.
func (f *File) Require(spec string, params ...interface{}) string {
	spec = fmt.Sprintf(spec, params...)
	str, err := f.Get("", spec)
	if err != nil {
		panic(err)
	}
	return str
}

// Child retrieves a child node from the specified node as follows:
//   .mapkey   - Get the key 'mapkey' of the Node, which must be a Map
//   [idx]     - Choose the index from the current Node, which must be a List
//
// The above selectors may be applied recursively, and each successive selector
// applies to the result of the previous selector.  For convenience, a "." is
// implied as the first character if the first character is not a "." or "[".
// The node tree is walked from the given node, considering each token of the
// above format.  If a node along the evaluation path is not found, an error is
// returned. If a node is not the proper type, an error is returned.  If the
// final node is not a Scalar, an error is returned.
func Child(root Node, spec string, params ...interface{}) (Node, error) {
	spec = fmt.Sprintf(spec, params...)
	if len(spec) == 0 {
		return root, nil
	}

	if first := spec[0]; first != '.' && first != '[' {
		spec = "." + spec
	}

	var recur func(Node, string, string) (Node, error)
	recur = func(n Node, last, s string) (Node, error) {

		if len(s) == 0 {
			return n, nil
		}

		if n == nil {
			return nil, &NodeNotFound{
				Full: spec,
				Spec: last,
			}
		}

		// Extract the next token
		delim := 1 + strings.IndexAny(s[1:], ".[")
		if delim <= 0 {
			delim = len(s)
		}
		tok := s[:delim]
		remain := s[delim:]

		switch s[0] {
		case '[':
			s, ok := n.(List)
			if !ok {
				return nil, &NodeTypeMismatch{
					Node:     n,
					Expected: "yaml.List",
					Full:     spec,
					Spec:     last,
					Token:    tok,
				}
			}

			if tok[0] == '[' && tok[len(tok)-1] == ']' {
				if num, err := strconv.Atoi(tok[1 : len(tok)-1]); err == nil {
					if num >= 0 && num < len(s) {
						return recur(s[num], last+tok, remain)
					}
				}
			}
			return nil, &NodeNotFound{
				Full: spec,
				Spec: last + tok,
			}
		default:
			m, ok := n.(Map)
			if !ok {
				return nil, &NodeTypeMismatch{
					Node:     n,
					Expected: "yaml.Map",
					Full:     spec,
					Spec:     last,
					Token:    tok,
				}
			}

			n, ok = m[tok[1:]]
			if !ok {
				return nil, &NodeNotFound{
					Full: spec,
					Spec: last + tok,
				}
			}
			return recur(n, last+tok, remain)
		}
	}
	return recur(root, "", spec)
}

func (e *NodeNotFound) Error() string {
	return fmt.Sprintf("yaml: %s: %q not found", e.Full, e.Spec)
}

func (e *NodeTypeMismatch) Error() string {
	return fmt.Sprintf("yaml: %s: type mismatch: %q is %T, want %s (at %q)",
		e.Full, e.Spec, e.Node, e.Expected, e.Token)
}

// Key returns the value associeted with the key in the map.
func (node Map) Key(key string) Node {
	return node[key]
}

func (node Map) write(out io.Writer, firstind, nextind int) {
	indent := bytes.Repeat([]byte{' '}, nextind)
	ind := firstind

	width := 0
	scalarkeys := []string{}
	objectkeys := []string{}
	for key, value := range node {
		if _, ok := value.(Scalar); ok {
			if swid := len(key); swid > width {
				width = swid
			}
			scalarkeys = append(scalarkeys, key)
			continue
		}
		objectkeys = append(objectkeys, key)
	}
	sort.Strings(scalarkeys)
	sort.Strings(objectkeys)

	for _, key := range scalarkeys {
		value := node[key].(Scalar)
		out.Write(indent[:ind])
		fmt.Fprintf(out, "%-*s %s\n", width+1, key+":", string(value))
		ind = nextind
	}
	for _, key := range objectkeys {
		out.Write(indent[:ind])
		if node[key] == nil {
			fmt.Fprintf(out, "%s: <nil>\n", key)
			continue
		}
		fmt.Fprintf(out, "%s:\n", key)
		ind = nextind
		node[key].write(out, ind+2, ind+2)
	}
}

// Get the number of items in the List.
func (node List) Len() int {
	return len(node)
}

// Get the idx'th item from the List.
func (node List) Item(idx int) Node {
	if idx >= 0 && idx < len(node) {
		return node[idx]
	}
	return nil
}

func (node List) write(out io.Writer, firstind, nextind int) {
	indent := bytes.Repeat([]byte{' '}, nextind)
	ind := firstind

	for _, value := range node {
		out.Write(indent[:ind])
		fmt.Fprint(out, "- ")
		ind = nextind
		value.write(out, 0, ind+2)
	}
}

// String returns the string represented by this Scalar.
func (node Scalar) String() string { return string(node) }

func (node Scalar) write(out io.Writer, ind, _ int) {
	fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", ind), string(node))
}

// Render returns a string of the node as a YAML document.  Note that
// Scalars will have a newline appended if they are rendered directly.
func Render(node Node) string {
	buf := bytes.NewBuffer(nil)
	node.write(buf, 0, 0)
	return buf.String()
}

// Parse returns a root-level Node parsed from the lines read from r.  In
// general, this will be done for you by one of the File constructors.
func Parse(r io.Reader) (node Node, err error) {
	lb := &lineBuffer{
		Reader: bufio.NewReader(r),
	}

	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = r
			case string:
				err = errors.New(r)
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	node = parseNode(lb, 0, nil)
	return
}

type lineReader interface {
	Next(minIndent int) *indentedLine
}

type indentedLine struct {
	lineno int
	indent int
	line   []byte
}

func (line *indentedLine) String() string {
	return fmt.Sprintf("%2d: %s%s", line.indent,
		strings.Repeat(" ", 0*line.indent), string(line.line))
}

func parseNode(r lineReader, ind int, initial Node) (node Node) {
	first := true
	node = initial

	// read lines
	for {
		line := r.Next(ind)
		if line == nil {
			break
		}

		if len(line.line) == 0 {
			continue
		}

		if first {
			ind = line.indent
			first = false
		}

		types := []int{}
		pieces := []string{}

		var inlineValue func([]byte)
		inlineValue = func(partial []byte) {
			// TODO(kevlar): This can be a for loop now
			vtyp, brk := getType(partial)
			begin, end := partial[:brk], partial[brk:]

			if vtyp == typMapping {
				end = end[1:]
			}
			end = bytes.TrimLeft(end, " ")

			switch vtyp {
			case typScalar:
				types = append(types, typScalar)
				pieces = append(pieces, strings.Trim(strings.Trim(string(bytes.TrimSpace(end)), ","), "\""))
				return
			case typMapping:
				types = append(types, typMapping)
				pieces = append(pieces, strings.TrimSpace(string(begin)))

				trimmed := bytes.TrimSpace(end)
				if len(trimmed) == 1 && trimmed[0] == '|' {
					text := ""

					for {
						l := r.Next(1)
						if l == nil {
							break
						}

						s := string(l.line)
						s = strings.TrimSpace(s)
						if len(s) == 0 {
							break
						}
						text = text + "\n" + s
					}

					types = append(types, typScalar)
					pieces = append(pieces, string(text))
					return
				}
				inlineValue(end)
			case typSequence:
				types = append(types, typSequence)
				pieces = append(pieces, "-")
				inlineValue(end)
			}
		}

		inlineValue(line.line)
		var prev Node

		// Nest inlines
		for len(types) > 0 {
			last := len(types) - 1
			typ, piece := types[last], pieces[last]

			var current Node
			if last == 0 {
				current = node
			}
			//child := parseNode(r, line.indent+1, typUnknown) // TODO allow scalar only

			// Add to current node
			switch typ {
			case typScalar: // last will be == nil
				if _, ok := current.(Scalar); current != nil && !ok {
					panic("cannot append scalar to non-scalar node")
				}
				if current != nil {
					current = Scalar(piece) + " " + current.(Scalar)
					break
				}
				current = Scalar(piece)
			case typMapping:
				var mapNode Map
				var ok bool
				var child Node

				// Get the current map, if there is one
				if mapNode, ok = current.(Map); current != nil && !ok {
					_ = current.(Map) // panic
				} else if current == nil {
					mapNode = make(Map)
				}

				if _, inlineMap := prev.(Scalar); inlineMap && last > 0 {
					current = Map{
						piece: prev,
					}
					break
				}

				child = parseNode(r, line.indent+1, prev)
				mapNode[piece] = child
				current = mapNode

			case typSequence:
				var listNode List
				var ok bool
				var child Node

				// Get the current list, if there is one
				if listNode, ok = current.(List); current != nil && !ok {
					_ = current.(List) // panic
				} else if current == nil {
					listNode = make(List, 0)
				}

				if _, inlineList := prev.(Scalar); inlineList && last > 0 {
					current = List{
						prev,
					}
					break
				}

				child = parseNode(r, line.indent+1, prev)
				listNode = append(listNode, child)
				current = listNode

			}

			if last < 0 {
				last = 0
			}
			types = types[:last]
			pieces = pieces[:last]
			prev = current
		}

		node = prev
	}
	return
}

func getType(line []byte) (typ, split int) {
	if len(line) == 0 {
		return
	}

	if line[0] == '-' {
		typ = typSequence
		split = 1
		return
	}

	typ = typScalar

	if line[0] == ' ' || line[0] == '"' {
		return
	}

	// the first character is real
	// need to iterate past the first word
	// things like "foo:" and "foo :" are mappings
	// everything else is a scalar

	idx := bytes.IndexAny(line, " \":")
	if idx < 0 {
		return
	}

	if line[idx] == '"' {
		return
	}

	if line[idx] == ':' {
		typ = typMapping
		split = idx
	} else if line[idx] == ' ' {
		// we have a space
		// need to see if its all spaces until a :
		for i := idx; i < len(line); i++ {
			switch ch := line[i]; ch {
			case ' ':
				continue
			case ':':
				// only split on colons followed by a space
				if i+1 < len(line) && line[i+1] != ' ' {
					continue
				}

				typ = typMapping
				split = i
				break
			default:
				break
			}
		}
	}

	if typ == typMapping && split+1 < len(line) && line[split+1] != ' ' {
		typ = typScalar
		split = 0
	}

	return
}

// lineReader implementations

type lineBuffer struct {
	*bufio.Reader
	readLines int
	pending   *indentedLine
}

func (lb *lineBuffer) Next(min int) (next *indentedLine) {
	if lb.pending == nil {
		var (
			read []byte
			more bool
			err  error
		)

		l := new(indentedLine)
		l.lineno = lb.readLines
		more = true
		for more {
			read, more, err = lb.ReadLine()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				panic(err)
			}
			l.line = append(l.line, read...)
		}
		lb.readLines++

		for _, ch := range l.line {
			switch ch {
			case ' ':
				l.indent += 1
				continue
			default:
			}
			break
		}
		l.line = l.line[l.indent:]

		// Ignore blank lines and comments.
		if len(l.line) == 0 || l.line[0] == '#' {
			return lb.Next(min)
		}

		lb.pending = l
	}
	next = lb.pending
	if next.indent < min {
		return nil
	}
	lb.pending = nil
	return
}

type lineSlice []*indentedLine

func (ls *lineSlice) Next(min int) (next *indentedLine) {
	if len(*ls) == 0 {
		return nil
	}
	next = (*ls)[0]
	if next.indent < min {
		return nil
	}
	*ls = (*ls)[1:]
	return
}

func (ls *lineSlice) Push(line *indentedLine) {
	*ls = append(*ls, line)
}
