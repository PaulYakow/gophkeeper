// Package views реализация терминального интерфейса пользователя (TUI) для клиентского приложения.
package views

import (
	"context"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/PaulYakow/gophkeeper/cmd/client/config"
	"github.com/PaulYakow/gophkeeper/internal/client/controller"
	"github.com/PaulYakow/gophkeeper/internal/entity"
)

const (
	mainMenu     = "main"
	unitsMenu    = "units"
	pairsPage    = "pairs"
	cardsPage    = "cards"
	notesPage    = "notes"
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
	header *tview.TextView
	body   *tview.Pages
	footer *tview.Flex

	mainMenu *tview.List

	unitsMenu *tview.List

	pairsPage *tview.Flex
	pairsList *tview.List
	pairsInfo *tview.TextView

	cardsPage *tview.Flex
	cardsList *tview.List
	cardInfo  *tview.TextView

	notesPage *tview.Flex
	notesList *tview.List
	noteInfo  *tview.TextView

	signForm     *tview.Form
	registerFail *tview.Modal
}

// View обеспечивает взаимодействие TUI и Controller'a
type View struct {
	ctrl *controller.Controller
	cfg  *config.Config
	tui  *ui
}

// New создаёт объект View.
func New(ctrl *controller.Controller, cfg *config.Config) (v *View) {
	v = &View{
		ctrl: ctrl,
		cfg:  cfg,
		tui: &ui{
			Application: tview.NewApplication(),
			body:        tview.NewPages(),
			signForm:    tview.NewForm(),
		},
	}

	v.createHeader()

	v.createMainMenu()
	v.createRegisterFail()
	v.createUnitsMenu()
	v.createPairsPage()
	v.createCardsPage()
	v.createNotesPage()

	v.createFooter()
	v.createRoot()

	v.tui.EnableMouse(true)
	return
}

// Run запускает TUI клиента.
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
			v.switchToMainMenu()
		}

		if err != nil {
			v.tui.body.SwitchToPage(registerFail)
			return
		}

		v.ctrl.Token = token
		v.switchToUnitsMenu()
	})

	v.tui.signForm.AddButton("Cancel", func() {
		v.switchToMainMenu()
	})
}

func (v *View) createMainMenu() {
	v.tui.mainMenu = tview.NewList().
		AddItem("Register", "Sign up new user", 'r', func() {
			v.tui.signForm.Clear(true)
			v.callSignForm(register)
			v.setHeader("Register")
			v.tui.body.SwitchToPage(signForm)
		}).
		AddItem("Login", "Sign in with exist user", 'l', func() {
			v.tui.signForm.Clear(true)
			v.callSignForm(login)
			v.setHeader("Login")
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
			v.switchToMainMenu()
		}).
		SetBackgroundColor(tcell.ColorLightCoral)

	v.tui.body.AddPage(registerFail, v.tui.registerFail, true, false)
}

func (v *View) createUnitsMenu() {
	v.tui.unitsMenu = tview.NewList().
		AddItem("Pairs", "show login/password pairs", 'r', func() {
			v.getPairsList()
			v.setHeader("Pairs (press ESC to exit)")
			v.tui.body.SwitchToPage(pairsPage)
		}).
		AddItem("Notes", "show arbitrary text data", 'l', func() {
			v.getNotesList()
			v.setHeader("Notes (press ESC to exit)")
			v.tui.body.SwitchToPage(notesPage)
		}).
		AddItem("Cards", "show bank cards data", 'c', func() {
			v.getCardsList()
			v.setHeader("Cards (press ESC to exit)")
			v.tui.body.SwitchToPage(cardsPage)
		}).
		AddItem("Binary", "show arbitrary binary data", 'b', nil).
		AddItem("Back", "... to main menu", ' ', func() {
			v.switchToMainMenu()
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

	v.tui.pairsPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			v.switchToUnitsMenu()
		}
		return event
	})

	v.tui.body.AddPage(pairsPage, v.tui.pairsPage, true, false)
}

func (v *View) getPairsList() {
	pairs, err := v.ctrl.Pairs.ViewAllPairs(context.Background(), v.ctrl.Token)
	if err != nil {
		v.switchToUnitsMenu()
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

func (v *View) createCardsPage() {
	v.tui.cardsList = tview.NewList().ShowSecondaryText(false)
	v.tui.cardInfo = tview.NewTextView()

	v.tui.cardsPage = tview.NewFlex().
		AddItem(v.tui.cardsList, 0, 1, true).
		AddItem(v.tui.cardInfo, 0, 3, false)

	v.tui.cardsPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			v.switchToUnitsMenu()
		}
		return event
	})

	v.tui.body.AddPage(cardsPage, v.tui.cardsPage, true, false)
}

func (v *View) getCardsList() {
	cards, err := v.ctrl.Cards.ViewAllCards(context.Background(), v.ctrl.Token)
	if err != nil {
		v.switchToUnitsMenu()
		return
	}

	v.tui.cardsList.Clear()
	for _, card := range cards {
		v.tui.cardsList.AddItem(strconv.Itoa(card.ID), "", ' ', nil)
	}

	v.tui.cardsList.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		v.setCardInfo(cards[index])
	})
}

func (v *View) setCardInfo(card entity.BankDTO) {
	var sb strings.Builder

	v.tui.cardInfo.Clear()
	sb.WriteString(card.CardHolder)
	sb.WriteString("\n")
	sb.WriteString(card.Number)
	sb.WriteString("\n")
	sb.WriteString(card.ExpirationDate)
	sb.WriteString("\n")
	sb.WriteString(card.Metadata)
	sb.WriteString("\n")

	v.tui.cardInfo.SetText(sb.String())
}

func (v *View) createNotesPage() {
	v.tui.notesList = tview.NewList().ShowSecondaryText(false)
	v.tui.noteInfo = tview.NewTextView()

	v.tui.notesPage = tview.NewFlex().
		AddItem(v.tui.notesList, 0, 1, true).
		AddItem(v.tui.noteInfo, 0, 3, false)

	v.tui.notesPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			v.switchToUnitsMenu()
		}
		return event
	})

	v.tui.body.AddPage(notesPage, v.tui.notesPage, true, false)
}

func (v *View) getNotesList() {
	notes, err := v.ctrl.Notes.ViewAllNotes(context.Background(), v.ctrl.Token)
	if err != nil {
		v.switchToUnitsMenu()
		return
	}

	v.tui.notesList.Clear()
	for _, note := range notes {
		v.tui.notesList.AddItem(strconv.Itoa(note.ID), "", ' ', nil)
	}

	v.tui.notesList.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		v.setNoteInfo(notes[index])
	})
}

func (v *View) setNoteInfo(note entity.TextDTO) {
	var sb strings.Builder

	v.tui.cardInfo.Clear()
	sb.WriteString(note.Note)
	sb.WriteString("\n")
	sb.WriteString(note.Metadata)
	sb.WriteString("\n")

	v.tui.noteInfo.SetText(sb.String())
}

func (v *View) createHeader() {
	v.tui.header = tview.NewTextView()
	v.tui.header.SetBorder(true)
	v.tui.header.SetText("Main menu")
}

func (v *View) setHeader(text string) {
	v.tui.header.SetText(text)
}

func (v *View) createFooter() {
	clientInfo := tview.NewTextView().
		SetText("version: " + v.cfg.App.Version)

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
		AddItem(v.tui.header, 0, 1, false).
		AddItem(v.tui.body, 0, 3, true).
		AddItem(v.tui.footer, 0, 1, false)
}

func (v *View) switchToMainMenu() {
	v.setHeader("Main menu")
	v.tui.body.SwitchToPage(mainMenu)
}

func (v *View) switchToUnitsMenu() {
	v.setHeader("Resources")
	v.tui.body.SwitchToPage(unitsMenu)
}
