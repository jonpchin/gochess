package mud

type Equipment struct {
	Weapon   Weapon
	Sidearm  Sidearm
	Shield   Shield
	Helmet   Helmet
	Torso    Torso
	Belt     Belt
	Arms     Arms
	Legs     Legs
	Shoes    Shoes
	Ring     Ring
	Floating Floating
}

type Weapon struct {
	Object
	Strength int
	Effect   string // Not all weapons have special effects
}

type Sidearm struct {
	Object
	Strength int
	Effect   string // Not all weapons have special effects
}

type Helmet struct {
	Armor
}

type Shield struct {
	Armor
}

type Torso struct {
	Armor
}

type Belt struct {
	Armor
}

type Arms struct {
	Armor
}

type Legs struct {
	Armor
}

type Shoes struct {
	Armor
}

type Ring struct {
	Armor
}

type Floating struct {
	Armor
}

type Armor struct {
	Object
	SharpProtection int // Sharp object physical defense such as spears, axe, daggers
	BluntProtection int // Blunt objct physical defense such as clubs, mace, brass knuckles
	Resistance      int // Magical defense
	Effect          string
}

type Dagger struct {
	Weapon
}

type Axe struct {
	Weapon
}
