Mud = {};

Mud.Equipment = {
    Weapon:   Mud.Object,
    Sidearm:  Mud.Object,
    Shield:   Mud.Object,
    Helmet:   Mud.Object,
    Torso:    Mud.Object,
    Belt:     Mud.Object,
    Arms:     Mud.Object,
    Legs:     Mud.Object,
    Shoes:    Mud.Object,
    Ring:     Mud.Object,
    Floating: Mud.Object
} 

Mud.Object = {
    Type:           "",
    Name:           "",
    Description:    "",
    Weight:          1,
    Value:           0,
    Location:       "",
    Strength:        0,
    Intelligence:    0,
    Wisdom:          0,
    Effect:         "",
    SharpProtection: 0,
    BluntProtection: 0,
    Resistance:      0
}

Mud.Coordinate = {
    Row: 5,
    Col: 5,
    Level: 5
}
