package item

// Status represents an action taken on an item.
type Status string

const (
	Read  Status = "read"  // Leave blank if unread.
	Skip         = "skip"  // For items not worth reading right now.
	Notes        = "notes" // Present if notes were taken on this item.
)
