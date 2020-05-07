package vars

// Actions & Status
const (
	EventStatusOpen     = "event_open"
	EventStatusClosed   = "event_closed"
	TodoStatusOpen      = "todo_open"
	TodoStatusPostponed = "todo_postponed"
	TodoStatusCanceled  = "todo_canceled"
	TodoStatusFinished  = "todo_finished"
	TodoStatusOverdue   = "todo_overdue"

	ActionCreateTodo = "create_todo" // +1
	ActionFinishTodo = "finish_todo" // +1
	ActionPushTodo   = "push_todo"   // -1
	ActionCopyTodo   = "copy_todo"   // +0
	ActionPullTodo   = "pull_todo"   // +0
	ActionCancelTodo = "cancel_todo" // -2
	ActionRemoveTodo = "remove_todo" // -10
	ActionRenewTodo  = "renew_todo"  // -10
	ActionEditTodo   = "edit_todo"   // +0
)