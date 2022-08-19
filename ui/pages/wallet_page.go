package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

const (
	WalletPageID = "Wallet"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	NavDrawerWidth          = unit.Dp(160)
	NavDrawerMinimizedWidth = unit.Dp(72)
)

//
//type MainPage struct {
//	*app.MasterPage
//
//	*handlers.Load
//	*listeners.SyncProgressListener
//	*listeners.TxAndBlockNotificationListener
//	*listeners.ProposalNotificationListener
//	ctx                  context.Context
//	ctxCancel            context.CancelFunc
//	drawerNav            components.NavDrawer
//	bottomNavigationBar  components.BottomNavigationBar
//	floatingActionButton components.BottomNavigationBar
//
//	hideBalanceItem HideBalanceItem
//
//	sendPage    *send.Page   // reuse value to keep data persistent onresume.
//	receivePage *ReceivePage // pointer to receive page. to avoid duplication.
//
//	refreshExchangeRateBtn *decredmaterial.Clickable
//	darkmode               *decredmaterial.Clickable
//	openWalletSelector     *decredmaterial.Clickable
//
//	// page state variables
//	totalBalance dcrutil.Amount
//}
