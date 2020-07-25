package uci

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Engine struct {
	stdin    *bufio.Writer
	stdout   *bufio.Scanner
	Meta     Meta
	Side     int
	StartPos string
}

// Meta data about the Engine
type Meta struct {
	Name    string   // Name of Engine
	Author  string   // Author of the Engine
	Options []Option // Available options for this Engine
}

// Available options to set on this Engine
type Option struct {
	Name    string      // Name of option
	Type    string      // Type of option
	Default interface{} // Default value of option
	Min     int         // Min value of option
	Max     int         // Max value of option
	Vars    []string    // Enum vars for option
}

// Options for creating a new game
type NewGameOpts struct {
	Variant struct {
		Key string
	}
	InitialFen string
	Moves string
	Side     int    // Which side should the Engine play as. Must be uci.W or uci.B
}

// Options to pass when looking for best move
type GoOpts struct {
	SearchMoves string // <move1> .... <movei>. restrict search to this moves only
	Ponder      bool   // start searching in pondering mode
	Wtime       int    // number of ms white has left
	Btime       int    // number of ms black has left
	Winc        int    // number of ms white increases by each move
	Binc        int    // number of ms black increases by each move
	MovesToGo   int    // number of moves until next time control
	Depth       int    // maximum search depth
	Nodes       int    // maximum search nodes
	Mate        int    // search for mate in x moves
	MoveTime    int    // number of ms exactly to search for
}

// Response from searching
type GoResp struct {
	Bestmove string // Best move the Engine could find
	Ponder   string // Ponder move
}

const (
	White int = 0 // Side to play as. White
	Black int = 1 // Side to play as. Black
)

var execCommand = exec.Command

// Create a new Engine. Requires the path to a uci chess Engine such as stockfish
func NewEngine(path string) (*Engine, error) {
	eng := Engine{}
	cmd := execCommand(path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	eng.stdin = bufio.NewWriter(stdin)
	eng.stdout = bufio.NewScanner(stdout)
	eng.Meta = eng.getMeta()
	return &eng, nil
}

func (eng *Engine) getMeta() (meta Meta) {
	eng.send("uci")
	lines := eng.receive("uciok")

	namePrefix := "id name "
	authorPrefix := "id author "
	optionPrefix := "option "
	for _, line := range lines {
		if strings.HasPrefix(line, namePrefix) {
			meta.Name = strings.TrimPrefix(line, namePrefix)
		} else if strings.HasPrefix(line, authorPrefix) {
			meta.Author = strings.TrimPrefix(line, authorPrefix)
		} else if strings.HasPrefix(line, optionPrefix) {
			meta.Options = append(meta.Options, newOption(strings.TrimPrefix(line, optionPrefix)))
		}
	}
	return meta
}

func getOption(line, regex string) interface{} {
	rr := regexp.MustCompile(regex)
	results := rr.FindStringSubmatch(line)

	if len(results) == 2 {
		return results[1]
	}
	return nil
}

func newOption(line string) (option Option) {
	option.Name, _ = getOption(line, `name (.*) type`).(string)
	option.Type, _ = getOption(line, `type (\w+)`).(string)
	option.Default = getOption(line, `default (\w+)`)
	minStr, _ := getOption(line, `min (\w+)`).(string)
	maxStr, _ := getOption(line, `max (\w+)`).(string)
	option.Min, _ = strconv.Atoi(minStr)
	option.Max, _ = strconv.Atoi(maxStr)

	varRegex := regexp.MustCompile(`var (\w+)`)

	vars := []string{}
	for _, v := range varRegex.FindAllStringSubmatch(line, -1) {
		vars = append(vars, v[1])
	}
	option.Vars = vars

	return
}

// Pass an option to the underlying Engine
func (eng *Engine) SetOption(name string, value interface{}) bool {
	for _, option := range eng.Meta.Options {
		if option.Name == name {
			var v string
			switch value.(type) {
			case string:
				v, _ = value.(string)
			case int:
				vv, _ := value.(int)
				if (vv < option.Min) {
					vv = option.Min
				} else if (vv > option.Max) {
					vv = option.Max
				}
				v = strconv.Itoa(vv)
			case bool:
				vv, _ := value.(bool)
				if vv {
					v = "true"
				} else {
					v = "false"
				}
			}
			eng.send("setoption name " + name + " value " + v)
			return true
		}
	}
	return false
}

// Check if Engine is ready to start receiving commands
func (eng *Engine) IsReady() bool {
	eng.send("isready")
	lines := eng.receive("readyok")
	return lines[0] == "readyok"
}

func (eng *Engine) send(input string) {
	_, err := eng.stdin.WriteString(input + "\n")
	if err == nil {
		eng.stdin.Flush()
	}
}

func (eng *Engine) receive(stopPrefix string) (lines []string) {
	scanner := eng.stdout
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if strings.HasPrefix(line, stopPrefix) {
			break
		}
	}
	if err := eng.stdout.Err(); err != nil {
		log.Fatal("reading standard input:", err)
	}
	return
}

// Start a new game. Only one game should be played at a time
func (eng *Engine) NewGame(opts NewGameOpts) {
	if opts.Variant.Key == "chess960" || opts.Variant.Key == "fromPosition" {
		eng.SetOption("UCI_Variant", "chess")
		eng.SetOption("UCI_Chess960", "true")
	} else {
		if opts.Variant.Key == "threeCheck" {
			eng.SetOption("UCI_Variant", "3check")
		} else if opts.Variant.Key == "standard" {
			eng.SetOption("UCI_Variant", "chess")
		} else {
			eng.SetOption("UCI_Variant", strings.ToLower(opts.Variant.Key))
		}
		eng.SetOption("UCI_Chess960", "false")
	}
	if opts.InitialFen == "" || opts.InitialFen == "startpos" {
		eng.StartPos = "position startpos moves "
	} else {
		eng.StartPos = "position fen " + opts.InitialFen + " moves "
	}
	eng.Position(opts.Moves)
	eng.Side = opts.Side
}

// Set the position of the game after the initial start position. Algrebriac notiation, e.g `e2e4 e7e6`
func (eng *Engine) Position(moves string) {
	eng.send(eng.StartPos + moves)
}

func addOpt(name string, value int) string {
	if value > 0 {
		return name + " " + strconv.Itoa(value) + " "
	}
	return ""
}

// Search for the bestmove
func (eng *Engine) Go(opts GoOpts) GoResp {
	goCmd := "go "
	if opts.Ponder {
		goCmd += "ponder "
	}
	goCmd += addOpt("wtime", opts.Wtime)
	goCmd += addOpt("btime", opts.Btime)
	goCmd += addOpt("winc", opts.Winc)
	goCmd += addOpt("binc", opts.Binc)
	goCmd += addOpt("movestogo", opts.MovesToGo)
	goCmd += addOpt("depth", opts.Depth)
	goCmd += addOpt("nodes", opts.Nodes)
	goCmd += addOpt("mate", opts.Mate)
	goCmd += addOpt("movetime", opts.MoveTime)

	eng.send(goCmd)
	lines := eng.receive("bestmove")
	words := strings.Split(lines[len(lines)-1], " ")

	bestmove := ""
	ponder := ""
	if len(words) >= 2 {
		bestmove = words[1]
	}
	if len(words) >= 4 {
		ponder = words[3]
	}

	return GoResp{
		Bestmove: bestmove,
		Ponder:   ponder,
	}
}

// Quit the Engine. Engine struct cannot be used after this command has been sent
func (eng *Engine) Quit() {
	eng.send("quit")
	eng.stdin = nil
	eng.stdout = nil
	eng.Meta = Meta{}
	eng.Side = 0
	eng.StartPos = ""
}
