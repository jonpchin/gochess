package mud

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"

	"github.com/jonpchin/gochess/gostuff"
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

// Read equipment files to update actual count
var TOTAL_BELTS = 100
var TOTAL_HELMETS = 100
var TOTAL_LEGS = 100
var TOTAL_SHIELDS = 100
var TOTAL_SHOES = 100
var TOTAL_TORSO = 100
var TOTAL_WEAPONS = 100

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
				belt.Value = 2*(belt.SharpProtection+1)*(belt.BluntProtection+1) - belt.Weight
				belts = append(belts, belt)
			} else if fileName == "daggers.txt" {
				var weapon Weapon
				weapon.Type = "dagger"
				weapon.Name = name
				weapon.Description = "This is a dagger."
				weapon.Weight = rand.Intn(WEAPON_MAX_WEIGHT)
				weapon.DamageType = "sharp"
				weapon.Strength = rand.Intn(WEAPON_MAX_STAT) + 1

				// About %25 of weapons will have the ability to cause afflictions
				weightedEffect := rand.Intn(3)
				weapon.Effect = "None"
				if weightedEffect == 0 {
					weapon.Effect = WEAPON_EFFECTS[rand.Intn(len(WEAPON_EFFECTS)-1)]
				}
				weapon.Value = WEAPON_MAX_WEIGHT*weapon.Strength - weapon.Weight

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
				shoe.Value = 2*(shoe.SharpProtection+1)*(shoe.BluntProtection+1)*(shoe.Resistance+1) - shoe.Weight
				shoes = append(shoes, shoe)
			} else if fileName == "helmets.txt" {
				var helmet Helmet
				helmet.Type = "helmet"
				helmet.Name = name
				helmet.Description = "This is a helmet."
				helmet.Weight = rand.Intn(HELMET_MAX_WEIGHT)
				helmet.SharpProtection = rand.Intn(HELMET_MAX_STAT)
				helmet.BluntProtection = rand.Intn(HELMET_MAX_STAT)
				helmet.Resistance = rand.Intn(HELMET_MAX_STAT)
				helmet.Value = 2*(helmet.SharpProtection+1)*(helmet.BluntProtection+1)*(helmet.Resistance+1) - helmet.Weight
				helmets = append(helmets, helmet)
			} else if fileName == "legs.txt" {
				var leg Legs
				leg.Name = name
				leg.Type = "leg"
				leg.Description = "This is a legging."
				leg.Weight = rand.Intn(LEG_MAX_WEIGHT)
				leg.SharpProtection = rand.Intn(LEG_MAX_STAT)
				leg.BluntProtection = rand.Intn(LEG_MAX_STAT)
				leg.Resistance = rand.Intn(LEG_MAX_STAT)
				leg.Value = 2*(leg.SharpProtection+1)*(leg.BluntProtection+1)*(leg.Resistance+1) - leg.Weight
				legs = append(legs, leg)
			} else if fileName == "shields.txt" {
				var shield Shield
				shield.Name = name
				shield.Type = "shield"
				shield.Description = "This is a shield."
				shield.Weight = rand.Intn(SHIELD_MAX_WEIGHT)
				shield.SharpProtection = rand.Intn(SHIELD_MAX_STAT)
				shield.BluntProtection = rand.Intn(SHIELD_MAX_STAT)
				shield.Resistance = rand.Intn(SHIELD_MAX_STAT)
				shield.Value = 2*(shield.SharpProtection+1)*(shield.BluntProtection+1)*(shield.Resistance+1) - shield.Weight
				shields = append(shields, shield)
			} else if fileName == "torso.txt" {
				var torso Torso
				torso.Name = name
				torso.Type = "torso"
				torso.Description = "This is a torso."
				torso.Weight = rand.Intn(TORSO_MAX_WEIGHT)
				torso.SharpProtection = rand.Intn(TORSO_MAX_STAT)
				torso.BluntProtection = rand.Intn(TORSO_MAX_STAT)
				torso.Resistance = rand.Intn(TORSO_MAX_STAT)
				torso.Value = 2*(torso.SharpProtection+1)*(torso.BluntProtection+1)*(torso.Resistance+1) - torso.Weight
				torsos = append(torsos, torso)
			}
		}
	}

	TOTAL_BELTS = len(belts)
	TOTAL_HELMETS = len(helmets)
	TOTAL_LEGS = len(legs)
	TOTAL_SHIELDS = len(shields)
	TOTAL_SHOES = len(shoes)
	TOTAL_TORSO = len(torsos)

	// Sort equipment by value so their drop frequency can be scaled with lower value drops occuring more often
	sort.Slice(belts, func(i, j int) bool { return belts[i].Value < belts[j].Value })
	sort.Slice(helmets, func(i, j int) bool { return helmets[i].Value < helmets[j].Value })
	sort.Slice(legs, func(i, j int) bool { return legs[i].Value < legs[j].Value })
	sort.Slice(shoes, func(i, j int) bool { return shoes[i].Value < shoes[j].Value })

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
	err = ioutil.WriteFile(equipmentFolder+"/helmets.txt", helmetFinal, 0644)
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
	err = ioutil.WriteFile(equipmentFolder+"/torsos.txt", torsosFinal, 0644)
	if err != nil {
		fmt.Println("GenerateEquipmentStats torsos 2", err)
	}
}

func getRandomBelt() Belt {

	var equipmentFile = "mud/equipment/belts.txt"
	var belt Belt

	if gostuff.IsFileExist(equipmentFile) == false {
		fmt.Println("File does not exist for ", equipmentFile)
		return belt
	}

	equipment, err := ioutil.ReadFile(equipmentFile)
	if err != nil {
		fmt.Println(err)
	}
	var data []map[string]interface{}

	err = json.Unmarshal(equipment, &data)
	if err != nil {
		fmt.Println(err)
	}

	beltWeights := []int{16, 5, 3, 1}
	itemsInEachSection := TOTAL_BELTS / len(beltWeights)
	sumOfWeights := 0
	for _, weight := range beltWeights {
		sumOfWeights += (weight * itemsInEachSection)
	}

	randomBeltWeight := rand.Intn(sumOfWeights)

	for index, equipment := range data {

		randomBeltWeight -= beltWeights[(index / itemsInEachSection)]

		if randomBeltWeight <= 0 {
			belt.Type = equipment["Type"].(string)
			belt.Name = equipment["Name"].(string)
			belt.Description = equipment["Description"].(string)
			belt.Weight = int(equipment["Weight"].(float64))
			belt.Value = int(equipment["Value"].(float64))
			belt.SharpProtection = int(equipment["SharpProtection"].(float64))
			belt.BluntProtection = int(equipment["BluntProtection"].(float64))
			belt.Resistance = int(equipment["Resistance"].(float64))
			return belt
		}
	}
	return belt
}

func getRandomHelmet() Helmet {
	var equipmentFile = "mud/equipment/helmets.txt"
	var helmet Helmet

	if gostuff.IsFileExist(equipmentFile) == false {
		fmt.Println("File does not exist for ", equipmentFile)
		return helmet
	}

	equipment, err := ioutil.ReadFile(equipmentFile)
	if err != nil {
		fmt.Println(err)
	}
	var data []map[string]interface{}

	err = json.Unmarshal(equipment, &data)
	if err != nil {
		fmt.Println(err)
	}

	helmetWeights := []int{16, 5, 3, 1}
	itemsInEachSection := TOTAL_HELMETS / len(helmetWeights)
	sumOfWeights := 0
	for _, weight := range helmetWeights {
		sumOfWeights += (weight * itemsInEachSection)
	}

	randomHelmetWeight := rand.Intn(sumOfWeights)

	for index, equipment := range data {

		randomHelmetWeight -= helmetWeights[(index / itemsInEachSection)]

		if randomHelmetWeight <= 0 {

			helmet.Type = equipment["Type"].(string)
			helmet.Name = equipment["Name"].(string)
			helmet.Description = equipment["Description"].(string)
			helmet.Weight = int(equipment["Weight"].(float64))
			helmet.Value = int(equipment["Value"].(float64))
			helmet.SharpProtection = int(equipment["SharpProtection"].(float64))
			helmet.BluntProtection = int(equipment["BluntProtection"].(float64))
			helmet.Resistance = int(equipment["Resistance"].(float64))
			return helmet
		}
	}
	return helmet
}

func getRandomLegs() Legs {

	var equipmentFile = "mud/equipment/legs.txt"
	var legs Legs

	if gostuff.IsFileExist(equipmentFile) == false {
		fmt.Println("File does not exist for ", equipmentFile)
		return legs
	}

	equipment, err := ioutil.ReadFile(equipmentFile)
	if err != nil {
		fmt.Println(err)
	}
	var data []map[string]interface{}

	err = json.Unmarshal(equipment, &data)
	if err != nil {
		fmt.Println(err)
	}

	legWeights := []int{16, 5, 3, 1}
	itemsInEachSection := TOTAL_LEGS / len(legWeights)
	sumOfWeights := 0
	for _, weight := range legWeights {
		sumOfWeights += (weight * itemsInEachSection)
	}

	randomLegWeights := rand.Intn(sumOfWeights)

	for index, equipment := range data {

		randomLegWeights -= legWeights[(index / itemsInEachSection)]

		if randomLegWeights <= 0 {

			legs.Type = equipment["Type"].(string)
			legs.Name = equipment["Name"].(string)
			legs.Description = equipment["Description"].(string)
			legs.Weight = int(equipment["Weight"].(float64))
			legs.Value = int(equipment["Value"].(float64))
			legs.SharpProtection = int(equipment["SharpProtection"].(float64))
			legs.BluntProtection = int(equipment["BluntProtection"].(float64))
			legs.Resistance = int(equipment["Resistance"].(float64))
			return legs
		}
	}
	return legs
}

func getRandomShield() Shield {

	var equipmentFile = "mud/equipment/shields.txt"
	var shield Shield

	if gostuff.IsFileExist(equipmentFile) == false {
		fmt.Println("File does not exist for ", equipmentFile)
		return shield
	}

	equipment, err := ioutil.ReadFile(equipmentFile)
	if err != nil {
		fmt.Println(err)
	}
	var data []map[string]interface{}

	err = json.Unmarshal(equipment, &data)
	if err != nil {
		fmt.Println(err)
	}

	shieldWeights := []int{16, 5, 3, 1}
	itemsInEachSection := TOTAL_SHIELDS / len(shieldWeights)
	sumOfWeights := 0
	for _, weight := range shieldWeights {
		sumOfWeights += (weight * itemsInEachSection)
	}

	randomShieldWeights := rand.Intn(sumOfWeights)

	for index, equipment := range data {

		randomShieldWeights -= shieldWeights[(index / itemsInEachSection)]

		if randomShieldWeights <= 0 {

			shield.Type = equipment["Type"].(string)
			shield.Name = equipment["Name"].(string)
			shield.Description = equipment["Description"].(string)
			shield.Weight = int(equipment["Weight"].(float64))
			shield.Value = int(equipment["Value"].(float64))
			shield.SharpProtection = int(equipment["SharpProtection"].(float64))
			shield.BluntProtection = int(equipment["BluntProtection"].(float64))
			shield.Resistance = int(equipment["Resistance"].(float64))
			return shield
		}
	}
	return shield
}

func getRandomShoes() Shoes {

	var equipmentFile = "mud/equipment/shoes.txt"
	var shoes Shoes

	if gostuff.IsFileExist(equipmentFile) == false {
		fmt.Println("File does not exist for ", equipmentFile)
		return shoes
	}

	equipment, err := ioutil.ReadFile(equipmentFile)
	if err != nil {
		fmt.Println(err)
	}
	var data []map[string]interface{}

	err = json.Unmarshal(equipment, &data)
	if err != nil {
		fmt.Println(err)
	}

	shoeWeights := []int{16, 5, 3, 1}
	itemsInEachSection := TOTAL_SHOES / len(shoeWeights)
	sumOfWeights := 0
	for _, weight := range shoeWeights {
		sumOfWeights += (weight * itemsInEachSection)
	}

	randomShoeWeights := rand.Intn(sumOfWeights)

	for index, equipment := range data {

		randomShoeWeights -= shoeWeights[(index / itemsInEachSection)]

		if randomShoeWeights <= 0 {

			shoes.Type = equipment["Type"].(string)
			shoes.Name = equipment["Name"].(string)
			shoes.Description = equipment["Description"].(string)
			shoes.Weight = int(equipment["Weight"].(float64))
			shoes.Value = int(equipment["Value"].(float64))
			shoes.SharpProtection = int(equipment["SharpProtection"].(float64))
			shoes.BluntProtection = int(equipment["BluntProtection"].(float64))
			shoes.Resistance = int(equipment["Resistance"].(float64))
			return shoes
		}
	}
	return shoes
}

func getRandomTorso() Torso {

	var equipmentFile = "mud/equipment/torsos.txt"
	var torso Torso

	if gostuff.IsFileExist(equipmentFile) == false {
		fmt.Println("File does not exist for ", equipmentFile)
		return torso
	}

	equipment, err := ioutil.ReadFile(equipmentFile)
	if err != nil {
		fmt.Println(err)
	}
	var data []map[string]interface{}

	err = json.Unmarshal(equipment, &data)
	if err != nil {
		fmt.Println(err)
	}

	torsoWeights := []int{16, 5, 3, 1}
	itemsInEachSection := TOTAL_TORSO / len(torsoWeights)
	sumOfWeights := 0
	for _, weight := range torsoWeights {
		sumOfWeights += (weight * itemsInEachSection)
	}

	randomTorsoWeights := rand.Intn(sumOfWeights)

	for index, equipment := range data {

		randomTorsoWeights -= torsoWeights[(index / itemsInEachSection)]

		if randomTorsoWeights <= 0 {

			torso.Type = equipment["Type"].(string)
			torso.Name = equipment["Name"].(string)
			torso.Description = equipment["Description"].(string)
			torso.Weight = int(equipment["Weight"].(float64))
			torso.Value = int(equipment["Value"].(float64))
			torso.SharpProtection = int(equipment["SharpProtection"].(float64))
			torso.BluntProtection = int(equipment["BluntProtection"].(float64))
			torso.Resistance = int(equipment["Resistance"].(float64))
			return torso
		}
	}
	return torso
}

func getRandomWeapon() Weapon {

	var equipmentFile = "mud/equipment/weapons.txt"
	var weapon Weapon

	if gostuff.IsFileExist(equipmentFile) == false {
		fmt.Println("File does not exist for ", equipmentFile)
		return weapon
	}

	equipment, err := ioutil.ReadFile(equipmentFile)
	if err != nil {
		fmt.Println(err)
	}
	var data []map[string]interface{}

	err = json.Unmarshal(equipment, &data)
	if err != nil {
		fmt.Println(err)
	}

	weaponWeights := []int{16, 5, 3, 1}
	itemsInEachSection := TOTAL_WEAPONS / len(weaponWeights)
	sumOfWeights := 0
	for _, weight := range weaponWeights {
		sumOfWeights += (weight * itemsInEachSection)
	}

	randomWeaponWeights := rand.Intn(sumOfWeights)

	for index, equipment := range data {

		randomWeaponWeights -= weaponWeights[(index / itemsInEachSection)]

		if randomWeaponWeights <= 0 {

			weapon.Type = equipment["Type"].(string)
			weapon.Name = equipment["Name"].(string)
			weapon.Description = equipment["Description"].(string)
			weapon.Weight = int(equipment["Weight"].(float64))
			weapon.Value = int(equipment["Value"].(float64))
			weapon.Strength = int(equipment["Strength"].(float64))
			weapon.Effect = equipment["Effect"].(string)
			return weapon
		}
	}
	return weapon
}
