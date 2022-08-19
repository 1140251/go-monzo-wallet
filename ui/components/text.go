package components

import (
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Text struct {
	material.LabelStyle
}

func (t *Theme) H1(txt string) Text {
	return t.labelWithDefaultColor(Text{material.H1(t.Base, txt)})
}

func (t *Theme) H2(txt string) Text {
	return t.labelWithDefaultColor(Text{material.H2(t.Base, txt)})
}

func (t *Theme) H3(txt string) Text {
	return t.labelWithDefaultColor(Text{material.H2(t.Base, txt)})
}

func (t *Theme) H4(txt string) Text {
	return t.labelWithDefaultColor(Text{material.H4(t.Base, txt)})
}

func (t *Theme) H5(txt string) Text {
	return t.labelWithDefaultColor(Text{material.H5(t.Base, txt)})
}

func (t *Theme) H6(txt string) Text {
	return t.labelWithDefaultColor(Text{material.H6(t.Base, txt)})
}

func (t *Theme) Body1(txt string) Text {
	return t.labelWithDefaultColor(Text{material.Body1(t.Base, txt)})
}

func (t *Theme) Body2(txt string) Text {
	return t.labelWithDefaultColor(Text{material.Body2(t.Base, txt)})
}

func (t *Theme) Caption(txt string) Text {
	return t.labelWithDefaultColor(Text{material.Caption(t.Base, txt)})
}

func (t *Theme) ErrorLabel(txt string) Text {
	label := t.Caption(txt)
	label.Color = t.Color.Danger
	return label
}

func (t *Theme) Text(size unit.Sp, txt string) Text {
	return t.labelWithDefaultColor(Text{material.Label(t.Base, size, txt)})
}

func (t *Theme) labelWithDefaultColor(l Text) Text {
	l.Color = t.Color.DeepBlue
	return l
}
