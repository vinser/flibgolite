package parser

// Based on: https://github.com/dschnelldavis/parse-full-name
// 1. Make go-routine proof by removing racondition on globals
// 2. Make faster by using maps for lookups instead of Dijksta's "L" algorythm
// 3. Fix case sensitivity issue.
// 4. Fix speed by generating regular expressions only once.

// See: https://github.com/pschlump/names.git

import (
	"log"
	"regexp"
	"strings"

	"github.com/vinser/flibgolite/pkg/model"
)

type ParsedName struct {
	Title  string
	First  string
	Middle string
	Last   string
	Nick   string
	Suffix string
}

var (
	suffixListLookup = map[string]bool{
		"2":       true,
		"cfp":     true,
		"chfc":    true,
		"clu":     true,
		"d.c.":    true,
		"d.o.":    true,
		"esq":     true,
		"esquire": true,
		"ii":      true,
		"iii":     true,
		"iv":      true,
		"j.d.":    true,
		"jnr":     true,
		"jr":      true,
		"ll.m.":   true,
		"m.d.":    true,
		"md":      true,
		"p.c.":    true,
		"ph.d.":   true,
		"phd":     true,
		"snr":     true,
		"sr":      true,
		"v":       true,
	}

	prefixListLookup = map[string]bool{
		"a":       true,
		"ab":      true,
		"antune":  true,
		"ap":      true,
		"abu":     true,
		"al":      true,
		"alm":     true,
		"alt":     true,
		"bab":     true,
		"bäck":    true,
		"bar":     true,
		"bath":    true,
		"bat":     true,
		"beau":    true,
		"beck":    true,
		"ben":     true,
		"berg":    true,
		"bet":     true,
		"bin":     true,
		"bint":    true,
		"birch":   true,
		"björk":   true,
		"björn":   true,
		"bjur":    true,
		"da":      true,
		"dahl":    true,
		"dal":     true,
		"de":      true,
		"degli":   true,
		"dele":    true,
		"del":     true,
		"della":   true,
		"der":     true,
		"di":      true,
		"dos":     true,
		"du":      true,
		"e":       true,
		"ek":      true,
		"el":      true,
		"escob":   true,
		"esch":    true,
		"fleisch": true,
		"fitz":    true,
		"fors":    true,
		"gott":    true,
		"griff":   true,
		"haj":     true,
		"haug":    true,
		"holm":    true,
		"ibn":     true,
		"kauf":    true,
		"kil":     true,
		"koop":    true,
		"kvarn":   true,
		"la":      true,
		"le":      true,
		"lind":    true,
		"lönn":    true,
		"lund":    true,
		"mac":     true,
		"mhic":    true,
		"mic":     true,
		"mir":     true,
		"na":      true,
		"naka":    true,
		"neder":   true,
		"nic":     true,
		"ni":      true,
		"nin":     true,
		"nord":    true,
		"norr":    true,
		"ny":      true,
		"o":       true,
		"ua":      true,
		`ui\'`:    true,
		"öfver":   true,
		"ost":     true,
		"över":    true,
		"öz":      true,
		"papa":    true,
		"pour":    true,
		"quarn":   true,
		"skog":    true,
		"skoog":   true,
		"sten":    true,
		"stor":    true,
		"ström":   true,
		"söder":   true,
		"ter":     true,
		"tre":     true,
		"türk":    true,
		"van":     true,
		"väst":    true,
		"väster":  true,
		"vest":    true,
		"von":     true,
	}

	titleListLookup = map[string]bool{
		"a v m":             true,
		"admiraal":          true,
		"admiral":           true,
		"air cdre":          true,
		"air commodore":     true,
		"air marshal":       true,
		"air vice marshal":  true,
		"alderman":          true,
		"alhaji":            true,
		"ambassador":        true,
		"baron":             true,
		"barones":           true,
		"brig gen":          true,
		"brig general":      true,
		"brig":              true,
		"brigadier general": true,
		"brigadier":         true,
		"brother":           true,
		"canon":             true,
		"capt":              true,
		"captain":           true,
		"cardinal":          true,
		"cdr":               true,
		"chief":             true,
		"cik":               true,
		"cmdr":              true,
		"coach":             true,
		"col":               true,
		"colonel":           true,
		"commandant":        true,
		"commander":         true,
		"commissioner":      true,
		"commodore":         true,
		"comte":             true,
		"comtessa":          true,
		"congressman":       true,
		"conseiller":        true,
		"consul":            true,
		"conte":             true,
		"contessa":          true,
		"corporal":          true,
		"councillor":        true,
		"count":             true,
		"countess":          true,
		"crown prince":      true,
		"crown princess":    true,
		"dame":              true,
		"datin":             true,
		"dato":              true,
		"datuk seri":        true,
		"datuk":             true,
		"deacon":            true,
		"deaconess":         true,
		"dean":              true,
		"dhr":               true,
		"dipl ing":          true,
		"doctor":            true,
		"dott sa":           true,
		"dott":              true,
		"dr ing":            true,
		"dr":                true,
		"dra":               true,
		"drs":               true,
		"embajador":         true,
		"embajadora":        true,
		"en":                true,
		"encik":             true,
		"eng":               true,
		"eur ing":           true,
		"exma sra":          true,
		"exmo sr":           true,
		"f o":               true,
		"father":            true,
		"first lieutenant":  true,
		"first officer":     true,
		"flt lieut":         true,
		"flying officer":    true,
		"fr":                true,
		"frau":              true,
		"fraulein":          true,
		"fru":               true,
		"gen":               true,
		"generaal":          true,
		"general":           true,
		"governor":          true,
		"graaf":             true,
		"gravin":            true,
		"group captain":     true,
		"grp capt":          true,
		"h e dr":            true,
		"h h":               true,
		"h m":               true,
		"h r h":             true,
		"hajah":             true,
		"haji":              true,
		"hajim":             true,
		"her highness":      true,
		"her majesty":       true,
		"herr":              true,
		"high chief":        true,
		"his highness":      true,
		"his holiness":      true,
		"his majesty":       true,
		"hon":               true,
		"hr":                true,
		"hra":               true,
		"ing":               true,
		"ir":                true,
		"jonkheer":          true,
		"judge":             true,
		"justice":           true,
		"khun ying":         true,
		"kolonel":           true,
		"lady":              true,
		"lcda":              true,
		"lic":               true,
		"lieut cdr":         true,
		"lieut col":         true,
		"lieut gen":         true,
		"lieut":             true,
		"lord":              true,
		"m l":               true,
		"m r":               true,
		"m":                 true,
		"madame":            true,
		"mademoiselle":      true,
		"maj gen":           true,
		"major":             true,
		"master":            true,
		"mevrouw":           true,
		"miss":              true,
		"mlle":              true,
		"mme":               true,
		"monsieur":          true,
		"monsignor":         true,
		"mr":                true,
		"mrs":               true,
		"ms":                true,
		"mstr":              true,
		"nti":               true,
		"pastor":            true,
		"president":         true,
		"prince":            true,
		"princess":          true,
		"princesse":         true,
		"prinses":           true,
		"prof sir":          true,
		"prof":              true,
		"professor":         true,
		"puan sri":          true,
		"puan":              true,
		"rabbi":             true,
		"rear admiral":      true,
		"rev canon":         true,
		"rev dr":            true,
		"rev mother":        true,
		"rev":               true,
		"reverend":          true,
		"rva":               true,
		"senator":           true,
		"sergeant":          true,
		"sheikh":            true,
		"sheikha":           true,
		"sig na":            true,
		"sig ra":            true,
		"sig":               true,
		"sir":               true,
		"sister":            true,
		"sqn ldr":           true,
		"sr d":              true,
		"sr":                true,
		"sra":               true,
		"srta":              true,
		"sultan":            true,
		"tan sri dato":      true,
		"tan sri":           true,
		"tengku":            true,
		"teuku":             true,
		"than puying":       true,
		"the hon dr":        true,
		"the hon justice":   true,
		"the hon miss":      true,
		"the hon mr":        true,
		"the hon mrs":       true,
		"the hon ms":        true,
		"the hon sir":       true,
		"the very rev":      true,
		"toh puan":          true,
		"tun":               true,
		"vice admiral":      true,
		"viscount":          true,
		"viscountess":       true,
		"wg cdr":            true,
	}

	conjunctionListLookup = map[string]bool{
		"&":   true,
		"and": true,
		"et":  true,
		"e":   true,
		"of":  true,
		"the": true,
		"und": true,
		"y":   true,
	}
)

func AuthorByFullName(fullName string) *model.Author {
	author := &model.Author{}
	names := ParseFullName(fullName)
	fullName = names.Title + " " + names.First + " " + names.Middle + " " + names.Last + " " + names.Suffix + " (" + names.Nick + ")"
	author.Name = strings.TrimSpace(strings.TrimSuffix(fullName, " ()"))
	sortName := names.Last + ", " + names.First + " " + names.Middle + " (" + names.Nick + ")"
	author.Sort = strings.TrimSuffix(strings.TrimSpace(strings.TrimSuffix(sortName, " ()")), ",")
	return author
}

var reDelimGlued = regexp.MustCompile(`\p{Ll}\p{Lu}|[^ ]\(|,[^ \r\n]|\.[^ ,\r\n]|\)[^ ,\r\n]`)

// DelimitGluedName changes "CamelCase" to "Camel Case".
func DelimitGluedName(fullName string) string {
	fullName = strings.Trim(reDelimGlued.ReplaceAllStringFunc(
		fullName,
		func(s string) string {
			r := []rune(s)
			if len(r) == 2 {
				return string(r[:1]) + " " + string(r[1:])
			}
			return s
		}), " \t")
	return fullName
}

// Add comma after last name
func AddCommaAfterLastName(fullName string) string {
	if !strings.Contains(fullName, ",") {
		parts := strings.Split(fullName, " ")
		if !strings.Contains(parts[0], ",") {
			parts[0] += ","
			fullName = strings.Join(parts, " ")
		}
	}
	return fullName
}

var reNickName = regexp.MustCompile(`\s?[\'\"\(\[]([^\[\]\)\)\'\"]+)[\'\"\)\]]`)

// findNickName pulls out parenthesized nicknames and returns them.
func findNickName(fullName string) (partsFound []string, newFullName string) {
	newFullName = fullName
	matches := reNickName.FindAllStringSubmatch(newFullName, -1)
	for _, vv := range matches {
		partsFound = append(partsFound, vv[1])
		newFullName = strings.Replace(newFullName, vv[0], "", -1)
	}
	return
}

var reSplitName = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

func splitNameIntoParts(fullName string) (parts []string, comma []bool) {
	fullName = DelimitGluedName(fullName)
	fullName = reSplitName.ReplaceAllLiteralString(fullName, " ")
	parts = strings.Split(strings.Trim(fullName, " \t"), " ")
	for ii, vv := range parts {
		parts[ii] = strings.Trim(vv, " \t")
		comma = append(comma, strings.HasSuffix(vv, ","))
		if comma[ii] {
			parts[ii] = strings.TrimSuffix(vv, ",")
		}
	}
	return
}

// searchMapForParts will find and remove the items that end in '.' like Dr.
func searchMapForParts(list map[string]bool, nameParts *[]string, nameCommas *[]bool) (partsFound []string) {
	if db1 {
		log.Printf("*nameParts %v\n", *nameParts)
	}
	for _, namePart := range *nameParts {
		if namePart == "" {
			continue
		}

		partToCheck := strings.TrimSuffix(strings.ToLower(namePart), ".")
		if found, ok := list[partToCheck]; ok && found {
			partsFound = append(partsFound, namePart)
		}
	}

	if db1 {
		log.Printf("partsFound %s\n", partsFound)
	}
	for _, vv := range partsFound {
		if foundIndex := locationInArray(vv, *nameParts); foundIndex >= 0 {
			*nameParts = removeAt(*nameParts, foundIndex)
			if (*nameCommas)[foundIndex] && foundIndex != len(*nameCommas)-1 {
				*nameCommas = removeAt(*nameCommas, foundIndex+1)
			} else {
				*nameCommas = removeAt(*nameCommas, foundIndex)
			}
		}
	}
	return
}

func findTitles(nameParts *[]string, nameCommas *[]bool) []string {
	return searchMapForParts(titleListLookup, nameParts, nameCommas)
}

func findSuffixes(nameParts *[]string, nameCommas *[]bool) []string {
	return searchMapForParts(suffixListLookup, nameParts, nameCommas)
}

func joinPrefixes(nameParts *[]string, nameCommas *[]bool) {
	if len((*nameParts)) > 1 {
		for ii := len((*nameParts)) - 2; ii >= 0; ii-- {
			if np, ok := prefixListLookup[(*nameParts)[ii]]; ok && np {
				(*nameParts)[ii] = (*nameParts)[ii] + " " + (*nameParts)[ii+1]
				(*nameParts) = removeAt(*nameParts, ii+1)
				(*nameCommas) = removeAt(*nameCommas, ii)
			}
		}
	}
}

func joinConjunctions(nameParts *[]string, nameCommas *[]bool) {
	if len((*nameParts)) > 2 {
		for ii := len((*nameParts)) - 3; ii >= 0; ii-- {
			if found, ok := conjunctionListLookup[(*nameParts)[ii+1]]; ok && found {
				(*nameParts)[ii] = (*nameParts)[ii] + " " + (*nameParts)[ii+1] + " " + (*nameParts)[ii+2]
				(*nameParts) = append((*nameParts)[:ii+1], (*nameParts)[ii+3:]...)
				(*nameCommas) = append((*nameCommas)[:ii], (*nameCommas)[ii+2:]...)
				ii--
			}
		}
	}
}

func findExtraSuffixes(nameParts *[]string, nameCommas *[]bool) (extraSuffixes []string) {
	commasCount := 0
	for _, v := range *nameCommas {
		if v {
			commasCount++
		}
	}
	if commasCount > 1 {
		for ii := len((*nameParts)) - 1; ii >= 2; ii-- {
			if (*nameCommas)[ii] {
				extraSuffixes = append(extraSuffixes, (*nameParts)[ii])
				(*nameParts) = removeAt(*nameParts, ii)
				(*nameCommas) = removeAt(*nameCommas, ii)
			}
		}
	}
	return
}

func extractLastName(nameParts *[]string, nameCommas *[]bool) (lastname string) {
	posFirstComma := len((*nameParts)) - 1
	for ii, vv := range *nameCommas {
		if vv {
			posFirstComma = ii
		}
	}
	lastname = (*nameParts)[posFirstComma]
	(*nameParts) = removeAt(*nameParts, posFirstComma)
	(*nameCommas) = (*nameCommas)[:0]
	return
}

func extrctFirstName(nameParts *[]string, nameCommas *[]bool) (firstname string) {
	firstname = (*nameParts)[0]
	(*nameParts) = (*nameParts)[1:]
	return
}

func extractMiddleName(nameParts *[]string, nameCommas *[]bool) (middlename string) {
	middlename = strings.Join((*nameParts), " ")
	(*nameParts) = (*nameParts)[:0]
	return
}

func ParseFullName(fullName string) (parsedName ParsedName) {

	var nameParts []string
	var nameCommas []bool

	nick, fullName := findNickName(fullName)
	parsedName.Nick = strings.Join(nick, ",")            // remove and store nicknames
	nameParts, nameCommas = splitNameIntoParts(fullName) // split name to parts tracking commas

	if len(nameParts) > 1 {
		parsedName.Suffix = strings.Join(findSuffixes(&nameParts, &nameCommas), ", ") // remove and store suffixes
	}

	if len(nameParts) > 1 {
		parsedName.Title = strings.Join(findTitles(&nameParts, &nameCommas), ", ") // remove and store titles
	}

	if len(nameParts) > 1 {
		joinPrefixes(&nameParts, &nameCommas) // Join name prefixes to following names
	}

	if len(nameParts) > 1 {
		joinConjunctions(&nameParts, &nameCommas) // Join conjunctions to surrounding names
	}

	if len(nameParts) > 1 {
		extraSuffixes := findExtraSuffixes(&nameParts, &nameCommas) // stuff after commas is assumed to be additonal suffixes
		if len(extraSuffixes) > 0 {
			if parsedName.Suffix != "" {
				parsedName.Suffix += ", " + strings.Join(extraSuffixes, ", ")
			} else {
				parsedName.Suffix = strings.Join(extraSuffixes, ", ")
			}
		}
	}

	if len(nameParts) > 0 {
		parsedName.Last = extractLastName(&nameParts, &nameCommas) // remove last name and store it
	}

	if len(nameParts) > 0 {
		parsedName.First = extrctFirstName(&nameParts, &nameCommas) // remove first name and store it
	}

	if len(nameParts) > 0 {
		parsedName.Middle = extractMiddleName(&nameParts, &nameCommas) // Use all remaining parts as middle name
	}

	return
}

const db1 = false

// removeAt removes from `slice` the item at postion `pos`.  If pos is out of range it returns the original `slice`.
func removeAt[T any](slice []T, pos int) []T {
	if pos < 0 {
		return slice
	} else if pos >= len(slice) {
		return slice
	} else if pos == 0 {
		return slice[1:]
	} else if pos == len(slice)-1 {
		return slice[0:pos]
	}
	return append(slice[:pos], slice[pos+1:]...)
}

func locationInArray[T comparable](needle T, haystack []T) int {
	for ii, val := range haystack {
		if val == needle {
			return ii
		}
	}
	return -1
}
