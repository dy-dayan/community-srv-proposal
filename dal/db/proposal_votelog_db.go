package db

import (
	"encoding/json"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CProposalVoteLog = "proposal_vote_log"
)

type ProposalVoteLog struct {
	ID bson.ObjectId `bson:"_id"`
	ProposalID int64 `bson:"proposal_id"`
	Voter int64 `bson:"voter"`
	Option string `bson:"option"`
	CreatedAt int64 `bson:"created_at"`
	UpdatedAt int64 `bson:"updated_at"`
}

func UpsertProposalVoteLog(pid int64, voter int64, opt string) error {
	ses := defaultMgo.Copy()
	if ses == nil {
		return errors.New("mgo session is nil")
	}
	defer ses.Close()

	now := time.Now().Unix()

	query := bson.M{
		"proposal_id":pid,
		"voter":voter,
	}

	data := bson.M{
		"$set": bson.M{
			"option": opt,
			"updated_at": now,
		},
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	}

	change := mgo.Change{
		Update:data,
		Upsert:true,
		Remove:    false,
		ReturnNew: false,
	}

	_, err := ses.DB(DBProposal).C(CProposalVoteLog).Find(query).Apply(change, nil)
	return err
}

func GetProposalVoteInfo(pid int64, opts []string) (result []byte, err error) {
	ses := defaultMgo.Copy()
	if ses == nil {
		return nil, errors.New("mgo session is nil")
	}
	defer ses.Close()

	info := bson.M{}
	for _, v := range opts {
		info[v] = 0
	}

	group := bson.M{
		"_id": "$option",
		"count": bson.M{
			"$sum":1,
		},
	}

	query := bson.M{
		"proposal_id":pid,
		"$group": group,
	}

	ret := []struct{
		ID    map[string]string `bson:"_id,omitempty"`
		Count int
	}{}
	err = ses.DB(DBProposal).C(CProposalVoteLog).Pipe(query).All(&ret)
	if err != nil {
		return nil, err
	}

	for _, v := range ret {
		field := v.ID["option"]
		info[field] = v.Count
	}

	result, err = json.Marshal(info)
	if err != nil {
		return nil, err
	}

	return result, err
}