package inventory

// Weapon — базовая структура оружия
type Weapon struct {
	Item
	WeaponType  string     // "revolver", "rifle" и т.д.
	Damage      float64    // базовый урон за выстрел
	Range       float64    // эффективная дальность
	ReloadTime  float64    // время перезарядки в секундах
	AmmoType    string     // тип патрона, например "unitary"
	Capacity    int        // ёмкость барабана/магазина
	Hands       []SlotType // слоты, куда можно экипировать (SlotWeapon, SlotOffHand)
	TwoHanded   bool       // если true — занимает оба слота оружия
	Description string
}

// Пример стартового револьвера
var ColtNavy1851 = Weapon{
	Item: Item{
		ID:   "colt_navy_1851_unitary",
		Name: "Colt Navy 1851 (унитарный)",
		Type: "revolver",
	},
	WeaponType:  "revolver",
	Damage:      25,
	Range:       50,
	ReloadTime:  4.0,
	AmmoType:    "unitary",
	Capacity:    6,
	Hands:       []SlotType{SlotWeapon, SlotOffHand},
	TwoHanded:   false, // одноручное оружие
	Description: "Классический револьвер Colt 1851 Navy, модернизированный под унитарные патроны. Надёжный, но уже устаревший.",
}

// EquipWeapon экипирует оружие в выбранный слот
func (e *Equipment) EquipWeapon(w *Weapon, preferredSlot SlotType) {
	if e.Slots == nil {
		e.Slots = make(map[SlotType]*Item)
	}

	// Двуручное оружие — занимает сразу оба слота
	if w.TwoHanded {
		e.Slots[SlotWeapon] = &w.Item
		e.Slots[SlotOffHand] = &w.Item
		return
	}

	// Проверка, чтобы preferredSlot был разрешён для данного оружия
	valid := false
	for _, s := range w.Hands {
		if s == preferredSlot {
			valid = true
			break
		}
	}
	if !valid && len(w.Hands) > 0 {
		preferredSlot = w.Hands[0]
	}

	// Снимаем старый предмет из слота
	if oldItem, exists := e.Slots[preferredSlot]; exists {
		// здесь можно вернуть старое оружие в инвентарь, если нужно
		_ = oldItem
	}

	// Экипируем оружие
	e.Slots[preferredSlot] = &w.Item
}
