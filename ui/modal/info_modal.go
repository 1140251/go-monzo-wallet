package modal

import (
	"go-monzo-wallet/ui/components"
	"go-monzo-wallet/ui/handlers"
	"go-monzo-wallet/ui/renderers"
	"go-monzo-wallet/ui/values"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
)

const InfoID = "info_page"

type InfoModal struct {
	*handlers.Load
	*components.Modal
	*GenericPageModal

	enterKeyPressed bool

	dialogIcon *components.Icon

	dialogTitle    string
	subtitle       string
	customTemplate []layout.Widget
	customWidget   layout.Widget

	positiveButtonText    string
	positiveButtonClicked func(isChecked bool) bool
	btnPositve            components.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           components.Button

	checkbox      components.CheckBoxStyle
	mustBeChecked bool

	titleAlignment, btnAlignment layout.Direction
	materialLoader               material.LoaderStyle

	isCancelable bool
	isLoading    bool
}

func NewInfoModal(l *handlers.Load) *InfoModal {
	return NewInfoModalWithKey(l, "info_modal")
}

func NewInfoModalWithKey(l *handlers.Load, key string) *InfoModal {

	in := &InfoModal{
		Load:             l,
		Modal:            l.Theme.ModalFloatTitle(key),
		btnPositve:       l.Theme.OutlineButton(values.String(values.StrYes)),
		btnNegative:      l.Theme.OutlineButton(values.String(values.StrNo)),
		isCancelable:     true,
		isLoading:        false,
		btnAlignment:     layout.E,
		GenericPageModal: NewGenericPageModal(InfoID),
	}

	in.btnPositve.Font.Weight = text.Medium
	in.btnNegative.Font.Weight = text.Medium

	in.materialLoader = material.Loader(l.Theme.Base)

	return in
}

func (in *InfoModal) OnResume() {}

func (in *InfoModal) OnDismiss() {}

func (in *InfoModal) SetCancelable(min bool) *InfoModal {
	in.isCancelable = min
	return in
}

func (in *InfoModal) SetContentAlignment(title, btn layout.Direction) *InfoModal {
	in.titleAlignment = title
	in.btnAlignment = btn
	return in
}

func (in *InfoModal) Icon(icon *components.Icon) *InfoModal {
	in.dialogIcon = icon
	return in
}

func (in *InfoModal) CheckBox(checkbox components.CheckBoxStyle, mustBeChecked bool) *InfoModal {
	in.checkbox = checkbox
	in.mustBeChecked = mustBeChecked // determine if the checkbox must be selected to proceed
	return in
}

func (in *InfoModal) SetLoading(loading bool) {
	in.isLoading = loading
	in.Modal.SetDisabled(loading)
}

func (in *InfoModal) Title(title string) *InfoModal {
	in.dialogTitle = title
	return in
}

func (in *InfoModal) Body(subtitle string) *InfoModal {
	in.subtitle = subtitle
	return in
}

func (in *InfoModal) PositiveButton(text string, clicked func(isChecked bool) bool) *InfoModal {
	in.positiveButtonText = text
	in.positiveButtonClicked = clicked
	return in
}

func (in *InfoModal) PositiveButtonStyle(background, text color.NRGBA) *InfoModal {
	in.btnPositve.Background, in.btnPositve.Color = background, text
	return in
}

func (in *InfoModal) NegativeButton(text string, clicked func()) *InfoModal {
	in.negativeButtonText = text
	in.negativeButtonClicked = clicked
	return in
}

// for backwards compatibilty
func (in *InfoModal) SetupWithTemplate(template string) *InfoModal {
	title := in.dialogTitle
	subtitle := in.subtitle
	var customTemplate []layout.Widget
	switch template {
	case TransactionDetailsInfoTemplate:
		title = values.String(values.StrHowToCopy)
		customTemplate = transactionDetailsInfo(in.Theme)
	case SignMessageInfoTemplate:
		customTemplate = signMessageInfo(in.Theme)
	case VerifyMessageInfoTemplate:
		customTemplate = verifyMessageInfo(in.Theme)
	case PrivacyInfoTemplate:
		title = values.String(values.StrUseMixer)
		customTemplate = privacyInfo(in.Load)
	case SetupMixerInfoTemplate:
		customTemplate = setupMixerInfo(in.Theme)
	case WalletBackupInfoTemplate:
		customTemplate = backupInfo(in.Theme)
	}

	in.dialogTitle = title
	in.subtitle = subtitle
	in.customTemplate = customTemplate
	return in
}

func (in *InfoModal) UseCustomWidget(layout layout.Widget) *InfoModal {
	in.customWidget = layout
	return in
}

// HandleKeyEvent is called when a key is pressed on the current window.
// Satisfies the load.KeyEventHandler interface for receiving key events.
func (in *InfoModal) HandleKeyEvent(evt *key.Event) {
	if (evt.Name == key.NameReturn || evt.Name == key.NameEnter) && evt.State == key.Press {
		in.btnPositve.Click()
		in.ParentWindow().Reload()
	}
}

func (in *InfoModal) Handle() {
	for in.btnPositve.Clicked() {
		if in.isLoading {
			return
		}
		isChecked := false
		if in.checkbox.CheckBox != nil {
			isChecked = in.checkbox.CheckBox.Value
		}

		if in.positiveButtonClicked(isChecked) {
			in.Dismiss()
		}
	}

	for in.btnNegative.Clicked() {
		if !in.isLoading {
			in.Dismiss()
			in.negativeButtonClicked()
		}
	}

	if in.Modal.BackdropClicked(in.isCancelable) {
		if !in.isLoading {
			in.Dismiss()
		}
	}

	if in.checkbox.CheckBox != nil {
		if in.mustBeChecked {
			in.btnNegative.SetEnabled(in.checkbox.CheckBox.Value)
		}
	}
}

func (in *InfoModal) Layout(gtx layout.Context) D {
	icon := func(gtx C) D {
		if in.dialogIcon == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return in.dialogIcon.Layout(gtx, values.MarginPadding50)
			})
		})
	}

	checkbox := func(gtx C) D {
		if in.checkbox.CheckBox == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPaddingMinus5, Left: values.MarginPaddingMinus5}.Layout(gtx, func(gtx C) D {
			in.checkbox.TextSize = values.TextSize14
			in.checkbox.Color = in.Theme.Color.GrayText1
			in.checkbox.IconColor = in.Theme.Color.Gray2
			if in.checkbox.CheckBox.Value {
				in.checkbox.IconColor = in.Theme.Color.Primary
			}
			return in.checkbox.Layout(gtx)
		})
	}

	subtitle := func(gtx C) D {
		text := in.Theme.Body1(in.subtitle)
		text.Color = in.Theme.Color.GrayText2
		return text.Layout(gtx)
	}

	var w []layout.Widget

	// Every section of the dialog is optional
	if in.dialogIcon != nil {
		w = append(w, icon)
	}

	if in.dialogTitle != "" {
		w = append(w, in.titleLayout())
	}

	if in.subtitle != "" {
		w = append(w, subtitle)
	}

	if in.customTemplate != nil {
		w = append(w, in.customTemplate...)
	}

	if in.checkbox.CheckBox != nil {
		w = append(w, checkbox)
	}

	if in.customWidget != nil {
		w = append(w, in.customWidget)
	}

	if in.negativeButtonText != "" || in.positiveButtonText != "" {
		w = append(w, in.actionButtonsLayout())
	}

	return in.Modal.Layout(gtx, w)
}

func (in *InfoModal) titleLayout() layout.Widget {
	return func(gtx C) D {
		t := in.Theme.H6(in.dialogTitle)
		t.Font.Weight = text.SemiBold
		return in.titleAlignment.Layout(gtx, t.Layout)
	}
}

func (in *InfoModal) actionButtonsLayout() layout.Widget {
	return func(gtx C) D {
		return in.btnAlignment.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if in.negativeButtonText == "" || in.isLoading {
						return layout.Dimensions{}
					}

					in.btnNegative.Text = in.negativeButtonText
					gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding250)
					return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, in.btnNegative.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					if in.isLoading {
						return in.materialLoader.Layout(gtx)
					}

					if in.positiveButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnPositve.Text = in.positiveButtonText
					gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding250)
					return in.btnPositve.Layout(gtx)
				}),
			)
		})
	}
}

const (
	VerifyMessageInfoTemplate      = "VerifyMessageInfo"
	SignMessageInfoTemplate        = "SignMessageInfo"
	PrivacyInfoTemplate            = "PrivacyInfo"
	SetupMixerInfoTemplate         = "ConfirmSetupMixer"
	TransactionDetailsInfoTemplate = "TransactionDetailsInfoInfo"
	WalletBackupInfoTemplate       = "WalletBackupInfo"
	AllowUnmixedSpendingTemplate   = "AllowUnmixedSpending"
	TicketPriceErrorTemplate       = "TicketPriceError"
	SecurityToolsInfoTemplate      = "SecurityToolsInfo"
)

func verifyMessageInfo(th *components.Theme) []layout.Widget {
	text := values.StringF(values.StrVerifyMessageInfo, `<span style="text-color: gray">`, `<br />`, `<font color="success">`, `</font>`, `<font color="danger">`, `</font>`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func signMessageInfo(th *components.Theme) []layout.Widget {
	text := values.StringF(values.StrSignMessageInfo, `<span style="text-color: gray">`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func privacyInfo(l *handlers.Load) []layout.Widget {
	text := values.StringF(values.StrPrivacyInfo, `<span style="text-color: gray">`, `<br/><span style="font-weight: bold">`, `</span></br>`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, l.Theme).Layout,
	}
}

func setupMixerInfo(th *components.Theme) []layout.Widget {
	text := values.StringF(values.StrSetupMixerInfo, `<span style="text-color: gray">`, `<span style="font-weight: bold">`, `</span>`, `<span style="font-weight: bold">`, `</span>`, `<br> <span style="font-weight: bold">`, `</span></span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func transactionDetailsInfo(th *components.Theme) []layout.Widget {
	text := values.StringF(values.StrTxdetailsInfo, `<span style="text-color: gray">`, `<span style="text-color: primary">`, `</span>`, `</span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func backupInfo(th *components.Theme) []layout.Widget {
	text := values.StringF(values.StrBackupInfo, `<span style="text-color: danger"> <span style="font-weight: bold">`, `</span>`, `<span style="font-weight: bold">`, `</span>`, `<span style="font-weight: bold">`, `</span></span>`)
	return []layout.Widget{
		renderers.RenderHTML(text, th).Layout,
	}
}

func allowUnspendUnmixedAcct(l *handlers.Load) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := components.NewIcon(l.Theme.Icons.ActionInfo)
					ic.Color = l.Theme.Color.GrayText1
					return layout.Inset{Top: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
						return ic.Layout(gtx, values.MarginPadding18)
					})
				}),
				layout.Rigid(func(gtx C) D {
					text := values.StringF(values.StrAllowUnspendUnmixedAcct, `<span style="text-color: gray">`, `<br>`, `<span style="font-weight: bold">`, `</span>`, `</span>`)
					return renderers.RenderHTML(text, l.Theme).Layout(gtx)
				}),
			)
		},
	}
}
