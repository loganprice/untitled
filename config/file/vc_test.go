package file

type validTestStruct struct {
	Foo    string   `yaml:"foo" json:"foo,omitempty"`
	Bar    string   `yaml:"bar" json:"bar,omitempty"`
	FooBar []string `yaml:"foo_bar" json:"foo_bar,omitempty"`
}

const goodYAML = `
foo: Bar
bar: Baz
foo_bar:
- Foo
- Bar
`

const goodJSON = `
{
  "foo": "Bar",
  "bar": "Baz",
  "foo_bar": ["Foo", "Bar"]
}`

const goodTOML = `
foo = "Bar"
bar = "Baz"
foo_bar = [
	"Foo",
	"Bar"
]	
`
