package data

import (
	"math/rand"
	"strings"
)

var adjectives = []string{
	"able", "brave", "calm", "clean", "clear", "cold", "dark", "deep", "dry", "easy",
	"elated", "fair", "fast", "fine", "free", "friendly", "good", "hard", "harsh", "heavy",
	"honest", "hot", "huge", "innocent", "jolly", "keen", "kind", "lazy", "light", "little",
	"loyal", "mean", "modern", "neat", "nice", "open", "optimistic", "petty", "quick", "rare",
	"real", "reliable", "rough", "safe", "sharp", "shy", "silly", "simple", "slim", "smart",
	"soft", "strong", "sweet", "tender", "thin", "tight", "timid", "tough", "unique", "vain",
	"vague", "valid", "warm", "weak", "wet", "wild", "wise", "young", "zany", "zealous",
	"altruistic", "artsy", "bright", "broad", "bold", "bitter", "bouncy", "careful", "cheerful", "chilly",
	"chunky", "clumsy", "courageous", "crafty", "cute", "curious", "dense", "delicate", "difficult", "eager",
	"easygoing", "energetic", "elegant", "fiery", "flat", "fresh", "genuine", "graceful", "grim", "gritty",
	"happy", "harmless", "ideal", "intense", "inventive", "jumpy", "lively", "loving", "loud", "low",
	"messy", "mild", "mischievous", "needy", "new", "noble", "nostalgic", "precious", "racy", "reasonable",
	"righteous", "robust", "serene", "shiny", "sleek", "stiff", "sturdy", "swift", "tasty", "trustworthy",
}

var nouns = []string{
	"apple", "banana", "car", "dog", "elephant", "fish", "guitar", "house", "island", "jacket",
	"kangaroo", "lamp", "moon", "notebook", "orange", "pencil", "quilt", "river", "sun", "tree",
	"umbrella", "violin", "window", "xylophone", "yacht", "zebra", "airplane", "bicycle", "cat",
	"dolphin", "ear", "flame", "grape", "honey", "insect", "jungle", "key", "lemon", "mountain",
	"nail", "ocean", "planet", "queen", "rocket", "star", "telescope", "unicorn", "vampire", "whale",
	"x-ray", "yellow", "zoo", "lime", "butterfly", "cliff", "daisy", "eagle", "firefly", "giraffe",
	"hammock", "ice", "jewel", "kettle", "lava", "mango", "needle", "oasis", "parrot", "quail",
	"rock", "snow", "turtle", "vulture", "waterfall", "xenon", "yarn", "zinnia", "acorn", "bottle",
	"candle", "diamond", "eggplant", "feather", "gravel", "halo", "inbox", "necklace", "keychain",
	"leash", "mushroom", "net", "octopus", "piano", "quiver", "ring", "stone", "tea", "underwear",
	"vase", "wood", "pumpkin", "yogurt", "zenith", "albatross", "bee", "clover", "duck", "rat",
	"fern", "goose", "hummingbird", "ivory", "jaguar", "koala", "lily", "maple", "nut", "olive",
	"pear", "quince", "rooster", "seahorse", "toucan", "violet", "wombat", "mouse", "yak",
	"rose",
}

// GenerateName generates a human-readable name for device identification.
func GenerateName() string {
	return strings.ToUpper(adjectives[rand.Intn(len(adjectives))]) + "-" + strings.ToUpper(nouns[rand.Intn(len(nouns))])
}
