package store

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"fmt"
)

// 사용자 권한 조회
// param: 사용자번호
func (r *Repository) GetUserRoleListByUno(ctx context.Context, db Queryer, uno int64) ([]entity.UserRoleMap, error) {
	var list []entity.UserRoleMap

	query := `
		SELECT 
			USER_UNO,
			ROLE_CODE,
			JNO
		FROM IRIS_USER_ROLE_MAP
		WHERE USER_UNO = :1
		AND JNO != 0`

	if err := db.SelectContext(ctx, &list, query, uno); err != nil {
		return nil, fmt.Errorf("GetUserRoleListByUno fail: %w", err)
	}
	return list, nil
}

// 사용자 권한 조회
// param: 권한코드, 프로젝트번호
func (r *Repository) GetUserRoleListByCodeAndJno(ctx context.Context, db Queryer, code string, jno int64) ([]entity.UserRoleMap, error) {
	var list []entity.UserRoleMap

	query := `
		SELECT 
			USER_UNO,
			ROLE_CODE,
			JNO
		FROM IRIS_USER_ROLE_MAP
		WHERE ROLE_CODE = :1
		AND JNO = :2`

	if err := db.SelectContext(ctx, &list, query, code, jno); err != nil {
		return nil, fmt.Errorf("GetUserRoleListByCodeAndJno fail: %w", err)
	}
	return list, nil
}

// 사용자 권한 추가
func (r *Repository) AddUserRole(ctx context.Context, tx Execer, userRoles []entity.UserRoleMap) error {
	agent := utils.GetAgent()

	query := `
		INSERT INTO IRIS_USER_ROLE_MAP(USER_UNO, ROLE_CODE, JNO, REG_DATE, REG_AGENT, REG_USER, REG_UNO)
		VALUES (:1, :2, :3, SYSDATE, :4, :5, :6)`

	for _, userRole := range userRoles {
		if _, err := tx.ExecContext(ctx, query, userRole.UserUno, userRole.RoleCode, userRole.Jno, agent, userRole.RegUser, userRole.RegUno); err != nil {
			return fmt.Errorf("AddUserRole fail: %w", err)
		}
	}
	return nil
}

// 사용자 권한 삭제
func (r *Repository) RemoveUserRole(ctx context.Context, tx Execer, userRoles []entity.UserRoleMap) error {
	query := `
		DELETE FROM IRIS_USER_ROLE_MAP
		WHERE USER_UNO = :1
		AND ROLE_CODE = :2
		AND JNO = :3`

	for _, userRole := range userRoles {
		if _, err := tx.ExecContext(ctx, query, userRole.UserUno, userRole.RoleCode, userRole.Jno); err != nil {
			return fmt.Errorf("RemoveUserRole fail: %w", err)
		}
	}
	return nil
}
