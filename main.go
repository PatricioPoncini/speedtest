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
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#5A42BC")).
		Padding(0, 1).
		Bold(true)

	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#767676")).
		Width(12).
		Render

	dimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		Render

	serverInfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D0D0D0")).Bold(true)
	dlStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	ulStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C8AFF")).Bold(true)
	errStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444"))
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
		spinner: "⠋",
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
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
			case "⠋":
				m.spinner = "⠙"
			case "⠙":
				m.spinner = "⠹"
			case "⠹":
				m.spinner = "⠸"
			case "⠸":
				m.spinner = "⠼"
			case "⠼":
				m.spinner = "⠴"
			case "⠴":
				m.spinner = "⠦"
			case "⠦":
				m.spinner = "⠧"
			case "⠧":
				m.spinner = "⠇"
			case "⠇":
				m.spinner = "⠏"
			case "⠏":
				m.spinner = "⠋"
			default:
				m.spinner = "⠋"
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
		return appStyle.Render(errStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	s := titleStyle.Render("SPEEDTEST") + "\n\n"

	s += labelStyle("Server")
	if m.stage >= 1 {
		s += serverInfoStyle.Render(fmt.Sprintf("%s (%s)", m.server.Name, m.server.Country))
	} else {
		s += dimStyle(fmt.Sprintf("%s Finding best server...", m.spinner))
	}
	s += "\n"

	s += labelStyle("Download")
	if m.stage > 1 {
		mbps := (m.dlSpeed * 8) / 1000000.0
		s += dlStyle.Render(fmt.Sprintf("%.2f Mbps", mbps))
	} else if m.stage == 1 {
		s += dimStyle(fmt.Sprintf("%s Testing download...", m.spinner))
	} else {
		s += dimStyle("...")
	}
	s += "\n"

	s += labelStyle("Upload")
	if m.stage > 2 {
		mbps := (m.ulSpeed * 8) / 1000000.0
		s += ulStyle.Render(fmt.Sprintf("%.2f Mbps", mbps))
	} else if m.stage == 2 {
		s += dimStyle(fmt.Sprintf("%s Testing upload...", m.spinner))
	} else {
		s += dimStyle("...")
	}
	s += "\n\n"

	if m.loading {
		s += dimStyle("Press 'q' to quit")
	} else {
		s += dimStyle("Done. Press 'q' to exit.")
	}

	return appStyle.Render(s)
}

func findServerCmd() tea.Cmd {
	return func() tea.Msg {
		client := speedtest.New()

		serverList, err := client.FetchServers()
		if err != nil {
			return errMsg(fmt.Errorf("error fetching servers: %v", err))
		}

		targets, err := serverList.FindServer([]int{})
		if err != nil {
			return errMsg(fmt.Errorf("error selecting server: %v", err))
		}

		if len(targets) > 0 {
			return serverMsg(targets[0])
		}

		return errMsg(fmt.Errorf("no servers found"))
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
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
