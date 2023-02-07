package views

import (
	"context"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/PaulYakow/gophkeeper/internal/client/controller"
	"github.com/PaulYakow/gophkeeper/internal/entity"
)

const (
	mainMenu     = "main"
	unitsMenu    = "units"
	pairsPage    = "pairs"
	signForm     = "sign"
	registerFail = "login exist"
)

type Sign int

const (
	undefined Sign = iota
	register
	login
)

type ui struct {
	*tview.Application
	root   *tview.Flex
	body   *tview.Pages
	footer *tview.Flex

	mainMenu *tview.List

	unitsMenu *tview.List

	pairsPage *tview.Flex
	pairsList *tview.List
	pairsInfo *tview.TextView

	signForm     *tview.Form
	registerFail *tview.Modal
}

type View struct {
	ctrl *controller.Controller
	tui  *ui
}

func New(ctrl *controller.Controller) (v *View) {
	v = &View{
		ctrl: ctrl,
		tui: &ui{
			Application: tview.NewApplication(),
			body:        tview.NewPages(),
			signForm:    tview.NewForm(),
		},
	}

	v.createMainMenu()
	v.createRegisterFail()
	v.createUnitsMenu()
	v.createPairsPage()

	v.createFooter()
	v.createRoot()

	v.tui.EnableMouse(true)
	return
}

func (v *View) Run() {
	if err := v.tui.SetRoot(v.tui.root, true).Run(); err != nil {
		panic(err)
	}
}

func (v *View) callSignForm(signType Sign) {
	var regLogin, regPassword string
	v.tui.signForm.AddInputField("login", "", 20, nil, func(login string) {
		regLogin = login
	})

	v.tui.signForm.AddPasswordField("password", "", 20, '*', func(password string) {
		regPassword = password
	})

	v.tui.signForm.AddButton("OK", func() {
		var token string
		var err error

		switch signType {
		case register:
			token, err = v.ctrl.Auth.Register(context.Background(), regLogin, regPassword)
		case login:
			token, err = v.ctrl.Auth.Login(context.Background(), regLogin, regPassword)
		default:
			v.tui.body.SwitchToPage(mainMenu)
		}

		if err != nil {
			v.tui.body.SwitchToPage(registerFail)
			return
		}

		v.ctrl.Token = token
		v.tui.body.SwitchToPage(unitsMenu)
	})

	v.tui.signForm.AddButton("Cancel", func() {
		v.tui.body.SwitchToPage(mainMenu)
	})
}

func (v *View) createMainMenu() {
	v.tui.mainMenu = tview.NewList().
		AddItem("Register", "Sign up new user", 'r', func() {
			v.tui.signForm.Clear(true)
			v.callSignForm(register)
			v.tui.body.SwitchToPage(signForm)
		}).
		AddItem("Login", "Sign in with exist user", 'l', func() {
			v.tui.signForm.Clear(true)
			v.callSignForm(login)
			v.tui.body.SwitchToPage(signForm)
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			v.tui.Stop()
		})

	v.tui.body.AddPage(mainMenu, v.tui.mainMenu, true, true)
	v.tui.body.AddPage(signForm, v.tui.signForm, true, false)
}

func (v *View) createRegisterFail() {
	v.tui.registerFail = tview.NewModal().
		SetText("Login already exist!").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			v.tui.body.SwitchToPage(mainMenu)
		}).
		SetBackgroundColor(tcell.ColorLightCoral)

	v.tui.body.AddPage(registerFail, v.tui.registerFail, true, false)
}

func (v *View) createUnitsMenu() {
	v.tui.unitsMenu = tview.NewList().
		AddItem("Pairs", "show login/password pairs", 'r', func() {
			v.getPairsList()
			v.tui.body.SwitchToPage(pairsPage)
		}).
		AddItem("Notes", "show arbitrary text data", 'l', nil).
		AddItem("Cards", "show bank cards data", 'c', nil).
		AddItem("Binary", "show arbitrary binary data", 'b', nil).
		AddItem("Back", "... to main menu", ' ', func() {
			v.tui.body.SwitchToPage(mainMenu)
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			v.tui.Stop()
		})

	v.tui.body.AddPage(unitsMenu, v.tui.unitsMenu, true, false)
}

func (v *View) createPairsPage() {
	v.tui.pairsList = tview.NewList().ShowSecondaryText(false)
	v.tui.pairsInfo = tview.NewTextView()

	v.tui.pairsPage = tview.NewFlex().
		AddItem(v.tui.pairsList, 0, 1, true).
		AddItem(v.tui.pairsInfo, 0, 3, false)

	v.tui.body.AddPage(pairsPage, v.tui.pairsPage, true, false)
}

func (v *View) getPairsList() {
	pairs, err := v.ctrl.Pairs.ViewAllPairs(context.Background(), v.ctrl.Token)
	if err != nil {
		v.tui.body.SwitchToPage(unitsMenu)
		return
	}

	v.tui.pairsList.Clear()
	for _, pair := range pairs {
		v.tui.pairsList.AddItem(strconv.Itoa(pair.ID), "", ' ', nil)
	}

	v.tui.pairsList.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		v.setPairInfo(pairs[index])
	})
}

func (v *View) setPairInfo(pair entity.PairDTO) {
	var sb strings.Builder

	v.tui.pairsInfo.Clear()
	sb.WriteString(pair.Login)
	sb.WriteString("\n")
	sb.WriteString(pair.Password)
	sb.WriteString("\n")
	sb.WriteString(pair.Metadata)
	sb.WriteString("\n")

	v.tui.pairsInfo.SetText(sb.String())
}

func (v *View) createFooter() {
	clientInfo := tview.NewTextView().
		SetText("version: ")

	clientInfo.SetBorder(true).
		SetTitle("Client info")

	serverInfo := tview.NewTextView().
		SetText("version: \ntarget: ")
	serverInfo.SetBorder(true).
		SetTitle("Server info")

	v.tui.footer = tview.NewFlex().
		AddItem(clientInfo, 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(serverInfo, 0, 1, false)
}

func (v *View) createRoot() {
	v.tui.root = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(v.tui.body, 0, 3, true).
		AddItem(v.tui.footer, 0, 1, false)
}
