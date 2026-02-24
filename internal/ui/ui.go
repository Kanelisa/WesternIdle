package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	_ "github.com/jezek/xgb/res"

	"WesternIdle/internal/system"
	"sort"
)

//
// ---------- UI ELEMENT INTERFACE ----------
//

type UIElement interface {
	Draw(*ebiten.Image)
	HandleClick(mx, my int)
}

//
// ---------- BUTTON ----------
//

type UIButton struct {
	Text     string
	X, Y     float64
	Width    float64
	Height   float64
	Color    color.RGBA
	Font     text.Face
	ActionID string
	OnClick  func()

	Progress  float64
	Duration  float64
	Remaining float64
}

func (b *UIButton) IsHovered(mx, my int) bool {
	return float64(mx) >= b.X &&
		float64(mx) <= b.X+b.Width &&
		float64(my) >= b.Y &&
		float64(my) <= b.Y+b.Height
}

func (b *UIButton) Draw(screen *ebiten.Image) {

	mx, my := ebiten.CursorPosition()
	hovered := b.IsHovered(mx, my)

	// ---- 1. Базовый цвет ----
	base := b.Color

	// Лёгкое затемнение если действие выполняется
	if b.Progress > 0 {
		base = color.RGBA{
			R: b.Color.R - 20,
			G: b.Color.G - 20,
			B: b.Color.B - 20,
			A: 255,
		}
	}

	// Hover подсветка
	if hovered {
		base = color.RGBA{
			R: min(base.R+20, 255),
			G: min(base.G+20, 255),
			B: min(base.B+20, 255),
			A: 255,
		}
	}

	// ---- 2. Рисуем фон ----
	ebitenutil.DrawRect(screen, b.X, b.Y, b.Width, b.Height, base)

	// ---- 3. Прогресс (поверх фона, но под текстом) ----
	if b.Progress > 0 {

		progressWidth := b.Width * b.Progress

		ebitenutil.DrawRect(
			screen,
			b.X,
			b.Y,
			progressWidth,
			b.Height,
			color.RGBA{220, 140, 50, 180},
		)
	}

	// ---- 4. Рамка ----
	ebitenutil.DrawRect(screen, b.X, b.Y, b.Width, 1, color.RGBA{0, 0, 0, 255})
	ebitenutil.DrawRect(screen, b.X, b.Y+b.Height-1, b.Width, 1, color.RGBA{0, 0, 0, 255})
	ebitenutil.DrawRect(screen, b.X, b.Y, 1, b.Height, color.RGBA{0, 0, 0, 255})
	ebitenutil.DrawRect(screen, b.X+b.Width-1, b.Y, 1, b.Height, color.RGBA{0, 0, 0, 255})

	// ---- 5. Текст ----
	w, h := text.Measure(b.Text, b.Font, 1.0)
	tx := b.X + (b.Width-w)/2
	ty := b.Y + (b.Height-h)/2

	DrawText(
		screen,
		b.Text,
		b.Font,
		tx,
		ty,
		color.RGBA{255, 255, 255, 255},
	)

	// ---- 6. Таймер ----
	if b.Progress > 0 && b.Remaining > 0 {

		timeStr := fmt.Sprintf("%.1fс", b.Remaining)
		tw, th := text.Measure(timeStr, b.Font, 1.0)

		DrawText(
			screen,
			timeStr,
			b.Font,
			b.X+b.Width-tw-6,
			b.Y+(b.Height-th)/2,
			color.RGBA{255, 255, 255, 255},
		)
	}
}

func (b *UIButton) HandleClick(mx, my int) {
	if float64(mx) >= b.X && float64(mx) <= b.X+b.Width &&
		float64(my) >= b.Y && float64(my) <= b.Y+b.Height {

		if b.OnClick != nil {
			b.OnClick()
		}
	}
}

//
// ---------- PANEL ----------
//

type Panel struct {
	X, Y          float64
	Width, Height float64
	Color         color.RGBA
	Elements      []UIElement
}

func (p *Panel) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.X, p.Y, p.Width, p.Height, p.Color)
	for _, e := range p.Elements {
		e.Draw(screen)
	}
}

func (p *Panel) Layout() {
	currentY := p.Y + 10 // отступ сверху
	spacing := 8.0

	for _, el := range p.Elements {
		if cat, ok := el.(*UICategory); ok {
			// теперь currentY включает и хедер, и кнопки (если развернута)
			currentY = cat.Layout(currentY)
			currentY += spacing
		}
	}
}

func (p *Panel) HandleClick(mx, my int) {
	for _, e := range p.Elements {
		if cat, ok := e.(*UICategory); ok {
			cat.HandleClick(mx, my)
		} else {
			e.HandleClick(mx, my)
		}
	}
}

//
// ---------- UI ----------
//

func (ui *UI) buildTooltip(btn *UIButton) []string {

	info, ok := system.GetActionInfo(btn.ActionID)
	if !ok {
		return nil
	}

	var lines []string

	if !info.Instant {
		lines = append(lines,
			fmt.Sprintf("Длительность: %.1f сек", info.Duration),
		)
	}

	if len(info.Cost) > 0 {
		for res, val := range info.Cost {
			name := system.ResourceDisplayName(res)
			lines = append(lines,
				fmt.Sprintf("Стоимость: -%.1f %s", val, name),
			)
		}
	}

	if len(info.Reward) > 0 {
		for res, val := range info.Reward {
			name := system.ResourceDisplayName(res)
			lines = append(lines,
				fmt.Sprintf("Награда: +%.1f %s", val, name),
			)
		}
	}

	return lines
}

func (ui *UI) drawTooltip(screen *ebiten.Image) {

	if ui.HoveredButton == nil {
		return
	}

	info, ok := system.GetActionInfo(ui.HoveredButton.ActionID)
	if !ok {
		return
	}

	var lines []string

	// --- Название ---
	lines = append(lines, ui.HoveredButton.Text)

	// --- Разделитель ---
	lines = append(lines, "----------------")

	// --- Длительность ---
	if !info.Instant {
		lines = append(lines,
			fmt.Sprintf("Длительность: %.1f сек", info.Duration),
		)
	}

	// --- Стоимость ---
	for res, val := range info.Cost {
		name := system.ResourceDisplayName(res)
		lines = append(lines,
			fmt.Sprintf("Стоимость: -%.1f %s", val, name),
		)
	}

	// --- Награда ---
	for res, val := range info.Reward {
		name := system.ResourceDisplayName(res)
		lines = append(lines,
			fmt.Sprintf("Награда: +%.1f %s", val, name),
		)
	}

	// --- Описание ---
	if info.Description != "" {
		lines = append(lines, " ")
		lines = append(lines, info.Description)
	}

	// --- Размеры ---
	padding := 10.0
	lineHeight := 18.0

	maxWidth := 0.0
	for _, line := range lines {
		w, _ := text.Measure(line, ui.GameFont, 1.0)
		if w > maxWidth {
			maxWidth = w
		}
	}

	width := maxWidth + padding*2
	height := padding*2 + lineHeight*float64(len(lines))

	// --- Позиция (справа от кнопки) ---
	//x := ui.HoveredButton.X + ui.HoveredButton.Width + 15
	//y := ui.HoveredButton.Y
	xint, yint := ebiten.CursorPosition()
	var x, y float64 = float64(xint + 15), float64(yint)
	// --- Фон ---
	ebitenutil.DrawRect(
		screen,
		x,
		y,
		width,
		height,
		color.RGBA{25, 20, 15, 240},
	)

	// --- Рамка ---
	ebitenutil.DrawRect(screen, x, y, width, 1, color.RGBA{80, 60, 40, 255})
	ebitenutil.DrawRect(screen, x, y+height-1, width, 1, color.RGBA{80, 60, 40, 255})
	ebitenutil.DrawRect(screen, x, y, 1, height, color.RGBA{80, 60, 40, 255})
	ebitenutil.DrawRect(screen, x+width-1, y, 1, height, color.RGBA{80, 60, 40, 255})

	// --- Текст ---
	for i, line := range lines {
		DrawText(
			screen,
			line,
			ui.GameFont,
			x+padding,
			y+padding+float64(i)*lineHeight,
			color.RGBA{255, 255, 255, 255},
		)
	}
}

type UI struct {
	State          *system.GameState
	GameFont       text.Face
	LeftPanel      *Panel
	LocationsPanel *Panel
	CenterPanel    *Panel
	RightPanel     *Panel
	OnAction       func(id string, duration float64)
	Notification   *Notification
	lastLogIndex   int
	HoveredButton  *UIButton
	ActiveTab      string
}

func NewUI(state *system.GameState, font text.Face, onAction func(id string, duration float64)) *UI {
	ui := &UI{
		State:     state,
		GameFont:  font,
		OnAction:  onAction,
		ActiveTab: "main",
	}

	ui.initPanels()
	ui.buildLocationPanel()
	ui.buildCenterPanel()
	return ui
}

func (ui *UI) initPanels() {

	// ---------------- LEFT PANEL (Только вкладки) ----------------
	ui.LeftPanel = &Panel{
		X:      0,
		Y:      0,
		Width:  160,
		Height: 600,
		Color:  color.RGBA{60, 40, 20, 255},
	}

	mainBtn := &UIButton{
		Text:   "Главная",
		X:      10,
		Y:      20,
		Width:  ui.LeftPanel.Width - 20,
		Height: 40,
		Color:  color.RGBA{100, 60, 40, 255},
		Font:   ui.GameFont,
		OnClick: func() {
			ui.ActiveTab = "main"
			ui.buildLocationPanel() // показываем кнопки локаций
			ui.buildCenterPanel()   // перестраиваем категории действий
		},
	}

	invBtn := &UIButton{
		Text:   "Инвентарь",
		X:      10,
		Y:      70,
		Width:  ui.LeftPanel.Width - 20,
		Height: 40,
		Color:  color.RGBA{100, 60, 40, 255},
		Font:   ui.GameFont,
		OnClick: func() {
			ui.ActiveTab = "inventory"
			ui.buildLocationPanel() // очищаем панель локаций
			ui.buildCenterPanel()
		},
	}

	ui.LeftPanel.Elements = []UIElement{mainBtn, invBtn}

	// ---------------- LOCATIONS PANEL ----------------
	ui.LocationsPanel = &Panel{
		X:      ui.LeftPanel.X + ui.LeftPanel.Width,
		Y:      0,
		Width:  160,
		Height: 600,
		Color:  color.RGBA{70, 45, 25, 255},
	}

	// ---------------- CENTER PANEL (Действия) ----------------
	ui.CenterPanel = &Panel{
		X:      ui.LocationsPanel.X + ui.LocationsPanel.Width,
		Y:      0,
		Width:  400,
		Height: 600,
		Color:  color.RGBA{40, 25, 15, 255},
	}

	// ---------------- RIGHT PANEL (Ресурсы) ----------------
	ui.RightPanel = &Panel{
		X:      ui.CenterPanel.X + ui.CenterPanel.Width,
		Y:      0,
		Width:  900 - ui.LeftPanel.Width - ui.CenterPanel.Width - ui.LocationsPanel.Width,
		Height: 600,
		Color:  color.RGBA{50, 30, 15, 255},
	}

	// ---------------- INITIAL BUILD ----------------
	ui.ActiveTab = "main"

	// buildLocationPanel вызываем после того, как State уже содержит Locations
	ui.buildLocationPanel() // создаём кнопки локаций
	ui.buildCenterPanel()   // создаём категории действий
}

//------------Панель локаций----------------------------------

func (ui *UI) buildLocationPanel() {
	ui.LocationsPanel.Elements = nil

	if ui.ActiveTab != "main" {
		return
	}

	panelOffsetX := ui.LocationsPanel.X + 10
	panelWidth := ui.LocationsPanel.Width - 20

	category := &UICategory{
		Title:    "Локации",
		X:        panelOffsetX,
		Width:    panelWidth,
		HeaderH:  28,
		Expanded: true,
		Font:     ui.GameFont,
		Color:    color.RGBA{70, 40, 25, 255},
	}

	// Сортировка ID по Order
	locIDs := make([]string, 0, len(ui.State.Locations))
	for id := range ui.State.Locations {
		locIDs = append(locIDs, id)
	}

	sort.Slice(locIDs, func(i, j int) bool {
		return ui.State.Locations[locIDs[i]].Order < ui.State.Locations[locIDs[j]].Order
	})

	for _, id := range locIDs {
		loc := ui.State.Locations[id]
		locID := id // для замыкания

		btn := &UIButton{
			Text:   loc.Name,
			Width:  panelWidth,
			Height: 28,
			Font:   ui.GameFont,
		}

		// Подсветка текущей локации
		if ui.State.CurrentLocation != nil && ui.State.CurrentLocation.ID == id {
			btn.Color = color.RGBA{150, 100, 60, 255}
		} else {
			btn.Color = color.RGBA{100, 70, 40, 255}
		}

		btn.OnClick = func() {
			if ui.State.ChangeLocation(locID) {
				ui.buildCenterPanel()   // перестраиваем действия
				ui.buildLocationPanel() // перестроить подсветку кнопок
			}
		}

		category.Elements = append(category.Elements, btn)
	}

	ui.LocationsPanel.Elements = append(ui.LocationsPanel.Elements, category)
	ui.LocationsPanel.Layout()
}

// Новая вспомогательная функция: обновляет цвет кнопок в зависимости от текущей локации
func (ui *UI) updateLocationHighlight() {
	if len(ui.LocationsPanel.Elements) == 0 {
		return
	}
	cat, ok := ui.LocationsPanel.Elements[0].(*UICategory)
	if !ok {
		return
	}

	for _, el := range cat.Elements {
		if btn, ok := el.(*UIButton); ok {
			if ui.State.CurrentLocation != nil && btn.Text == ui.State.CurrentLocation.Name {
				btn.Color = color.RGBA{150, 100, 60, 255} // активная
			} else {
				btn.Color = color.RGBA{100, 70, 40, 255} // неактивная
			}
		}
	}
}

//------------Центральная панель действий----------------------

func (ui *UI) layoutCenterTwoColumns() {
	panel := ui.CenterPanel

	padding := 10.0
	spacing := 10.0
	columnGap := 10.0

	availableWidth := panel.Width - padding*2 - columnGap
	columnWidth := availableWidth / 2

	leftX := panel.X + padding
	rightX := leftX + columnWidth + columnGap

	currentYLeft := panel.Y + padding
	currentYRight := panel.Y + padding

	for _, el := range panel.Elements {
		cat, ok := el.(*UICategory)
		if !ok {
			continue
		}

		cat.Width = columnWidth

		if cat.Title == "Дуэли" {
			cat.X = rightX
			currentYRight = cat.Layout(currentYRight)
			currentYRight += spacing
		} else {
			cat.X = leftX
			currentYLeft = cat.Layout(currentYLeft)
			currentYLeft += spacing
		}
	}
}

func (ui *UI) buildCenterPanel() {
	ui.CenterPanel.Elements = nil // очищаем перед построением

	panelOffsetX := ui.CenterPanel.X + 10
	panelWidth := ui.CenterPanel.Width/2 - 5
	topY := ui.CenterPanel.Y + 10 // верхний отступ для первой строки категорий

	switch ui.ActiveTab {
	case "main":
		if ui.State.CurrentLocation == nil {
			return // если локация не выбрана, ничего не строим
		}

		// 1️⃣ Колонка слева
		longCategory := &UICategory{
			Title:    "Длительные действия",
			X:        panelOffsetX,
			Y:        topY, // верхняя позиция для первой категории
			Width:    panelWidth,
			HeaderH:  28,
			Expanded: true,
			Font:     ui.GameFont,
			Parent:   ui.CenterPanel,
		}

		instantCategory := &UICategory{
			Title:    "Мгновенные действия",
			X:        panelOffsetX,
			Y:        topY, // верхняя позиция; Layout для кнопок внутри категории сдвинет их вниз
			Width:    panelWidth,
			HeaderH:  28,
			Expanded: true,
			Font:     ui.GameFont,
			Parent:   ui.CenterPanel,
		}

		// 2️⃣ Колонка справа — Дуэли
		duelCategory := &UICategory{
			Title:    "Дуэли",
			X:        panelOffsetX + panelWidth + 5, // смещение вправо
			Y:        topY,                          // верхняя позиция параллельно первым категориям
			Width:    panelWidth,
			HeaderH:  28,
			Expanded: true,
			Font:     ui.GameFont,
			Parent:   ui.CenterPanel,
		}

		// создаём кнопки для текущей локации
		for _, id := range ui.State.CurrentLocation.AvailableActions {
			info, ok := system.GetActionInfo(id)
			if !ok {
				continue
			}

			btn := &UIButton{
				Text:     info.Name,
				ActionID: id,
				Width:    panelWidth,
				Height:   28,
				Color:    color.RGBA{120, 80, 50, 255},
				Font:     ui.GameFont,
			}

			switch info.Category {
			case "instant":
				btn.OnClick = func(id string) func() {
					return func() { ui.State.StartAction(id) }
				}(id)
				instantCategory.Elements = append(instantCategory.Elements, btn)

			case "long":
				btn.OnClick = func(id string, dur float64) func() {
					return func() { ui.OnAction(id, dur) }
				}(id, info.Duration)
				longCategory.Elements = append(longCategory.Elements, btn)

			case "duel":
				btn.Color = color.RGBA{160, 60, 60, 255}
				btn.OnClick = func(id string, dur float64) func() {
					return func() { ui.OnAction(id, dur) }
				}(id, info.Duration)
				duelCategory.Elements = append(duelCategory.Elements, btn)
			}
		}

		// добавляем категории только если они не пустые
		if len(longCategory.Elements) > 0 {
			ui.CenterPanel.Elements = append(ui.CenterPanel.Elements, longCategory)
		}
		if len(instantCategory.Elements) > 0 {
			ui.CenterPanel.Elements = append(ui.CenterPanel.Elements, instantCategory)
		}
		if len(duelCategory.Elements) > 0 {
			ui.CenterPanel.Elements = append(ui.CenterPanel.Elements, duelCategory)
		}

	case "inventory":
		invCategory := &UICategory{
			Title:    "Инвентарь",
			X:        panelOffsetX,
			Y:        topY,
			Width:    panelWidth,
			HeaderH:  28,
			Expanded: true,
			Font:     ui.GameFont,
			Parent:   ui.CenterPanel,
		}

		for res := range ui.State.Resources {
			text := fmt.Sprintf("%s: %.1f",
				system.ResourceDisplayName(res),
				ui.State.GetResource(res),
			)
			label := &UIButton{
				Text:   text,
				Width:  panelWidth,
				Height: 28,
				Color:  color.RGBA{80, 60, 40, 255},
				Font:   ui.GameFont,
			}
			invCategory.Elements = append(invCategory.Elements, label)
		}
		ui.CenterPanel.Elements = append(ui.CenterPanel.Elements, invCategory)
	}

	// расставляем кнопки внутри категорий по вертикали
	ui.layoutCenterTwoColumns()
}

//
// ---------- UPDATE ----------
//

func (ui *UI) Update() {

	// 1️⃣ Layout (пересчёт координат)
	ui.LeftPanel.Layout()
	ui.layoutCenterTwoColumns()

	// 2️⃣ Hover reset
	ui.HoveredButton = nil
	mx, my := ebiten.CursorPosition()

	// 3️⃣ Проверка hover
	for _, el := range ui.CenterPanel.Elements {
		if cat, ok := el.(*UICategory); ok && cat.Expanded {
			for _, child := range cat.Elements {
				if btn, ok := child.(*UIButton); ok {
					if btn.IsHovered(mx, my) {
						ui.HoveredButton = btn
					}
				}
			}
		}
	}

	// 4️⃣ Обновление прогресса кнопок
	for _, el := range ui.CenterPanel.Elements {
		if cat, ok := el.(*UICategory); ok {
			for _, child := range cat.Elements {
				if btn, ok := child.(*UIButton); ok {
					if ui.State.CurrentAction != nil && btn.ActionID == ui.State.CurrentAction.ID {
						progress := ui.State.CurrentAction.Progress / ui.State.CurrentAction.Duration
						if progress > 1 {
							progress = 1
						}
						btn.Progress = progress
						btn.Duration = ui.State.CurrentAction.Duration
						btn.Remaining = ui.State.CurrentAction.Duration - ui.State.CurrentAction.Progress
						if btn.Remaining < 0 {
							btn.Remaining = 0
						}
					} else {
						btn.Progress = 0
						btn.Remaining = 0
					}
				}
			}
		}
	}

	// 5️⃣ Обработка клика (ПОСЛЕ hover)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ui.LeftPanel.HandleClick(mx, my)
		if ui.ActiveTab == "main" {
			ui.LocationsPanel.HandleClick(mx, my)
		}
		ui.CenterPanel.HandleClick(mx, my)
	}

	// 6️⃣ Логи → уведомления
	if ui.State != nil && len(ui.State.Log) > ui.lastLogIndex {

		last := ui.State.Log[len(ui.State.Log)-1]

		ui.Notification = &Notification{
			Text:  last,
			Timer: 3.0,
		}

		ui.lastLogIndex = len(ui.State.Log)
	}

	// 7️⃣ Таймер уведомления
	if ui.Notification != nil {
		ui.Notification.Timer -= 1.0 / 60.0
		if ui.Notification.Timer <= 0 {
			ui.Notification = nil
		}
	}
}

// ---------- UI NOTIFICATIONS ----------

type Notification struct {
	Text  string
	Timer float64
}

// Отрисовка уведомления (вызывать в Draw после всех панелей)
func (ui *UI) drawNotification(screen *ebiten.Image) {
	if ui.Notification == nil {
		return
	}

	x, y := 250.0, 250.0
	width, height := 400.0, 40.0
	ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{50, 50, 50, 200})

	w, h := text.Measure(ui.Notification.Text, ui.GameFont, 1.0)
	tx := x + (width-w)/2
	ty := y + (height-h)/2
	DrawText(screen, ui.Notification.Text, ui.GameFont, tx, ty, color.RGBA{255, 255, 255, 255})
}

//
// ---------- PROGRESS BAR ----------
//

func (ui *UI) drawProgress(screen *ebiten.Image) {

	if ui.State.CurrentAction == nil {
		return
	}

	// Новое: проверяем, можно ли выполнить текущее действие
	if !ui.State.CanPerformAction(ui.State.CurrentAction.ID) {
		return
	}

	progress := ui.State.CurrentAction.Progress /
		ui.State.CurrentAction.Duration

	if progress > 1 {
		progress = 1
	}

	barX := 260.0
	barY := 100.0
	barWidth := 400.0
	barHeight := 25.0

	ebitenutil.DrawRect(screen,
		barX, barY,
		barWidth, barHeight,
		color.RGBA{80, 50, 30, 255},
	)

	ebitenutil.DrawRect(screen,
		barX, barY,
		barWidth*progress, barHeight,
		color.RGBA{200, 120, 40, 255},
	)

	remaining := ui.State.CurrentAction.Duration -
		ui.State.CurrentAction.Progress

	if remaining < 0 {
		remaining = 0
	}

	DrawText(screen,
		fmt.Sprintf("%.2f", remaining),
		ui.GameFont,
		barX+barWidth-60,
		barY+7.5,
		color.RGBA{255, 255, 255, 255},
	)
}

//
// ---------- RESOURCES PANEL ----------
//

func (ui *UI) drawResources(screen *ebiten.Image) {
	if ui.State == nil || ui.RightPanel == nil {
		return
	}

	padding := 10.0 // отступ внутри панели
	y := ui.RightPanel.Y + padding

	keys := make([]system.Resource, 0, len(ui.State.Resources))
	for r := range ui.State.Resources {
		keys = append(keys, r)
	}

	// сортировка по ресурсу
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, resource := range keys {
		current := ui.State.GetResource(resource)
		max := 0.0
		if ui.State.MaxResources != nil {
			max = ui.State.MaxResources[resource]
		}

		displayName := system.ResourceDisplayName(resource)

		var text string
		if max > 0 {
			text = fmt.Sprintf("%s: %.1f / %.1f", displayName, current, max)
		} else {
			text = fmt.Sprintf("%s: %.1f", displayName, current)
		}

		// позиция x привязана к панели
		x := ui.RightPanel.X + padding
		DrawText(
			screen,
			text,
			ui.GameFont,
			x,
			y,
			color.RGBA{255, 255, 255, 255},
		)
		y += 30 // расстояние между ресурсами
	}
}

//
//-----------КАТЕГОРИЗАЦИЯ----------
//

type UICategory struct {
	Title    string
	X, Y     float64
	Width    float64
	HeaderH  float64
	Expanded bool
	Font     text.Face
	Color    color.RGBA
	Elements []UIElement
	Parent   *Panel
}

func (c *UICategory) Draw(screen *ebiten.Image) {
	// Заголовок категории
	ebitenutil.DrawRect(screen,
		c.X, c.Y,
		c.Width, c.HeaderH,
		c.Color,
	)

	DrawText(screen,
		c.Title,
		c.Font,
		c.X+10,
		c.Y+5,
		color.RGBA{255, 255, 255, 255},
	)

	// Кнопки отрисовываем только если развернута категория
	if c.Expanded {
		for _, el := range c.Elements {
			el.Draw(screen)
		}
	}
}

func (c *UICategory) HandleClick(mx, my int) {
	// Клик по заголовку
	if float64(mx) >= c.X && float64(mx) <= c.X+c.Width &&
		float64(my) >= c.Y && float64(my) <= c.Y+c.HeaderH {

		c.Expanded = !c.Expanded
		// Обновляем Layout панели, чтобы сдвинуть все заголовки
		if c.Parent != nil {
			c.Parent.Layout()
		}
		return
	}

	// Клик по кнопкам
	if c.Expanded {
		for _, el := range c.Elements {
			el.HandleClick(mx, my)
		}
	}
}

func (c *UICategory) Layout(startY float64) float64 {
	c.Y = startY
	currentY := startY + c.HeaderH

	if c.Expanded {
		topPadding := 5.0 // отступ первой кнопки от хедера
		spacing := 5.0    // отступ между остальными кнопками

		for i, el := range c.Elements {
			if btn, ok := el.(*UIButton); ok {
				btn.X = c.X
				if i == 0 {
					btn.Y = currentY + topPadding
					currentY += btn.Height + topPadding + spacing
				} else {
					btn.Y = currentY
					currentY += btn.Height + spacing
				}
				btn.Width = c.Width
			}
		}
	}
	//Возвращает текущий Y для следующей категории
	return currentY
}

//
// ---------- TEXT HELPER ----------
//

func DrawText(screen *ebiten.Image,
	str string,
	font text.Face,
	x, y float64,
	clr color.RGBA) {

	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x, y)

	r, g, b, a := clr.RGBA()
	opts.ColorM.Scale(
		float64(r)/65535,
		float64(g)/65535,
		float64(b)/65535,
		float64(a)/65535,
	)

	text.Draw(screen, str, font, opts)
}

//
// ---------- DRAW ----------
//

func (ui *UI) Draw(screen *ebiten.Image) {

	ui.LeftPanel.Draw(screen)
	if ui.ActiveTab == "main" {
		ui.LocationsPanel.Draw(screen)
	}
	ui.CenterPanel.Draw(screen)
	ui.RightPanel.Draw(screen)

	ui.drawResources(screen)
	//ui.drawNotification(screen)

	ui.drawTooltip(screen)
}
