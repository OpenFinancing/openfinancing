package solar

// General Concept:
// This can be applied as a feedback system to actors that are part of projects (i.e. relate to Entities.go, such as contractors).
// It allows investors and recipients to give comment son the services provided by these entities in the project.
// TODO: build and improve this functionality and consider a reputation system to score services.

type Feedback struct {
	Content string
	// the content of the feedback, good / bad
	// maybe we could have a  rating system baked in? a star based rating system?
	// would be nice, idk
	From Entity
	// who gave the feedback?
	To Entity
	// regarding whom is this feedback about
	Date string
	// time at which this feedback was written
	RelatedContract []Project
	// the contract regarding which this feedback is directed at
}
