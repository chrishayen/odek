package tui

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

const (
	transitionDuration = 320 * time.Millisecond
	transitionFPS      = 60
	transitionPinLabel = " feature "
	pinFadeInStart     = 0.55
)

type transitionFrameMsg time.Time

func tickTransitionFrame() tea.Cmd {
	return tea.Tick(time.Second/transitionFPS, func(t time.Time) tea.Msg {
		return transitionFrameMsg(t)
	})
}

var featurePinStyle = lipgloss.NewStyle().
	Foreground(fgBright).
	Background(accent).
	Padding(0, 1).
	Bold(true)

func renderFeaturePin() string {
	return featurePinStyle.Render(transitionPinLabel)
}

type transitionModel struct {
	from      tea.Model
	to        tea.Model
	width     int
	height    int
	startedAt time.Time
	duration  time.Duration
	pin       string
	fromSnap  string
	toSnap    string
}

func newTransition(from, to tea.Model, width, height int, pin string) transitionModel {
	m := transitionModel{
		from:      from,
		to:        to,
		width:     width,
		height:    height,
		startedAt: time.Now(),
		duration:  transitionDuration,
		pin:       pin,
	}
	m.snapshot()
	return m
}

func (m *transitionModel) snapshot() {
	m.fromSnap = m.from.View().Content
	m.toSnap = m.to.View().Content
}

func (m transitionModel) Init() tea.Cmd {
	return tickTransitionFrame()
}

func (m transitionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.from, _ = m.from.Update(msg)
		m.to, _ = m.to.Update(msg)
		m.snapshot()
		return m, nil
	case transitionFrameMsg:
		if time.Since(m.startedAt) >= m.duration {
			return m.to, m.to.Init()
		}
		return m, tickTransitionFrame()
	case tea.KeyPressMsg:
		return m.to, m.to.Init()
	}
	return m, nil
}

// pinAnchor is satisfied by destination screens that want the transition's
// end-state icon to land at a specific row (not the bottom edge), so handoff
// is seamless.
type pinAnchor interface {
	PinRow() int
}

func (m transitionModel) View() tea.View {
	progress := float64(time.Since(m.startedAt)) / float64(m.duration)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	eased := easeOutCubic(progress)

	pinY := m.height - 1
	if pa, ok := m.to.(pinAnchor); ok {
		if r := pa.PinRow(); r >= 0 && r < m.height {
			pinY = r
		}
	}

	content := composeSlideCollapse(m.fromSnap, m.toSnap, m.pin, m.width, m.height, pinY, eased)
	v := tea.NewView(content)
	v.AltScreen = true
	v.BackgroundColor = bgMain
	return v
}

func easeOutCubic(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	u := 1 - t
	return 1 - u*u*u
}

func lerpInt(a, b int, t float64) int {
	return a + int(float64(b-a)*t+0.5)
}

// composeSlideCollapse returns a single rendered frame of the transition: the
// outgoing view shrinks into a small rectangle anchored to (x=0, y=pinY)
// while the incoming view slides in from the right. At t=1 the pin is drawn
// in the iconW×1 rectangle at row pinY, matching where the destination screen
// will draw its own pin — so the handoff is seamless.
func composeSlideCollapse(fromSnap, toSnap, pin string, width, height, pinY int, t float64) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	if pinY < 0 {
		pinY = 0
	}
	if pinY >= height {
		pinY = height - 1
	}

	iconW := min(max(ansi.StringWidth(pin), 1), width)

	outW := max(lerpInt(width, iconW, t), iconW)
	outH := max(lerpInt(height, 1, t), 1)
	bottomEdge := lerpInt(height-1, pinY, t)
	outY := max(bottomEdge-outH+1, 0)
	if outY+outH > height {
		outH = height - outY
	}
	slideLeft := lerpInt(width, 0, t)

	fromLines := padLines(splitLines(fromSnap), width)
	toLines := padLines(splitLines(toSnap), width)

	bgLine := lipgloss.NewStyle().Background(bgMain).Render(strings.Repeat(" ", width))

	rows := make([]string, height)
	for row := range height {
		line := bgLine

		visInc := width - slideLeft
		if visInc > 0 && row < len(toLines) {
			incCropped := ansi.Cut(toLines[row], 0, visInc)
			line = overlayAt(line, incCropped, slideLeft, width)
		}

		if row >= outY && row < outY+outH {
			origRow := row - outY
			if origRow >= 0 && origRow < len(fromLines) {
				outCropped := ansi.Cut(fromLines[origRow], 0, outW)
				line = overlayAt(line, outCropped, 0, width)
			}
		}

		rows[row] = line
	}

	if t >= pinFadeInStart && pinY >= 0 && pinY < len(rows) {
		rows[pinY] = overlayAt(rows[pinY], pin, 0, width)
	}

	return strings.Join(rows, "\n")
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func padLines(lines []string, width int) []string {
	bgPadStyle := lipgloss.NewStyle().Background(bgMain)
	out := make([]string, len(lines))
	for i, line := range lines {
		w := ansi.StringWidth(line)
		switch {
		case w == width:
			out[i] = line
		case w < width:
			out[i] = line + bgPadStyle.Render(strings.Repeat(" ", width-w))
		default:
			out[i] = ansi.Truncate(line, width, "")
		}
	}
	return out
}

// overlayAt paints `overlay` onto `base` starting at column `x`, returning a
// line of `canvasWidth` columns. ANSI-aware: it uses ansi.Cut to preserve
// color state on either side of the overlay.
func overlayAt(base, overlay string, x, canvasWidth int) string {
	if x >= canvasWidth || overlay == "" {
		return base
	}
	overlayW := ansi.StringWidth(overlay)
	if overlayW == 0 {
		return base
	}
	if x+overlayW > canvasWidth {
		overlayW = canvasWidth - x
		overlay = ansi.Truncate(overlay, overlayW, "")
	}
	left := ansi.Cut(base, 0, x)
	right := ansi.Cut(base, x+overlayW, canvasWidth)
	return left + overlay + right
}
