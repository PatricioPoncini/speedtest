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

## `🛠️ Installation & Usage`

Follow these steps to install the tool and run it from anywhere in your terminal.

### `Prerequisites`
Make sure you have **Go** installed on your system. You can download it [here](https://go.dev/dl/).

### Installation Steps (Linux/macOS)

1. **Clone the repository and navigate to the folder:**
   ```bash
   git clone [https://github.com/PatricioPoncini/speedtest.git](https://github.com/PatricioPoncini/speedtest.git)
   cd speedtest
   ```
   
2. Download dependencies and build the binary: This will create an executable file named speedtest.
   ```bash
   go mod tidy
   go build -o speedtest main.go
   ```
   
3. Move the binary to your system's PATH: This allows you to run the command globally.
   ```bash
   sudo mv speedtest /usr/local/bin/
   ```

Done! You can now open a new terminal window and simply run:
```bash
speedtest
```