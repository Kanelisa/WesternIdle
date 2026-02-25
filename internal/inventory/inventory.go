package inventory

//-----------Инвентарь---------------

type SlotType string

const (
	SlotHead      SlotType = "head"
	SlotBody      SlotType = "body"
	SlotLegs      SlotType = "legs"
	SlotWeapon    SlotType = "weapon"
	SlotOffHand   SlotType = "offhand"
	SlotAccessory SlotType = "accessory"
)

var AllSlots = []SlotType{
	SlotHead,
	SlotBody,
	SlotLegs,
	SlotWeapon,
	SlotOffHand,
	SlotAccessory,
}

func (s SlotType) DisplayName() string {
	switch s {
	case SlotHead:
		return "Голова"
	case SlotBody:
		return "Тело"
	case SlotLegs:
		return "Ноги"
	case SlotWeapon:
		return "Оружие"
	case SlotOffHand:
		return "Вторая рука"
	case SlotAccessory:
		return "Аксессуар"
	default:
		return string(s)
	}
}

type Inventory struct {
	Items map[string]*Item
}

func NewInventory() *Inventory {
	return &Inventory{
		Items: make(map[string]*Item),
	}
}

func (inv *Inventory) Add(item Item) {
	inv.Items[item.ID] = &item
}

func (inv *Inventory) Remove(id string, amount int) bool {
	if item, ok := inv.Items[id]; ok {
		if item.Quantity >= amount {
			item.Quantity -= amount
			if item.Quantity == 0 {
				delete(inv.Items, id)
			}
			return true
		}
	}
	return false
}

//------------Предметы---------------

type Item struct {
	ID       string
	Name     string
	Type     string
	Slot     SlotType // В какой слот можно надеть
	Quantity int
	// Бонусы (можно расширять)
	Damage  int
	Defense int
	Speed   int
}

//-----------Экипировка--------------

type Equipment struct {
	Slots map[SlotType]*Item
}

type EquipWeapon struct {
	Slots map[SlotType]*Weapon
}

func NewEquipment() *Equipment {
	return &Equipment{
		Slots: make(map[SlotType]*Item),
	}
}

func (e *Equipment) Equip(inv *Inventory, itemID string) bool {
	item, ok := inv.Items[itemID]
	if !ok {
		return false
	}

	if item.Slot == "" {
		return false // не экипируемый предмет
	}

	// Снимаем старый предмет из слота
	if oldItem, exists := e.Slots[item.Slot]; exists {
		inv.Add(*oldItem)
	}

	// Убираем предмет из инвентаря
	inv.Remove(itemID, 1)

	// Копируем предмет в слот
	e.Slots[item.Slot] = &Item{
		ID:      item.ID,
		Name:    item.Name,
		Slot:    item.Slot,
		Damage:  item.Damage,
		Defense: item.Defense,
		Speed:   item.Speed,
	}

	return true
}

func (e *Equipment) Unequip(inv *Inventory, slot SlotType) bool {
	item, ok := e.Slots[slot]
	if !ok {
		return false
	}

	inv.Add(*item)
	delete(e.Slots, slot)
	return true
}

func (e *Equipment) TotalStats() (damage, defense, speed int) {
	for _, item := range e.Slots {
		damage += item.Damage
		defense += item.Defense
		speed += item.Speed
	}
	return
}
