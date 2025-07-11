// Code generated by ent, DO NOT EDIT.

package prowsuites

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldID, id))
}

// JobID applies equality check predicate on the "job_id" field. It's identical to JobIDEQ.
func JobID(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldJobID, v))
}

// JobURL applies equality check predicate on the "job_url" field. It's identical to JobURLEQ.
func JobURL(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldJobURL, v))
}

// JobName applies equality check predicate on the "job_name" field. It's identical to JobNameEQ.
func JobName(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldJobName, v))
}

// SuiteName applies equality check predicate on the "suite_name" field. It's identical to SuiteNameEQ.
func SuiteName(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldSuiteName, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldName, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldStatus, v))
}

// ErrorMessage applies equality check predicate on the "error_message" field. It's identical to ErrorMessageEQ.
func ErrorMessage(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldErrorMessage, v))
}

// ExternalServicesImpact applies equality check predicate on the "external_services_impact" field. It's identical to ExternalServicesImpactEQ.
func ExternalServicesImpact(v bool) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldExternalServicesImpact, v))
}

// Time applies equality check predicate on the "time" field. It's identical to TimeEQ.
func Time(v float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldTime, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldCreatedAt, v))
}

// JobIDEQ applies the EQ predicate on the "job_id" field.
func JobIDEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldJobID, v))
}

// JobIDNEQ applies the NEQ predicate on the "job_id" field.
func JobIDNEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldJobID, v))
}

// JobIDIn applies the In predicate on the "job_id" field.
func JobIDIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldJobID, vs...))
}

// JobIDNotIn applies the NotIn predicate on the "job_id" field.
func JobIDNotIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldJobID, vs...))
}

// JobIDGT applies the GT predicate on the "job_id" field.
func JobIDGT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldJobID, v))
}

// JobIDGTE applies the GTE predicate on the "job_id" field.
func JobIDGTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldJobID, v))
}

// JobIDLT applies the LT predicate on the "job_id" field.
func JobIDLT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldJobID, v))
}

// JobIDLTE applies the LTE predicate on the "job_id" field.
func JobIDLTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldJobID, v))
}

// JobIDContains applies the Contains predicate on the "job_id" field.
func JobIDContains(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContains(FieldJobID, v))
}

// JobIDHasPrefix applies the HasPrefix predicate on the "job_id" field.
func JobIDHasPrefix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasPrefix(FieldJobID, v))
}

// JobIDHasSuffix applies the HasSuffix predicate on the "job_id" field.
func JobIDHasSuffix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasSuffix(FieldJobID, v))
}

// JobIDEqualFold applies the EqualFold predicate on the "job_id" field.
func JobIDEqualFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEqualFold(FieldJobID, v))
}

// JobIDContainsFold applies the ContainsFold predicate on the "job_id" field.
func JobIDContainsFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContainsFold(FieldJobID, v))
}

// JobURLEQ applies the EQ predicate on the "job_url" field.
func JobURLEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldJobURL, v))
}

// JobURLNEQ applies the NEQ predicate on the "job_url" field.
func JobURLNEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldJobURL, v))
}

// JobURLIn applies the In predicate on the "job_url" field.
func JobURLIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldJobURL, vs...))
}

// JobURLNotIn applies the NotIn predicate on the "job_url" field.
func JobURLNotIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldJobURL, vs...))
}

// JobURLGT applies the GT predicate on the "job_url" field.
func JobURLGT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldJobURL, v))
}

// JobURLGTE applies the GTE predicate on the "job_url" field.
func JobURLGTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldJobURL, v))
}

// JobURLLT applies the LT predicate on the "job_url" field.
func JobURLLT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldJobURL, v))
}

// JobURLLTE applies the LTE predicate on the "job_url" field.
func JobURLLTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldJobURL, v))
}

// JobURLContains applies the Contains predicate on the "job_url" field.
func JobURLContains(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContains(FieldJobURL, v))
}

// JobURLHasPrefix applies the HasPrefix predicate on the "job_url" field.
func JobURLHasPrefix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasPrefix(FieldJobURL, v))
}

// JobURLHasSuffix applies the HasSuffix predicate on the "job_url" field.
func JobURLHasSuffix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasSuffix(FieldJobURL, v))
}

// JobURLEqualFold applies the EqualFold predicate on the "job_url" field.
func JobURLEqualFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEqualFold(FieldJobURL, v))
}

// JobURLContainsFold applies the ContainsFold predicate on the "job_url" field.
func JobURLContainsFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContainsFold(FieldJobURL, v))
}

// JobNameEQ applies the EQ predicate on the "job_name" field.
func JobNameEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldJobName, v))
}

// JobNameNEQ applies the NEQ predicate on the "job_name" field.
func JobNameNEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldJobName, v))
}

// JobNameIn applies the In predicate on the "job_name" field.
func JobNameIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldJobName, vs...))
}

// JobNameNotIn applies the NotIn predicate on the "job_name" field.
func JobNameNotIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldJobName, vs...))
}

// JobNameGT applies the GT predicate on the "job_name" field.
func JobNameGT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldJobName, v))
}

// JobNameGTE applies the GTE predicate on the "job_name" field.
func JobNameGTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldJobName, v))
}

// JobNameLT applies the LT predicate on the "job_name" field.
func JobNameLT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldJobName, v))
}

// JobNameLTE applies the LTE predicate on the "job_name" field.
func JobNameLTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldJobName, v))
}

// JobNameContains applies the Contains predicate on the "job_name" field.
func JobNameContains(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContains(FieldJobName, v))
}

// JobNameHasPrefix applies the HasPrefix predicate on the "job_name" field.
func JobNameHasPrefix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasPrefix(FieldJobName, v))
}

// JobNameHasSuffix applies the HasSuffix predicate on the "job_name" field.
func JobNameHasSuffix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasSuffix(FieldJobName, v))
}

// JobNameEqualFold applies the EqualFold predicate on the "job_name" field.
func JobNameEqualFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEqualFold(FieldJobName, v))
}

// JobNameContainsFold applies the ContainsFold predicate on the "job_name" field.
func JobNameContainsFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContainsFold(FieldJobName, v))
}

// SuiteNameEQ applies the EQ predicate on the "suite_name" field.
func SuiteNameEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldSuiteName, v))
}

// SuiteNameNEQ applies the NEQ predicate on the "suite_name" field.
func SuiteNameNEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldSuiteName, v))
}

// SuiteNameIn applies the In predicate on the "suite_name" field.
func SuiteNameIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldSuiteName, vs...))
}

// SuiteNameNotIn applies the NotIn predicate on the "suite_name" field.
func SuiteNameNotIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldSuiteName, vs...))
}

// SuiteNameGT applies the GT predicate on the "suite_name" field.
func SuiteNameGT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldSuiteName, v))
}

// SuiteNameGTE applies the GTE predicate on the "suite_name" field.
func SuiteNameGTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldSuiteName, v))
}

// SuiteNameLT applies the LT predicate on the "suite_name" field.
func SuiteNameLT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldSuiteName, v))
}

// SuiteNameLTE applies the LTE predicate on the "suite_name" field.
func SuiteNameLTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldSuiteName, v))
}

// SuiteNameContains applies the Contains predicate on the "suite_name" field.
func SuiteNameContains(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContains(FieldSuiteName, v))
}

// SuiteNameHasPrefix applies the HasPrefix predicate on the "suite_name" field.
func SuiteNameHasPrefix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasPrefix(FieldSuiteName, v))
}

// SuiteNameHasSuffix applies the HasSuffix predicate on the "suite_name" field.
func SuiteNameHasSuffix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasSuffix(FieldSuiteName, v))
}

// SuiteNameEqualFold applies the EqualFold predicate on the "suite_name" field.
func SuiteNameEqualFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEqualFold(FieldSuiteName, v))
}

// SuiteNameContainsFold applies the ContainsFold predicate on the "suite_name" field.
func SuiteNameContainsFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContainsFold(FieldSuiteName, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContainsFold(FieldName, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContainsFold(FieldStatus, v))
}

// ErrorMessageEQ applies the EQ predicate on the "error_message" field.
func ErrorMessageEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldErrorMessage, v))
}

// ErrorMessageNEQ applies the NEQ predicate on the "error_message" field.
func ErrorMessageNEQ(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldErrorMessage, v))
}

// ErrorMessageIn applies the In predicate on the "error_message" field.
func ErrorMessageIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldErrorMessage, vs...))
}

// ErrorMessageNotIn applies the NotIn predicate on the "error_message" field.
func ErrorMessageNotIn(vs ...string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldErrorMessage, vs...))
}

// ErrorMessageGT applies the GT predicate on the "error_message" field.
func ErrorMessageGT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldErrorMessage, v))
}

// ErrorMessageGTE applies the GTE predicate on the "error_message" field.
func ErrorMessageGTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldErrorMessage, v))
}

// ErrorMessageLT applies the LT predicate on the "error_message" field.
func ErrorMessageLT(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldErrorMessage, v))
}

// ErrorMessageLTE applies the LTE predicate on the "error_message" field.
func ErrorMessageLTE(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldErrorMessage, v))
}

// ErrorMessageContains applies the Contains predicate on the "error_message" field.
func ErrorMessageContains(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContains(FieldErrorMessage, v))
}

// ErrorMessageHasPrefix applies the HasPrefix predicate on the "error_message" field.
func ErrorMessageHasPrefix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasPrefix(FieldErrorMessage, v))
}

// ErrorMessageHasSuffix applies the HasSuffix predicate on the "error_message" field.
func ErrorMessageHasSuffix(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldHasSuffix(FieldErrorMessage, v))
}

// ErrorMessageIsNil applies the IsNil predicate on the "error_message" field.
func ErrorMessageIsNil() predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIsNull(FieldErrorMessage))
}

// ErrorMessageNotNil applies the NotNil predicate on the "error_message" field.
func ErrorMessageNotNil() predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotNull(FieldErrorMessage))
}

// ErrorMessageEqualFold applies the EqualFold predicate on the "error_message" field.
func ErrorMessageEqualFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEqualFold(FieldErrorMessage, v))
}

// ErrorMessageContainsFold applies the ContainsFold predicate on the "error_message" field.
func ErrorMessageContainsFold(v string) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldContainsFold(FieldErrorMessage, v))
}

// ExternalServicesImpactEQ applies the EQ predicate on the "external_services_impact" field.
func ExternalServicesImpactEQ(v bool) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldExternalServicesImpact, v))
}

// ExternalServicesImpactNEQ applies the NEQ predicate on the "external_services_impact" field.
func ExternalServicesImpactNEQ(v bool) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldExternalServicesImpact, v))
}

// ExternalServicesImpactIsNil applies the IsNil predicate on the "external_services_impact" field.
func ExternalServicesImpactIsNil() predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIsNull(FieldExternalServicesImpact))
}

// ExternalServicesImpactNotNil applies the NotNil predicate on the "external_services_impact" field.
func ExternalServicesImpactNotNil() predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotNull(FieldExternalServicesImpact))
}

// TimeEQ applies the EQ predicate on the "time" field.
func TimeEQ(v float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldTime, v))
}

// TimeNEQ applies the NEQ predicate on the "time" field.
func TimeNEQ(v float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldTime, v))
}

// TimeIn applies the In predicate on the "time" field.
func TimeIn(vs ...float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldTime, vs...))
}

// TimeNotIn applies the NotIn predicate on the "time" field.
func TimeNotIn(vs ...float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldTime, vs...))
}

// TimeGT applies the GT predicate on the "time" field.
func TimeGT(v float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldTime, v))
}

// TimeGTE applies the GTE predicate on the "time" field.
func TimeGTE(v float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldTime, v))
}

// TimeLT applies the LT predicate on the "time" field.
func TimeLT(v float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldTime, v))
}

// TimeLTE applies the LTE predicate on the "time" field.
func TimeLTE(v float64) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldTime, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldLTE(FieldCreatedAt, v))
}

// CreatedAtIsNil applies the IsNil predicate on the "created_at" field.
func CreatedAtIsNil() predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldIsNull(FieldCreatedAt))
}

// CreatedAtNotNil applies the NotNil predicate on the "created_at" field.
func CreatedAtNotNil() predicate.ProwSuites {
	return predicate.ProwSuites(sql.FieldNotNull(FieldCreatedAt))
}

// HasProwSuites applies the HasEdge predicate on the "prow_suites" edge.
func HasProwSuites() predicate.ProwSuites {
	return predicate.ProwSuites(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProwSuitesTable, ProwSuitesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProwSuitesWith applies the HasEdge predicate on the "prow_suites" edge with a given conditions (other predicates).
func HasProwSuitesWith(preds ...predicate.Repository) predicate.ProwSuites {
	return predicate.ProwSuites(func(s *sql.Selector) {
		step := newProwSuitesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ProwSuites) predicate.ProwSuites {
	return predicate.ProwSuites(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ProwSuites) predicate.ProwSuites {
	return predicate.ProwSuites(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ProwSuites) predicate.ProwSuites {
	return predicate.ProwSuites(sql.NotPredicates(p))
}
