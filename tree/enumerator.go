package tree

type Enumerator func(children Children, index int) string

func DefaultEnumerator(children Children, index int) string {
	if children.Length()-1 == index {
		return "└──"
	}
	return "├──"
}

func RoundedEnumerator(children Children, index int) string {
	if children.Length()-1 == index {
		return "╰──"
	}
	return "├──"
}

type Indenter func(children Children, index int) string

func DefaultIndenter(children Children, index int) string {
	if children.Length()-1 == index {
		return "   "
	}
	return "│  "
}
