package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Smileys struct {
	Id   int       `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
	Time time.Time `json:"time,omitempty"`
	Tags []string  `json:"tags,omitempty"`
}
type SmileysCount struct {
	Name  string `json:"name,omitempty"`
	Id    int    `json:"id,omitempty"`
	Count int    `json:"count,omitempty"`
}

type SmileysModel struct {
	DB *sql.DB
}

func (m SmileysModel) GetSmileys() ([]*Smileys, error) {

	query := ` SELECT * FROM smiley `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return []*Smileys{}, err
	}
	defer rows.Close()
	smileys := []*Smileys{}
	for rows.Next() {
		var smiley Smileys
		err := rows.Scan(&smiley.Id, &smiley.Name)
		if err != nil {
			return []*Smileys{}, err
		}

		smileys = append(smileys, &smiley)
	}
	if err = rows.Err(); err != nil {
		return []*Smileys{}, err
	}
	return smileys, nil
}

func (m SmileysModel) GetUserSmileys(userID string) ([]*Smileys, error) {

	query := ` SELECT smiley.id , smiley.name, user_smiley.granted_at, user_smiley.tags 
	FROM smiley
	JOIN user_smiley ON smiley.id = user_smiley.smiley_id
	WHERE user_smiley.user_id = $1; `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	smileys := []*Smileys{}
	for rows.Next() {
		var smiley Smileys
		err := rows.Scan(&smiley.Id, &smiley.Name, &smiley.Time, pq.Array(&smiley.Tags))
		if err != nil {
			return nil, err
		}
		smileys = append(smileys, &smiley)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return smileys, nil
}

func (m SmileysModel) InsertUserSmileys(userID string, smiley Smileys) error {
	query := `
	INSERT INTO user_smiley (user_id, smiley_id, granted_at, tags)
	VALUES ($1, $2, $3, $4) `

	args := []any{userID, smiley.Id, time.Now(), pq.Array(smiley.Tags)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m SmileysModel) GetUserSmileysCount(userID string, interval int) ([]*SmileysCount, error) {
	query := fmt.Sprintf(`
    SELECT smiley.name, smiley.id, COUNT(*) AS count
    FROM smiley
    JOIN user_smiley ON smiley.id = user_smiley.smiley_id
    WHERE user_smiley.user_id = $1
    AND user_smiley.granted_at >= NOW() - INTERVAL '%d days'
    GROUP BY smiley.name, smiley.id; `, interval)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	smileys := []*SmileysCount{}
	for rows.Next() {
		var smiley SmileysCount
		err := rows.Scan(&smiley.Name, &smiley.Id, &smiley.Count)
		if err != nil {
			return nil, err
		}
		smileys = append(smileys, &smiley)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return smileys, nil
}