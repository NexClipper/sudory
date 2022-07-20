package vanilla_test

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

func Test_Equal(t *testing.T) {
	c := vanilla.Equal("foo", "bar").Parse()
	t.Log(c.Query(), c.Args())
}

func Test_IsNull(t *testing.T) {
	c := vanilla.IsNull("foo").Parse()
	t.Log(c.Query(), c.Args())
}

func Test_In(t *testing.T) {
	c := vanilla.In("foo", 1, 2, 3, 4).Parse()
	t.Log(c.Query(), c.Args())
}

func Test_GreaterThan(t *testing.T) {
	c := vanilla.GreaterThan("foo", 1).Parse()
	t.Log(c.Query(), c.Args())
}

func Test_GreaterThanEqual(t *testing.T) {
	c := vanilla.GreaterThanEqual("foo", 1).Parse()
	t.Log(c.Query(), c.Args())
}

func Test_LessThan(t *testing.T) {
	c := vanilla.LessThan("foo", 1).Parse()
	t.Log(c.Query(), c.Args())
}

func Test_LessThanEqual(t *testing.T) {
	c := vanilla.LessThanEqual("foo", 1).Parse()
	t.Log(c.Query(), c.Args())
}

func Test_And(t *testing.T) {
	c := vanilla.And(vanilla.Equal("c1", 1), vanilla.Equal("c2", 2)).Parse()
	t.Log(c.Query(), c.Args())
}

func Test_Or(t *testing.T) {
	c := vanilla.Or(vanilla.Equal("c1", 1), vanilla.Equal("c2", 2)).Parse()
	t.Log(c.Query(), c.Args())
}

func Test_AndOr(t *testing.T) {
	c := vanilla.And(vanilla.Equal("c1", 1), vanilla.Equal("c2", 2),
		vanilla.Or(vanilla.Equal("c3", 3), vanilla.Equal("c4", 4))).Parse()
	t.Log(c.Query(), c.Args())
}

func Test_Asc(t *testing.T) {
	c := vanilla.Asc("a1", "b1", "c1").Parse()
	t.Log(c.Order())
}

func Test_Desc(t *testing.T) {
	c := vanilla.Desc("a2", "b2", "c2").Parse()
	t.Log(c.Order())
}

func Test_AscDesc(t *testing.T) {
	c := vanilla.Asc("a1", "b1", "c1").Combine(
		vanilla.Desc("a2", "b2", "c2"),
		vanilla.Asc("a3", "b3", "c3"),
	).Parse()
	t.Log(c.Order())
}
