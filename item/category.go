package item

// Categories could also be bit flags. This would make them easier to keep track of and avoids duplicates.

// These are strings I'm currently using to organize my links. Subject to change.
const (
	// Advice is for links that talk about life advice, philosophy, mentality, etc.
	Advice string = "advice"

	// Career is for anything related to finding, keeping, and succeeding in a job.
	Career string = "career"

	// Hobby is for anything fun.
	Hobby string = "hobby"

	// Money is for investment advice and personal finance tips.
	Money string = "money"

	// Opinion is for opinion pieces on any topic.
	Opinion string = "opinion"

	// Product is for any miscellaneous software/hardware.
	Product string = "product"

	// Project is for project ideas.
	Project string = "project"

	// Reference is for any reference material like books, courses, and videos.
	Reference string = "reference"

	// Skill is for any useful skills, for example "being focused" or "drawing".
	Skill string = "skill"

	// Tip is for small but useful ideas or routines to accomplish some task (e.g. "method for cleaning glasses" or "morning routine").
	Tip string = "tip"
)
