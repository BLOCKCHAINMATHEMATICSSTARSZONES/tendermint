package types

import (
	"time"

	"github.com/tendermint/tendermint/libs/bytes"
	tmtime "github.com/tendermint/tendermint/types/time"
)

// Canonical* wraps the structs in types for amino encoding them for use in SignBytes / the Signable interface.

// TimeFormat is used for generating the sigs
const TimeFormat = time.RFC3339Nano

type CanonicalBlockID struct {
	Hash        bytes.HexBytes
	PartsHeader CanonicalPartSetHeader
}

type CanonicalPartSetHeader struct {
	Hash  bytes.HexBytes
	Total int
}

type CanonicalProposal struct {
	Type      SignedMsgType // type alias for byte
	Height    int64         `binary:"fixed64"`
	Round     int64         `binary:"fixed64"`
	POLRound  int64         `binary:"fixed64"`
	BlockID   CanonicalBlockID
	Timestamp time.Time
	ChainID   string
	Data      bytes.HexBytes // [peppermint] add data in proposal to that signers can sign
}

type CanonicalVote struct {
	Type      SignedMsgType // type alias for byte
	Height    int64         `binary:"fixed64"`
	Round     int64         `binary:"fixed64"`
	BlockID   CanonicalBlockID
	Timestamp time.Time
	ChainID   string
}

// [peppermint] create RLP vote to decode in contract
type CanonicalRLPVote struct {
	ChainID string
	Type    byte
	Height  uint
	Round   uint
	Data    []byte
}

//-----------------------------------
// Canonicalize the structs

func CanonicalizeBlockID(blockID BlockID) CanonicalBlockID {
	return CanonicalBlockID{
		Hash:        blockID.Hash,
		PartsHeader: CanonicalizePartSetHeader(blockID.PartsHeader),
	}
}

func CanonicalizePartSetHeader(psh PartSetHeader) CanonicalPartSetHeader {
	return CanonicalPartSetHeader{
		psh.Hash,
		psh.Total,
	}
}

func CanonicalizeProposal(chainID string, proposal *Proposal) CanonicalProposal {
	return CanonicalProposal{
		Type:      ProposalType,
		Height:    proposal.Height,
		Round:     int64(proposal.Round), // cast int->int64 to make amino encode it fixed64 (does not work for int)
		POLRound:  int64(proposal.POLRound),
		BlockID:   CanonicalizeBlockID(proposal.BlockID),
		Timestamp: proposal.Timestamp,
		ChainID:   chainID,
	}
}

func CanonicalizeVote(chainID string, vote *Vote) CanonicalRLPVote {
	// return CanonicalVote{
	// 	Type:      vote.Type,
	// 	Height:    vote.Height,
	// 	Round:     int64(vote.Round), // cast int->int64 to make amino encode it fixed64 (does not work for int)
	// 	BlockID:   CanonicalizeBlockID(vote.BlockID),
	// 	Timestamp: vote.Timestamp,
	// 	ChainID:   chainID,
	// }
	// TODO ensure that removing Timestamp and BlockID has no security issues
	return CanonicalRLPVote{
		ChainID: chainID,
		Type:    byte(vote.Type),
		Height:  uint(vote.Height),
		Round:   uint(vote.Round),
		Data:    vote.Data,
	}
}

// CanonicalTime can be used to stringify time in a canonical way.
func CanonicalTime(t time.Time) string {
	// Note that sending time over amino resets it to
	// local time, we need to force UTC here, so the
	// signatures match
	return tmtime.Canonical(t).Format(TimeFormat)
}
