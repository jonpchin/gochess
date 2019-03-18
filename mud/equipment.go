package mud

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

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

// List of effects a weapon can afflict on a target
var WEAPON_EFFECTS = []string{
	"shock",
	"shiver",
	"burn",
	"weaken",
	"poison",
	"absorb", // Absorbs health
}

// TODO: Generate location for equipment
func GenerateEquipmentStats() {

	const generatedEquipmentFolder = "mud/equipment/generated"
	const equipmentFolder = "mud/equipment"

	files, err := ioutil.ReadDir(generatedEquipmentFolder)
	if err != nil {
		log.Fatal(err)
	}

	var belts []Belt
	var weapons []Weapon
	var shoes []Shoes
	var helmets []Helmet
	var legs []Legs
	var shields []Shield
	var torsos []Torso

	BELT_MAX_STAT := 15
	BELT_MAX_WEIGHT := 10
	WEAPON_MAX_STAT := 25
	WEAPON_MAX_WEIGHT := 20
	SHOES_MAX_STAT := 14
	SHOES_MAX_WEIGHT := 12
	HELMET_MAX_STAT := 15
	HELMET_MAX_WEIGHT := 11
	LEG_MAX_STAT := 16
	LEG_MAX_WEIGHT := 15
	SHIELD_MAX_STAT := 18
	SHIELD_MAX_WEIGHT := 20
	TORSO_MAX_STAT := 20
	TORSO_MAX_WEIGHT := 25

	for _, f := range files {
		equipmentTxtFile, err := os.Open(generatedEquipmentFolder + "/" + f.Name())
		defer equipmentTxtFile.Close()

		if err != nil {
			fmt.Println("GenerateEquipmentStats 0", err)
			return
		}

		scanner := bufio.NewScanner(equipmentTxtFile)

		for scanner.Scan() {
			name := scanner.Text()
			fileName := f.Name()
			if fileName == "belts.txt" {
				var belt Belt
				belt.Type = "belt"
				belt.Name = name
				belt.Description = "This is a belt."
				belt.Weight = rand.Intn(BELT_MAX_WEIGHT)
				belt.SharpProtection = rand.Intn(BELT_MAX_STAT)
				belt.BluntProtection = rand.Intn(BELT_MAX_STAT)
				belt.Resistance = rand.Intn(BELT_MAX_STAT)
				belt.Value = 2*belt.SharpProtection*belt.BluntProtection - belt.Weight
				belts = append(belts, belt)
			} else if fileName == "daggers.txt" {
				var weapon Weapon
				weapon.Type = "dagger"
				weapon.Name = name
				weapon.Description = "This is a dagger."
				weapon.Weight = rand.Intn(WEAPON_MAX_WEIGHT)
				weapon.DamageType = "sharp"
				weapon.Strength = rand.Intn(WEAPON_MAX_STAT)

				// About %25 of weapons will have the ability to cause afflictions
				weightedEffect := rand.Intn(3)
				weapon.Effect = "None"
				if weightedEffect == 0 {
					weapon.Effect = WEAPON_EFFECTS[rand.Intn(len(WEAPON_EFFECTS)-1)]
				}
				weapon.Value = 3*weapon.Strength - weapon.Weight

				weapons = append(weapons, weapon)
			} else if fileName == "boots.txt" {
				var shoe Shoes
				shoe.Type = "shoes"
				shoe.Name = name
				shoe.Description = "This is a pair of shoes."
				shoe.Weight = rand.Intn(SHOES_MAX_WEIGHT)
				shoe.SharpProtection = rand.Intn(SHOES_MAX_STAT)
				shoe.BluntProtection = rand.Intn(SHOES_MAX_STAT)
				shoe.Resistance = rand.Intn(SHOES_MAX_STAT)
				shoe.Value = 2*shoe.SharpProtection*shoe.BluntProtection*shoe.Resistance - shoe.Weight
				shoes = append(shoes, shoe)
			} else if fileName == "helmets.txt" {
				var helmet Helmet
				helmet.Name = name
				helmet.Description = "This is a helmet."
				helmet.Weight = rand.Intn(HELMET_MAX_WEIGHT)
				helmet.SharpProtection = rand.Intn(HELMET_MAX_STAT)
				helmet.BluntProtection = rand.Intn(HELMET_MAX_STAT)
				helmet.Resistance = rand.Intn(HELMET_MAX_STAT)
				helmet.Value = 2*helmet.SharpProtection*helmet.BluntProtection*helmet.Resistance - helmet.Weight
				helmets = append(helmets, helmet)
			} else if fileName == "legs.txt" {
				var leg Legs
				leg.Name = name
				leg.Description = "This is a legging."
				leg.Weight = rand.Intn(LEG_MAX_WEIGHT)
				leg.SharpProtection = rand.Intn(LEG_MAX_STAT)
				leg.BluntProtection = rand.Intn(LEG_MAX_STAT)
				leg.Resistance = rand.Intn(LEG_MAX_STAT)
				leg.Value = 2*leg.SharpProtection*leg.BluntProtection*leg.Resistance - leg.Weight
				legs = append(legs, leg)
			} else if fileName == "shields.txt" {
				var shield Shield
				shield.Name = name
				shield.Description = "This is a shield."
				shield.Weight = rand.Intn(SHIELD_MAX_WEIGHT)
				shield.SharpProtection = rand.Intn(SHIELD_MAX_STAT)
				shield.BluntProtection = rand.Intn(SHIELD_MAX_STAT)
				shield.Resistance = rand.Intn(SHIELD_MAX_STAT)
				shield.Value = 2*shield.SharpProtection*shield.BluntProtection*shield.Resistance - shield.Weight
				shields = append(shields, shield)
			} else if fileName == "torso.txt" {
				var torso Torso
				torso.Name = name
				torso.Description = "This is a torso."
				torso.Weight = rand.Intn(TORSO_MAX_WEIGHT)
				torso.SharpProtection = rand.Intn(TORSO_MAX_STAT)
				torso.BluntProtection = rand.Intn(TORSO_MAX_STAT)
				torso.Resistance = rand.Intn(TORSO_MAX_STAT)
				torso.Value = 2*torso.SharpProtection*torso.BluntProtection*torso.Resistance - torso.Weight
				torsos = append(torsos, torso)
			}
		}

	}

	beltsFinal, err := json.Marshal(belts)
	if err != nil {
		fmt.Println("GenerateEquipmentStats belts 1", err)
	}
	err = ioutil.WriteFile(equipmentFolder+"/belts.txt", beltsFinal, 0644)
	if err != nil {
		fmt.Println("GenerateEquipmentStats belts 2", err)
	}

	weaponsFinal, err := json.Marshal(weapons)
	if err != nil {
		fmt.Println("GenerateEquipmentStats weapons 1 ", err)
	}
	err = ioutil.WriteFile(equipmentFolder+"/weapons.txt", weaponsFinal, 0644)
	if err != nil {
		fmt.Println("GenerateEquipmentStats weapons", err)
	}

	shoesFinal, err := json.Marshal(shoes)
	if err != nil {
		fmt.Println("GenerateEquipmentStats shoes 1", err)
	}
	err = ioutil.WriteFile(equipmentFolder+"/shoes.txt", shoesFinal, 0644)
	if err != nil {
		fmt.Println("GenerateEquipmentStats shoes 2", err)
	}

	helmetFinal, err := json.Marshal(helmets)
	if err != nil {
		fmt.Println("GenerateEquipmentStats helmet 1", err)
	}
	err = ioutil.WriteFile(equipmentFolder+"/helmet.txt", helmetFinal, 0644)
	if err != nil {
		fmt.Println("GenerateEquipmentStats helmet 2", err)
	}

	legsFinal, err := json.Marshal(legs)
	if err != nil {
		fmt.Println("GenerateEquipmentStats legs 1", err)
	}
	err = ioutil.WriteFile(equipmentFolder+"/legs.txt", legsFinal, 0644)
	if err != nil {
		fmt.Println("GenerateEquipmentStats legs 2", err)
	}

	shieldsFinal, err := json.Marshal(shields)
	if err != nil {
		fmt.Println("GenerateEquipmentStats shields 1", err)
	}
	err = ioutil.WriteFile(equipmentFolder+"/shields.txt", shieldsFinal, 0644)
	if err != nil {
		fmt.Println("GenerateEquipmentStats shields 2", err)
	}

	torsosFinal, err := json.Marshal(torsos)
	if err != nil {
		fmt.Println("GenerateEquipmentStats torsos 1", err)
	}
	err = ioutil.WriteFile(equipmentFolder+"/torso.txt", torsosFinal, 0644)
	if err != nil {
		fmt.Println("GenerateEquipmentStats torsos 2", err)
	}
}

func generateBeltStats(beltName string) {
	const equipmentFolder = "mud/equipment"

}
