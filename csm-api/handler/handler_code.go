package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-03-18
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct, func: 코드 조회
type HandlerCode struct {
	Service service.ServiceCode
}

func (h *HandlerCode) ListByPCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pCode := r.URL.Query().Get("p_code")
	if pCode == "" {
		BadRequestResponse(ctx, w)
		return
	}

	list, err := h.Service.GetCodeList(ctx, pCode)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.Codes `json:"list"`
	}{List: *list}

	SuccessValuesResponse(ctx, w, values)
}

func (h *HandlerCode) ListCodeTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	codeTrees, err := h.Service.GetCodeTree(ctx)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		CodeTrees entity.CodeTrees `json:"code_trees"`
	}{CodeTrees: *codeTrees}

	SuccessValuesResponse(ctx, w, values)
}

func (h *HandlerCode) Merge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 데이터 파싱
	code := entity.Code{}
	if err := json.NewDecoder(r.Body).Decode(&code); err != nil {
		FailResponse(ctx, w, err)
	}

	err := h.Service.MergeCode(ctx, code)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}
