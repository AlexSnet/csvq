package parser

import (
	"testing"

	"github.com/mithrandie/csvq/lib/value"
)

func TestFormatFieldIdentifier(t *testing.T) {
	var e QueryExpression = NewStringValue("str")
	expect := "str"
	result := FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = NewDatetimeValueFromString("2006-01-02 15:04:05 -08:00")
	expect = "2006-01-02T15:04:05-08:00"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = NewIntegerValue(1)
	expect = "1"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = FieldReference{Column: Identifier{Literal: "column1"}}
	expect = "column1"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}

	e = ColumnNumber{View: Identifier{Literal: "table1"}, Number: value.NewInteger(1)}
	expect = "table1.1"
	result = FormatFieldIdentifier(e)
	if result != expect {
		t.Errorf("field identifier = %q, want %q for %#v", result, expect, e)
	}
}
