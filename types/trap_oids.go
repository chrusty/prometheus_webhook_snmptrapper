package types

import (
	snmpgo "github.com/k-sone/snmpgo"
)

type TrapOIDs struct {
	FiringTrap   *snmpgo.Oid
	RecoveryTrap *snmpgo.Oid
	Instance     *snmpgo.Oid
	Service      *snmpgo.Oid
	Location     *snmpgo.Oid
	Severity     *snmpgo.Oid
	Description  *snmpgo.Oid
	JobName      *snmpgo.Oid
	TimeStamp    *snmpgo.Oid
}
