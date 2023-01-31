package models

type LiteralBin byte

const (
	B_null LiteralBin = iota
	B_a
	B_activate
	B_an
	B_and
	B_any
	B_are
	B_at
	B_blue
	B_breakpoints
	B_by
	B_change
	B_cloning
	B_close_par
	B_comma
	B_conclude
	B_create
	B_deactivate
	B_degrees
	B_delete
	B_different
	B_enabled
	B_end
	B_equal
	B_equal_sym
	B_first_task
	B_focus
	B_for
	B_from
	B_greater
	B_green
	B_halt
	B_hide
	B_if
	B_inform
	B_initially
	B_insert
	B_instance
	B_invoke
	B_is
	B_less
	B_move
	B_named
	B_of
	B_on
	B_open_par
	B_operator
	B_or
	B_parent
	B_red
	B_remove
	B_rotate
	B_set
	B_show
	B_start
	B_task_queue
	B_than
	B_that
	B_the
	B_then
	B_to
	B_transfer
	B_unconditionally
	B_when
	B_whenever
	B_whose
	B_with
	B_yellow
)

var LiteralBinStr = map[string]LiteralBin{
	"":                B_null,
	"a":               B_a,
	"activate":        B_activate,
	"an":              B_an,
	"and":             B_and,
	"any":             B_any,
	"are":             B_are,
	"at":              B_at,
	"blue":            B_blue,
	"breakpoints":     B_breakpoints,
	"by":              B_by,
	"change":          B_change,
	"cloning":         B_cloning,
	")":               B_close_par,
	",":               B_comma,
	"conclude":        B_conclude,
	"create":          B_create,
	"deactivate":      B_deactivate,
	"degrees":         B_degrees,
	"delete":          B_delete,
	"different":       B_different,
	"enabled":         B_enabled,
	"end":             B_end,
	"equal":           B_equal,
	"=":               B_equal_sym,
	"first-task":      B_first_task,
	"focus":           B_focus,
	"for":             B_for,
	"from":            B_from,
	"greater":         B_greater,
	"green":           B_green,
	"halt":            B_halt,
	"hide":            B_hide,
	"if":              B_if,
	"inform":          B_inform,
	"initially":       B_initially,
	"insert":          B_insert,
	"instance":        B_instance,
	"invoke":          B_invoke,
	"is":              B_is,
	"less":            B_less,
	"move":            B_move,
	"named":           B_named,
	"of":              B_of,
	"on":              B_on,
	"(":               B_open_par,
	"operator":        B_operator,
	"or":              B_or,
	"parent":          B_parent,
	"red":             B_red,
	"remove":          B_remove,
	"rotate":          B_rotate,
	"set":             B_set,
	"show":            B_show,
	"start":           B_start,
	"task-queue":      B_task_queue,
	"than":            B_than,
	"that":            B_that,
	"the":             B_the,
	"then":            B_then,
	"to":              B_to,
	"transfer":        B_transfer,
	"unconditionally": B_unconditionally,
	"when":            B_when,
	"whenever":        B_whenever,
	"whose":           B_whose,
	"with":            B_with,
	"yellow":          B_yellow,
}

func (me LiteralBin) String() string {
	return string(me)
}

func (me LiteralBin) Size() int {
	return len(LiteralBinStr)
}
