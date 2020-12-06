package database

import (
	"encoding/json"
	"github.com/pkg/errors"

	edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/openx/consts"
)

// Save inserts a User object into the database
func (a *User) Save() error {
	return edb.Save(consts.DbDir+consts.DbName, UserBucket, a, a.Index)
}

// RetrieveUser retrieves a User from the database
func RetrieveUser(key int) (User, error) {
	var user User
	x, err := edb.Retrieve(consts.DbDir+consts.DbName, UserBucket, key)
	if err != nil {
		return user, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(x, &user)
	if user.Index == 0 {
		return user, errors.New("Error while retrieving user")
	}
	return user, err
}

// RetrieveAllUsers gets a list of all Users in the database
func RetrieveAllUsers() ([]User, error) {
	var arr []User
	x, err := edb.RetrieveAllKeys(consts.DbDir+consts.DbName, UserBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users")
	}
	for _, value := range x {
		var temp User
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.New("error while unmarshalling json, quitting")
		}
		arr = append(arr, temp)
	}

	return arr, nil
}

// RetrieveAllUsersLim gets the number of users in the bucket
func RetrieveAllUsersLim() (int, error) {
	return edb.RetrieveAllKeysLim(consts.DbDir+consts.DbName, UserBucket)
}

// TopReputationUsers gets a list of users sorted by descending reputation
func TopReputationUsers() ([]User, error) {
	// these reputation functions should mostly be used by the frontend through the
	// RPC to display to other users what other users' reputation is.
	arr, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}
	for i := range arr {
		for j := range arr {
			if arr[i].Reputation > arr[j].Reputation {
				tmp := arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}
	return arr, nil
}

// RetrieveAllAdmins retrieves a list of all admisn from the database
func RetrieveAllAdmins() ([]User, error) {
	var arr []User
	users, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Admin {
			var x User
			x.Index = user.Index
			x.Admin = user.Admin
			x.Image = user.Image
			x.Name = user.Name
			x.Username = user.Username
			x.Email = user.Email
			x.Country = user.Country
			arr = append(arr, x)
		}
	}

	return arr, nil
}

// ValidatePwhash validates a username / pwhash combination
func ValidatePwhash(name string, pwhash string) (User, error) {
	var dummy User
	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		// Running locally without token confirmation

		//if !user.Conf {
		//	continue
		//}
		if user.Username == name && user.Pwhash == pwhash {
			return user, nil
		}
	}
	return dummy, errors.New("could not find user with requested credentials")
}

// ValidatePwhashReg validates a username / pwhash combination during registration
func ValidatePwhashReg(name string, pwhash string) (User, error) {
	var dummy User
	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Username == name && user.Pwhash == pwhash {
			return user, nil
		}
	}
	return dummy, errors.New("could not find user with requested credentials")
}

// ValidateAccessToken validates a username / accessToken combination
func ValidateAccessToken(name string, accessToken string) (User, error) {
	var dummy User

	if len(accessToken) > consts.AccessTokenLength || len(accessToken) < consts.AccessTokenLength {
		return dummy, errors.New("incorrect token length")
	}

	users, err := RetrieveAllUsers()
	if err != nil {
		return dummy, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		// Running locally without token confirmation
		//if !user.Conf {
		//	continue
		//}
		if user.Username == name {
			dummy = user
			break
		}
	}

	timeNow := utils.Unix()
	for storedToken, timeout := range dummy.AccessToken {
		if storedToken == accessToken && timeNow-timeout < consts.AccessTokenLife {
			return dummy, nil
		}
	}

	return dummy, errors.New("could not find user with requested credentials")
}
