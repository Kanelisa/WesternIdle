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
		Name: "Colt Navy 1851",
		Type: "revolver",
	},
	WeaponType:  "revolver",
	Damage:      25,
	Range:       50,
	ReloadTime:  4.0,
	AmmoType:    "unitary_revolver",
	Capacity:    6,
	Hands:       []SlotType{SlotWeapon, SlotOffHand},
	TwoHanded:   false, // одноручное оружие
	Description: "Классический револьвер Colt 1851 Navy, модернизированный под унитарные патроны. Надёжный, но уже устаревший.",
}

// EquipWeapon экипирует оружие в выбранный слот
func (e *Equipment) EquipWeapon(inv *Inventory, w *Weapon, preferredSlot SlotType) {

	if e.Slots == nil {
		e.Slots = make(map[SlotType]*Item)
	}

	// --- 1. Проверка допустимого слота ---
	valid := false
	for _, s := range w.Hands {
		if s == preferredSlot {
			valid = true
			break
		}
	}
	if !valid {
		preferredSlot = w.Hands[0]
	}

	// --- 2. Если оружие уже экипировано в другой руке ---
	for slot, equipped := range e.Slots {
		if equipped.ID == w.ID {

			// Если это тот же самый слот — ничего не делаем
			if slot == preferredSlot {
				return
			}

			// Удаляем из старого слота (перемещение, НЕ в инвентарь)
			delete(e.Slots, slot)
		}
	}

	// --- 3. Если новое оружие двуручное ---
	if w.TwoHanded {

		// Снимаем оба слота
		for _, slot := range []SlotType{SlotWeapon, SlotOffHand} {
			if oldItem, exists := e.Slots[slot]; exists {
				inv.Add(*oldItem)
				delete(e.Slots, slot)
			}
		}

		e.Slots[SlotWeapon] = &w.Item
		e.Slots[SlotOffHand] = &w.Item
		return
	}

	// --- 4. Если стоит двуручное — снять полностью ---
	if main, ok := e.Slots[SlotWeapon]; ok {
		if off, ok2 := e.Slots[SlotOffHand]; ok2 && main.ID == off.ID {
			inv.Add(*main)
			delete(e.Slots, SlotWeapon)
			delete(e.Slots, SlotOffHand)
		}
	}

	// --- 5. Если в выбранном слоте есть другое оружие ---
	if oldItem, exists := e.Slots[preferredSlot]; exists {
		inv.Add(*oldItem)
		delete(e.Slots, preferredSlot)
	}

	// --- 6. Экипируем ---
	e.Slots[preferredSlot] = &w.Item
}

var WeaponRegistry = map[string]*Weapon{
	ColtNavy1851.ID: &ColtNavy1851,
}

func GetWeaponByID(id string) *Weapon {
	if w, ok := WeaponRegistry[id]; ok {
		return w
	}
	return nil
}
