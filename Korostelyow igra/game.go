package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

type Item struct {
	Name    string
	Attack  int
	Defence int
	Heal    int
}

type Player struct {
	Name      string
	HP        int
	MaxHP     int
	Strength  int
	Inventory []Item
	Weapon    Item
	Armor     Item
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:      name,
		HP:        100,
		MaxHP:     100,
		Strength:  10,
		Inventory: []Item{},
	}
}

func (p *Player) Attack() int {
	dmg := p.Strength + rand.Intn(10)
	if p.Weapon.Name != "" {
		dmg += p.Weapon.Attack
	}
	return dmg
}

func (p *Player) Defence() int {
	if p.Armor.Name != "" {
		return p.Armor.Defence
	}
	return 0
}

func (p *Player) TakeDamage(dmg int) {
	def := p.Defence()
	actual := dmg - def
	if actual < 0 {
		actual = 0
	}
	p.HP -= actual
	if p.HP < 0 {
		p.HP = 0
	}
	fmt.Printf("%s получает %d урона", p.Name, actual)
	if def > 0 {
		fmt.Printf(" (броня -%d)", def)
	}
	fmt.Printf(" ❤️ %d/%d\n", p.HP, p.MaxHP)
}

func (p *Player) ShowInventory() {
	fmt.Println("\n📦 ИНВЕНТАРЬ")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
	if len(p.Inventory) == 0 {
		fmt.Println("  Пусто")
		return
	}
	for i, item := range p.Inventory {
		fmt.Printf("%d. %s", i, item.Name)
		if item.Attack > 0 {
			fmt.Printf(" ⚔️ +%d", item.Attack)
		}
		if item.Defence > 0 {
			fmt.Printf(" 🛡️ +%d", item.Defence)
		}
		if item.Heal > 0 {
			fmt.Printf(" 💊 +%d HP", item.Heal)
		}
		fmt.Println()
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
}

func (p *Player) ShowEquipment() {
	fmt.Println("\n⚔️ ЭКИПИРОВКА")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
	if p.Weapon.Name != "" {
		fmt.Printf("  Оружие: %s (атака +%d)\n", p.Weapon.Name, p.Weapon.Attack)
	} else {
		fmt.Println("  Оружие: не надето")
	}
	if p.Armor.Name != "" {
		fmt.Printf("  Броня: %s (защита +%d)\n", p.Armor.Name, p.Armor.Defence)
	} else {
		fmt.Println("  Броня: не надета")
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
}

func (p *Player) UseItem(index int) {
	if index < 0 || index >= len(p.Inventory) {
		fmt.Println("❌ Неверный индекс")
		return
	}
	item := p.Inventory[index]

	if item.Heal > 0 {
		p.HP += item.Heal
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
		fmt.Printf("💊 Использован %s, восстановлено %d HP\n", item.Name, item.Heal)
		p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
	} else if item.Attack > 0 {
		if p.Weapon.Name != "" {
			p.Inventory = append(p.Inventory, p.Weapon)
			fmt.Printf("🔁 Снято: %s\n", p.Weapon.Name)
		}
		p.Weapon = item
		fmt.Printf("⚔️ Надето: %s (атака +%d)\n", item.Name, item.Attack)
		p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
	} else if item.Defence > 0 {
		if p.Armor.Name != "" {
			p.Inventory = append(p.Inventory, p.Armor)
			fmt.Printf("🔁 Снято: %s\n", p.Armor.Name)
		}
		p.Armor = item
		fmt.Printf("🛡️ Надета: %s (защита +%d)\n", item.Name, item.Defence)
		p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
	}
}

func RandomItem() Item {
	items := []Item{
		{"M41A Pulse Rifle", 8, 0, 0},
		{"M240 Flamethrower", 12, 0, 0},
		{"Combat Knife", 5, 0, 0},
		{"Smartgun", 10, 0, 0},
		{"USCM Armor", 0, 5, 0},
		{"Tactical Vest", 0, 3, 0},
		{"Medkit", 0, 0, 30},
		{"Stimpack", 0, 0, 20},
	}
	return items[rand.Intn(len(items))]
}

type Enemy struct {
	Name     string
	HP       int
	Strength int
	Item     Item
}

func NewEnemy() *Enemy {
	names := []string{"Facehugger", "Drone Xenomorph", "Warrior Xenomorph", "Praetorian"}
	return &Enemy{
		Name:     names[rand.Intn(len(names))],
		HP:       40 + rand.Intn(60),
		Strength: 8 + rand.Intn(12),
		Item:     RandomItem(),
	}
}

func (e *Enemy) Attack() int {
	return e.Strength + rand.Intn(10)
}

func (e *Enemy) TakeDamage(dmg int) {
	e.HP -= dmg
	if e.HP < 0 {
		e.HP = 0
	}
	fmt.Printf("👾 %s получает %d урона! ❤️ %d\n", e.Name, dmg, e.HP)
}

func FightPVE(p *Player, e *Enemy) bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\n☠️ БОЙ: %s VS %s ☠️\n", p.Name, e.Name)

	for p.HP > 0 && e.HP > 0 {
		fmt.Printf("\n❤️ %s: %d/%d | 👾 %s: %d\n", p.Name, p.HP, p.MaxHP, e.Name, e.HP)
		fmt.Println("1. Атаковать")
		fmt.Println("2. Инвентарь")
		fmt.Println("3. Экипировка")
		fmt.Print("Выбор: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "2":
			p.ShowInventory()
			fmt.Print("Использовать предмет (индекс, -1 назад): ")
			scanner.Scan()
			var idx int
			fmt.Sscan(scanner.Text(), &idx)
			if idx >= 0 {
				p.UseItem(idx)
			}
			continue
		case "3":
			p.ShowEquipment()
			fmt.Print("Снять предмет (0-оружие,1-броня): ")
			scanner.Scan()
			var idx int
			fmt.Sscan(scanner.Text(), &idx)
			if idx == 0 && p.Weapon.Name != "" {
				p.Inventory = append(p.Inventory, p.Weapon)
				p.Weapon = Item{}
				fmt.Println("Оружие снято")
			} else if idx == 1 && p.Armor.Name != "" {
				p.Inventory = append(p.Inventory, p.Armor)
				p.Armor = Item{}
				fmt.Println("Броня снята")
			}
			continue
		}

		dmg := p.Attack()
		e.TakeDamage(dmg)
		fmt.Printf("🔫 Урон %d\n", dmg)

		if e.HP <= 0 {
			fmt.Printf("\n🏆 Трофей: %s\n", e.Item.Name)
			p.Inventory = append(p.Inventory, e.Item)
			return true
		}

		dmg = e.Attack()
		p.TakeDamage(dmg)
		fmt.Printf("👾 %s нанес %d урона\n", e.Name, dmg)
	}
	return false
}

func Campaign() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🎬 ПРОЛОГ: LV-426")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Корабль 'Прометей' совершил аварийную посадку на спутнике LV-426.")
	fmt.Println("Сигнал бедствия привел вас на заброшенную колонию Вэйланд-Ютани.")
	fmt.Println("Внутри - тишина. Стены покрыты слизью. Вентиляция шипит.")
	fmt.Println("Вы - элитный боец USCM. Ваша задача - найти выживших...")
	fmt.Println("\n📟 СООБЩЕНИЕ ОТ КОМАНДОВАНИЯ:")
	fmt.Println("   'В колонии обнаружены неизвестные формы жизни. Будьте осторожны.'")
	fmt.Println("   'Докладывайте о любом контакте. Удачи, морпех.'")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Print("Нажмите Enter, чтобы продолжить...")
	scanner.Scan()

	p := NewPlayer("Рипли")
	p.Inventory = append(p.Inventory, RandomItem())
	p.Inventory = append(p.Inventory, RandomItem())
	fmt.Println("\n📦 Вы получили стартовое снаряжение!")
	p.ShowInventory()
	fmt.Print("Нажмите Enter...")
	scanner.Scan()

	// Бой 1
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🔍 ИССЛЕДОВАНИЕ КОЛОНИИ")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Вы идете по темному коридору. Датчик движения пищит.")
	fmt.Println("Из вентиляции выпрыгивает ЧУЖОЙ ДРОН!")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Print("Нажмите Enter...")
	scanner.Scan()

	e1 := NewEnemy()
	if !FightPVE(p, e1) {
		fmt.Println("\n💀 Игра окончена... Колония потеряна.")
		return
	}

	p.HP = p.MaxHP
	fmt.Println("\n💊 Медики восстановили ваше здоровье. Вы готовы к следующей миссии.")
	fmt.Print("Нажмите Enter...")
	scanner.Scan()

	// Бой 2
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🏭 ВХОД В РЕАКТОРНУЮ")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Вы спускаетесь в реакторный отсек. Воздух тяжелый, висит туман.")
	fmt.Println("Слышен звук шагов. Из-за угла выходит КСЕНОМОРФ-ВОИН!")
	fmt.Println("Он крупнее и быстрее обычного. Готовьтесь!")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Print("Нажмите Enter...")
	scanner.Scan()

	e2 := NewEnemy()
	if !FightPVE(p, e2) {
		fmt.Println("\n💀 Вы погибли в реакторной...")
		return
	}

	p.HP = p.MaxHP
	fmt.Println("\n💊 Эвакуационный модуль восстановил силы.")
	fmt.Print("Нажмите Enter...")
	scanner.Scan()

	// Бой 3 - Королева
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("👑 ЛОГОВО КОРОЛЕВЫ")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Вы нашли главное гнездо. В центре - ОГРОМНАЯ КОРОЛЕВА ЧУЖИХ.")
	fmt.Println("Она охраняет яйца. Ее челюсти щелкают, хвост бьет по полу.")
	fmt.Println("Это последний бой. Вся колония зависит от вас!")
	fmt.Println("ЗА ВЕЙЛАНД-ЮТАНИ! ЗА ЧЕЛОВЕЧЕСТВО!")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Print("Нажмите Enter...")
	scanner.Scan()

	e3 := NewEnemy()
	e3.Name = "👑 КОРОЛЕВА КСЕНОМОРФОВ 👑"
	e3.HP = 120
	e3.Strength = 15

	if !FightPVE(p, e3) {
		fmt.Println("\n💀 Королева уничтожила вас... Но, может быть, кто-то еще найдет это сообщение?")
		return
	}

	// Эпилог
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🏆 ЭПИЛОГ: СПАСЕНИЕ")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Королева повержена. Гнездо уничтожено.")
	fmt.Println("Вы запускаете самоуничтожение колонии и успеваете на эвакуационный шаттл.")
	fmt.Println("Сидя в кресле пилота, вы смотрите, как LV-426 взрывается.")
	fmt.Println("На экране появляется сообщение:")
	fmt.Println("   'Миссия выполнена. Вы - герой. Добро пожаловать домой, морпех.'")
	fmt.Println("\n🏅 ВЫ ВЫЖИЛИ! 🏅")
	fmt.Println("Конец игры.")
	fmt.Println(strings.Repeat("=", 50))
}

func PvPLocal() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n=== PVP НА ОДНОМ КОМПЕ ===")
	fmt.Print("Имя игрока 1: ")
	scanner.Scan()
	p1 := NewPlayer(scanner.Text())
	fmt.Print("Имя игрока 2: ")
	scanner.Scan()
	p2 := NewPlayer(scanner.Text())

	p1.Inventory = append(p1.Inventory, RandomItem())
	p2.Inventory = append(p2.Inventory, RandomItem())

	curr, other := p1, p2

	for p1.HP > 0 && p2.HP > 0 {
		fmt.Printf("\n🔥 ХОД %s\n", curr.Name)
		fmt.Printf("%s HP: %d | %s HP: %d\n", p1.Name, p1.HP, p2.Name, p2.HP)
		fmt.Println("1. Атаковать")
		fmt.Println("2. Инвентарь")
		fmt.Println("3. Написать сообщение")
		fmt.Print("Выбор: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "2":
			curr.ShowInventory()
			fmt.Print("Использовать предмет (индекс): ")
			scanner.Scan()
			var idx int
			fmt.Sscan(scanner.Text(), &idx)
			if idx >= 0 {
				curr.UseItem(idx)
			}
			curr, other = other, curr
			continue
		case "3":
			fmt.Print("Сообщение: ")
			scanner.Scan()
			msg := scanner.Text()
			fmt.Printf("\n💬 %s: %s\n", curr.Name, msg)
			continue
		}

		fmt.Printf("\n%s, отвернись!\n", other.Name)
		fmt.Scanln()

		dmg := curr.Attack()
		other.TakeDamage(dmg)
		fmt.Printf("%s нанес %d урона!\n", curr.Name, dmg)

		curr, other = other, curr
	}

	if p1.HP <= 0 {
		fmt.Printf("\n🏆 ПОБЕДИЛ %s!\n", p2.Name)
	} else {
		fmt.Printf("\n🏆 ПОБЕДИЛ %s!\n", p1.Name)
	}
}

// ==================== СЕТЬ - РАБОТАЕТ НОРМАЛЬНО ====================

type GameServer struct {
	conn1    net.Conn
	conn2    net.Conn
	name1    string
	name2    string
	p1       *Player
	p2       *Player
	waiting1 bool
	waiting2 bool
}

func NewGameServer() *GameServer {
	return &GameServer{}
}

func (s *GameServer) Start(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer ln.Close()
	fmt.Printf("🚀 Сервер на %s. Жду игроков...\n", port)

	s.conn1, _ = ln.Accept()
	fmt.Println("✅ Игрок 1 подключился")
	s.conn2, _ = ln.Accept()
	fmt.Println("✅ Игрок 2 подключился")

	buf := make([]byte, 1024)
	n, _ := s.conn1.Read(buf)
	s.name1 = strings.TrimSpace(string(buf[:n]))
	n, _ = s.conn2.Read(buf)
	s.name2 = strings.TrimSpace(string(buf[:n]))

	s.p1 = NewPlayer(s.name1)
	s.p2 = NewPlayer(s.name2)

	s.conn1.Write([]byte(fmt.Sprintf("👾 Противник: %s\n", s.name2)))
	s.conn2.Write([]byte(fmt.Sprintf("👾 Противник: %s\n", s.name1)))
	s.conn1.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━\n"))
	s.conn1.Write([]byte("1 - атаковать\n2 - инвентарь\n3 - сообщение\n"))
	s.conn2.Write([]byte("━━━━━━━━━━━━━━━━━━━━━━\n"))
	s.conn2.Write([]byte("1 - атаковать\n2 - инвентарь\n3 - сообщение\n"))

	go s.handleClient1()
	go s.handleClient2()
	s.gameLoop()
}

func (s *GameServer) handleClient1() {
	r := bufio.NewReader(s.conn1)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.TrimSpace(msg)
		if msg == "1" {
			s.waiting1 = true
		} else if msg == "2" {
			s.showInventoryTo(s.p1, s.conn1)
		} else if strings.HasPrefix(msg, "3") {
			chatMsg := strings.TrimPrefix(msg, "3")
			chatMsg = strings.TrimSpace(chatMsg)
			if chatMsg != "" {
				s.conn2.Write([]byte(fmt.Sprintf("\n💬 %s: %s\n", s.name1, chatMsg)))
			}
		}
	}
}

func (s *GameServer) handleClient2() {
	r := bufio.NewReader(s.conn2)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.TrimSpace(msg)
		if msg == "1" {
			s.waiting2 = true
		} else if msg == "2" {
			s.showInventoryTo(s.p2, s.conn2)
		} else if strings.HasPrefix(msg, "3") {
			chatMsg := strings.TrimPrefix(msg, "3")
			chatMsg = strings.TrimSpace(chatMsg)
			if chatMsg != "" {
				s.conn1.Write([]byte(fmt.Sprintf("\n💬 %s: %s\n", s.name2, chatMsg)))
			}
		}
	}
}

func (s *GameServer) showInventoryTo(p *Player, conn net.Conn) {
	if len(p.Inventory) == 0 {
		conn.Write([]byte("📦 Инвентарь пуст\n"))
		return
	}
	conn.Write([]byte("📦 ИНВЕНТАРЬ:\n"))
	for i, item := range p.Inventory {
		conn.Write([]byte(fmt.Sprintf("%d. %s", i, item.Name)))
		if item.Attack > 0 {
			conn.Write([]byte(fmt.Sprintf(" ⚔️ +%d", item.Attack)))
		}
		if item.Defence > 0 {
			conn.Write([]byte(fmt.Sprintf(" 🛡️ +%d", item.Defence)))
		}
		if item.Heal > 0 {
			conn.Write([]byte(fmt.Sprintf(" 💊 +%d", item.Heal)))
		}
		conn.Write([]byte("\n"))
	}
}

func (s *GameServer) gameLoop() {
	for s.p1.HP > 0 && s.p2.HP > 0 {
		s.waiting1 = false
		s.conn1.Write([]byte("\n⚔️ ТВОЙ ХОД! (1-атака,2-инв,3-чат): "))
		
		timeout := time.After(60 * time.Second)
		select {
		case <-timeout:
			s.conn1.Write([]byte("\n❌ Таймаут! Ход переходит...\n"))
			s.conn2.Write([]byte(fmt.Sprintf("\n⏰ %s не ответил, ход передан\n", s.name1)))
			continue
		default:
			for !s.waiting1 {
				time.Sleep(100 * time.Millisecond)
			}
		}
		
		if s.waiting1 {
			dmg := s.p1.Attack()
			s.p2.HP -= dmg
			if s.p2.HP < 0 {
				s.p2.HP = 0
			}
			s.conn1.Write([]byte(fmt.Sprintf("💥 Ты нанес %d урона!\n", dmg)))
			s.conn2.Write([]byte(fmt.Sprintf("💥 %s нанес %d урона! Твой HP: %d\n", s.name1, dmg, s.p2.HP)))
		}
		
		if s.p2.HP <= 0 {
			break
		}
		
		s.waiting2 = false
		s.conn2.Write([]byte("\n⚔️ ТВОЙ ХОД! (1-атака,2-инв,3-чат): "))
		
		timeout2 := time.After(60 * time.Second)
		select {
		case <-timeout2:
			s.conn2.Write([]byte("\n❌ Таймаут! Ход переходит...\n"))
			s.conn1.Write([]byte(fmt.Sprintf("\n⏰ %s не ответил, ход передан\n", s.name2)))
			continue
		default:
			for !s.waiting2 {
				time.Sleep(100 * time.Millisecond)
			}
		}
		
		if s.waiting2 {
			dmg := s.p2.Attack()
			s.p1.HP -= dmg
			if s.p1.HP < 0 {
				s.p1.HP = 0
			}
			s.conn2.Write([]byte(fmt.Sprintf("💥 Ты нанес %d урона!\n", dmg)))
			s.conn1.Write([]byte(fmt.Sprintf("💥 %s нанес %d урона! Твой HP: %d\n", s.name2, dmg, s.p1.HP)))
		}
	}
	
	if s.p1.HP <= 0 {
		s.conn1.Write([]byte("\n💀 ТЫ ПРОИГРАЛ!\n"))
		s.conn2.Write([]byte("\n🏆 ТЫ ПОБЕДИЛ!\n"))
	} else {
		s.conn1.Write([]byte("\n🏆 ТЫ ПОБЕДИЛ!\n"))
		s.conn2.Write([]byte("\n💀 ТЫ ПРОИГРАЛ!\n"))
	}
	
	s.conn1.Close()
	s.conn2.Close()
	fmt.Println("\n👽 Бой окончен!")
}

func RunClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}
	defer conn.Close()
	
	fmt.Println("✅ Подключено к серверу!")
	fmt.Print("Введите имя: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()
	conn.Write([]byte(name + "\n"))
	
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("\n❌ Соединение потеряно")
				os.Exit(0)
			}
			fmt.Print(string(buf[:n]))
		}
	}()
	
	fmt.Println("\n💬 КОМАНДЫ:")
	fmt.Println("   1 - атаковать")
	fmt.Println("   2 - инвентарь")
	fmt.Println("   3 текст - сообщение (например: 3 привет)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
	
	for scanner.Scan() {
		text := scanner.Text()
		conn.Write([]byte(text + "\n"))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("👽 ALIEN: USCM MISSION 👽")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. 🎬 Одиночная кампания (Прометей)")
	fmt.Println("2. 🔥 PvP на одном компьютере")
	fmt.Println("3. 🌐 Сетевой PvP (клиент)")
	fmt.Println("4. 🖥️ Запустить сервер")
	fmt.Println("5. ❌ Выход")
	fmt.Print("\nВыбор: ")

	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		Campaign()
	case "2":
		PvPLocal()
	case "3":
		fmt.Print("Адрес (localhost:8888): ")
		scanner.Scan()
		addr := scanner.Text()
		if addr == "" {
			addr = "localhost:8888"
		}
		RunClient(addr)
	case "4":
		fmt.Print("Порт (8888): ")
		scanner.Scan()
		port := scanner.Text()
		if port == "" {
			port = "8888"
		}
		NewGameServer().Start(port)
	case "5":
		fmt.Println("До свидания!")
	default:
		fmt.Println("Неверный выбор")
	}
}