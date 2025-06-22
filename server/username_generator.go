package main

import (
	"fmt"
	"math/rand"
	"time"
)

var adjectives = []string{
	"happy", "brave", "clever", "bright", "swift", "gentle", "mighty", "silent",
	"golden", "silver", "crimson", "azure", "emerald", "violet", "amber", "coral",
	"dancing", "laughing", "singing", "jumping", "running", "flying", "swimming", "climbing",
	"mysterious", "magical", "ancient", "modern", "tiny", "giant", "smooth", "rough",
	"shiny", "sparkly", "fuzzy", "fluffy", "spiky", "round", "square", "curved",
	"wild", "tame", "bold", "shy", "loud", "quiet", "fast", "slow",
	"hot", "cold", "warm", "cool", "fresh", "old", "new", "young",
	"wise", "silly", "funny", "serious", "cheerful", "grumpy", "sleepy", "awake",
	"hungry", "full", "thirsty", "satisfied", "excited", "calm", "nervous", "confident",
	"lucky", "unlucky", "strong", "weak", "tall", "short", "big", "small",
	"electric", "cosmic", "stellar", "lunar", "solar", "oceanic", "mountain", "forest",
	"desert", "arctic", "tropical", "windy", "stormy", "sunny", "cloudy", "misty",
}

var nouns = []string{
	"cat", "dog", "bird", "fish", "rabbit", "turtle", "hamster", "gecko",
	"lion", "tiger", "bear", "wolf", "fox", "deer", "eagle", "hawk",
	"dolphin", "whale", "shark", "octopus", "jellyfish", "starfish", "crab", "lobster",
	"butterfly", "dragonfly", "bee", "ant", "spider", "ladybug", "cricket", "firefly",
	"tree", "flower", "grass", "moss", "fern", "cactus", "rose", "daisy",
	"mountain", "river", "lake", "ocean", "beach", "forest", "desert", "valley",
	"star", "moon", "sun", "cloud", "rainbow", "thunder", "lightning", "snow",
	"book", "pen", "paper", "computer", "phone", "camera", "watch", "lamp",
	"chair", "table", "bed", "couch", "pillow", "blanket", "cup", "plate",
	"car", "bike", "train", "plane", "boat", "rocket", "balloon", "kite",
	"pizza", "cake", "cookie", "apple", "banana", "orange", "grape", "cherry",
	"music", "song", "dance", "art", "paint", "brush", "canvas", "sculpture",
	"game", "toy", "puzzle", "ball", "doll", "robot", "castle", "bridge",
	"key", "door", "window", "mirror", "clock", "bell", "drum", "guitar",
	"magic", "wizard", "dragon", "unicorn", "fairy", "knight", "princess", "treasure",
	"adventure", "journey", "quest", "mystery", "secret", "surprise", "gift", "party",
}

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GenerateRandomUsername() string {
	adjective := adjectives[rng.Intn(len(adjectives))]
	noun := nouns[rng.Intn(len(nouns))]
	return fmt.Sprintf("%s-%s", adjective, noun)
}
