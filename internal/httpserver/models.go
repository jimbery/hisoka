package httpserver

type voteType string

const (
	DubVote voteType = "dub"
	SubVote voteType = "sub"
)

type voteInput struct {
	MalID    int      `json:"mal_id"`
	VoteType voteType `json:"vote_type"`
}
