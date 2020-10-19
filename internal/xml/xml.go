package xml

import (
	"bytes"
	"encoding/xml"
	goxml "encoding/xml"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type StartElement struct {
	Name xml.Name
	Attr []XMLAttr
}

type XMLAttr struct {
	Name  xml.Name
	Value []byte
}

type bufioReader struct {
	buf      []byte
	rd       io.Reader // reader provided by the client
	r, w     int       // buf read and write positions
	err      error
	nextByte int
}

const maxConsecutiveEmptyReads = 100
const defaultBufSize = 4096

func (b *bufioReader) readErr() error {
	err := b.err
	b.err = nil
	return err
}

func (b *bufioReader) ReadByte() (byte, error) {
	if b.nextByte >= 0 {
		bt := byte(b.nextByte)
		b.nextByte = -1
		return bt, nil
	}
	for b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		b.fill() // buffer is empty
	}
	b.r++
	return b.buf[b.r-1], nil
}

// fill reads a new chunk into the buffer.
func (b *bufioReader) fill() {
	// Slide existing data to beginning.
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}

	if b.w >= len(b.buf) {
		panic("bufio: tried to fill full buffer")
	}

	// Read new data: try a limited number of times.
	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		n, err := b.rd.Read(b.buf[b.w:])
		b.w += n
		if err != nil {
			b.err = err
			return
		}
		if n > 0 {
			return
		}
	}
	b.err = io.ErrNoProgress
}

const nameCacheSize = 8

// A Decoder represents an XML parser reading a particular input stream.
// The parser assumes that its input is encoded in UTF-8.
type Decoder struct {
	OnStart func(StartElement)
	OnEnd   func(xml.EndElement)
	OnChar  func(xml.CharData)

	names     map[[nameCacheSize]byte]string
	r         *bufioReader
	buf       bytes.Buffer
	stk       *stack
	free      *stack
	needClose bool
	toClose   goxml.Name
	ns        map[string]string
	err       error
	attrPool  []XMLAttr
	strPool   []bytes.Buffer
}

// NewDecoder creates a new XML parser reading from r.
// If r does not implement io.ByteReader, NewDecoder will
// do its own buffering.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		ns:       make(map[string]string),
		names:    make(map[[nameCacheSize]byte]string),
		attrPool: make([]XMLAttr, 10),
		strPool:  make([]bytes.Buffer, 10),
		r: &bufioReader{
			buf:      make([]byte, defaultBufSize),
			rd:       r,
			nextByte: -1,
		},
	}
	return d
}

func (d *Decoder) handleStartElement(t StartElement) {
	for _, a := range t.Attr {
		if a.Name.Space == xmlnsPrefix {
			v, ok := d.ns[a.Name.Local]
			d.pushNs(a.Name.Local, v, ok)
			d.ns[a.Name.Local] = string(a.Value)
		}
		if a.Name.Space == "" && a.Name.Local == xmlnsPrefix {
			// Default space for untagged names
			v, ok := d.ns[""]
			d.pushNs("", v, ok)
			d.ns[""] = string(a.Value)
		}
	}

	d.translate(&t.Name, true)
	for i := range t.Attr {
		d.translate(&t.Attr[i].Name, false)
	}
	d.pushElement(t.Name)
	if d.OnStart != nil {
		d.OnStart(t)
	}
}

func (d *Decoder) handleEndElement(t goxml.EndElement) bool {
	d.translate(&t.Name, true)
	if d.popElement(&t) {
		if d.OnEnd != nil {
			d.OnEnd(t)
		}
		return true
	}
	return false
}

const (
	xmlURL      = "http://www.w3.org/XML/1998/namespace"
	xmlnsPrefix = "xmlns"
	xmlPrefix   = "xml"
)

// Apply name space translation to name n.
// The default name space (for Space=="")
// applies only to element names, not to attribute names.
func (d *Decoder) translate(n *goxml.Name, isElementName bool) {
	switch {
	case n.Space == xmlnsPrefix:
		return
	case n.Space == "" && !isElementName:
		return
	case n.Space == xmlPrefix:
		n.Space = xmlURL
	case n.Space == "" && n.Local == xmlnsPrefix:
		return
	}
	if v, ok := d.ns[n.Space]; ok {
		n.Space = v
	}
}

// Parsing state - stack holds old name space translations
// and the current set of open elements. The translations to pop when
// ending a given tag are *below* it on the stack, which is
// more work but forced on us by XML.
type stack struct {
	next *stack
	kind int
	name goxml.Name
	ok   bool
}

const (
	stkStart = iota
	stkNs
	stkEOF
)

func (d *Decoder) push(kind int) *stack {
	s := d.free
	if s != nil {
		d.free = s.next
	} else {
		s = new(stack)
	}
	s.next = d.stk
	s.kind = kind
	d.stk = s
	return s
}

func (d *Decoder) pop() *stack {
	s := d.stk
	if s != nil {
		d.stk = s.next
		s.next = d.free
		d.free = s
	}
	return s
}

// Record that after the current element is finished
// (that element is already pushed on the stack)
// Token should return EOF until popEOF is called.
func (d *Decoder) pushEOF() {
	// Walk down stack to find Start.
	// It might not be the top, because there might be stkNs
	// entries above it.
	start := d.stk
	for start.kind != stkStart {
		start = start.next
	}
	// The stkNs entries below a start are associated with that
	// element too; skip over them.
	for start.next != nil && start.next.kind == stkNs {
		start = start.next
	}
	s := d.free
	if s != nil {
		d.free = s.next
	} else {
		s = new(stack)
	}
	s.kind = stkEOF
	s.next = start.next
	start.next = s
}

// Undo a pushEOF.
// The element must have been finished, so the EOF should be at the top of the stack.
func (d *Decoder) popEOF() bool {
	if d.stk == nil || d.stk.kind != stkEOF {
		return false
	}
	d.pop()
	return true
}

// Record that we are starting an element with the given name.
func (d *Decoder) pushElement(name goxml.Name) {
	s := d.push(stkStart)
	s.name = name
}

// Record that we are changing the value of ns[local].
// The old value is url, ok.
func (d *Decoder) pushNs(local string, url string, ok bool) {
	s := d.push(stkNs)
	s.name.Local = local
	s.name.Space = url
	s.ok = ok
}

// Creates a SyntaxError.
func (d *Decoder) syntaxError(msg string) error {
	return &goxml.SyntaxError{Msg: msg}
}

// Record that we are ending an element with the given name.
// The name must match the record at the top of the stack,
// which must be a pushElement record.
// After popping the element, apply any undo records from
// the stack to restore the name translations that existed
// before we saw this element.
func (d *Decoder) popElement(t *goxml.EndElement) bool {
	s := d.pop()
	name := t.Name
	switch {
	case s == nil || s.kind != stkStart:
		d.err = d.syntaxError("unexpected end element </" + name.Local + ">")
		return false
	case s.name.Local != name.Local:
		d.err = d.syntaxError("element <" + s.name.Local + "> closed by </" + name.Local + ">")
		return false
	case s.name.Space != name.Space:
		d.err = d.syntaxError("element <" + s.name.Local + "> in space " + s.name.Space +
			"closed by </" + name.Local + "> in space " + name.Space)
		return false
	}

	// Pop stack until a Start or EOF is on the top, undoing the
	// translations that were associated with the element we just closed.
	for d.stk != nil && d.stk.kind != stkStart && d.stk.kind != stkEOF {
		s := d.pop()
		if s.ok {
			d.ns[s.name.Local] = s.name.Space
		} else {
			delete(d.ns, s.name.Local)
		}
	}

	return true
}

func (d *Decoder) RawToken() error {
	if d.stk != nil && d.stk.kind == stkEOF {
		return io.EOF
	}
	if d.err != nil {
		if d.err == io.EOF && d.stk != nil && d.stk.kind != stkEOF {
			d.err = d.syntaxError("unexpected EOF")
		}
		return d.err
	}
	if d.needClose {
		// The last element we read was self-closing and
		// we returned just the StartElement half.
		// Return the EndElement half now.
		d.needClose = false
		if !d.handleEndElement(goxml.EndElement{Name: d.toClose}) {
			return d.err
		}
		return nil
	}

	b, ok := d.getc()
	if !ok {
		if d.err == io.EOF && d.stk != nil && d.stk.kind != stkEOF {
			d.err = d.syntaxError("unexpected EOF")
		}
		return d.err
	}

	if b != '<' {
		// Text section.
		d.ungetc(b)
		data := d.text(-1)
		if data == nil {
			return d.err
		}
		if d.OnChar != nil {
			d.OnChar(goxml.CharData(data))
		}
		return nil
	}

	if b, ok = d.getc(); !ok {
		d.mustNotEOF()
		return d.err
	}
	switch b {
	case '/':
		// </: End element
		var name goxml.Name
		if name, ok = d.nsname(); !ok {
			if d.err == nil {
				d.err = d.syntaxError("expected element name after </")
			}
			return d.err
		}
		d.space()
		if b, ok = d.getc(); !ok {
			d.mustNotEOF()
			return d.err
		}
		if b != '>' {
			d.err = d.syntaxError("invalid characters between </" + name.Local + " and >")
			return d.err
		}
		if !d.handleEndElement(goxml.EndElement{Name: name}) {
			return d.err
		}
		return nil

	case '?':
		// <?: Processing instruction.
		if _, ok = d.name(); !ok {
			if d.err == nil {
				d.err = d.syntaxError("expected target name after <?")
			}
			return d.err
		}
		d.space()
		var b0 byte
		for {
			if b, ok = d.getc(); !ok {
				d.mustNotEOF()
				return d.err
			}
			d.buf.WriteByte(b)
			if b0 == '?' && b == '>' {
				break
			}
			b0 = b
		}
		return nil
	}

	// Must be an open element like <a href="foo">
	d.ungetc(b)

	var (
		name  goxml.Name
		empty bool
	)
	if name, ok = d.nsname(); !ok {
		if d.err == nil {
			d.err = d.syntaxError("expected element name after <")
		}
		return d.err
	}

	i := 0
	for {
		d.space()
		if b, ok = d.getc(); !ok {
			d.mustNotEOF()
			return d.err
		}
		if b == '/' {
			empty = true
			if b, ok = d.getc(); !ok {
				d.mustNotEOF()
				return d.err
			}
			if b != '>' {
				d.err = d.syntaxError("expected /> in element")
				return d.err
			}
			break
		}
		if b == '>' {
			break
		}
		d.ungetc(b)

		a := XMLAttr{}
		if a.Name, ok = d.nsname(); !ok {
			if d.err == nil {
				d.err = d.syntaxError("expected attribute name in element")
			}
			return d.err
		}
		d.space()
		if b, ok = d.getc(); !ok {
			d.mustNotEOF()
			return d.err
		}
		if b != '=' {
			d.err = d.syntaxError("attribute name without = in element")
			return d.err
		}
		d.space()
		data := d.attrval()
		if data == nil {
			return d.err
		}
		i++
		if len(d.attrPool) < i {
			d.attrPool = append(d.attrPool, make([]XMLAttr, len(d.attrPool))...)
			d.strPool = append(d.strPool, make([]bytes.Buffer, len(d.strPool))...)
		}
		bldr := &d.strPool[i-1]
		bldr.Reset()
		bldr.Write(data)
		a.Value = bldr.Bytes()
		d.attrPool[i-1] = a
	}
	if empty {
		d.needClose = true
		d.toClose = name
	}
	d.handleStartElement(StartElement{Name: name, Attr: d.attrPool[:i]})
	return nil
}

func (d *Decoder) attrval() []byte {
	b, ok := d.getc()
	if !ok {
		d.mustNotEOF()
		return nil
	}
	// Handle quoted attribute values
	if b == '"' || b == '\'' {
		return d.text(int(b))
	}
	// Handle unquoted attribute values for strict parsers
	d.err = d.syntaxError("unquoted or missing attribute value in element")
	return nil
}

// Skip spaces if any
func (d *Decoder) space() {
	for {
		b, ok := d.getc()
		if !ok {
			return
		}
		switch b {
		case ' ', '\r', '\n', '\t':
		default:
			d.ungetc(b)
			return
		}
	}
}

// Read a single byte.
// If there is no byte to read, return ok==false
// and leave the error in d.err.
func (d *Decoder) getc() (b byte, ok bool) {
	b, d.err = d.r.ReadByte()
	ok = d.err == nil
	return
}

// Must read a single byte.
// If there is no byte to read,
// set d.err to SyntaxError("unexpected EOF")
// and return ok==false
func (d *Decoder) mustNotEOF() {
	if d.err == io.EOF {
		d.err = d.syntaxError("unexpected EOF")
	}
	return
}

// Unread a single byte.
func (d *Decoder) ungetc(b byte) {
	d.r.nextByte = int(b)
}

var entity = map[string]int{
	"lt":   '<',
	"gt":   '>',
	"amp":  '&',
	"apos": '\'',
	"quot": '"',
}

// Read plain text section (XML calls it character data).
// If quote >= 0, we are in a quoted string and need to find the matching quote.
// If cdata == true, we are in a <![CDATA[ section and need to find ]]>.
// On failure return nil and leave the error in d.err.
func (d *Decoder) text(quote int) []byte {
	var b1 byte
	d.buf.Reset()
Input:
	for {
		b, ok := d.getc()
		if !ok {
			break Input
		}

		// Stop reading text if we see a <.
		if b == '<' {
			if quote >= 0 {
				d.err = d.syntaxError("unescaped < inside quoted string")
				return nil
			}
			d.ungetc('<')
			break Input
		}
		if quote >= 0 && b == byte(quote) {
			break Input
		}
		if b == '&' {
			// Read escaped character expression up to semicolon.
			// XML in all its glory allows a document to define and use
			// its own character names with <!ENTITY ...> directives.
			// Parsers are required to recognize lt, gt, amp, apos, and quot
			// even if they have not been declared.
			before := d.buf.Len()
			d.buf.WriteByte('&')
			var ok bool
			var text string
			var haveText bool
			if b, ok = d.getc(); !ok {
				d.mustNotEOF()
				return nil
			}
			if b == '#' {
				d.buf.WriteByte(b)
				if b, ok = d.getc(); !ok {
					d.mustNotEOF()
					return nil
				}
				base := 10
				if b == 'x' {
					base = 16
					d.buf.WriteByte(b)
					if b, ok = d.getc(); !ok {
						d.mustNotEOF()
						return nil
					}
				}
				start := d.buf.Len()
				for '0' <= b && b <= '9' ||
					base == 16 && 'a' <= b && b <= 'f' ||
					base == 16 && 'A' <= b && b <= 'F' {
					d.buf.WriteByte(b)
					if b, ok = d.getc(); !ok {
						d.mustNotEOF()
						return nil
					}
				}
				if b != ';' {
					d.ungetc(b)
				} else {
					s := string(d.buf.Bytes()[start:])
					d.buf.WriteByte(';')
					n, err := strconv.ParseUint(s, base, 64)
					if err == nil && n <= unicode.MaxRune {
						text = string(rune(n))
						haveText = true
					}
				}
			} else {
				d.ungetc(b)
				if !d.readName() {
					if d.err != nil {
						return nil
					}
				}
				if b, ok = d.getc(); !ok {
					d.mustNotEOF()
					return nil
				}
				if b != ';' {
					d.ungetc(b)
				} else {
					name := d.buf.Bytes()[before+1:]
					d.buf.WriteByte(';')
					if isName(name) {
						s := string(name)
						if r, ok := entity[s]; ok {
							text = string(rune(r))
							haveText = true
						}
					}
				}
			}

			if haveText {
				d.buf.Truncate(before)
				d.buf.Write([]byte(text))
				b1 = 0
				continue Input
			}
			ent := string(d.buf.Bytes()[before:])
			if ent[len(ent)-1] != ';' {
				ent += " (no semicolon)"
			}
			d.err = d.syntaxError("invalid character entity " + ent)
			return nil
		}

		// We must rewrite unescaped \r and \r\n into \n.
		if b == '\r' {
			d.buf.WriteByte('\n')
		} else if b1 == '\r' && b == '\n' {
			// Skip \r\n--we already wrote \n.
		} else {
			d.buf.WriteByte(b)
		}

		b1 = b
	}
	return d.buf.Bytes()
}

// Get name space name: name with a : stuck in the middle.
// The part before the : is the name space identifier.
func (d *Decoder) nsname() (name goxml.Name, ok bool) {
	s, ok := d.name()
	if !ok {
		return
	}
	i := strings.IndexByte(s, ':')
	if i < 0 {
		name.Local = s
	} else {
		name.Space = s[0:i]
		name.Local = s[i+1:]
	}
	return
}

// Get name: /first(first|second)*/
// Do not set d.err if the name is missing (unless unexpected EOF is received):
// let the caller provide better context.
func (d *Decoder) name() (s string, ok bool) {
	d.buf.Reset()
	if !d.readName() {
		return "", false
	}

	// Now we check the characters.
	b := d.buf.Bytes()
	if len(b) == 1 {
		switch b[0] {
		case 'x':
			return "x", true
		case 'y':
			return "y", true
		case 'z':
			return "z", true
		case 'u':
			return "u", true
		case 'v':
			return "v", true
		}
	} else if len(b) == 2 {
		if b[0] == 'v' {
			switch b[1] {
			case '1':
				return "v1", true
			case '2':
				return "v2", true
			case '3':
				return "v3", true
			}
		} else if b[0] == 'p' {
			switch b[1] {
			case '1':
				return "p1", true
			case '2':
				return "p2", true
			case '3':
				return "p3", true
			}
		}
	} else if string(b) == "vertex" {
		return "vertex", true
	} else if string(b) == "triangle" {
		return "triangle", true
	} else if len(b) == 3 && string(b) == "pid" {
		return "pid", true
	} else if len(b) == 5 && string(b) == "color" {
		return "color", true
	}
	var arr [nameCacheSize]byte
	if len(b) <= nameCacheSize {
		copy(arr[:], b)
		if s, ok = d.names[arr]; ok {
			return s, ok
		}
	}
	if !isName(b) {
		d.err = d.syntaxError("invalid XML name: " + string(b))
		return "", false
	}
	s = string(b)
	if len(b) <= nameCacheSize {
		copy(arr[:], b)
		d.names[arr] = s
	}
	return s, true
}

// Read a name and append its bytes to d.buf.
// The name is delimited by any single-byte character not valid in names.
// All multi-byte characters are accepted; the caller must check their validity.
func (d *Decoder) readName() bool {
	const runSelf = 0x80 // characters below RuneSelf are represented as themselves in a single byte.
	var b byte
	var ok bool
	if b, ok = d.getc(); !ok {
		d.mustNotEOF()
		return false
	}
	if b < runSelf && !isNameByte(b) {
		d.ungetc(b)
		return false
	}
	d.buf.WriteByte(b)

	for {
		if b, ok = d.getc(); !ok {
			d.mustNotEOF()
			return false
		}
		if b < runSelf && !isNameByte(b) {
			d.ungetc(b)
			break
		}
		d.buf.WriteByte(b)
	}
	return true
}

func isNameByte(c byte) bool {
	return 'A' <= c && c <= 'Z' ||
		'a' <= c && c <= 'z' ||
		'0' <= c && c <= '9' ||
		c == '_' || c == ':' || c == '.' || c == '-'
}
func isName(s []byte) bool {
	return len(s) != 0
}
