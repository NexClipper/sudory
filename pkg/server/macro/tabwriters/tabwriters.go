package tabwriters

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type Writer struct {
	keys []string
	tabwriter.Writer
}

func NewWriter(output io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) *Writer {
	w := &Writer{}
	w.Writer.Init(output, minwidth, tabwidth, padding, padchar, flags)
	return w
}

const missingValue = "(MISSING)"

func (writer *Writer) PrintKeyValue(keysAndValues ...interface{}) {
	convert_values := func(v ...interface{}) []string {
		rst := make([]string, 0, len(v))
		for _, v := range v {
			b := &bytes.Buffer{}

			// The type checks are sorted so that more frequently used ones
			// come first because that is then faster in the common
			// cases. In Kubernetes, ObjectRef (a Stringer) is more common
			// than plain strings
			// (https://github.com/kubernetes/kubernetes/pull/106594#issuecomment-975526235).
			switch v := v.(type) {
			case fmt.Stringer:
				writeStringValue(b, true, stringerToString(v))
			case string:
				writeStringValue(b, true, v)
			case error:
				writeStringValue(b, true, v.Error())
			case []byte:
				// In https://github.com/kubernetes/klog/pull/237 it was decided
				// to format byte slices with "%+q". The advantages of that are:
				// - readable output if the bytes happen to be printable
				// - non-printable bytes get represented as unicode escape
				//   sequences (\uxxxx)
				//
				// The downsides are that we cannot use the faster
				// strconv.Quote here and that multi-line output is not
				// supported. If developers know that a byte array is
				// printable and they want multi-line output, they can
				// convert the value to string before logging it.
				// b.WriteByte('=')
				b.WriteString(fmt.Sprintf("%+q", v))
			default:
				writeStringValue(b, false, fmt.Sprintf("%+v", v))
			}

			rst = append(rst, b.String())
		}

		return rst
	}

	var keys []string = make([]string, 0, len(keysAndValues))
	var values []interface{} = make([]interface{}, 0, len(keysAndValues))

	for i := 0; i < len(keysAndValues); i += 2 {
		var v interface{}
		k := keysAndValues[i]
		if i+1 < len(keysAndValues) {
			v = keysAndValues[i+1]
		} else {
			v = missingValue
		}

		// Keys are assumed to be well-formed according to
		// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/migration-to-structured-logging.md#name-arguments
		// for the sake of performance. Keys with spaces,
		// special characters, etc. will break parsing.
		var strkey string
		if k_, ok := k.(string); ok {
			// Avoid one allocation when the key is a string, which
			// normally it should be.
			strkey = k_
		} else {
			strkey = fmt.Sprintf("%s", k)
		}

		keys = append(keys, strkey)

		values = append(values, v)
	}

	if !writer.CompareKeys(keys...) {
		//flush
		writer.Flush()

		//replace keys
		writer.keys = keys
		//write new names
		writer.Write([]byte(strings.Join(keys, "\t")))
		writer.Write([]byte("\n"))

	}

	//write values
	writer.Write([]byte(strings.Join(convert_values(values...), "\t")))
	writer.Write([]byte("\n"))
}

func (writer *Writer) CompareKeys(keys ...string) bool {
	compare_string := func(a, b string) bool {
		return a == b
	}

	if len(writer.keys) != len(keys) {
		return false
	}

	for i := range writer.keys {
		a := writer.keys[i]
		b := keys[i]
		if !compare_string(a, b) {
			return false
		}
	}
	return true
}

func stringerToString(s fmt.Stringer) (ret string) {
	defer func() {
		if err := recover(); err != nil {
			ret = "nil"
		}
	}()
	ret = s.String()
	return
}

func writeStringValue(b *bytes.Buffer, quote bool, v string) {
	data := []byte(v)
	index := bytes.IndexByte(data, '\n')
	if index == -1 {
		// b.WriteByte('=')
		// if quote {
		// 	// Simple string, quote quotation marks and non-printable characters.
		// 	b.WriteString(strconv.Quote(v))
		// 	return
		// }
		// Non-string with no line breaks.
		b.WriteString(v)
		return
	}

	// Complex multi-line string, show as-is with indention like this:
	// I... "hello world" key=<
	// <tab>line 1
	// <tab>line 2
	//  >
	//
	// Tabs indent the lines of the value while the end of string delimiter
	// is indented with a space. That has two purposes:
	// - visual difference between the two for a human reader because indention
	//   will be different
	// - no ambiguity when some value line starts with the end delimiter
	//
	// One downside is that the output cannot distinguish between strings that
	// end with a line break and those that don't because the end delimiter
	// will always be on the next line.
	// b.WriteString("=<\n")
	b.WriteString("=<")
	for index != -1 {
		// b.WriteByte('\t')
		b.Write(data[0 : index+1])
		data = data[index+1:]
		index = bytes.IndexByte(data, '\n')
	}
	if len(data) == 0 {
		// String ended with line break, don't add another.
		b.WriteString(" >")
	} else {
		// No line break at end of last line, write rest of string and
		// add one.
		// b.WriteByte('\t')
		b.Write(data)
		// b.WriteString("\n >")
		b.WriteString(" >")
	}
}
