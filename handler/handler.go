package handler

import (
	"context"
	"github.com/dy-dayan/community-srv-proposal/dal/db"
	"github.com/dy-gopkg/kit"
	"github.com/sirupsen/logrus"
	srv "github.com/dy-dayan/community-srv-proposal/idl/dayan/community/srv-proposal"
	atomicid "github.com/dy-dayan/community-srv-proposal/idl/dayan/common/srv-atomicid"
	"github.com/dy-dayan/community-srv-proposal/idl"
	"gopkg.in/mgo.v2"
)

type Handler struct {
}

// 新建议案
func (h *Handler) NewProposal(ctx context.Context, req *srv.NewProposalReq, rsp *srv.NewProposalResp) error {
	rsp.BaseResp = &base.Resp{
		Code: int32(base.CODE_OK),
	}
	//检测入参
	if len(req.Title) == 0 || len(req.Content) == 0 || len(req.Options) == 0 || req.CommunityID < 0{
		logrus.Debug("invalid parameters")
		rsp.BaseResp.Code = int32(base.CODE_INVALID_PARAMETER)
		rsp.BaseResp.Msg = "invalid parameters"
		return nil
	}

	// 请求获取一个proposal id
	cl := atomicid.NewAtomicIDService("dayan.common.srv.atomicid", kit.Client())
	req1 := &atomicid.GetIDReq{Label: "dayan.community.srv.proposal.proposal_id"}
	rsp1, err := cl.GetID(ctx, req1)
	if err != nil {
		logrus.Errorf("atomicid.GetID error:%v", err)
		return err
	}

	if rsp1.BaseResp.Code != int32(base.CODE_OK) {
		logrus.Warnf("atomicid.GetID resp code:%v, msg:%s", rsp1.BaseResp.Code, rsp1.BaseResp.Msg)
		rsp.BaseResp = rsp1.BaseResp
		return nil
	}

	err = db.InsertProposal(rsp1.Id, req.CommunityID, req.Title, req.Content, req.Creator, req.Imgs, req.Options)
	if err != nil{
		logrus.Warnf("db.InsertProposal error:%v", err)
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
		return nil
	}

	rsp.ProposalID = rsp1.Id

	return nil
}

func (h *Handler) GetProposalInfo(ctx context.Context, req *srv.GetProposalInfoReq, rsp *srv.GetProposalInfoResp) error {
	rsp.BaseResp = &base.Resp{
		Code: int32(base.CODE_OK),
	}

	if req.ProposalID < 0 {
		logrus.Debug("invalid parameters")
		rsp.BaseResp.Code = int32(base.CODE_INVALID_PARAMETER)
		rsp.BaseResp.Msg = "invalid parameters"
		return nil
	}

	p, err := db.GetProposal(req.ProposalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			logrus.Warnf("no such proposal:%d", req.ProposalID)
			rsp.BaseResp.Code = int32(base.CODE_INVALID_PARAMETER)
			rsp.BaseResp.Msg = "no such proposal"
			return nil
		}

		logrus.Warnf("db.GetProposal error:%v", err)
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
		return nil
	}

	rsp.ProposalID = req.ProposalID
	rsp.Title = p.Title
	rsp.Content = p.Content
	rsp.Creator = p.Creator
	rsp.CreatedAt = p.CreatedAt
	rsp.UpdatedAt = p.UpdatedAt

	res, err := db.GetProposalVoteInfo(req.ProposalID, p.Options)
	if err != nil {
		logrus.Warnf("db.GetProposalVoteInfo error:%v", err)
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
		return nil
	}
	rsp.VoteInfo = res

	return nil
}

// 只能修改议案标题和内容
func (h *Handler) ModifyProposal(ctx context.Context, req *srv.ModifyProposalReq, rsp *srv.ModifyProposalResp) error {
	rsp.BaseResp = &base.Resp{
		Code: int32(base.CODE_OK),
	}

	if len(req.Title) == 0 && len(req.Content) == 0 {
		logrus.Debug("invalid parameters")
		rsp.BaseResp.Code = int32(base.CODE_INVALID_PARAMETER)
		rsp.BaseResp.Msg = "invalid parameters"
		return nil
	}

	p, err := db.GetProposal(req.ProposalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			logrus.Warnf("no such proposal:%d", req.ProposalID)
			rsp.BaseResp.Code = int32(base.CODE_INVALID_PARAMETER)
			rsp.BaseResp.Msg = "no such proposal"
			return nil
		}

		logrus.Warnf("db.GetProposal error:%v", err)
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
		return nil
	}

	// 只更新标题和内容
	err = db.UpdateProposalTitleAndContent(p.ID, req.Title, req.Content)
	if err != nil {
		logrus.Warnf("db.UpdateProposalTitleAndContent error:%v", err)
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
		return nil
	}

	return nil
}

func (h *Handler) CommentProposal(ctx context.Context, req *srv.CommentProposalReq, rsp *srv.CommentProposalResp) error {
	rsp.BaseResp = &base.Resp{
		Code: int32(base.CODE_OK),
	}


	return nil
}


func (h *Handler) GetComments(ctx context.Context, req *srv.GetCommentsReq, rsp *srv.GetCommentsResp) error {


	return nil
}

func (h *Handler) VoteProposal(ctx context.Context, req *srv.VoteProposalReq, rsp *srv.VoteProposalResp) error {
	rsp.BaseResp = &base.Resp{
		Code: int32(base.CODE_OK),
	}

	// 查询议案
	p, err := db.GetProposal(req.ProposalID)
	if err != nil {
		if err == mgo.ErrNotFound {
			logrus.Warnf("no such proposal:%d", req.ProposalID)
			rsp.BaseResp.Code = int32(base.CODE_INVALID_PARAMETER)
			rsp.BaseResp.Msg = "no such proposal"
			return nil
		}

		logrus.Warnf("db.GetProposal error:%v", err)
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
		return nil
	}

	// 检测是否存在该选项
	var has bool = false
	for _, opt := range p.Options {
		if opt == req.Option {
			has = true
			break
		}
	}

	if has == false {
		rsp.BaseResp.Code = int32(base.CODE_INVALID_PARAMETER)
		rsp.BaseResp.Msg = "no such proposal option"
		return nil
	}

	// 更新投票记录
	err = db.UpsertProposalVoteLog(req.ProposalID, req.Voter, req.Option)
	if err != nil {
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
	}

	return nil
}

func (h *Handler) GetProposalList(ctx context.Context, req *srv.GetProposalListReq, rsp *srv.GetProposalListResp) error {
	rsp.BaseResp = &base.Resp{
		Code: int32(base.CODE_OK),
	}

	ret, err := db.GetProposalList(req.CommunityID, req.State, req.Page)
	if err != nil {
		logrus.Warnf("db.GetProposalList error:%v", err)
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
		return nil
	}

	for _, p := range ret {
		proposal := &srv.ProposalInfo{
			ID:                   p.ID,
			Title:                p.Title,
		}
		rsp.Proposals = append(rsp.Proposals, proposal)
	}
	return nil
}