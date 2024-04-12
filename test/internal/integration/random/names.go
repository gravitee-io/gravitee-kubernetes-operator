// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package random

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
)

var (
	safeRandom = os.Getenv("SAFE_RANDOM") == env.TrueString

	dbFileNAme = "names.db"

	eol = "\n"

	intRandMax = 99

	dbFile *os.File
)

var (
	left = [...]string{
		"admiring",
		"adoring",
		"affectionate",
		"agitated",
		"amazing",
		"angry",
		"awesome",
		"beautiful",
		"blissful",
		"bold",
		"boring",
		"brave",
		"busy",
		"charming",
		"clever",
		"compassionate",
		"competent",
		"condescending",
		"confident",
		"cool",
		"cranky",
		"crazy",
		"dazzling",
		"determined",
		"distracted",
		"dreamy",
		"eager",
		"ecstatic",
		"elastic",
		"elated",
		"elegant",
		"eloquent",
		"epic",
		"exciting",
		"fervent",
		"festive",
		"flamboyant",
		"focused",
		"friendly",
		"frosty",
		"funny",
		"gallant",
		"gifted",
		"goofy",
		"gracious",
		"great",
		"happy",
		"hardcore",
		"heuristic",
		"hopeful",
		"hungry",
		"infallible",
		"inspiring",
		"intelligent",
		"interesting",
		"jolly",
		"jovial",
		"keen",
		"kind",
		"laughing",
		"loving",
		"lucid",
		"magical",
		"modest",
		"musing",
		"mystifying",
		"naughty",
		"nervous",
		"nice",
		"nifty",
		"nostalgic",
		"objective",
		"optimistic",
		"peaceful",
		"pedantic",
		"pensive",
		"practical",
		"priceless",
		"quirky",
		"quizzical",
		"recursing",
		"relaxed",
		"reverent",
		"romantic",
		"sad",
		"serene",
		"sharp",
		"silly",
		"sleepy",
		"stoic",
		"strange",
		"stupefied",
		"suspicious",
		"sweet",
		"tender",
		"thirsty",
		"trusting",
		"unruffled",
		"upbeat",
		"vibrant",
		"vigilant",
		"vigorous",
		"wizardly",
		"wonderful",
		"xenodochial",
		"youthful",
		"zealous",
		"zen",
	}

	right = [...]string{

		"agnesi",

		"albattani",

		"allen",

		"almeida",

		"antonelli",

		"archimedes",

		"ardinghelli",

		"aryabhata",

		"austin",

		"babbage",

		"banach",

		"banzai",

		"bardeen",

		"bartik",

		"bassi",

		"beaver",

		"bell",

		"benz",

		"bhabha",

		"bhaskara",

		"black",

		"blackburn",

		"blackwell",

		"bohr",

		"booth",

		"borg",

		"bose",

		"bouman",

		"boyd",

		"brahmagupta",

		"brattain",

		"brown",

		"buck",

		"burnell",

		"cannon",

		"carson",

		"cartwright",

		"carver",

		"cerf",

		"chandrasekhar",

		"chaplygin",

		"chatelet",

		"chatterjee",

		"chaum",

		"chebyshev",

		"clarke",

		"cohen",

		"colden",

		"cori",

		"cray",

		"curie",

		"curran",

		"darwin",

		"davinci",

		"dewdney",

		"dhawan",

		"diffie",

		"dijkstra",

		"dirac",

		"driscoll",

		"dubinsky",

		"easley",

		"edison",

		"einstein",

		"elbakyan",

		"elgamal",

		"elion",

		"ellis",

		"engelbart",

		"euclid",

		"euler",

		"faraday",

		"feistel",

		"fermat",

		"fermi",

		"feynman",

		"franklin",

		"gagarin",

		"galileo",

		"galois",

		"ganguly",

		"gates",

		"gauss",

		"germain",

		"goldberg",

		"goldstine",

		"goldwasser",

		"golick",

		"goodall",

		"gould",

		"greider",

		"grothendieck",

		"haibt",

		"hamilton",

		"haslett",

		"hawking",

		"heisenberg",

		"hellman",

		"hermann",

		"herschel",

		"hertz",

		"heyrovsky",

		"hodgkin",

		"hofstadter",

		"hoover",

		"hopper",

		"hugle",

		"hypatia",

		"ishizaka",

		"jackson",

		"jang",

		"jemison",

		"jennings",

		"jepsen",

		"johnson",

		"joliot",

		"jones",

		"kalam",

		"kapitsa",

		"kare",

		"keldysh",

		"keller",

		"kepler",

		"khayyam",

		"khorana",

		"kilby",

		"kirch",

		"knuth",

		"kowalevski",

		"lalande",

		"lamarr",

		"lamport",

		"leakey",

		"leavitt",

		"lederberg",

		"lehmann",

		"lewin",

		"lichterman",

		"liskov",

		"lovelace",

		"lumiere",

		"mahavira",

		"margulis",

		"matsumoto",

		"maxwell",

		"mayer",

		"mccarthy",

		"mcclintock",

		"mclaren",

		"mclean",

		"mcnulty",

		"meitner",

		"mendel",

		"mendeleev",

		"meninsky",

		"merkle",

		"mestorf",

		"mirzakhani",

		"montalcini",

		"moore",

		"morse",

		"moser",

		"murdock",

		"napier",

		"nash",

		"neumann",

		"newton",

		"nightingale",

		"nobel",

		"noether",

		"northcutt",

		"noyce",

		"panini",

		"pare",

		"pascal",

		"pasteur",

		"payne",

		"perlman",

		"pike",

		"poincare",

		"poitras",

		"proskuriakova",

		"ptolemy",

		"raman",

		"ramanujan",

		"rhodes",

		"ride",

		"ritchie",

		"robinson",

		"roentgen",

		"rosalind",

		"rubin",

		"saha",

		"sammet",

		"sanderson",

		"satoshi",

		"shamir",

		"shannon",

		"shaw",

		"shirley",

		"shockley",

		"shtern",

		"sinoussi",

		"snyder",

		"solomon",

		"spence",

		"stonebraker",

		"sutherland",

		"swanson",

		"swartz",

		"swirles",

		"taussig",

		"tesla",

		"tharp",

		"thompson",

		"torvalds",

		"tu",

		"turing",

		"varahamihira",

		"vaughan",

		"villani",

		"visvesvaraya",

		"volhard",

		"wescoff",

		"wilbur",

		"wiles",

		"williams",

		"williamson",

		"wilson",

		"wing",

		"wozniak",

		"wright",

		"wu",

		"yalow",

		"yonath",

		"zhukovsky",

		"brassely",

		"elamrani",

		"geraud",

		"compiegne",

		"ahmadpour",

		"cordier",

		"michaux",

		"tschacher",

		"beauchemin",

		"cheggour",

		"chamfroy",

		"fernandez",

		"cambier",

		"santos",

		"enachi",

		"fernandes",

		"timoska",

		"stojanovski",

		"veljanoski",

		"gjorgievski",

		"netkov",

		"cusnieux",

		"pacaud",

		"leleu",

		"khelifi",

		"devaux",

		"lamirand",

		"tavernier",

		"avenier",

		"waller",

		"giovaresco",

		"maisse",

		"pisicchio",

		"haeyaert",
	}
)

func GetName() string {
	name := strings.ReplaceAll(
		left[rand.Intn(len(left))]+"_"+right[rand.Intn(len(right))], "_", "-", //nolint:gosec // this is safe
	)
	if safeRandom {
		return withSafeRandom(name)
	}
	return name
}

func GetSuffix() string {
	return "-" + GetName()
}

// This is only used in CI to avoid failing tests because of conflicts on the cluster.
func withSafeRandom(name string) string {
	if exists(name) {
		name = fmt.Sprintf("%s-%02d", name, rand.Intn(intRandMax)) //nolint:gosec // this is safe
		return withSafeRandom(name)
	}
	_, err := dbFile.WriteString(name + eol)
	if err != nil {
		panic(err)
	}
	return name
}

func exists(name string) bool {
	content, err := os.ReadFile(dbFile.Name())
	if err != nil {
		panic(err)
	}
	known := strings.Split(string(content), eol)
	for _, n := range known {
		if n == name {
			return true
		}
	}
	return false
}

func init() {
	if safeRandom {
		filePath := path.Join(os.TempDir(), dbFileNAme)
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		dbFile = file
	}
}
