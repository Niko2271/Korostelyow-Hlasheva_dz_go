package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strings"
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
	fmt.Println("ПРАВИЛА:")
	fmt.Println("1. Когда ходит соперник - ОТВЕРНИТЕСЬ!")
	fmt.Println("2. Вводите выбор, когда подойдёт очередь")
	fmt.Println("3. Не подсматривайте!")
	fmt.Println("4. ЧАТ: введите 'чат сообщение' для общения")
	fmt.Println("════════════════════════════════════════════")

	fmt.Println("\nНажмите Enter чтобы начать...")
	fmt.Scanln()

	// Канал для сообщений чата
	chatMessages := make(chan string, 10)
	stopChat := make(chan bool)
	
	// Запускаем горутину для отображения чата
	go func() {
		for {
			select {
			case msg := <-chatMessages:
				fmt.Printf("\n[ЧАТ] %s\n", msg)
				fmt.Print("> ")
			case <-stopChat:
				return
			}
		}
	}()

	round := 1
	for players[0].IsAlive() && players[1].IsAlive() {
		var choices [2]struct {
			block int
			hit   int
		}

		fmt.Printf("\n=== РАУНД %d ===\n", round)
		
		// Ход первого игрока
		fmt.Printf("\n %s, ваш ход (%s отвернись!)\n", players[0].GetName(), players[1].GetName())
		fmt.Println("(Введите 'чат сообщение' для отправки в чат)")
		fmt.Println("Нажмите Enter когда готовы...")
		
		// Проверяем, не хочет ли игрок написать в чат
		var input string
		fmt.Scanln(&input)
		
		if strings.HasPrefix(input, "чат") {
			msg := strings.TrimSpace(strings.TrimPrefix(input, "чат"))
			if msg != "" {
				chatMessages <- fmt.Sprintf("%s: %s", players[0].GetName(), msg)
			}
			fmt.Println("Нажмите Enter чтобы продолжить...")
			fmt.Scanln()
		}

		choices[0].block = players[0].Block()
		choices[0].hit = players[0].Hit()

		// Ход второго игрока
		fmt.Printf("\n %s, ваш ход (%s отвернись!)\n", players[1].GetName(), players[0].GetName())
		fmt.Println("(Введите 'чат сообщение' для отправки в чат)")
		fmt.Println("Нажмите Enter когда готовы...")
		
		fmt.Scanln(&input)
		
		if strings.HasPrefix(input, "чат") {
			msg := strings.TrimSpace(strings.TrimPrefix(input, "чат"))
			if msg != "" {
				chatMessages <- fmt.Sprintf("%s: %s", players[1].GetName(), msg)
			}
			fmt.Println("Нажмите Enter чтобы продолжить...")
			fmt.Scanln()
		}

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

	stopChat <- true

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

func showPrologue() {
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
}

func playStoryMode() {
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
}

func playPVPMode() {
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
}

func showMainMenu() int {
	fmt.Println("\n══════════════════════════════════════════════════════════")
	fmt.Println("                    ГЛАВНОЕ МЕНЮ")
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println("1. Начать сюжетную кампанию")
	fmt.Println("2. Пропустить сюжет и начать игру")
	fmt.Println("3. Режим PvP")
	fmt.Println("4. Выйти из игры")
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Print("Выбор: ")

	var choice int
	fmt.Scan(&choice)
	return choice
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		choice := showMainMenu()

		switch choice {
		case 1:
			showPrologue()
			playStoryMode()
		case 2:
			fmt.Println("\nПропускаем сюжет... Начинаем игру!")
			playStoryMode()
		case 3:
			playPVPMode()
		case 4:
			fmt.Println("\nСпасибо за игру! До свидания!")
			return
		default:
			fmt.Println("\nНеверный выбор. Попробуйте снова.")
			continue
		}

		fmt.Println("\n══════════════════════════════════════════════════════════")
		fmt.Println("Хотите вернуться в главное меню?")
		fmt.Println("1. Да")
		fmt.Println("2. Нет (выйти из игры)")
		fmt.Print("Выбор: ")

		var backToMenu int
		fmt.Scan(&backToMenu)

		if backToMenu != 1 {
			fmt.Println("\nСпасибо за игру! До свидания!")
			break
		}
	}
}

func StartServer() {
	ln, _ := net.Listen("tcp", ":8080")
	defer ln.Close()
	fmt.Println("Сервер запущен, ждём клиента...")
	fmt.Println("Чат доступен! Используйте 'чат сообщение' для общения")

	conn, _ := ln.Accept()
	defer conn.Close()
	
	player1 := &Player{Name: "Игрок 1 (хост)", HP: 100, MaxHP: 100, Strength: 10}
	player2 := &Player{Name: "Игрок 2", HP: 100, MaxHP: 100, Strength: 10}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Канал для сообщений чата
	chatMessages := make(chan string, 10)
	stopChat := make(chan bool)
	
	// Запускаем горутину для отображения чата
	go func() {
		for {
			select {
			case msg := <-chatMessages:
				fmt.Printf("\n[ЧАТ] %s\n", msg)
				fmt.Print("> ")
			case <-stopChat:
				return
			}
		}
	}()
	
	// Запускаем горутину для получения сообщений чата от клиента
	go func() {
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if strings.HasPrefix(msg, "ЧАТ:") {
				chatMsg := strings.TrimPrefix(msg, "ЧАТ:")
				chatMessages <- strings.TrimSpace(chatMsg)
			}
		}
	}()

	var hit2, block2 int
	round := 1

	for player1.IsAlive() && player2.IsAlive() {
		fmt.Printf("\n=== РАУНД %d ===\n", round)
		fmt.Println("\n=== ВАШ ХОД (игрок 1) ===")
		fmt.Println("(Для чата введите: чат сообщение)")
		
		// Проверяем, не хочет ли игрок написать в чат
		var input string
		fmt.Scanln(&input)
		
		if strings.HasPrefix(input, "чат") {
			msg := strings.TrimSpace(strings.TrimPrefix(input, "чат"))
			if msg != "" {
				chatMsg := fmt.Sprintf("%s: %s", player1.Name, msg)
				chatMessages <- chatMsg
				writer.WriteString("ЧАТ:" + chatMsg + "\n")
				writer.Flush()
			}
			fmt.Println("Нажмите Enter чтобы продолжить...")
			fmt.Scanln()
		}
		
		block1 := player1.Block()
		hit1 := player1.Hit()

		writer.WriteString(fmt.Sprintf("%d %d\n", hit1, block1))
		writer.Flush()

		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Соединение разорвано")
			break
		}

		// Проверяем, не сообщение ли это чата
		if strings.HasPrefix(data, "ЧАТ:") {
			continue
		}

		fmt.Sscanf(data, "%d %d", &hit2, &block2)

		fmt.Println("\n--- РЕЗУЛЬТАТЫ РАУНДА ---")
		fmt.Printf("%s защищает %d, бьёт в %d\n", player1.Name, block1, hit1)
		fmt.Printf("%s защищает %d, бьёт в %d\n", player2.Name, block2, hit2)

		damage1 := 0
		damage2 := 0
		
		if hit1 != block2 {
			damage1 = player1.GetStrength()
			player2.TakeDamage(damage1)
			fmt.Printf("%s попадает! Урон: %d\n", player1.Name, damage1)
		} else {
			fmt.Printf("%s блокирует удар %s!\n", player2.Name, player1.Name)
		}

		if hit2 != block1 {
			damage2 = player2.GetStrength()
			player1.TakeDamage(damage2)
			fmt.Printf("%s попадает! Урон: %d\n", player2.Name, damage2)
		} else {
			fmt.Printf("%s блокирует удар %s!\n", player1.Name, player2.Name)
		}

		resultMsg := fmt.Sprintf("%d %d %d %d %d %d\n", 
			player1.HP, player2.HP, hit1, block1, damage1, damage2)
		writer.WriteString(resultMsg)
		writer.Flush()

		fmt.Printf("\n%s: %d/%d HP\n", player1.Name, player1.HP, player1.MaxHP)
		fmt.Printf("%s: %d/%d HP\n", player2.Name, player2.HP, player2.MaxHP)
		
		round++
	}

	stopChat <- true

	var winner string
	if player1.IsAlive() {
		winner = player1.Name
	} else {
		winner = player2.Name
	}
	
	fmt.Println("\n════════════════════════════════════════════")
	fmt.Printf("ПОБЕДИТЕЛЬ: %s!\n", winner)
	fmt.Println("════════════════════════════════════════════")
	
	writer.WriteString("GAME_OVER " + winner + "\n")
	writer.Flush()
}

func StartClient() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		fmt.Println("Убедитесь что сервер запущен")
		return
	}
	defer conn.Close()

	player := &Player{Name: "Игрок 2", HP: 100, MaxHP: 100, Strength: 10}
	enemyName := "Игрок 1 (хост)"

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Канал для сообщений чата
	chatMessages := make(chan string, 10)
	stopChat := make(chan bool)
	
	// Запускаем горутину для отображения чата
	go func() {
		for {
			select {
			case msg := <-chatMessages:
				fmt.Printf("\n[ЧАТ] %s\n", msg)
				fmt.Print("> ")
			case <-stopChat:
				return
			}
		}
	}()
	
	// Запускаем горутину для получения сообщений чата от сервера
	go func() {
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if strings.HasPrefix(msg, "ЧАТ:") {
				chatMsg := strings.TrimPrefix(msg, "ЧАТ:")
				chatMessages <- strings.TrimSpace(chatMsg)
			}
		}
	}()

	fmt.Printf("\nПодключились к серверу! Вы - %s\n", player.Name)
	fmt.Println("Чат доступен! Используйте 'чат сообщение' для общения")
	fmt.Println("Ожидаем начала игры...")

	for player.IsAlive() {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Соединение разорвано")
			break
		}

		if len(data) >= 9 && data[:9] == "GAME_OVER" {
			var winner string
			fmt.Sscanf(data, "GAME_OVER %s", &winner)
			fmt.Println("\n════════════════════════════════════════════")
			fmt.Printf("ПОБЕДИТЕЛЬ: %s!\n", winner)
			fmt.Println("════════════════════════════════════════════")
			break
		}

		// Проверяем, не сообщение ли это чата
		if strings.HasPrefix(data, "ЧАТ:") {
			continue
		}

		var enemyHit, enemyBlock int
		fmt.Sscanf(data, "%d %d", &enemyHit, &enemyBlock)

		fmt.Println("\n=== ВАШ ХОД ===")
		fmt.Println("(Для чата введите: чат сообщение)")
		
		// Проверяем, не хочет ли игрок написать в чат
		var input string
		fmt.Scanln(&input)
		
		if strings.HasPrefix(input, "чат") {
			msg := strings.TrimSpace(strings.TrimPrefix(input, "чат"))
			if msg != "" {
				chatMsg := fmt.Sprintf("%s: %s", player.Name, msg)
				chatMessages <- chatMsg
				writer.WriteString("ЧАТ:" + chatMsg + "\n")
				writer.Flush()
			}
			fmt.Println("Нажмите Enter чтобы продолжить...")
			fmt.Scanln()
		}
		
		block := player.Block()
		hit := player.Hit()

		writer.WriteString(fmt.Sprintf("%d %d\n", hit, block))
		writer.Flush()

		resultData, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Соединение разорвано")
			break
		}

		var playerHP, enemyHP, enemyHit2, enemyBlock2, damageToEnemy, damageToPlayer int
		fmt.Sscanf(resultData, "%d %d %d %d %d %d", 
			&playerHP, &enemyHP, &enemyHit2, &enemyBlock2, &damageToEnemy, &damageToPlayer)

		player.HP = playerHP

		fmt.Println("\n--- РЕЗУЛЬТАТЫ РАУНДА ---")
		fmt.Printf("%s защищает %d, бьёт в %d\n", enemyName, enemyBlock, enemyHit)
		fmt.Printf("%s защищает %d, бьёт в %d\n", player.Name, block, hit)

		if damageToEnemy > 0 {
			fmt.Printf("%s попадает! Урон: %d\n", player.Name, damageToEnemy)
		} else {
			fmt.Printf("%s блокирует удар %s!\n", enemyName, player.Name)
		}

		if damageToPlayer > 0 {
			fmt.Printf("%s попадает! Урон: %d\n", enemyName, damageToPlayer)
		} else {
			fmt.Printf("%s блокирует удар %s!\n", player.Name, enemyName)
		}

		fmt.Printf("\n%s: %d/%d HP\n", player.Name, player.HP, player.MaxHP)
		fmt.Printf("%s: %d HP\n", enemyName, enemyHP)
	}

	stopChat <- true
}
