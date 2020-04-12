package xml

import (
	"bytes"
	"encoding/xml"
	goxml "encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type bufioReader struct {
	buf  []byte
	rd   io.Reader // reader provided by the client
	r, w int       // buf read and write positions
	err  error
}

const minReadBufferSize = 16
const maxConsecutiveEmptyReads = 100
const defaultBufSize = 4096

func (b *bufioReader) readErr() error {
	err := b.err
	b.err = nil
	return err
}

func (b *bufioReader) ReadByte() (byte, error) {
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

// A Decoder represents an XML parser reading a particular input stream.
// The parser assumes that its input is encoded in UTF-8.
type Decoder struct {
	OnStart func(xml.StartElement)
	OnEnd   func(xml.EndElement)
	OnChar  func(xml.CharData)
	// Entity can be used to map non-standard entity names to string replacements.
	// The parser behaves as if these standard mappings are present in the map,
	// regardless of the actual map content:
	//
	//	"lt": "<",
	//	"gt": ">",
	//	"amp": "&",
	//	"apos": "'",
	//	"quot": `"`,
	Entity map[string]string

	// DefaultSpace sets the default name space used for unadorned tags,
	// as if the entire XML stream were wrapped in an element containing
	// the attribute xmlns="DefaultSpace".
	DefaultSpace string
	names        map[[8]byte]string
	r            *bufioReader
	buf          bytes.Buffer
	stk          *stack
	free         *stack
	needClose    bool
	toClose      goxml.Name
	nextByte     int
	ns           map[string]string
	err          error
	offset       int64
	attrPool     []xml.Attr
}

// NewDecoder creates a new XML parser reading from r.
// If r does not implement io.ByteReader, NewDecoder will
// do its own buffering.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		ns:       make(map[string]string),
		names:    make(map[[8]byte]string),
		nextByte: -1,
		attrPool: make([]xml.Attr, 10),
		r: &bufioReader{
			buf: make([]byte, defaultBufSize),
			rd:  r,
		},
	}
	return d
}

func (d *Decoder) handleStartElement(t goxml.StartElement) {
	for _, a := range t.Attr {
		if a.Name.Space == xmlnsPrefix {
			v, ok := d.ns[a.Name.Local]
			d.pushNs(a.Name.Local, v, ok)
			d.ns[a.Name.Local] = a.Value
		}
		if a.Name.Space == "" && a.Name.Local == xmlnsPrefix {
			// Default space for untagged names
			v, ok := d.ns[""]
			d.pushNs("", v, ok)
			d.ns[""] = a.Value
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
	} else if n.Space == "" {
		n.Space = d.DefaultSpace
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

	if b, ok = d.mustgetc(); !ok {
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
		if b, ok = d.mustgetc(); !ok {
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
		var target string
		if target, ok = d.name(); !ok {
			if d.err == nil {
				d.err = d.syntaxError("expected target name after <?")
			}
			return d.err
		}
		d.space()
		d.buf.Reset()
		var b0 byte
		for {
			if b, ok = d.mustgetc(); !ok {
				return d.err
			}
			d.buf.WriteByte(b)
			if b0 == '?' && b == '>' {
				break
			}
			b0 = b
		}
		data := d.buf.Bytes()
		data = data[0 : len(data)-2] // chop ?>

		if target == "xml" {
			content := string(data)
			ver := procInst("version", content)
			if ver != "" && ver != "1.0" {
				d.err = fmt.Errorf("xml: unsupported version %q; only version 1.0 is supported", ver)
				return d.err
			}
			enc := procInst("encoding", content)
			if enc != "" && enc != "utf-8" && enc != "UTF-8" && !strings.EqualFold(enc, "utf-8") {
				d.err = fmt.Errorf("xml: encoding %q declared but Decoder.CharsetReader is nil", enc)
				return d.err
			}
		}
		return nil // goxml.ProcInst{Target: target, Inst: data}, nil

	case '!':
		// <!: Maybe comment, maybe CDATA.
		if b, ok = d.mustgetc(); !ok {
			return d.err
		}
		switch b {
		case '-': // <!-
			// Probably <!-- for a comment.
			if b, ok = d.mustgetc(); !ok {
				return d.err
			}
			if b != '-' {
				d.err = d.syntaxError("invalid sequence <!- not part of <!--")
				return d.err
			}
			// Look for terminator.
			d.buf.Reset()
			var b0, b1 byte
			for {
				if b, ok = d.mustgetc(); !ok {
					return d.err
				}
				d.buf.WriteByte(b)
				if b0 == '-' && b1 == '-' {
					if b != '>' {
						d.err = d.syntaxError(
							`invalid sequence "--" not allowed in comments`)
						return d.err
					}
					break
				}
				b0, b1 = b1, b
			}
			data := d.buf.Bytes()
			data = data[0 : len(data)-3] // chop -->
			return nil

		case '[': // <![
			// Probably <![CDATA[.
			for i := 0; i < 6; i++ {
				if b, ok = d.mustgetc(); !ok {
					return d.err
				}
				if b != "CDATA["[i] {
					d.err = d.syntaxError("invalid <![ sequence")
					return d.err
				}
			}
			// CDATA not supported
			return d.err
		}

		// Probably a directive: <!DOCTYPE ...>, <!ENTITY ...>, etc.
		// We don't care, but accumulate for caller. Quoted angle
		// brackets do not count for nesting.
		d.buf.Reset()
		d.buf.WriteByte(b)
		inquote := uint8(0)
		depth := 0
		for {
			if b, ok = d.mustgetc(); !ok {
				return d.err
			}
			if inquote == 0 && b == '>' && depth == 0 {
				break
			}
		HandleB:
			d.buf.WriteByte(b)
			switch {
			case b == inquote:
				inquote = 0

			case inquote != 0:
				// in quotes, no special action

			case b == '\'' || b == '"':
				inquote = b

			case b == '>' && inquote == 0:
				depth--

			case b == '<' && inquote == 0:
				// Look for <!-- to begin comment.
				s := "!--"
				for i := 0; i < len(s); i++ {
					if b, ok = d.mustgetc(); !ok {
						return d.err
					}
					if b != s[i] {
						for j := 0; j < i; j++ {
							d.buf.WriteByte(s[j])
						}
						depth++
						goto HandleB
					}
				}

				// Remove < that was written above.
				d.buf.Truncate(d.buf.Len() - 1)

				// Look for terminator.
				var b0, b1 byte
				for {
					if b, ok = d.mustgetc(); !ok {
						return d.err
					}
					if b0 == '-' && b1 == '-' && b == '>' {
						break
					}
					b0, b1 = b1, b
				}
			}
		}
		return nil // goxml.Directive(d.buf.Bytes()), nil
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
		if b, ok = d.mustgetc(); !ok {
			return d.err
		}
		if b == '/' {
			empty = true
			if b, ok = d.mustgetc(); !ok {
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

		a := goxml.Attr{}
		if a.Name, ok = d.nsname(); !ok {
			if d.err == nil {
				d.err = d.syntaxError("expected attribute name in element")
			}
			return d.err
		}
		d.space()
		if b, ok = d.mustgetc(); !ok {
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
		a.Value = string(data)
		i++
		if len(d.attrPool) < i {
			d.attrPool = append(d.attrPool, make([]xml.Attr, len(d.attrPool))...)
		}
		d.attrPool[i-1] = a
	}
	if empty {
		d.needClose = true
		d.toClose = name
	}
	d.handleStartElement(goxml.StartElement{Name: name, Attr: d.attrPool[:i]})
	return nil
}

func (d *Decoder) attrval() []byte {
	b, ok := d.mustgetc()
	if !ok {
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
	if d.err != nil {
		return 0, false
	}
	if d.nextByte >= 0 {
		b = byte(d.nextByte)
		d.nextByte = -1
	} else {
		b, d.err = d.r.ReadByte()
		if d.err != nil {
			return 0, false
		}
	}
	d.offset++
	return b, true
}

// InputOffset returns the input stream byte offset of the current decoder position.
// The offset gives the location of the end of the most recently returned token
// and the beginning of the next token.
func (d *Decoder) InputOffset() int64 {
	return d.offset
}

// Must read a single byte.
// If there is no byte to read,
// set d.err to SyntaxError("unexpected EOF")
// and return ok==false
func (d *Decoder) mustgetc() (b byte, ok bool) {
	if b, ok = d.getc(); !ok && d.err == io.EOF {
		d.err = d.syntaxError("unexpected EOF")
	}
	return
}

// Unread a single byte.
func (d *Decoder) ungetc(b byte) {
	d.nextByte = int(b)
	d.offset--
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
			if b, ok = d.mustgetc(); !ok {
				return nil
			}
			if b == '#' {
				d.buf.WriteByte(b)
				if b, ok = d.mustgetc(); !ok {
					return nil
				}
				base := 10
				if b == 'x' {
					base = 16
					d.buf.WriteByte(b)
					if b, ok = d.mustgetc(); !ok {
						return nil
					}
				}
				start := d.buf.Len()
				for '0' <= b && b <= '9' ||
					base == 16 && 'a' <= b && b <= 'f' ||
					base == 16 && 'A' <= b && b <= 'F' {
					d.buf.WriteByte(b)
					if b, ok = d.mustgetc(); !ok {
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
						text = string(n)
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
				if b, ok = d.mustgetc(); !ok {
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
							text = string(r)
							haveText = true
						} else if d.Entity != nil {
							text, haveText = d.Entity[s]
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
	return name, true
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
	var arr [8]byte
	b := d.buf.Bytes()
	if len(b) <= 8 {
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
	if len(b) <= 8 {
		d.names[arr] = s
	}
	return s, true
}

// Read a name and append its bytes to d.buf.
// The name is delimited by any single-byte character not valid in names.
// All multi-byte characters are accepted; the caller must check their validity.
func (d *Decoder) readName() (ok bool) {
	const runSelf = 0x80 // characters below RuneSelf are represented as themselves in a single byte.
	var b byte
	if b, ok = d.mustgetc(); !ok {
		return
	}
	if b < runSelf && !isNameByte(b) {
		d.ungetc(b)
		return false
	}
	d.buf.WriteByte(b)

	for {
		if b, ok = d.mustgetc(); !ok {
			return
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

// procInst parses the `param="..."` or `param='...'`
// value out of the provided string, returning "" if not found.
func procInst(param, s string) string {
	// TODO: this parsing is somewhat lame and not exact.
	// It works for all actual cases, though.
	param = param + "="
	idx := strings.Index(s, param)
	if idx == -1 {
		return ""
	}
	v := s[idx+len(param):]
	if v == "" {
		return ""
	}
	if v[0] != '\'' && v[0] != '"' {
		return ""
	}
	idx = strings.IndexRune(v[1:], rune(v[0]))
	if idx == -1 {
		return ""
	}
	return v[1 : idx+1]
}
