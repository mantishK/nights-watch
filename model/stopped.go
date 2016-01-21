package model

import "plivo/nights-watch/config"

type Stopped struct {
	AuthID string
	From   string
	To     string
}

//Adds the stopped fields into db
func (s *Stopped) Add() (err error) {
	_, err = config.DB.Exec(`
    INSERT INTO stopped VALUES ($1, $2, $3)
    `, s.AuthID, s.From, s.To)
	return
}

//Deletes the row from DB
func (s *Stopped) Delete() (err error) {
	_, err = config.DB.Exec(`
    DELETE FROM stopped 
    WHERE auth_id = $1 AND from_number = $2 AND to_number = $3
    `, s.AuthID, s.From, s.To)
	return
}

//Returns if the row exists or not
func (s *Stopped) Exists() (exists bool, err error) {
	err = config.DB.QueryRow(`
    SELECT EXISTS 
      (SELECT 1 
      FROM stopped 
      WHERE auth_id = $1 AND from_number = $2 AND to_number = $3)
  `, s.AuthID, s.From, s.To).Scan(&exists)
	return
}
