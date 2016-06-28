package bqs

type EntityType struct {
	Id    int
	Name  string
	Class string
}

var (
	TYPE_WARRIOR = EntityType{1, "warrior", "player"}
	// Mobs
	TYPE_RAT = EntityType{2, "RAT", "mob"}
	TYPE_SKELETON = EntityType{3, "SKELETON", "mob"}
	TYPE_GOBLIN = EntityType{4, "GOBLIN", "mob"}
	TYPE_OGRE = EntityType{5, "OGRE", "mob"}
	TYPE_SPECTRE = EntityType{6, "SPECTRE", "mob"}
	TYPE_CRAB = EntityType{7, "CRAB", "mob"}
	TYPE_BAT = EntityType{8, "BAT", "mob"}
	TYPE_WIZARD = EntityType{9, "WIZARD", "mob"}
	TYPE_EYE = EntityType{10, "EYE", "mob"}
	TYPE_SNAKE = EntityType{11, "SNAKE", "mob"}
	TYPE_SKELETON2 = EntityType{12, "SKELETON2", "mob"}
	TYPE_BOSS = EntityType{13, "BOSS", "mob"}
	TYPE_DEATHKNIGHT = EntityType{14, "DEATHKNIGHT", "mob"}

	// Weapons
	TYPE_SWORD1 = EntityType{60, "SWORD1", "weapon"}
	TYPE_SWORD2 = EntityType{61, "SWORD2", "weapon"}
	TYPE_REDSWORD = EntityType{62, "REDSWORD", "weapon"}
	TYPE_GOLDENSWORD = EntityType{63, "GOLDENSWORD", "weapon"}
	TYPE_MORNINGSTAR = EntityType{64, "MORNINGSTAR", "weapon"}
	TYPE_AXE = EntityType{65, "AXE", "weapon"}
	TYPE_BLUESWORD = EntityType{66, "BLUESWORD", "weapon"}

	// Armors
	TYPE_FIREFOX = EntityType{20, "FIREFOX", "armor"}
	TYPE_CLOTHARMOR = EntityType{21, "CLOTHARMOR", "armor"}
	TYPE_LEATHERARMOR = EntityType{22, "LEATHERARMOR", "armor"}
	TYPE_MAILARMOR = EntityType{23, "MAILARMOR", "armor"}
	TYPE_PLATEARMOR = EntityType{24, "PLATEARMOR", "armor"}
	TYPE_REDARMOR = EntityType{25, "REDARMOR", "armor"}
	TYPE_GOLDENARMOR = EntityType{26, "GOLDENARMOR", "armor"}

	// Objects
	TYPE_FLASK = EntityType{35, "FLASK", "object"}
	TYPE_BURGER = EntityType{36, "BURGER", "object"}
	TYPE_CHEST = EntityType{37, "CHEST", "object"}
	TYPE_FIREPOTION = EntityType{38, "FIREPOTION", "object"}
	TYPE_CAKE = EntityType{39, "CAKE", "object"}

	// NPCs
	TYPE_GUARD = EntityType{40, "GUARD", "npc"}
	TYPE_KING = EntityType{41, "KING", "npc"}
	TYPE_OCTOCAT = EntityType{42, "OCTOCAT", "npc"}
	TYPE_VILLAGEGIRL = EntityType{43, "VILLAGEGIRL", "npc"}
	TYPE_VILLAGER = EntityType{44, "VILLAGER", "npc"}
	TYPE_PRIEST = EntityType{45, "PRIEST", "npc"}
	TYPE_SCIENTIST = EntityType{46, "SCIENTIST", "npc"}
	TYPE_AGENT = EntityType{47, "AGENT", "npc"}
	TYPE_RICK = EntityType{48, "RICK", "npc"}
	TYPE_NYAN = EntityType{49, "NYAN", "npc"}
	TYPE_SORCERER = EntityType{50, "SORCERER", "npc"}
	TYPE_BEACHNPC = EntityType{51, "BEACHNPC", "npc"}
	TYPE_FORESTNPC = EntityType{52, "FORESTNPC", "npc"}
	TYPE_DESERTNPC = EntityType{53, "DESERTNPC", "npc"}
	TYPE_LAVANPC = EntityType{54, "LAVANPC", "npc"}
	TYPE_CODER = EntityType{55, "CODER", "npc"}
)

var EntityTypeMap = map[string]EntityType{
	"rat": TYPE_RAT,
	"skeleton": TYPE_SKELETON,
	"skeleton2": TYPE_SKELETON2,
	"goblin": TYPE_GOBLIN,
	"ogre": TYPE_OGRE,

	"spectre": TYPE_SPECTRE,
	"deathknight": TYPE_DEATHKNIGHT,
	"crab": TYPE_CRAB,
	"snake": TYPE_SNAKE,
	"bat": TYPE_BAT,
	"wizard": TYPE_WIZARD,
	"eye": TYPE_EYE,
	"boss": TYPE_BOSS,

	"sword1": TYPE_SWORD1,
	"sword2": TYPE_SWORD2,
	"axe": TYPE_AXE,
	"redsword": TYPE_REDSWORD,
	"bluesword": TYPE_BLUESWORD,
	"goldensword": TYPE_GOLDENSWORD,
	"morningstar": TYPE_MORNINGSTAR,

	"firefox": TYPE_FIREFOX,
	"clotharmor": TYPE_CLOTHARMOR,
	"leatherarmor": TYPE_LEATHERARMOR,
	"mailarmor": TYPE_MAILARMOR,
	"platearmor": TYPE_PLATEARMOR,
	"redarmor": TYPE_REDARMOR,
	"goldenarmor": TYPE_GOLDENARMOR,

	"flask": TYPE_FLASK,
	"cake": TYPE_CAKE,
	"burger": TYPE_BURGER,
	"chest": TYPE_CHEST,
	"firepotion": TYPE_FIREPOTION,

	"guard": TYPE_GUARD,
	"villagegirl": TYPE_VILLAGEGIRL,
	"villager": TYPE_VILLAGER,
	"coder": TYPE_CODER,
	"scientist": TYPE_SCIENTIST,
	"priest": TYPE_PRIEST,
	"king": TYPE_KING,
	"rick": TYPE_RICK,
	"nyan": TYPE_NYAN,
	"sorcerer": TYPE_SORCERER,
	"agent": TYPE_AGENT,
	"octocat": TYPE_OCTOCAT,
	"beachnpc": TYPE_BEACHNPC,
	"forestnpc": TYPE_FORESTNPC,
	"desertnpc": TYPE_DESERTNPC,
	"lavanpc": TYPE_LAVANPC,
}

var RankedArmors = map[int]int{
	TYPE_CLOTHARMOR.Id: 1,
	TYPE_LEATHERARMOR.Id: 2,
	TYPE_MAILARMOR.Id: 3,
	TYPE_PLATEARMOR.Id: 4,
	TYPE_REDARMOR.Id: 5,
	TYPE_GOLDENARMOR.Id: 6,
}

var RankedWeapon = map[int]int{
	TYPE_SWORD1.Id: 1,
	TYPE_SWORD2.Id: 2,
	TYPE_AXE.Id: 3,
	TYPE_MORNINGSTAR.Id: 4,
	TYPE_BLUESWORD.Id: 5,
	TYPE_REDSWORD.Id: 6,
	TYPE_GOLDENSWORD.Id: 7,
}

func GetEntityTypeByString(kind string) EntityType {
	return EntityTypeMap[kind]
}

func GetWeaponRank(weaponId int) int {
	return RankedWeapon[weaponId]
}

func GetArmorRank(armorId int) int {
	return RankedWeapon[armorId]
}