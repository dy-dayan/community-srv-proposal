package handler

import (
	"context"
	"github.com/dy-dayan/community-srv-proposal/dal/db"
	"github.com/sirupsen/logrus"
	srv "github.com/dy-dayan/community-srv-proposal/idl/dayan/community/srv-proposal"
	"github.com/dy-dayan/community-srv-proposal/idl"
)

type Handler struct {
}

// 新建议案
func (h *Handler) NewProposal(ctx context.Context, req *srv.NewProposalReq, rsp *srv.NewProposalResp) error {
	rsp.BaseResp = &base.Resp{
		Code: int32(base.CODE_OK),
	}
	if len(req.Title) == 0 || len(req.Content) == 0 || len(req.Options) == 0 {
		logrus.Debug("invalid parameters")
		rsp.BaseResp.Code = int32(base.CODE_INVALID_PARAMETER)
		rsp.BaseResp.Msg = "invalid parameters"
		return nil
	}

	err := db.InsertProposal(req.Title, req.Content, req.Creator, req.Imgs, req.Options)
	if err != nil{
		logrus.Warnf("db.InsertProposal error:%v", err)
		rsp.BaseResp.Code = int32(base.CODE_DATA_EXCEPTION)
		rsp.BaseResp.Msg = err.Error()
		return nil
	}

	return nil
}

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


	return nil
}