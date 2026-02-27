package system

import (
	"WesternIdle/internal/inventory"
	"time"
)

// ---------------- Game State ----------------

type GameState struct {
	//--------Ресурсы---------
	Resources    map[Resource]float64
	MaxResources map[Resource]float64
	//--------Действия--------
	CurrentAction *ActiveAction
	Log           []string
	//--------Локации----------
	Locations       map[string]*Location
	CurrentLocation *Location
	//--------Сейв----------
	LastSaveTime time.Time
	//------Инвентарь-------
	Inventory *inventory.Inventory
	Equipment *inventory.Equipment
}

// ---------------- Start Action ----------------

func (s *GameState) StartAction(id string) bool {

	info, ok := GetActionInfo(id)
	if !ok {
		return false
	}

	// Проверка ресурсов
	if !s.CanPerformAction(id) {
		s.LogEvent("Недостаточно ресурсов!")
		return false
	}

	// ---- Instant действие ----
	if info.Instant {

		// Списание
		for res, cost := range info.Cost {
			s.AddResource(res, -cost)
		}

		// Resolver?
		if info.Resolver != nil {
			info.Resolver(s)
		} else {
			for res, reward := range info.Reward {
				s.AddResource(res, reward)
			}
			s.LogEvent(info.Name + " выполнено.")
		}

		return true
	}

	// ---- Если уже выполняется действие ----
	if s.CurrentAction != nil {
		return false
	}

	s.CurrentAction = &ActiveAction{
		ID:       id,
		Duration: info.Duration,
		Progress: 0,
	}

	if s.CurrentLocation != nil {
		allowed := false
		for _, a := range s.CurrentLocation.AvailableActions {
			if a == id {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}

	return true
}

//-----------------Смена локаций---------------------

func (s *GameState) ChangeLocation(id string) bool {

	if s.CurrentAction != nil {
		s.LogEvent("Нельзя сменить локацию во время действия.")
		return false
	}

	loc, ok := s.Locations[id]
	if !ok {
		return false
	}

	s.CurrentLocation = loc
	s.LogEvent("Вы переместились в: " + loc.Name)

	return true
}

// ---------------- Проверка ресурсов ----------------

func (s *GameState) CanPerformAction(id string) bool {

	info, ok := GetActionInfo(id)
	if !ok {
		return false
	}

	for res, cost := range info.Cost {
		if s.GetResource(res) < cost {
			return false
		}
	}

	for res := range info.Reward {
		if s.GetResource(res) >= s.MaxResources[res] {
			return false
		}
	}

	return true
}

// ---------------- Update Loop ----------------

func (s *GameState) Update(delta float64) {

	if s.CurrentAction == nil {
		return
	}

	s.CurrentAction.Progress += delta

	if s.CurrentAction.Progress >= s.CurrentAction.Duration {
		s.ResolveAction()
		s.CurrentAction = nil
	}
}

// ---------------- Resolve Action ----------------

func (s *GameState) ResolveAction() {

	if s.CurrentAction == nil {
		return
	}

	info, ok := GetActionInfo(s.CurrentAction.ID)
	if !ok {
		return
	}

	// Списание стоимости
	for res, cost := range info.Cost {
		s.AddResource(res, -cost)
	}

	// Кастомная логика
	if info.Resolver != nil {
		info.Resolver(s)
		return
	}

	// Стандартная награда
	for res, reward := range info.Reward {
		s.AddResource(res, reward)
	}

	s.LogEvent(info.Name + " завершено.")
}

// ---------------- Resources ----------------

func (s *GameState) AddResource(res Resource, value float64) {

	if s.Resources == nil {
		s.Resources = DefaultResources()
	}
	if s.MaxResources == nil {
		s.MaxResources = DefaultMaxResources()
	}

	s.Resources[res] += value

	// Не ниже 0
	if s.Resources[res] < 0 {
		s.Resources[res] = 0
	}

	// Ограничение сверху
	max := s.MaxResources[res]
	if max > 0 && s.Resources[res] > max {
		s.Resources[res] = max
	}
}

func (s *GameState) GetResource(res Resource) float64 {
	if s.Resources == nil {
		return 0
	}
	return s.Resources[res]
}

// ---------------- Logging ----------------

func (s *GameState) LogEvent(msg string) {
	s.Log = append(s.Log, msg)
	if len(s.Log) > 50 {
		s.Log = s.Log[1:]
	}
}

// ---------------- Game Load ----------------

func LoadGame() *GameState {
	// Создаём пустое состояние
	state := &GameState{
		Resources:    DefaultResources(),
		MaxResources: DefaultMaxResources(),
		Log:          []string{},
		LastSaveTime: time.Now(),
		Locations:    make(map[string]*Location),
		Inventory:    inventory.NewInventory(),
		Equipment:    inventory.NewEquipment(),
	}

	// Инициализируем локации
	InitLocations(state)
	state.Inventory.Add(inventory.ColtNavy1851.Item)

	// Первая локация по умолчанию
	if loc, ok := state.Locations["dust_town"]; ok {
		state.CurrentLocation = loc
	}

	return state
}
