package values

import (
	"gioui.org/widget"
	"go-monzo-wallet/ui/assets"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
)

type Icons struct {
	ContentAdd, NavigationCheck, NavigationMore, ActionCheckCircle, ActionInfo, NavigationArrowBack,
	NavigationArrowForward, ActionCheck, ChevronRight, NavigationCancel, NavMoreIcon,
	ImageBrightness1, ContentClear, DropDownIcon, Cached, ContentRemove, ConcealIcon, RevealIcon,
	SearchIcon, PlayIcon *widget.Icon

	MonzoLogo, SuccessIcon, FailedIcon, RedAlert image.Image
}

func (i *Icons) StandardMaterialIcons() *Icons {
	icon := MustIcon(widget.NewIcon(icons.ActionInfo))
	i.ActionInfo = icon

	return i
}

func (i *Icons) DefaultIcons() *Icons {
	decredIcons := assets.Icons

	i.StandardMaterialIcons()

	i.MonzoLogo = decredIcons["monzo_logo"]
	i.SuccessIcon = decredIcons["success_check"]
	i.FailedIcon = decredIcons["crossmark_red"]
	i.RedAlert = decredIcons["red_alert"]

	i.ImageBrightness1 = MustIcon(widget.NewIcon(icons.ImageBrightness1))
	return i
}

func (i *Icons) DarkModeIcons() *Icons {
	return i
}

func MustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
