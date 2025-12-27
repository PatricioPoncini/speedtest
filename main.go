package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/showwin/speedtest-go/speedtest"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0"))
	dlStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	ulStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C8AFF")).Bold(true)
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

type model struct {
	server  *speedtest.Server
	dlSpeed float64
	ulSpeed float64
	loading bool
	stage   int
	err     error
	spinner string
}

type serverMsg *speedtest.Server
type dlMsg float64
type ulMsg float64
type errMsg error
type tickMsg time.Time

func initialModel() model {
	return model{
		loading: true,
		stage:   0,
		spinner: "|",
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh, se rompió todo: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(findServerCmd(), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tickMsg:
		if m.stage < 3 {
			switch m.spinner {
			case "|":
				m.spinner = "/"
			case "/":
				m.spinner = "-"
			case "-":
				m.spinner = "\\"
			case "\\":
				m.spinner = "|"
			}
			return m, tickCmd()
		}

	case serverMsg:
		m.server = msg
		m.stage = 1
		return m, testDownloadCmd(m.server)

	case dlMsg:
		m.dlSpeed = float64(msg)
		m.stage = 2
		return m, testUploadCmd(m.server)

	case ulMsg:
		m.ulSpeed = float64(msg)
		m.stage = 3
		m.loading = false
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return errStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	s := "\n" + titleStyle.Render("🚀 SPEEDTEST") + "\n\n"

	if m.stage >= 1 {
		s += fmt.Sprintf("📡 Server: %s (%s)\n", m.server.Name, m.server.Country)
	} else {
		s += infoStyle.Render(fmt.Sprintf("%s Buscando el mejor servidor...", m.spinner)) + "\n"
	}

	if m.stage > 1 {
		s += fmt.Sprintf("⬇️  Bajada: %s\n", dlStyle.Render(fmt.Sprintf("%.2f Mbps", m.dlSpeed)))
	} else if m.stage == 1 {
		s += infoStyle.Render(fmt.Sprintf("%s Midiendo bajada...", m.spinner)) + "\n"
	}

	if m.stage > 2 {
		s += fmt.Sprintf("⬆️  Subida: %s\n", ulStyle.Render(fmt.Sprintf("%.2f Mbps", m.ulSpeed)))
	} else if m.stage == 2 {
		s += infoStyle.Render(fmt.Sprintf("%s Midiendo subida...", m.spinner)) + "\n"
	}

	s += "\n" + infoStyle.Render("(Presioná 'q' para salir)") + "\n"
	return s
}

func findServerCmd() tea.Cmd {
	return func() tea.Msg {
		client := speedtest.New()

		serverList, err := client.FetchServers()
		if err != nil {
			return errMsg(fmt.Errorf("error buscando servers: %v", err))
		}

		targets, err := serverList.FindServer([]int{})
		if err != nil {
			return errMsg(fmt.Errorf("error seleccionando server: %v", err))
		}

		if len(targets) > 0 {
			return serverMsg(targets[0])
		}

		return errMsg(fmt.Errorf("no se encontraron servidores"))
	}
}

func testDownloadCmd(s *speedtest.Server) tea.Cmd {
	return func() tea.Msg {
		s.DownloadTest()
		return dlMsg(s.DLSpeed)
	}
}

func testUploadCmd(s *speedtest.Server) tea.Cmd {
	return func() tea.Msg {
		s.UploadTest()
		return ulMsg(s.ULSpeed)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
