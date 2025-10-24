package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/M0F3/docstore-api/internal/database"
	"github.com/M0F3/docstore-api/internal/models"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

type ctxKey string

const UserKey ctxKey = "user"

func AttachUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		u, ok := claims["user"]
		b, _ := json.Marshal(u)

		var c models.User
		err := json.Unmarshal(b, &c)
		if err != nil {
			w.WriteHeader(401)
			return 
		}
		if ok {
			ctx := context.WithValue(r.Context(), UserKey, c)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUser(ctx context.Context) (models.User, bool) {
	val := ctx.Value(UserKey)
	if user, ok := val.(models.User); ok {
		return user, true
	}
	return models.User{}, false
}

const DatabaseSessionKey ctxKey = "database"

func AttachDatabaseConnection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := database.GetConnection().Acquire(r.Context())

		if err != nil {
			panic(err)
		}
		defer c.Release()
		ctx := context.WithValue(r.Context(), DatabaseSessionKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AttachDatabaseConnectionWithSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := database.GetConnection().Acquire(r.Context())

		if err != nil {
			http.Error(w, "failed to acquire DB connection", http.StatusInternalServerError)
			return
		}

		user, ok := GetUser(r.Context())

		if !ok {
			log.Println(err)
			http.Error(w, "failed to get user", http.StatusUnauthorized)
			return
		}

		if _, err := c.Conn().Exec(r.Context(),fmt.Sprintf("SET app.current_tenant_id = '%s'", user.TenantId)); err != nil {
    		log.Println(err)
			c.Release()
			w.WriteHeader(401)
			return
		}
		go func() {
			<-r.Context().Done()
			c.Exec(r.Context(), "UNSET app.current_tenant_id")
			c.Release()
	}()
		ctx := context.WithValue(r.Context(), DatabaseSessionKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetDatabaseConnectionFromContext(ctx context.Context) (*pgxpool.Conn, bool) {
	val := ctx.Value(DatabaseSessionKey)
	if c, ok := val.(*pgxpool.Conn); ok {
		return c, true
	}
	return nil, false
}