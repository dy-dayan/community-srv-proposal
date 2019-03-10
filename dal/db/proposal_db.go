package db

import (
	"errors"
	"github.com/dy-dayan/community-srv-proposal/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"time"
)

var (
	defaultMgo *mgo.Session
)

func Mgo() *mgo.Session {
	return defaultMgo
}

func Init() {
	dialInfo := &mgo.DialInfo{
		Addrs:     util.DefaultMgoConf.Addr,
		Direct:    false,
		Timeout:   time.Second * 3,
		PoolLimit: util.DefaultMgoConf.PoolLimit,
		Username:  util.DefaultMgoConf.Username,
		Password:  util.DefaultMgoConf.Password,
	}

	ses, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logrus.Fatalf("dail mgo server error:%v", err)
	}

	ses.SetMode(mgo.Monotonic, true)
	defaultMgo = ses
}

type Proposal struct {
	ID          int64    `bson:"_id"`
	CommunityID int64    `bson:"community_id"`
	Title       string   `bson:"title"`
	Content     string   `bson:"content"`
	Imgs        []string `bson:"imgs"`
	Options     []string `bson:"options"`
	Creator     int64    `bson:"creator"`
	State       int32    `bson:"state"`
	Result      string   `bson:"result"`
	CreatedAt   int64    `bson:"created_at"`
	UpdatedAt   int64    `bson:"updated_at"`
}

var (
	DBProposal = "dayan_community_proposal"
	CProposal  = "proposal"
)

const (
	ProposalState_Unknown  int32 = iota // 未知
	ProposalState_Ediion                // 编辑中
	ProposalState_Commited              // 已提交，投票中
	ProposalState_End                   // 结束
)

func InsertProposal(pid, cid int64, title, content string, creator int64, imgs, opts []string) error {
	now := time.Now().Unix()
	data := &Proposal{
		ID:          pid,
		CommunityID: cid,
		Title:       title,
		Content:     content,
		Imgs:        imgs,
		Options:     opts,
		Creator:     creator,
		State:       ProposalState_Ediion,
		Result:      "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ses := defaultMgo.Copy()
	if ses == nil {
		return errors.New("mgo session is nil")
	}
	defer ses.Close()

	return ses.DB(DBProposal).C(CProposal).Insert(data)
}

func GetProposal(id int64) (*Proposal, error) {
	query := bson.M{
		"_id": id,
	}
	ses := defaultMgo.Copy()
	if ses == nil {
		return nil, errors.New("mgo session is nil")
	}
	defer ses.Close()

	ret := &Proposal{}
	err := ses.DB(DBProposal).C(CProposal).Find(query).One(ret)
	if err != nil {
		return nil, err
	}

	return ret, err
}

func UpdateProposalTitleAndContent(id int64, title, content string) error {
	query := bson.M{
		"_id": id,
	}
	ses := defaultMgo.Copy()
	if ses == nil {
		return errors.New("mgo session is nil")
	}
	defer ses.Close()

	data := bson.M{}
	if len(title) != 0 {
		data["title"] = title
	}
	if len(content) != 0 {
		data["content"] = content
	}

	change := mgo.Change{
		Update: bson.M{
			"$set": data,
		},
		Upsert:    false,
		Remove:    false,
		ReturnNew: false,
	}

	_, err := ses.DB(DBProposal).C(CProposal).Find(query).Apply(change, nil)
	return err
}

func GetProposalList(cid int64, state, page int32) ([]*Proposal, error) {
	ses := defaultMgo.Copy()
	if ses == nil {
		return nil, errors.New("mgo session is nil")
	}
	defer ses.Close()

	limit := 30
	skip := int(page) * (limit)

	query := bson.M{
		"community_id": cid,
	}

	ret := []*Proposal{}
	err := ses.DB(DBProposal).C(CProposal).Find(query).Limit(limit).Skip(skip).All(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
