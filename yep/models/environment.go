// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/npiganeau/yep/yep/tools"
)

// Environment is a sql transaction, performed by a user with
// some context data.
type Environment struct {
	cr      *sqlx.Tx
	uid     int64
	context tools.Context
}

/*
Cr return a pointer to the transaction of the Environment
*/
func (env *Environment) Cr() *sqlx.Tx {
	return env.cr
}

/*
Uid returns the user id of the Environment
*/
func (env *Environment) Uid() int64 {
	return env.uid
}

/*
Context returns the Context of the Environment
*/
func (env *Environment) Context() tools.Context {
	return env.context
}

/*
WithContext returns a new Environment with its context updated by ctx.
If replace is true, then the context is replaced by the given ctx instead of
being updated.
*/
func (env *Environment) WithContext(ctx tools.Context, replace ...bool) *Environment {
	if len(replace) > 0 && replace[0] {
		return NewEnvironment(env.cr, env.uid, ctx)
	}
	newCtx := env.context
	for key, value := range ctx {
		newCtx[key] = value
	}
	return NewEnvironment(env.cr, env.uid, newCtx)
}

/*
Sudo returns a new Environment with the given userId or the superuser id if not specified
*/
func (env *Environment) Sudo(userId ...int64) *Environment {
	var uid int64
	if len(userId) > 0 {
		uid = userId[0]
	} else {
		uid = 1
	}
	return NewEnvironment(env.cr, uid, env.context)
}

///*
//Create creates a new record in database from the given data and returns a recordSet
//Data must be a struct pointer.
//*/
//func (env *Environment) Create(data interface{}) RecordSet {
//	if err := checkStructPtr(data); err != nil {
//		panic(fmt.Errorf("Create error: %s", err))
//	}
//	_, err := env.cr.Insert(data)
//	if err != nil {
//		panic(fmt.Errorf("Create error: %s", err))
//	}
//	rs := newRecordStructFromData(env, data)
//	rs.updateStoredFields(data)
//	rs.computeFields(data)
//	return rs
//}
//
///*
//Sync writes the given data to database.
//data must be a struct pointer that has been originally populated by RecordSet.ReadOne()
//or RecordSet.ReadAll().
//*/
//func (env Environment) Sync(data interface{}, cols ...string) bool {
//	if err := checkStructPtr(data); err != nil {
//		panic(fmt.Errorf("<Environment.Sync>: %s", err))
//	}
//	rs := newRecordStructFromData(env, data)
//	params := structToMap(data)
//	res := rs.Call("Write", params).(bool)
//	return res
//}

/*
NewEnvironment returns a new Environment with the given parameters.
*/
func NewEnvironment(cr *sqlx.Tx, uid int64, context ...tools.Context) *Environment {
	var ctx tools.Context
	if len(context) > 0 {
		ctx = context[0]
	}
	env := &Environment{
		cr:      cr,
		uid:     uid,
		context: ctx,
	}
	return env
}

func NewCursorEnvironment(uid int64, context ...tools.Context) *Environment {
	cr := DB.MustBegin()
	return NewEnvironment(cr, uid, context...)
}

/*
Pool returns an empty RecordSet from the given table name string or struct pointer
*/
func (env *Environment) Pool(tableNameOrStructPtr interface{}) RecordSet {
	return NewRecordSet(env, tableNameOrStructPtr)
}
