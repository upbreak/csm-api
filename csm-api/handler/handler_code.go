package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
	"strconv"
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

// func: 코드 조회
// @param
// - p_code: 부모 코드
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

// func: 코드트리 조회
// @param
// -
func (h *HandlerCode) ListCodeTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pCode := r.URL.Query().Get("p_code")

	codeTrees, err := h.Service.GetCodeTree(ctx, pCode)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		CodeTrees entity.CodeTrees `json:"code_trees"`
	}{CodeTrees: *codeTrees}

	SuccessValuesResponse(ctx, w, values)
}

// func: 코드 수정 및 생성
// @param
// - code
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

// func: 코드 삭제
// @param
// - idx: 삭제할 코드 pk
func (h *HandlerCode) Remove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idx, err := strconv.ParseInt(r.PathValue("idx"), 10, 64)
	if err != nil {
		BadRequestResponse(ctx, w)
		return
	}

	err = h.Service.RemoveCode(ctx, idx)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}
	SuccessResponse(ctx, w)
}

// func: 코드순서 변경
// @param
// - codeSorts
func (h *HandlerCode) SortNoModify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 데이터 파싱
	codeSorts := entity.CodeSorts{}
	if err := json.NewDecoder(r.Body).Decode(&codeSorts); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	err := h.Service.ModifySortNo(ctx, codeSorts)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)

}

// func: 코드 중복 검사
// @param
// - code
func (h *HandlerCode) DuplicateByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//데이터 파싱
	code := r.URL.Query().Get("code")
	value, err := h.Service.DuplicateCheckCode(ctx, code)

	if err != nil {
		FailResponse(ctx, w, err)
		return
	}
	SuccessValuesResponse(ctx, w, value)
}
