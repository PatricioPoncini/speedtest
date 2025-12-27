# `Go Speedtest CLI 🚀`

A minimal, beautiful, and fast command-line interface for testing internet speed written in Go.

This tool uses [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI (Terminal User Interface) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling, providing a modern and clean look compared to traditional speedtest tools.

## `✨ Features`

- **TUI Interface:** Clean, bordered dashboard layout.
- **Real-time Feedback:** Visual spinner and stage updates.
- **Accurate Metrics:**
    - **Ping/Latency**
    - **Download Speed**
    - **Upload Speed**
- **Smart Server Selection:** Automatically finds the closest and best server for testing.
- **Lightweight:** Compiled into a single binary with no external runtime dependencies.

## `📦 Tech Stack`
- **Language**: [Go (Golang)](https://go.dev/)
- **TUI Framework**: [Bubble Tea](https://go.dev/)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Speedtest Engine**: [speedtest-go](https://github.com/showwin/speedtest-go)