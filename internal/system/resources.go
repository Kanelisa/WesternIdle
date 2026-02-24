package system

//
// ---------- КОНСТАНТЫ РЕСУРСОВ ----------
//

// Resource — тип ресурса
type Resource string

const (
	ResourceMoney    Resource = "money"
	ResourceAmmo     Resource = "ammo"
	ResourceAccuracy Resource = "accuracy"
)

//
// ---------- РЕЕСТР РЕСУРСОВ ----------
//

// DefaultResources возвращает стартовые ресурсы игрока
func DefaultResources() map[Resource]float64 {
	return map[Resource]float64{
		ResourceMoney:    0,
		ResourceAmmo:     0,
		ResourceAccuracy: 1,
	}
}

func DefaultMaxResources() map[Resource]float64 {
	return map[Resource]float64{
		ResourceMoney:    100,
		ResourceAmmo:     30,
		ResourceAccuracy: 100,
	}
}

//
// ---------- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ----------
//

// ResourceDisplayName — красивое имя ресурса для UI
func ResourceDisplayName(r Resource) string {
	switch r {
	case ResourceMoney:
		return "Деньги"
	case ResourceAmmo:
		return "Патроны"
	case ResourceAccuracy:
		return "Меткость"
	default:
		return string(r)
	}
}

func ResourceOrder() []Resource {
	return []Resource{
		ResourceMoney,
		ResourceAmmo,
		ResourceAccuracy,
	}
}
