// This file defines the RequireStudentAccess middleware, which authorizes access to a student's resources based on the user's role (admin, teacher, student requesting their own profile, or guardian requesting their child's profile).

package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
	"github.com/openschool-org/openschool/internal/models"
)

// RequireStudentAccess aborts with 403 unless the caller is an admin, a teacher, the student themself, or a guardian of the student named by the :id URL parameter.
func RequireStudentAccess(pool *pgxpool.Pool) gin.HandlerFunc {
	queries := db.New(pool)
	return func(c *gin.Context) {
		userID, err := UserIDFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}

		userRoleList, ok := rolesFromContext(c)
		if !ok {
			return
		}

		if slices.Contains(userRoleList, models.RoleAdmin) || slices.Contains(userRoleList, models.RoleTeacher) {
			c.Next()
			return
		}

		studentIDStr := c.Param("id")
		if studentIDStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing student id"})
			return
		}
		studentID, err := uuid.Parse(studentIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
			return
		}

		if slices.Contains(userRoleList, models.RoleStudent) {
			student, err := queries.GetStudentByUserID(c.Request.Context(), pgtype.UUID{Bytes: userID, Valid: true})
			if err == nil && student.ID == studentID {
				c.Next()
				return
			}
		}

		if slices.Contains(userRoleList, models.RoleParent) {
			isGuardian, err := queries.IsGuardianOfStudent(c.Request.Context(), db.IsGuardianOfStudentParams{
				UserID:    pgtype.UUID{Bytes: userID, Valid: true},
				StudentID: studentID,
			})
			if err == nil && isGuardian {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied for this student resource"})
	}
}
