package mud

type Object struct {
	Type        string // Used to identify what type of item such as potion, dagger, shield, etc
	Name        string
	Description string
	Weight      int        // Players can only carry so much weight in their inventory
	Value       string     // 50% of value is what shop owners are willing to pay for this item
	Location    Coordinate // Used to locate where the object is in the world
}

type Potion struct {
	Object
	Effect string // Side effect given to player when potion is quaffed such as +4 strength
}

func placeObject() {

}
