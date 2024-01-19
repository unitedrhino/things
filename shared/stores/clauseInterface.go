package stores

import "gorm.io/gorm/clause"

const AuthModify = "authModify:%v"

type Opt int64

const (
	Create Opt = iota + 1
	Update
	Delete
	Select
)

type clauseInterface struct {
}

func (sd clauseInterface) Name() string {
	return ""
}

func (sd clauseInterface) Build(clause.Builder) {

}

func (sd clauseInterface) MergeClause(*clause.Clause) {

}
