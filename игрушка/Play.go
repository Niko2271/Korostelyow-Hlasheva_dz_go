package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type Character interface {
	Hit() int
	Block() int
	GetName() string
	GetHP() int
	GetStrength() int
	TakeDamage(damage int)
	IsAlive() bool
}

type Item struct {
	Type    string
	Attack  int
	Defence int
	PlusHP  int
	Name    string
}

type Player struct {
	Name      string
	HP        int
	MaxHP     int
	Strength  int
	hit       int
	block     int
	Inventory []Item
	Equipment []Item
}

func (p *Player) Hit() int {
	var hit int
	fmt.Print("Куда бьешь? 1-руки, 2-ноги, 3-голова: ")
	fmt.Scan(&hit)
	if hit < 1 || hit > 3 {
		hit = 1
	}
	return hit
}

func (p *Player) Block() int {
	var block int
	fmt.Print("Что защищаешь? 1-крылья(руки), 2-ноги, 3-корпус: ")
	fmt.Scan(&block)
	if block < 1 || block > 3 {
		block = 1
	}
	return block
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) GetHP() int {
	return p.HP
}

func (p *Player) GetStrength() int {
	bonus := 0
	for _, item := range p.Equipment {
		if item.Type == "оружие" {
			bonus += item.Attack
		}
	}
	return p.Strength + bonus
}

func (p *Player) TakeDamage(damage int) {
	defence := 0
	for _, item := range p.Equipment {
		if item.Type == "броня" {
			defence += item.Defence
		}
	}
	actualDamage := damage - defence
	if actualDamage < 1 {
		actualDamage = 1
	}
	p.HP -= actualDamage
	if p.HP < 0 {
		p.HP = 0
	}
}

func (p *Player) IsAlive() bool {
	return p.HP > 0
}

func (p *Player) AddItem(item Item) {
	p.Inventory = append(p.Inventory, item)
	fmt.Printf(" Получен предмет: %s\n", item.Name)
}

func (p *Player) ShowInventory() {
	fmt.Println("\n ИНВЕНТАРЬ:")
	if len(p.Inventory) == 0 {
		fmt.Println("   Пусто")
		return
	}
	for i, item := range p.Inventory {
		fmt.Printf("   %d. %s", i+1, item.Name)
		if item.Type == "оружие" {
			fmt.Printf(" (+%d атаки)", item.Attack)
		} else if item.Type == "броня" {
			fmt.Printf(" (+%d защиты)", item.Defence)
		} else if item.Type == "зелье" {
			fmt.Printf(" (+%d HP)", item.PlusHP)
		}
		fmt.Println()
	}
}

func (p *Player) ShowEquipment() {
	fmt.Println("\n ЭКИПИРОВКА:")
	if len(p.Equipment) == 0 {
		fmt.Println("   Нет надетых предметов")
		return
	}
	for _, item := range p.Equipment {
		fmt.Printf("   • %s", item.Name)
		if item.Type == "оружие" {
			fmt.Printf(" (+%d атаки)", item.Attack)
		} else if item.Type == "броня" {
			fmt.Printf(" (+%d защиты)", item.Defence)
		}
		fmt.Println()
	}
}

func (p *Player) EquipItem(itemNum int) {
	if itemNum < 1 || itemNum > len(p.Inventory) {
		fmt.Println(" Нет такого предмета")
		return
	}

	item := p.Inventory[itemNum-1]
	for _, equipped := range p.Equipment {
		if equipped.Type == item.Type {
			fmt.Printf(" Уже надет %s (%s)\n", equipped.Name, equipped.Type)
			return
		}
	}

	p.Equipment = append(p.Equipment, item)
	p.Inventory = append(p.Inventory[:itemNum-1], p.Inventory[itemNum:]...)
	fmt.Printf("Надето: %s\n", item.Name)

	if item.Type == "зелье" && item.PlusHP > 0 {
		p.HP += item.PlusHP
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
		fmt.Printf("Восстановлено %d HP (теперь %d/%d)\n", item.PlusHP, p.HP, p.MaxHP)
	}
}

type Enemy struct {
	Name     string
	HP       int
	Strength int
	hit      int
	block    int
	Item     Item
}

func (e *Enemy) Hit() int {
	return rand.Intn(3) + 1
}

func (e *Enemy) Block() int {
	return rand.Intn(3) + 1
}

func (e *Enemy) GetName() string {
	return e.Name
}

func (e *Enemy) GetHP() int {
	return e.HP
}

func (e *Enemy) GetStrength() int {
	return e.Strength
}

func (e *Enemy) TakeDamage(damage int) {
	e.HP -= damage
	if e.HP < 0 {
		e.HP = 0
	}
}

func (e *Enemy) IsAlive() bool {
	return e.HP > 0
}

func FightPvE(player *Player, enemy *Enemy) bool {
	fmt.Printf("\n БОЙ: %s против %s \n", player.GetName(), enemy.GetName())
	fmt.Println("════════════════════════════════════════════")

	round := 1

	for player.IsAlive() && enemy.IsAlive() {
		fmt.Printf("\n—— Раунд %d ——\n", round)
		fmt.Printf("%s: %d/%d HP\n", player.GetName(), player.GetHP(), player.MaxHP)
		fmt.Printf("%s: %d HP\n", enemy.GetName(), enemy.GetHP())

		playerBlock := player.Block()
		playerHit := player.Hit()
		enemyHit := enemy.Hit()
		enemyBlock := enemy.Block()

		fmt.Printf("\n%s защищает %d, бьёт в %d\n", player.GetName(), playerBlock, playerHit)
		fmt.Printf("%s защищает %d, бьёт в %d\n", enemy.GetName(), enemyBlock, enemyHit)

		if playerHit != enemyBlock {
			damage := player.GetStrength()
			enemy.TakeDamage(damage)
			fmt.Printf(" %s попадает! Нанесено %d урона\n", player.GetName(), damage)
			fmt.Printf("   У %s осталось %d HP\n", enemy.GetName(), enemy.GetHP())
		} else {
			fmt.Printf(" %s блокирует удар!\n", enemy.GetName())
		}

		if enemy.IsAlive() && enemyHit != playerBlock {
			damage := enemy.GetStrength()
			player.TakeDamage(damage)
			fmt.Printf(" %s попадает! Нанесено %d урона\n", enemy.GetName(), damage)
			fmt.Printf("   У %s осталось %d HP\n", player.GetName(), player.GetHP())
		} else if enemy.IsAlive() {
			fmt.Printf("  %s блокирует удар!\n", player.GetName())
		}

		round++
	}

	fmt.Println("\n════════════════════════════════════════════")
	if player.IsAlive() {
		fmt.Printf(" ПОБЕДА! %s побеждает!\n", player.GetName())
		if enemy.Item.Name != "" {
			player.AddItem(enemy.Item)
		}
		return true
	} else {
		fmt.Printf(" ПОРАЖЕНИЕ! %s побеждает!\n", enemy.GetName())
		return false
	}
}

func FightPvP(players [2]*Player) {
	fmt.Println("\n════════════════════════════════════════════")
	fmt.Println("           РЕЖИМ PvP - ГОРЯЧИЙ СТУЛ")
	fmt.Println("════════════════════════════════════════════")
	fmt.Println("📜 ПРАВИЛА:")
	fmt.Println("1. Когда ходит соперник - ОТВЕРНИТЕСЬ!")
	fmt.Println("2. Вводите выбор, когда подойдёт очередь")
	fmt.Println("3. Не подсматривайте!")
	fmt.Println("════════════════════════════════════════════")

	fmt.Println("\nНажмите Enter чтобы начать...")
	fmt.Scanln()

	round := 1
	for players[0].IsAlive() && players[1].IsAlive() {
		var choices [2]struct {
			block int
			hit   int
		}

		fmt.Printf("\n=== РАУНД %d ===\n", round)
		fmt.Printf("\n %s, ваш ход (%s отвернись!)\n", players[0].GetName(), players[1].GetName())
		fmt.Println("Нажмите Enter когда готовы...")
		fmt.Scanln()

		choices[0].block = players[0].Block()
		choices[0].hit = players[0].Hit()

		fmt.Printf("\n %s, ваш ход (%s отвернись!)\n", players[1].GetName(), players[0].GetName())
		fmt.Println("Нажмите Enter когда готовы...")
		fmt.Scanln()

		choices[1].block = players[1].Block()
		choices[1].hit = players[1].Hit()

		fmt.Println("\n=== РЕЗУЛЬТАТЫ РАУНДА ===")
		fmt.Printf("%s защищает %d, бьёт в %d\n",
			players[0].GetName(), choices[0].block, choices[0].hit)
		fmt.Printf("%s защищает %d, бьёт в %d\n",
			players[1].GetName(), choices[1].block, choices[1].hit)

		if choices[0].hit != choices[1].block {
			damage := players[0].GetStrength()
			players[1].TakeDamage(damage)
			fmt.Printf(" %s попадает! Нанесено %d урона\n", players[0].GetName(), damage)
		} else {
			fmt.Printf(" %s блокирует удар!\n", players[1].GetName())
		}

		if choices[1].hit != choices[0].block {
			damage := players[1].GetStrength()
			players[0].TakeDamage(damage)
			fmt.Printf(" %s попадает! Нанесено %d урона\n", players[1].GetName(), damage)
		} else {
			fmt.Printf("  %s блокирует удар!\n", players[0].GetName())
		}

		fmt.Printf("\n%s: %d/%d HP\n", players[0].GetName(), players[0].GetHP(), players[0].MaxHP)
		fmt.Printf("%s: %d/%d HP\n", players[1].GetName(), players[1].GetHP(), players[1].MaxHP)

		round++
	}

	fmt.Println("\n════════════════════════════════════════════")
	if players[0].IsAlive() {
		fmt.Printf("ПОБЕДИТЕЛЬ: %s!\n", players[0].GetName())
	} else {
		fmt.Printf("ПОБЕДИТЕЛЬ: %s!\n", players[1].GetName())
	}
}

func manageBetweenBattles(player *Player) {
	for {
		fmt.Println("\n════════════════════════════════════════════")
		fmt.Println("         МЕНЮ")
		fmt.Println("════════════════════════════════════════════")
		fmt.Println("1. Продолжить игру")
		fmt.Println("2. Посмотреть инвентарь")
		fmt.Println("3. Посмотреть экипировку")
		fmt.Println("4. Надеть предмет из инвентаря")
		fmt.Println("════════════════════════════════════════════")
		fmt.Print("Выбор: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			return
		case 2:
			player.ShowInventory()
		case 3:
			player.ShowEquipment()
		case 4:
			player.ShowInventory()
			if len(player.Inventory) > 0 {
				fmt.Print("Номер предмета для экипировки: ")
				var itemChoice int
				fmt.Scan(&itemChoice)
				player.EquipItem(itemChoice)
			} else {
				fmt.Println("Инвентарь пуст")
			}
		default:
			fmt.Println("Неверный выбор")
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println("                     ПРЕДИСЛОВИЕ")
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println("В древние времена, на окраине цветущего города Аркании,")
	fmt.Println("жил старый маг Эльдриан. Он искал ученика для передачи знаний.")
	fmt.Println()
	fmt.Println("Долго бродил он, пока не встретил Лилиан — девочку с добрым сердцем.")
	fmt.Println("Годы шли, она росла, впитывала магию. Когда маг умер, он оставил")
	fmt.Println("ей древний амулет — Сердце Дракона.")
	fmt.Println()
	fmt.Println("Лилиан помогала горожанам, но мэр Гаррик, полный зависти,")
	fmt.Println("обвинил её в колдовстве и поджёг дом.")
	fmt.Println()
	fmt.Println("В огне Лилиан произнесла заклинание... и превратилась в дракона.")
	fmt.Println("Теперь вы — дракон-хранитель. Защитите Арканию!")
	fmt.Println("══════════════════════════════════════════════════════════")

	player := &Player{
		Name:     "Дракон-Лилиан",
		HP:       120,
		MaxHP:    120,
		Strength: 15}

	player.Equipment = append(player.Equipment, Item{
		Name:   "Когти дракона",
		Type:   "оружие",
		Attack: 5,
	})

	enemies := []*Enemy{
		{
			Name:     "Стражник Гаррика",
			HP:       50,
			Strength: 8,
			Item: Item{
				Name:    "Щит стражника",
				Type:    "броня",
				Defence: 3,
			},
		},
		{
			Name:     "Капитан стражи",
			HP:       75,
			Strength: 12,
			Item: Item{
				Name:   "Меч капитана",
				Type:   "оружие",
				Attack: 8,
			},
		},
		{
			Name:     "Мэр Гаррик",
			HP:       95,
			Strength: 15,
			Item: Item{
				Name:   "Амулет прощения",
				Type:   "зелье",
				PlusHP: 50,
			},
		},
	}
	victory := true

	for i, enemy := range enemies {
		switch i {
		case 0:
			fmt.Println("\n══════════════════════════════════════════════════════════")
			fmt.Println("                     ГЛАВА 1: ПРОБУЖДЕНИЕ")
			fmt.Println("══════════════════════════════════════════════════════════")
			fmt.Println("Вы открываете глаза. Вокруг — дым и пепел.")
			fmt.Println("Ваш дом, ваша жизнь — всё в огне.")
			fmt.Println("Сердце Дракона на груди пульсирует тёплым светом.")
			fmt.Println("Пришло время защитить то, что дорого.")
			fmt.Println("══════════════════════════════════════════════════════════")

		case 1:
			fmt.Println("\n══════════════════════════════════════════════════════════")
			fmt.Println("                ГЛАВА 2: УЛИЦЫ АРКАНИИ")
			fmt.Println("══════════════════════════════════════════════════════════")
			fmt.Println("Город в панике. Люди бегут, не понимая, что происходит.")
			fmt.Println("Гаррик собрал стражу — они ищут вас.")
			fmt.Println("«Дракон должен умереть!» — кричит он.")
			fmt.Println("Но вы знаете правду. И будете сражаться за неё.")
			fmt.Println("══════════════════════════════════════════════════════════")

		case 2:
			fmt.Println("\n══════════════════════════════════════════════════════════")
			fmt.Println("                  ГЛАВА 3: У ПОДНОЖИЯ ЗАМКА")
			fmt.Println("══════════════════════════════════════════════════════════")
			fmt.Println("Перед вами — замок Гаррика. Его личная гвардия ждёт.")
			fmt.Println("Это последнее препятствие на пути к правде.")
			fmt.Println("Силы на исходе, но сдаваться нельзя.")
			fmt.Println("Судьба Аркании решается здесь и сейчас.")
			fmt.Println("══════════════════════════════════════════════════════════")
		}

		if !FightPvE(player, enemy) {
			victory = false
			break
		}

		if i < len(enemies)-1 && player.IsAlive() {
			manageBetweenBattles(player)
		}
	}

	fmt.Println("\n══════════════════════════════════════════════════════════")
	fmt.Println("                        ЭПИЛОГ")
	fmt.Println("══════════════════════════════════════════════════════════")

	if victory {
		fmt.Println("\nГаррик повержен. Правда восторжествовала.")
		fmt.Println("Лилиан вернула человеческий облик, но сила дракона")
		fmt.Println("осталась с ней. Она стала хранительницей Аркании,")
		fmt.Println("и город зажил в мире и процветании.")
		fmt.Println("\nИГРА ОКОНЧЕНА. ВЫ ПОБЕДИЛИ! ")
	} else {
		fmt.Println("\nАркания пала. Огонь поглотил город.")
		fmt.Println("Легенда о драконе-хранителе стала предостережением")
		fmt.Println("для будущих поколений.")
		fmt.Println("\n ИГРА ОКОНЧЕНА. ВЫ ПРОИГРАЛИ.")
	}
	fmt.Println("══════════════════════════════════════════════════════════")

	fmt.Println("\n══════════════════════════════════════════════════════════")
	fmt.Println("           ХОТИТЕ ПОПРОБОВАТЬ РЕЖИМ PvP?")
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println("1. Да")
	fmt.Println("2. Нет (выйти из игры)")
	fmt.Print("Выбор: ")

	var wantPvP int
	fmt.Scan(&wantPvP)

	if wantPvP == 1 {
		fmt.Println("\n════════════════════════════════════════════")
		fmt.Println("           ВЫБЕРИТЕ РЕЖИМ PvP")
		fmt.Println("════════════════════════════════════════════")
		fmt.Println("1. PvP локально (на одном компьютере)")
		fmt.Println("2. СОЗДАТЬ игру по сети (сервер)")
		fmt.Println("3. ПРИСОЕДИНИТЬСЯ к игре по сети (клиент)")
		fmt.Print("Выбор: ")

		var mode int
		fmt.Scan(&mode)

		if mode == 1 {
			player1 := &Player{
				Name:     "Игрок 1",
				HP:       100,
				MaxHP:    100,
				Strength: 10,
			}
			player2 := &Player{
				Name:     "Игрок 2",
				HP:       100,
				MaxHP:    100,
				Strength: 10,
			}
			player1.Equipment = append(player1.Equipment, Item{
				Name:   "Меч",
				Type:   "оружие",
				Attack: 5,
			})
			player2.Equipment = append(player2.Equipment, Item{
				Name:   "Топор",
				Type:   "оружие",
				Attack: 7,
			})
			FightPvP([2]*Player{player1, player2})

		} else if mode == 2 {
			StartServer()

		} else if mode == 3 {
			StartClient()
		}
	} else {
		fmt.Println("\nСпасибо за игру!")
	}
}

func StartServer() {
	ln, _ := net.Listen("tcp", ":8080")
	defer ln.Close()
	fmt.Println("Сервер запущен, ждём клиента...")

	conn, _ := ln.Accept()
	defer conn.Close()
	player1 := &Player{Name: "Игрок 1 (хост)", HP: 100, MaxHP: 100, Strength: 10}
	player2 := &Player{Name: "Игрок 2", HP: 100, MaxHP: 100, Strength: 10}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var hit2, block2 int

	for player1.IsAlive() && player2.IsAlive() {
		fmt.Println("\n=== ВАШ ХОД (игрок 1) ===")
		block1 := player1.Block()
		hit1 := player1.Hit()

		writer.WriteString(fmt.Sprintf("%d %d\n", hit1, block1))
		writer.Flush()

		data, _ := reader.ReadString('\n')
		fmt.Sscanf(data, "%d %d", &hit2, &block2)

		fmt.Println("\n--- РЕЗУЛЬТАТЫ РАУНДА ---")

		if hit1 != block2 {
			damage := player1.GetStrength()
			player2.TakeDamage(damage)
			fmt.Printf("%s попадает! Урон: %d\n", player1.Name, damage)
		} else {
			fmt.Printf("%s блокирует!\n", player2.Name)
		}

		if hit2 != block1 {
			damage := player2.GetStrength()
			player1.TakeDamage(damage)
			fmt.Printf("%s попадает! Урон: %d\n", player2.Name, damage)
		} else {
			fmt.Printf("%s блокирует!\n", player1.Name)
		}
		fmt.Printf("%s: %d/%d HP\n", player1.Name, player1.HP, player1.MaxHP)
		fmt.Printf("%s: %d/%d HP\n", player2.Name, player2.HP, player2.MaxHP)
	}
}

func StartClient() {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()

	player := &Player{Name: "Игрок 2", HP: 100, MaxHP: 100, Strength: 10}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for player.IsAlive() {
		data, _ := reader.ReadString('\n')
		var enemyHit, enemyBlock int
		fmt.Sscanf(data, "%d %d", &enemyHit, &enemyBlock)

		fmt.Println("\n=== ВАШ ХОД (игрок 2) ===")
		block := player.Block()
		hit := player.Hit()

		writer.WriteString(fmt.Sprintf("%d %d\n", hit, block))
		writer.Flush()

	}
}
