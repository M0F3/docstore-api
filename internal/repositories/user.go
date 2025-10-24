package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/M0F3/docstore-api/internal/middleware"
	"github.com/M0F3/docstore-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
	adminDb *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool, adminDb *pgxpool.Pool) *UserRepository {
	return &UserRepository{db, adminDb}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, first_name, last_name, tenant_id, email, password_hash FROM base.users WHERE email=$1 LIMIT 1`
	err := r.adminDb.QueryRow(ctx, query, email).Scan(&user.ID,  &user.FirstName, &user.LastName, &user.TenantId, &user.Email, &user.PasswordHash)
	return user, err
}

func (r *UserRepository) CreateTenant(ctx context.Context, t models.Tenant) (*models.Tenant, error) {
	_, e := r.adminDb.Exec(ctx, "INSERT INTO base.tenants (id, name) VALUES($1, $2)", t.ID.String(), &t.Name)
	if e != nil {
		return nil, e
	}

	return &t, nil
}


func (r *UserRepository) DeleteTenant(ctx context.Context, tenantId string) error {
	_, e := r.adminDb.Exec(ctx, "DELETE FROM base.tenants WHERE id = $1", tenantId)
	if e != nil {
		return e
	}

	return nil
}

func (r *UserRepository) CreateUser(ctx context.Context, t models.User) (*models.User, error) {
	_, e := r.adminDb.Exec(ctx, "INSERT INTO base.users (id, first_name, last_name, tenant_id, email, password_hash) VALUES($1, $2, $3, $4, $5, $6)", uuid.New().String(), &t.FirstName, &t.LastName, &t.TenantId, &t.Email, &t.PasswordHash)
	if e != nil {
		return nil, e
	}

	return &t, nil
}


func (r *UserRepository) GetUserWithTenantById(ctx context.Context, id string) (*models.UserDTO, error) {
	user := &models.UserDTO{}
	query := `SELECT u.id, u.first_name, u.last_name, u.email, jsonb_build_object('id', t.id, 'name', t.name) tenant FROM base.users u LEFT JOIN base.tenants t ON  t.id = u.tenant_id WHERE u.id=$1 LIMIT 1`
	err := r.adminDb.QueryRow(ctx, query, id).Scan(&user.ID,  &user.FirstName, &user.LastName, &user.Email, &user.Tenant)
	if err != nil {
		log.Fatalln("In repo")
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) ListUsers(ctx context.Context) ([]models.ListUserDTO, error) {
	c, ok := middleware.GetDatabaseConnectionFromContext(ctx)
	if !ok {
		return []models.ListUserDTO{}, errors.New("could not get database connection")
	}

	rows, err := c.Query(ctx, "SELECT id, first_name, last_name, email FROM base.users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.ListUserDTO

	for rows.Next() {
		var u models.ListUserDTO
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}