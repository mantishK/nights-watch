package model

import (
	"database/sql"
	"plivo/nights-watch/config"
	"strconv"
)

//Returns Account Id based on Auth ID
func GetAccountID(authID string) (accountID int, err error) {
	accountIDstr := ""
	err = config.DB.QueryRow(`
    SELECT id
    FROM account
    WHERE auth_id = $1`, authID).Scan(&accountIDstr)
	if err != nil {
		return
	}
	accountID, _ = strconv.Atoi(accountIDstr)
	return
}

//Returns Minimum used phone No that is not blocked
func GetMinUsedPhoneNo(toPhNo, authID string, accountID int, max int) (tx *sql.Tx, minUsedNo string, err error) {
	tx, err = config.DB.Begin()
	if err != nil {
		return
	}
	err = tx.QueryRow(`
    SELECT p.number 
    FROM phone_number AS p 
    LEFT JOIN stopped AS s ON s.auth_id = $1 AND s.from_number = p.number AND s.to_number = $2 
    WHERE p.account_id = $3 AND p.count < $4 AND s.auth_id IS NULL 
    ORDER BY p.count ASC LIMIT 1`, authID, toPhNo, accountID, max).Scan(&minUsedNo)
	return
}

//Increments the phone no used count
func IncPhoneUsedCount(phoneNo, authID string, tx *sql.Tx) (err error) {
	_, err = tx.Exec(`
    UPDATE phone_number 
    SET count = count + 1 
    WHERE number = $1 AND account_id = 
      (SELECT id FROM account WHERE auth_id = $2)
    `, phoneNo, authID)
	tx.Commit()
	return
}

//Decrements the phone no used count
func DecPhoneUsedCount(phoneNo, authID string) (err error) {
	_, err = config.DB.Exec(`
    UPDATE phone_number 
    SET count = count - 1 
    WHERE number = $1 AND account_id = 
      (SELECT id FROM account WHERE auth_id = $2)
    `, phoneNo, authID)
	return
}

//Returns Auth ID when given a Phone No
func GetAuthID(phoneNo string) (authID string, err error) {
	err = config.DB.QueryRow(`
    SELECT a.auth_id 
    FROM phone_number AS p 
    JOIN account AS a 
    ON p.account_id = a.id 
    WHERE p.number = $1
    `, phoneNo).Scan(&authID)
	return
}
