package query

import (
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/value"
)

type aggregateTests struct {
	List   []value.Primary
	Result value.Primary
}

var countTests = []aggregateTests{
	{
		List: []value.Primary{
			value.NewInteger(2),
			value.NewInteger(1),
			value.NewInteger(1),
			value.NewNull(),
			value.NewInteger(4),
			value.NewNull(),
		},
		Result: value.NewInteger(4),
	},
	{
		List: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewInteger(0),
	},
}

func TestCount(t *testing.T) {
	for _, v := range countTests {
		r := Count(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("count list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var maxTests = []aggregateTests{
	{
		List: []value.Primary{
			value.NewInteger(2),
			value.NewInteger(1),
			value.NewInteger(1),
			value.NewNull(),
			value.NewInteger(4),
			value.NewNull(),
		},
		Result: value.NewInteger(4),
	},
	{
		List: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestMax(t *testing.T) {
	for _, v := range maxTests {
		r := Max(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("max list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var minTests = []aggregateTests{
	{
		List: []value.Primary{
			value.NewInteger(2),
			value.NewInteger(1),
			value.NewInteger(1),
			value.NewNull(),
			value.NewInteger(4),
			value.NewNull(),
		},
		Result: value.NewInteger(1),
	},
	{
		List: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestMin(t *testing.T) {
	for _, v := range minTests {
		r := Min(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("min list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var sumTests = []aggregateTests{
	{
		List: []value.Primary{
			value.NewInteger(2),
			value.NewInteger(1),
			value.NewInteger(1),
			value.NewNull(),
			value.NewInteger(4),
			value.NewNull(),
		},
		Result: value.NewInteger(8),
	},
	{
		List: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestSum(t *testing.T) {
	for _, v := range sumTests {
		r := Sum(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("sum list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var avgTests = []aggregateTests{
	{
		List: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewNull(),
			value.NewInteger(4),
			value.NewNull(),
		},
		Result: value.NewInteger(2),
	},
	{
		List: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestAvg(t *testing.T) {
	for _, v := range avgTests {
		r := Avg(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("avg list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var medianTests = []aggregateTests{
	{
		List: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(4),
			value.NewInteger(6),
			value.NewNull(),
			value.NewInteger(1),
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewNull(),
		},
		Result: value.NewFloat(1.5),
	},
	{
		List: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(4),
			value.NewInteger(6),
			value.NewNull(),
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewNull(),
		},
		Result: value.NewInteger(2),
	},
	{
		List: []value.Primary{
			value.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			value.NewDatetime(time.Date(2012, 2, 5, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: value.NewInteger(time.Date(2012, 2, 4, 9, 18, 15, 0, GetTestLocation()).Unix()),
	},
	{
		List: []value.Primary{
			value.NewNull(),
		},
		Result: value.NewNull(),
	},
}

func TestMedian(t *testing.T) {
	for _, v := range medianTests {
		r := Median(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("median list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var listAggTests = []struct {
	List      []value.Primary
	Separator string
	Result    value.Primary
}{
	{
		List: []value.Primary{
			value.NewString("str1"),
			value.NewString("str3"),
			value.NewNull(),
			value.NewString("str2"),
			value.NewString("str1"),
			value.NewString("str2"),
		},
		Separator: ",",
		Result:    value.NewString("str1,str3,str2,str1,str2"),
	},
	{
		List: []value.Primary{
			value.NewNull(),
		},
		Separator: ",",
		Result:    value.NewNull(),
	},
}

func TestListAgg(t *testing.T) {
	for _, v := range listAggTests {
		r := ListAgg(v.List, v.Separator)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("Listagg list = %s: separator = %q, result = %s, want %s", v.List, v.Separator, r, v.Result)
		}
	}
}

var jsonAggTests = []struct {
	List   []value.Primary
	Result value.Primary
}{
	{
		List:   []value.Primary{},
		Result: value.NewString("[]"),
	},
	{
		List: []value.Primary{
			value.NewString("str3"),
			value.NewNull(),
			value.NewString("str2"),
		},
		Result: value.NewString("[\"str3\",null,\"str2\"]"),
	},
}

func TestJsonAgg(t *testing.T) {
	for _, v := range jsonAggTests {
		r := JsonAgg(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("JsonAgg list = %s, result = %s, want %s", v.List, r, v.Result)
		}
	}
}
