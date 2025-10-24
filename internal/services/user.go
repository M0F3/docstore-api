package services

import (
	"context"
	"log"
	"time"

	"github.com/M0F3/docstore-api/internal/auth"
	"github.com/M0F3/docstore-api/internal/models"
	"github.com/M0F3/docstore-api/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repositories.UserRepository
}


func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Register(ctx context.Context, data models.RegisterUserPayload) (*models.RegisteredUserDTO, error) {
	t, err := s.Repo.CreateTenant(ctx, models.Tenant{ID: uuid.New(), Name: data.TenantName})
	if err != nil {
		return nil, err
	}
    hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 4)
	if err != nil {
		_ = s.Repo.DeleteTenant(ctx, t.ID.String())
		log.Println(err)
		return nil, err
	}
	log.Println("Tenant Created")
	u, err := s.Repo.CreateUser(ctx, models.User{FirstName: data.FirstName, LastName: data.LastName, Email: data.Email, TenantId: t.ID, PasswordHash: string(hash)})
	if err != nil {
		erro := s.Repo.DeleteTenant(ctx, t.ID.String())
		if erro != nil {
			log.Println(erro)
			return nil, erro
		}
		log.Println(err)
		return nil, err
	}
	return &models.RegisteredUserDTO{Tenant: models.Tenant{ID: t.ID, Name: t.Name}, FirstName: u.FirstName, LastName: u.LastName, Email: u.Email, ID: u.ID}, nil
}

func (s *UserService) Login(ctx context.Context, data models.LoginPayload) (*models.SuccessfulLogin, error) {
	u, err := s.Repo.FindByEmail(ctx, data.Email)
	if err != nil {
		log.Fatalln("After Repo 1")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(data.Password))

	if err != nil {
		return nil, err
	}

	userResponse, err := s.Repo.GetUserWithTenantById(ctx, u.ID.String())
	
	if err != nil {
		log.Fatalln("After Repo")
		return nil, err
	}

	t, err := auth.GenerateToken(*u, 24 *time.Hour)

	if err != nil {
		log.Fatalln("After Token")
		return nil, err
	}

	return &models.SuccessfulLogin{Token: t, User: *userResponse }, nil
}

func (s UserService) ListUsers(ctx context.Context) ([]models.ListUserDTO, error) {
	return s.Repo.ListUsers(ctx)
}