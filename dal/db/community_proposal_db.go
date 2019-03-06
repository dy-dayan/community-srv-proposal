package db

import (
	"errors"
	"github.com/dy-dayan/community-srv-proposal/util/config"
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
		Addrs:     uconfig.DefaultMgoConf.Addr,
		Direct:    false,
		Timeout:   time.Second * 3,
		PoolLimit: uconfig.DefaultMgoConf.PoolLimit,
		Username:  uconfig.DefaultMgoConf.Username,
		Password:  uconfig.DefaultMgoConf.Password,
	}

	ses, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logrus.Fatalf("dail mgo server error:%v", err)
	}

	ses.SetMode(mgo.Monotonic, true)
	defaultMgo = ses
}

type Proposal struct {
	ID        bson.ObjectId  `bson:"_id"`
	Title     string  `bson:"title"`
	Content   string `bson:"content"`
	Imgs      []string `bson:"imgs"`
	Options   []string `bson:"options"`
	Creator   int64 `bson:"creator"`
	State     int32 `bson:"state"`
	Result    string `bson:"result"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

var (
	DBProposal = "dayan_community_proposal"
	CProposal  = "proposal"
)

const (
	ProposalState_Unknown = iota // 未知
	ProposalState_Ediion  // 编辑中
	ProposalState_Commited // 已提交，投票中
	ProposalState_End // 结束
)

func InsertProposal(title, content string, creator int64, imgs, opts []string) error {
	now := time.Now().Unix()
	data := &Proposal{
		Title:     title,
		Content:   content,
		Imgs:      imgs,
		Options:   opts,
		Creator:   creator,
		State:     ProposalState_Ediion,
		Result:    "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	ses := defaultMgo.Copy()
	if ses == nil {
		return errors.New("mgo session is nil")
	}
	defer ses.Close()

	return ses.DB(DBProposal).C(CProposal).Insert(data)
}



func GetProposal(id string) error {
	query := bson.M{
		"_id":bson.ObjectIdHex(id),
	}
	ses := defaultMgo.Copy()
	if ses == nil {
		return errors.New("mgo session is nil")
	}
	defer ses.Close()

	return ses.DB(DBProposal).C(CProposal).Insert(data)
}
