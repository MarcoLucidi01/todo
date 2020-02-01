package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

const FILENAME string = ".todo" // stored in user's home directory

const (
	COMPLETEPREFIX   string = "[x] "
	INCOMPLETEPREFIX string = "[ ] "
)

type Action int // an operation that can be performed on the list
const (
	ADD             Action = iota // desc...
	EDIT                          // -e id desc...
	MARKCOMPLETE                  // -c id...
	MARKINCOMPLETE                // -i id...
	PRINTALL                      // -c
	PRINTINCOMPLETE               // default action
	REMOVE                        // -r id...
	REMOVECOMPLETE                // -r
	REPLACE                       // -e id /sub/rep/
	SWAP                          // -s id id
)

type Args struct { // arguments from command line
	action   Action // action to be performed
	ids      []int  // parsed ids
	desc     string // new todo description
	sub, rep string // substring and replacement for REPLACE action
}

type Todo struct { // todo list item
	desc     string // todo description
	complete bool   // whether the todo was completed or not
}

var args Args   // parsed arguments
var list []Todo // in memory todo list

func main() {
	parseArgs()
	loadList()
	checkIds()
	doAction()
	if args.action != PRINTALL && args.action != PRINTINCOMPLETE {
		printIncomplete()
		saveList()
	}
}

func parseArgs() {
	flag.Usage = printUsage
	c := flag.Bool("c", false, "")
	i := flag.Bool("i", false, "")
	e := flag.Bool("e", false, "")
	s := flag.Bool("s", false, "")
	r := flag.Bool("r", false, "")
	flag.Parse()

	if *c {
		if flag.Arg(0) == "" {
			args.action = PRINTALL
		} else {
			args.action = MARKCOMPLETE
			args.ids = parseIds(flag.Args(), -1)
		}

	} else if *i {
		args.action = MARKINCOMPLETE
		args.ids = parseIds(flag.Args(), -1)

	} else if *e {
		args.action = EDIT
		args.ids = parseIds(flag.Args(), 1)
		args.desc = parseDesc(flag.Args()[1:])

		split := strings.Split(args.desc, "/")
		if len(split) == 4 && split[0] == "" && split[3] == "" { // /sub/rep/
			args.action = REPLACE
			args.sub = split[1]
			args.rep = split[2]
		}

	} else if *s {
		args.action = SWAP
		args.ids = parseIds(flag.Args(), 2)

	} else if *r {
		if flag.Arg(0) == "" {
			args.action = REMOVECOMPLETE
		} else {
			args.action = REMOVE
			args.ids = parseIds(flag.Args(), -1)
		}

	} else if flag.Arg(0) == "" {
		args.action = PRINTINCOMPLETE

	} else {
		args.action = ADD
		args.desc = parseDesc(flag.Args())
	}
}

func printUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), `usage: %v [-c|-i|-e|-s|-r|-h] [id...] [desc...]
command line todo list
todos are stored at %v

  desc...          add new todo
  -c               print also completed todos
  -c id...         mark specified todos as complete
  -i id...         mark specified todos as incomplete
  -e id desc...    edit description of specified todo
  -e id /sub/rep/  replace substring sub with rep in description of specified todo
  -s id id         swap position of specified todos
  -r               remove completed todos
  -r id...         remove specified todos
  -h               show usage message

repo: https://github.com/MarcoLucidi01/todo
`, os.Args[0], filepath())
}

func filepath() string {
	usr, err := user.Current()
	if err != nil {
		die("unable to get user's home directory")
	}
	return usr.HomeDir + "/" + FILENAME
}

func parseIds(ss []string, expected int) []int {
	var ids []int
	for i, s := range ss {
		if expected >= 0 && i == expected {
			break
		}
		id, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			die("invalid id \"%v\"", s)
		}
		ids = append(ids, int(id))
	}
	if expected >= 0 && len(ids) != expected {
		die("expected %v ids but got %v", expected, len(ids))
	}
	return ids
}

func parseDesc(ss []string) string {
	if len(ss) == 0 {
		die("missing description")
	}
	return strings.Join(ss, " ")
}

func loadList() {
	f, err := os.OpenFile(filepath(), os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		die(err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		desc := scanner.Text()
		// lines without PREFIX are assumed to be incomplete todos
		complete := strings.HasPrefix(desc, COMPLETEPREFIX)
		if complete || strings.HasPrefix(desc, INCOMPLETEPREFIX) {
			desc = desc[len(COMPLETEPREFIX):] // both PREFIX have same len
		}
		list = append(list, Todo{desc, complete})
	}
}

func checkIds() {
	max := len(list) - 1
	for _, id := range args.ids {
		if id < 0 || id > max {
			die("invalid id %v", id)
		}
	}
}

func doAction() {
	switch args.action {
	case ADD:
		add()
	case EDIT:
		edit()
	case MARKCOMPLETE:
		markComplete()
	case MARKINCOMPLETE:
		markIncomplete()
	case PRINTALL:
		printAll()
	case PRINTINCOMPLETE:
		printIncomplete()
	case REMOVE:
		remove()
	case REMOVECOMPLETE:
		removeComplete()
	case REPLACE:
		replace()
	case SWAP:
		swap()
	default:
		die("invalid action %v", args.action) // not reached
	}
}

func add() {
	list = append(list, Todo{args.desc, false})
}

func edit() {
	list[args.ids[0]].desc = args.desc
}

func markComplete() {
	for _, id := range args.ids {
		list[id].complete = true
	}
}

func markIncomplete() {
	for _, id := range args.ids {
		list[id].complete = false
	}
}

func printAll() {
	for id, todo := range list {
		prefix := INCOMPLETEPREFIX
		if todo.complete {
			prefix = COMPLETEPREFIX
		}
		fmt.Printf("%v %v%v\n", id, prefix, todo.desc)
	}
}

func printIncomplete() {
	for id, todo := range list {
		if !todo.complete {
			fmt.Printf("%v %v%v\n", id, INCOMPLETEPREFIX, todo.desc)
		}
	}
}

func remove() {
	removeIf(func(id int, todo Todo) bool {
		for _, aid := range args.ids {
			if id == aid {
				return todo.complete || askYesNo("todo %v is incomplete. Remove it?", id)
			}
		}
		return false
	})
}

func removeComplete() {
	if askYesNo("remove all completed todos?") {
		removeIf(func(id int, todo Todo) bool {
			return todo.complete
		})
	}
}

func removeIf(predicate func(int, Todo) bool) {
	var newList []Todo
	for id, todo := range list {
		if !predicate(id, todo) {
			newList = append(newList, todo)
		}
	}
	list = newList
}

func askYesNo(question string, a ...interface{}) bool {
	fmt.Printf(question, a...)
	fmt.Print(" [y/N]: ")
	var ans string
	_, err := fmt.Scanln(&ans)
	return err == nil && (strings.EqualFold(ans, "y") || strings.EqualFold(ans, "yes"))
}

func replace() {
	id := args.ids[0]
	list[id].desc = strings.Replace(list[id].desc, args.sub, args.rep, -1)
}

func swap() {
	id1, id2 := args.ids[0], args.ids[1]
	tmp := list[id1]
	list[id1] = list[id2]
	list[id2] = tmp
}

func saveList() {
	f, err := os.Create(filepath())
	if err != nil {
		die(err.Error())
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, todo := range list {
		if todo.complete {
			w.WriteString(COMPLETEPREFIX)
		} else {
			w.WriteString(INCOMPLETEPREFIX)
		}
		w.WriteString(todo.desc)
		w.WriteByte('\n')
	}
	w.Flush()
}

func die(reason string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, reason, a...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}
