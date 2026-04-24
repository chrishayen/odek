package tui

import (
	"context"
	"image/color"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"

	"shotgun.dev/odek/internal/decomposer"
	openai "shotgun.dev/odek/openai"
)

var taglineStyle = lipgloss.NewStyle().
	Foreground(dim).
	Italic(true)

var frameStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(border).
	Padding(1, 4)

// 70s VHS tape stripe palette
var gradStops []colorful.Color

func init() {
	hexes := []string{
		"#F9D050", // yellow
		"#F6B830", // amber
		"#F29A1E", // orange
		"#EC7A1E", // deep orange
		"#E05E22", // red-orange
		"#D4443C", // red
		"#B53A2E", // deep red
		"#863520", // rust
		"#5A2A10", // brown
	}
	for _, h := range hexes {
		c, _ := colorful.Hex(h)
		gradStops = append(gradStops, c)
	}
}

type model struct {
	width       int
	height      int
	help        help.Model
	ctx         context.Context
	client      *openai.Client
	decomposer  *decomposer.Decomposer
	kanjiOffset int
	logoX       int
	logoY       int
	logoVX      int
	logoVY      int
}

type landingTickMsg struct{}

func landingTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg { return landingTickMsg{} })
}

func (m model) Init() tea.Cmd {
	return landingTick()
}

func (m *model) bounceLogo() {
	_, lh, lw := logoMask()
	maxKanji := m.width / 2
	maxX := maxKanji - lw
	maxY := m.height - lh
	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}
	if m.logoVX == 0 && m.logoVY == 0 {
		m.logoVX, m.logoVY = 1, 1
	}
	m.logoX += m.logoVX
	m.logoY += m.logoVY
	if m.logoX <= 0 {
		m.logoX = 0
		m.logoVX = 1
	}
	if m.logoX >= maxX {
		m.logoX = maxX
		m.logoVX = -1
	}
	if m.logoY <= 0 {
		m.logoY = 0
		m.logoVY = 1
	}
	if m.logoY >= maxY {
		m.logoY = maxY
		m.logoVY = -1
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		return m, nil
	case landingTickMsg:
		m.kanjiOffset++
		m.bounceLogo()
		return m, landingTick()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "n", "N":
			cf := newCreateFeatureModel(m.ctx, m.client, m.decomposer, m.width, m.height)
			if m.width >= splitPaneMinWidth {
				_, rightW := splitWidths(m.width)
				right := newFeatureDecompModel(rightW, m.height, nil, cf.state)
				split := newSplitPaneModel(cf, right, m.width, m.height)
				return split, split.Init()
			}
			return cf, cf.Init()
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func renderStripes(text string, stops []colorful.Color) string {
	lines := strings.Split(text, "\n")
	n := len(stops) - 1

	var out strings.Builder
	for i, line := range lines {
		t := float64(i) / float64(len(lines)-1) * float64(n)
		idx := int(t)
		if idx >= n {
			idx = n - 1
		}
		frac := t - float64(idx)
		c := stops[idx].BlendLuv(stops[idx+1], frac)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex())).Bold(true)
		out.WriteString(style.Render(line))
		if i < len(lines)-1 {
			out.WriteRune('\n')
		}
	}
	return out.String()
}

func renderGradientOnBg(text string, stops []colorful.Color, bg string, totalWidth int) string {
	lines := strings.Split(text, "\n")
	bgColor := lipgloss.Color(bg)

	maxLen := 0
	for _, line := range lines {
		runes := []rune(line)
		if len(runes) > maxLen {
			maxLen = len(runes)
		}
	}
	if maxLen == 0 {
		return text
	}

	n := len(stops) - 1
	bgStyle := lipgloss.NewStyle().Background(bgColor)
	var out strings.Builder
	for i, line := range lines {
		runes := []rune(line)
		for j, r := range runes {
			if r == ' ' {
				out.WriteString(bgStyle.Render(" "))
				continue
			}
			t := float64(j) / float64(maxLen) * float64(n)
			idx := int(t)
			if idx >= n {
				idx = n - 1
			}
			frac := t - float64(idx)
			c := stops[idx].BlendLuv(stops[idx+1], frac)
			out.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex())).Background(bgColor).Bold(true).Render(string(r)))
		}
		pad := totalWidth - len(runes)
		if pad > 0 {
			out.WriteString(bgStyle.Render(strings.Repeat(" ", pad)))
		}
		if i < len(lines)-1 {
			out.WriteRune('\n')
		}
	}
	return out.String()
}

// odekMask is hand-crafted at kanji resolution: each character represents one
// kanji (2 screen cells wide × 1 screen row tall). '#' lights up a kanji with
// the warm logo gradient; space leaves it muted. Letter gaps are ≥2 kanji wide
// so they remain visible against the scrolling field.
var odekMask = []string{
	`  ####    #######   #######  ##    ##`,
	` ######   ########  #######  ##   ## `,
	`##    ##  ##    ##  ##       ##  ##  `,
	`##    ##  ##    ##  ##       ## ##   `,
	`##    ##  ##    ##  ######   ####    `,
	`##    ##  ##    ##  ######   ####    `,
	`##    ##  ##    ##  ##       ## ##   `,
	`##    ##  ##    ##  ##       ##  ##  `,
	` ######   ########  #######  ##   ## `,
	`  ####    #######   #######  ##    ##`,
}

// logoMask returns a 2D boolean grid from odekMask. h is row count (screen
// rows); w is the width in KANJI positions (not screen cells).
func logoMask() (mask [][]bool, h, w int) {
	h = len(odekMask)
	for _, row := range odekMask {
		if rc := len([]rune(row)); rc > w {
			w = rc
		}
	}
	mask = make([][]bool, h)
	for y, row := range odekMask {
		mask[y] = make([]bool, w)
		for x, r := range []rune(row) {
			if r != ' ' && r != 0 {
				mask[y][x] = true
			}
		}
	}
	return
}

// logoGradientRows computes one warm-gradient color per logo row by walking
// gradStops the same way renderStripes does. It also returns a lighter
// companion palette (each stop blended 50% toward white) used for help-line
// chars that fall under the logo mask — so the help bar pops against the
// surrounding logo stripes.
func logoGradientRows(h int) (normal, light []color.Color) {
	normal = make([]color.Color, h)
	light = make([]color.Color, h)
	n := len(gradStops) - 1
	white := colorful.Color{R: 1, G: 1, B: 1}
	for i := range normal {
		t := float64(i) / float64(h-1) * float64(n)
		idx := int(t)
		if idx >= n {
			idx = n - 1
		}
		frac := t - float64(idx)
		c := gradStops[idx].BlendLuv(gradStops[idx+1], frac)
		normal[i] = lipgloss.Color(c.Hex())
		light[i] = lipgloss.Color(c.BlendLuv(white, 0.55).Hex())
	}
	return
}

// buildHelpLine returns parallel slices describing the embedded help line:
// the rune, whether the rune belongs to a key glyph (bold), and whether the
// rune is a separator bullet. Length is even (padded with a trailing space if
// needed) so the help line spans whole kanji slots.
func buildHelpLine() (runes []rune, isKey, isSep []bool) {
	items := []struct{ k, d string }{{"n", "new"}, {"q", "quit"}}
	for i, it := range items {
		if i > 0 {
			for _, r := range "  •  " {
				runes = append(runes, r)
				isKey = append(isKey, false)
				isSep = append(isSep, r == '•')
			}
		}
		for _, r := range it.k {
			runes = append(runes, r)
			isKey = append(isKey, true)
			isSep = append(isSep, false)
		}
		runes = append(runes, ' ')
		isKey = append(isKey, false)
		isSep = append(isSep, false)
		for _, r := range it.d {
			runes = append(runes, r)
			isKey = append(isKey, false)
			isSep = append(isSep, false)
		}
	}
	if len(runes)%2 == 1 {
		runes = append(runes, ' ')
		isKey = append(isKey, false)
		isSep = append(isSep, false)
	}
	return
}

// renderKanjiField builds a w×h screen of scrolling kanji. Each row drifts
// horizontally at a uniform speed (alternating direction by parity). Kanji
// inside the bouncing ODEK logo mask get the warm VHS gradient; the rest stay
// muted gray. The help line is embedded near the bottom center, and its chars
// also pick up the warm gradient when the logo mask passes over them.
func renderKanjiField(w, h, offset, logoKX0, logoY0 int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	mask, lh, lw := logoMask()
	rowColors, lightColors := logoGradientRows(lh)

	baseStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
	bgStyle := lipgloss.NewStyle().Background(bgMain)
	litStyles := make([]lipgloss.Style, lh)
	litHelpStyles := make([]lipgloss.Style, lh)
	for i := range rowColors {
		litStyles[i] = lipgloss.NewStyle().Foreground(rowColors[i]).Background(bgMain).Bold(true)
		litHelpStyles[i] = lipgloss.NewStyle().Foreground(lightColors[i]).Background(bgMain).Bold(true)
	}

	hlKeyStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain).Bold(true)
	hlDescStyle := lipgloss.NewStyle().Foreground(fgBody).Background(bgMain)
	hlSepStyle := lipgloss.NewStyle().Foreground(fgDim).Background(bgMain)

	hlRunes, hlKey, hlSep := buildHelpLine()
	hlW := len(hlRunes)
	const hlPadX = 2
	hlBlockW := hlW + hlPadX*2
	hlBlockStart := (w - hlBlockW) / 2
	hlBlockStart &^= 1 // snap to even so kanji alignment on both sides survives
	if hlBlockStart < 0 {
		hlBlockStart = 0
	}
	hlStart := hlBlockStart + hlPadX
	hlBlockEnd := hlBlockStart + hlBlockW
	hlRow := h - 3
	if hlRow < 0 {
		hlRow = 0
	}

	var out strings.Builder
	for y := 0; y < h; y++ {
		rowDir := 1
		if y%2 == 1 {
			rowDir = -1
		}
		// 80ms tick → 12.5 ticks/sec. speed=1 → 1 char/sec = 1 char / 12 ticks.
		const speed = 6
		rowOff := rowDir * speed * offset / 12

		x := 0
		for x+2 <= w {
			kanjiIdx := x / 2
			lit := false
			litRow := 0
			if ly := y - logoY0; ly >= 0 && ly < lh {
				if kx := kanjiIdx - logoKX0; kx >= 0 && kx < lw && mask[ly][kx] {
					lit = true
					litRow = ly
				}
			}

			if y == hlRow && x >= hlBlockStart && x+2 <= hlBlockEnd {
				for i := 0; i < 2; i++ {
					cx := x + i
					var r rune = ' '
					var style lipgloss.Style
					inText := cx >= hlStart && cx < hlStart+hlW
					if inText {
						hi := cx - hlStart
						r = hlRunes[hi]
					}
					switch {
					case lit && inText:
						style = litHelpStyles[litRow]
					case !inText:
						style = bgStyle
					case hlKey[cx-hlStart]:
						style = hlKeyStyle
					case hlSep[cx-hlStart]:
						style = hlSepStyle
					default:
						style = hlDescStyle
					}
					out.WriteString(style.Render(string(r)))
				}
				x += 2
				continue
			}

			r := kanjiAt(y, kanjiIdx+rowOff)
			if lit {
				out.WriteString(litStyles[litRow].Render(string(r)))
			} else {
				out.WriteString(baseStyle.Render(string(r)))
			}
			x += 2
		}
		if x < w {
			out.WriteString(bgStyle.Render(" "))
		}
		if y < h-1 {
			out.WriteRune('\n')
		}
	}
	return out.String()
}

func (m model) View() tea.View {
	field := renderKanjiField(m.width, m.height, m.kanjiOffset, m.logoX, m.logoY)
	placed := lipgloss.Place(m.width, m.height,
		lipgloss.Left, lipgloss.Top,
		field,
	)
	v := tea.NewView(placed)
	v.AltScreen = true
	v.BackgroundColor = bgMain
	return v
}

func Run(ctx context.Context, client *openai.Client, dec *decomposer.Decomposer) error {
	teaModel := model{help: newHelpModel(), ctx: ctx, client: client, decomposer: dec}
	p := tea.NewProgram(teaModel)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
