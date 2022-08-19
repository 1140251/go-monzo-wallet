package pages

import (
	"context"
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-monzo-wallet/internal"
	"go-monzo-wallet/ui/components"
	"go-monzo-wallet/ui/handlers"
	"go-monzo-wallet/ui/modal"
	"go-monzo-wallet/ui/values"
	"golang.org/x/oauth2"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	StartPageID = "start_page"
	cfgFile     = "config.json"
)

type startPage struct {
	*handlers.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*modal.GenericPageModal

	loading   bool
	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	listLock        sync.Mutex
	scrollContainer *widget.List

	mainAccountsList internal.Accounts

	shadowBox    *components.Shadow
	accountsList *components.ClickableList

	wallectSelected func()
}

func NewStartPage(l *handlers.Load) handlers.Page {
	sp := &startPage{
		Load:             l,
		GenericPageModal: modal.NewGenericPageModal(StartPageID),
		loading:          true,
		scrollContainer: &widget.List{
			List: layout.List{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			},
		},
		shadowBox: l.Theme.Shadow(),
	}

	sp.accountsList = l.Theme.NewClickableList(layout.Vertical)

	return sp
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (sp *startPage) OnNavigatedTo() {

	if sp.WL.LoadedWallet() {
		sp.loading = false
	} else {
		err := sp.openWallet()
		if err != nil {

			startupPasswordModal := modal.NewInfoModal(sp.Load).
				Title(values.String(values.StrUnlockWithPassword)).
				NegativeButton(values.String(values.StrExit), func() {
					sp.WL.Shutdown()
					os.Exit(0)
				})
			sp.ParentWindow().ShowModal(startupPasswordModal)
		}

	}

}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (sp *startPage) HandleUserInteractions() {

	sp.listLock.Lock()
	mainWalletList := sp.mainAccountsList
	sp.listLock.Unlock()

	if ok, selectedItem := sp.accountsList.ItemClicked(); ok {
		sp.WL.SelectedAccount = mainWalletList[selectedItem]
		sp.wallectSelected()
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (sp *startPage) OnNavigatedFrom() {}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (sp *startPage) Layout(gtx values.C) values.D {
	if sp.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return sp.layoutMobile(gtx)
	}
	return sp.layoutDesktop(gtx)
}

// Desktop layout
func (sp *startPage) layoutDesktop(gtx values.C) values.D {
	gtx.Constraints.Min = gtx.Constraints.Max // use maximum height & width
	return layout.Flex{
		Alignment: layout.Middle,
		Axis:      layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx values.C) values.D {
			return sp.loadingSection(gtx)
		}),
	)
}

func (sp *startPage) loadingSection(gtx values.C) values.D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X // use maximum width
	if sp.loading {
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	} else {
		gtx.Constraints.Min.Y = (gtx.Constraints.Max.Y * 65) / 100 // use 65% of view height
	}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx values.C) values.D {
			return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx values.C) values.D {
					return layout.Center.Layout(gtx, func(gtx values.C) values.D {
						return components.NewImage(sp.Theme.Icons.MonzoLogo).LayoutSize(gtx, values.MarginPadding150)
						return values.D{}
					})
				}),
				layout.Rigid(func(gtx values.C) values.D {
					if sp.loading {
						loadStatus := sp.Theme.Text(values.TextSize20, values.String(values.StrLoading))
						if sp.WL.LoadedWallet() {
							loadStatus.Text = values.String(values.StrOpeningWallet)
						}

						return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, loadStatus.Layout)
					}
					pageContent := []func(gtx values.C) values.D{
						sp.Theme.Text(values.TextSize20, values.String(values.StrSelectWalletToOpen)).Layout,
						sp.walletSection, // wallet list layout
					}

					gtx.Constraints.Min = gtx.Constraints.Max
					return components.UniformPadding(gtx, func(gtx values.C) values.D {
						gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding550)
						list := &layout.List{
							Axis: layout.Vertical,
						}

						return layout.Center.Layout(gtx, func(gtx values.C) values.D {
							return list.Layout(gtx, len(pageContent), func(gtx values.C, i int) values.D {
								return layout.Inset{Top: values.MarginPadding26}.Layout(gtx, func(gtx values.C) values.D {
									return pageContent[i](gtx)
								})
							})
						})
					})
				}),
			)
		}),
	)
}

// Mobile layout
func (sp *startPage) layoutMobile(gtx values.C) values.D {
	gtx.Constraints.Min = gtx.Constraints.Max // use maximum height & width
	return layout.Flex{
		Alignment: layout.Middle,
		Axis:      layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx values.C) values.D {
			return sp.loadingSection(gtx)
		}),
		//layout.Rigid(func(gtx values.C) values.D {
		//	if sp.loading {
		//		return values.D{}
		//	}
		//
		//	gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding350)
		//	return layout.Inset{
		//		Left:  values.MarginPadding24,
		//		Right: values.MarginPadding24,
		//	}.Layout(gtx, sp.addWalletButton.Layout)
		//}),
	)
}

func (sp *startPage) openWallet() error {
	cfg, err := initConfig()
	if err != nil {
		logrus.Info("reading config:", err)
		// show err dialog
		return err
	}

	sp.wallectSelected = func() {
		//sp.ParentNavigator().ClearStackAndDisplay(NewWalletPage(sp.Load))
	}
	token, err := sp.WL.Connect(cfg)
	if err != nil {
		logrus.Info("connecting to monzo:", err)
		// show err dialog
		return err
	}

	startupPasswordModal := modal.NewInfoModal(sp.Load).
		Title(values.String(values.StrUnlockWithPassword)).
		Body(values.String(values.StrStartupPassword)).
		NegativeButton(values.String(values.StrExit), func() {
			os.Exit(0)
		})

	startupPasswordModal.PositiveButton(values.String(values.StrUnlock), func(bool) bool {
		err = retry(5, time.Second*3, func() error {
			err = sp.WL.FetchAccounts(token.AccessToken)
			if err != nil {

				logrus.Info("fetching accounts:", err)
				sp.Toast.NotifyError(err.Error())
				return err
			}

			accounts := sp.WL.AccountsList()

			sp.listLock.Lock()
			sp.mainAccountsList = accounts
			sp.listLock.Unlock()

			return nil
		})
		if err != nil {
			logrus.Info("retry error:", err)
			sp.Toast.NotifyError(err.Error())
			return false
		}

		sp.loading = false
		startupPasswordModal.SetLoading(false)
		startupPasswordModal.Dismiss()
		sp.ParentWindow().DismissModal(startupPasswordModal.ID())
		return true
	})
	sp.ParentWindow().ShowModal(startupPasswordModal)
	return nil

}

func initConfig() (*oauth2.Config, error) {

	// Use config file from the flag.
	viper.SetConfigFile(cfgFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg oauth2.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}

func (sp *startPage) walletList(gtx values.C) values.D {
	sp.listLock.Lock()
	mainWalletList := sp.mainAccountsList
	sp.listLock.Unlock()

	return sp.accountsList.Layout(gtx, len(mainWalletList), func(gtx values.C, i int) values.D {
		return sp.walletWrapper(gtx, mainWalletList[i])
	})
}

func (sp *startPage) walletSection(gtx values.C) values.D {
	walletSections := []func(gtx values.C) values.D{
		sp.walletList,
	}

	return sp.Theme.List(sp.scrollContainer).Layout(gtx, len(walletSections), func(gtx values.C, i int) values.D {
		return walletSections[i](gtx)
	})
}

func (sp *startPage) walletWrapper(gtx values.C, item *internal.Account) values.D {
	sp.shadowBox.SetShadowRadius(14)
	return components.LinearLayout{
		Width:      components.WrapContent,
		Height:     components.WrapContent,
		Padding:    layout.UniformInset(values.MarginPadding9),
		Background: sp.Theme.Color.Surface,
		Alignment:  layout.Middle,
		Shadow:     sp.shadowBox,
		Margin:     layout.UniformInset(values.MarginPadding5),
		Border:     components.Border{Radius: components.NewRadius(14)},
	}.Layout(gtx,
		layout.Rigid(func(gtx values.C) values.D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx values.C) values.D {
					return sp.Theme.Text(values.TextSize16, item.AccountNumber).Layout(gtx)
				}),
				layout.Rigid(func(gtx values.C) values.D {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(sp.syncStatusIcon),
						layout.Rigid(func(gtx values.C) values.D {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx values.C) values.D {
									ic := components.NewIcon(sp.Theme.Icons.ImageBrightness1)
									ic.Color = sp.Theme.Color.Gray1
									return layout.Inset{
										Left:  values.MarginPadding7,
										Right: values.MarginPadding7,
									}.Layout(gtx, func(gtx values.C) values.D {
										return ic.Layout(gtx, values.MarginPadding4)
									})
								}),
							)

						}),
					)
				}),
			)
		}),
		layout.Flexed(1, func(gtx values.C) values.D {
			balanceLabel := sp.Theme.Body1(fmt.Sprint(roundFloat(item.Balance, 2)))
			balanceLabel.Color = sp.Theme.Color.GrayText2
			return layout.Inset{
				Right: values.MarginPadding10,
			}.Layout(gtx, func(gtx values.C) values.D {
				return layout.E.Layout(gtx, balanceLabel.Layout)
			})
		}),
	)
}

func (sp *startPage) syncStatusIcon(gtx values.C) values.D {
	var (
		syncStatusIcon *components.Image
		syncStatus     string
	)

	switch {
	case sp.WL.LoadedWallet():
		syncStatusIcon = components.NewImage(sp.Theme.Icons.SuccessIcon)
		syncStatus = values.String(values.StrSynced)
	default:
		syncStatusIcon = components.NewImage(sp.Theme.Icons.FailedIcon)
		syncStatus = values.String(values.StrWalletNotSynced)
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(syncStatusIcon.Layout16dp),
		layout.Rigid(func(gtx values.C) values.D {
			return layout.Inset{
				Left: values.MarginPadding5,
			}.Layout(gtx, sp.Theme.Caption(syncStatus).Layout)
		}),
	)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
