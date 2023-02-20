package item

// Tags/Statuses could also be bit flags. This would make them easier to keep track of and avoids duplicates.

// Tag represents a label attached to an item for organizational purposes.
type Tag string

// These are tags I'm currently using to organize my links. Subject to change.
const (
	// Advice is for links that talk about life advice, philosophy, mentality, etc.
	Advice Tag = "advice"

	// Money is for investment advice and personal finance tips.
	Money Tag = "money"

	// Career is for anything related to finding, keeping, and succeeding in a job.
	Career Tag = "career"

	// Hobby is for anything fun.
	Hobby Tag = "hobby"

	// Project is for project ideas.
	Project Tag = "project"

	// Opinion is for opinion pieces on any topic.
	Opinion Tag = "opinion"

	// Reference is for any reference material like books, courses, and videos.
	Reference Tag = "reference"

	// Skill is for any useful skills, for example "being focused" or "drawing".
	Skill Tag = "skill"

	// Tip is for small but useful ideas or routines to accomplish some task (e.g. "method for cleaning glasses" or "morning routine").
	Tip Tag = "tip"

	// Present but unused.
	Favorite Tag = "favorite"
	Upvoted  Tag = "upvoted"
)
