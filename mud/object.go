package mud

type ObjectInfo struct {
	Type        string // Used to identify what type of item such as potion, dagger, shield, etc
	Name        string
	Description string
	Weight      int        // Players can only carry so much weight in their inventory
	Value       int        // 50% of value is what shop owners are willing to pay for this item
	Location    Coordinate // Used to locate where the object is in the world
}

type Weapon struct {
	ObjectInfo
	DamageType string // Blunt, sharp, resistance (magic)
	Strength   int
	Effect     string
}

// Magic users won't be able to use sidearms
type Sidearm struct {
	ObjectInfo
	DamageType string // Blunt, sharp, resistance (magic)
	Strength   int
	Effect     string
}

type Shield struct {
	ObjectInfo
	SharpProtection int
	BluntProtection int
	Resistance      int
}

type Helmet struct {
	ObjectInfo
	SharpProtection int
	BluntProtection int
	Resistance      int
}

type Torso struct {
	ObjectInfo
	SharpProtection int
	BluntProtection int
	Resistance      int
}

type Belt struct {
	ObjectInfo
	SharpProtection int
	BluntProtection int
	Resistance      int
}

type Arms struct {
	ObjectInfo
	SharpProtection int
	BluntProtection int
	Resistance      int
}

type Legs struct {
	ObjectInfo
	SharpProtection int
	BluntProtection int
	Resistance      int
}

type Shoes struct {
	ObjectInfo
	SharpProtection int
	BluntProtection int
	Resistance      int
}

type Ring struct {
	ObjectInfo
	Intelligence int
	Wisdom       int
	Effect       string
	Resistance   int
}

type Floating struct {
	ObjectInfo
	Intelligence int
	Wisdom       int
	Effect       string
}

type Item struct {
	ObjectInfo
}

type Potion struct {
	ObjectInfo
	Effect string
}

type Scroll struct {
	ObjectInfo
	Effect string
}
