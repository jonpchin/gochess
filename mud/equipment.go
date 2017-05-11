package mud

type Equipment struct {
	Weapon Weapon
}

type Weapon struct {
	Object
	Strength string
	Effect   string // Not all weapons have special effects
}

type Armor struct {
	Object
	Effect string
}

type Dagger struct {
	Weapon
}

type Axe struct {
	Weapon
}
