// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

import (
	"context"
	"errors"
	"net/http"

	v1 "github.com/sacloud/security-control-api-go/apis/v1"
)

type EvaluationRulesAPI interface {
	List(ctx context.Context, parameter v1.EvaluationRulesListParams) (*v1.EvaluationRulesListOK, error)
	Read(ctx context.Context, id string) (*v1.EvaluationRule, error)
	Update(ctx context.Context, id string, req *v1.EvaluationRuleInput) (*v1.EvaluationRule, error)
}

var _ EvaluationRulesAPI = (*evaluationRuleOp)(nil)

type evaluationRuleOp struct {
	client *v1.Client
}

func NewEvaluationRulesOp(client *v1.Client) EvaluationRulesAPI {
	return &evaluationRuleOp{client: client}
}

func (op *evaluationRuleOp) List(ctx context.Context, params v1.EvaluationRulesListParams) (*v1.EvaluationRulesListOK, error) {
	const methodName = "EvaluationRules.List"

	res, err := op.client.EvaluationRulesList(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.EvaluationRulesListOK:
		return r, nil
	case *v1.EvaluationRulesListApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesListApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesListApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesListApplicationJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesListApplicationJSONTooManyRequests:
		// TODO: Handle RetryAfter field
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, NewAPIError(methodName, r.StatusCode, errors.New(r.Response.Detail.Value))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *evaluationRuleOp) Read(ctx context.Context, id string) (*v1.EvaluationRule, error) {
	const methodName = "EvaluationRules.Read"

	res, err := op.client.EvaluationRulesRead(ctx, v1.EvaluationRulesReadParams{EvaluationRuleId: v1.EvaluationRuleID(id)})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.EvaluationRule:
		return r, nil
	case *v1.EvaluationRulesReadApplicationJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesReadApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesReadApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *evaluationRuleOp) Update(ctx context.Context, id string, req *v1.EvaluationRuleInput) (*v1.EvaluationRule, error) {
	const methodName = "EvaluationRules.Update"

	res, err := op.client.EvaluationRulesUpdate(ctx, req, v1.EvaluationRulesUpdateParams{EvaluationRuleId: v1.EvaluationRuleID(id)})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.EvaluationRule:
		return r, nil
	case *v1.EvaluationRulesUpdateApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesUpdateApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesUpdateApplicationJSONForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesUpdateApplicationJSONNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(string(r.Detail.Value)))
	case *v1.EvaluationRulesUpdateApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

type EvaluationRuleInputParams struct {
	ID                 string
	ServicePrincipalID string
	Targets            []string
	Enabled            bool
}

func SetupEvaluationRuleInput(params *EvaluationRuleInputParams) *v1.EvaluationRuleInput {
	switch params.ID {
	case "server-no-public-ip":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("server-no-public-ip"),
					ServerNoPublicIP: v1.ServerNoPublicIP{
						EvaluationRuleId: v1.ServerNoPublicIPEvaluationRuleId("server-no-public-ip"),
						Parameter: v1.NewOptEvaluationRuleParametersZonedEvaluationTarget(
							v1.EvaluationRuleParametersZonedEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
								Zones:              params.Targets,
							},
						),
					},
				},
			},
		}
	case "disk-encryption-enabled":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("disk-encryption-enabled"),
					DiskEncryptionEnabled: v1.DiskEncryptionEnabled{
						EvaluationRuleId: v1.DiskEncryptionEnabledEvaluationRuleId("disk-encryption-enabled"),
						Parameter: v1.NewOptEvaluationRuleParametersZonedEvaluationTarget(
							v1.EvaluationRuleParametersZonedEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
								Zones:              params.Targets,
							},
						),
					},
				},
			},
		}
	case "dba-encryption-enabled":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("dba-encryption-enabled"),
					DbaEncryptionEnabled: v1.DbaEncryptionEnabled{
						EvaluationRuleId: v1.DbaEncryptionEnabledEvaluationRuleId("dba-encryption-enabled"),
						Parameter: v1.NewOptEvaluationRuleParametersZonedEvaluationTarget(
							v1.EvaluationRuleParametersZonedEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
								Zones:              params.Targets,
							},
						),
					},
				},
			},
		}
	case "dba-no-public-ip":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("dba-no-public-ip"),
					DbaNoPublicIP: v1.DbaNoPublicIP{
						EvaluationRuleId: v1.DbaNoPublicIPEvaluationRuleId("dba-no-public-ip"),
						Parameter: v1.NewOptEvaluationRuleParametersZonedEvaluationTarget(
							v1.EvaluationRuleParametersZonedEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
								Zones:              params.Targets,
							},
						),
					},
				},
			},
		}
	case "objectstorage-bucket-acl-changed":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("objectstorage-bucket-acl-changed"),
					ObjectStorageBucketACLChanged: v1.ObjectStorageBucketACLChanged{
						EvaluationRuleId: v1.ObjectStorageBucketACLChangedEvaluationRuleId("objectstorage-bucket-acl-changed"),
						Parameter: v1.NewOptEvaluationRuleParametersEvaluationTarget(
							v1.EvaluationRuleParametersEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
							},
						),
					},
				},
			},
		}
	case "objectstorage-bucket-encryption-enabled":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("objectstorage-bucket-encryption-enabled"),
					ObjectStorageBucketEncryptionEnabled: v1.ObjectStorageBucketEncryptionEnabled{
						EvaluationRuleId: v1.ObjectStorageBucketEncryptionEnabledEvaluationRuleId("objectstorage-bucket-encryption-enabled"),
						Parameter: v1.NewOptEvaluationRuleParametersObjectStorageEvaluationTarget(
							v1.EvaluationRuleParametersObjectStorageEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
								Sites:              params.Targets,
							},
						),
					},
				},
			},
		}
	case "addon-datalake-no-public-access":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("addon-datalake-no-public-access"),
					AddonDatalakeNoPublicAccess: v1.AddonDatalakeNoPublicAccess{
						EvaluationRuleId: v1.AddonDatalakeNoPublicAccessEvaluationRuleId("addon-datalake-no-public-access"),
						Parameter: v1.NewOptEvaluationRuleParametersEvaluationTarget(
							v1.EvaluationRuleParametersEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
							},
						),
					},
				},
			},
		}
	case "addon-dwh-no-public-access":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("addon-dwh-no-public-access"),
					AddonDwhNoPublicAccess: v1.AddonDwhNoPublicAccess{
						EvaluationRuleId: v1.AddonDwhNoPublicAccessEvaluationRuleId("addon-dwh-no-public-access"),
						Parameter: v1.NewOptEvaluationRuleParametersEvaluationTarget(
							v1.EvaluationRuleParametersEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
							},
						),
					},
				},
			},
		}
	case "addon-threat-detection-enabled":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("addon-threat-detection-enabled"),
					AddonThreatDetectionEnabled: v1.AddonThreatDetectionEnabled{
						EvaluationRuleId: v1.AddonThreatDetectionEnabledEvaluationRuleId("addon-threat-detection-enabled"),
						Parameter: v1.NewOptEvaluationRuleParametersZonedEvaluationTarget(
							v1.EvaluationRuleParametersZonedEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
								Zones:              params.Targets,
							},
						),
					},
				},
			},
		}
	case "addon-threat-detections":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("addon-threat-detections"),
					AddonThreatDetections: v1.AddonThreatDetections{
						EvaluationRuleId: v1.AddonThreatDetectionsEvaluationRuleId("addon-threat-detections"),
					},
				},
			},
		}
	case "addon-vulnerability-detections":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("addon-vulnerability-detections"),
					AddonVulnerabilityDetections: v1.AddonVulnerabilityDetections{
						EvaluationRuleId: v1.AddonVulnerabilityDetectionsEvaluationRuleId("addon-vulnerability-detections"),
					},
				},
			},
		}
	case "elb-logging-enabled":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("elb-logging-enabled"),
					ELBLoggingEnabled: v1.ELBLoggingEnabled{
						EvaluationRuleId: v1.ELBLoggingEnabledEvaluationRuleId("elb-logging-enabled"),
						Parameter: v1.NewOptEvaluationRuleParametersEvaluationTarget(
							v1.EvaluationRuleParametersEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
							},
						),
					},
				},
			},
		}
	case "iam-member-operation-detected":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("iam-member-operation-detected"),
					IAMMemberOperationDetected: v1.IAMMemberOperationDetected{
						EvaluationRuleId: v1.IAMMemberOperationDetectedEvaluationRuleId("iam-member-operation-detected"),
						Parameter: v1.NewOptEvaluationRuleParametersEvaluationTarget(
							v1.EvaluationRuleParametersEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
							},
						),
					},
				},
			},
		}
	case "nosql-encryption-enabled":
		return &v1.EvaluationRuleInput{
			IsEnabled: params.Enabled,
			Rule: v1.EvaluationRuleUnion{
				OneOf: v1.EvaluationRuleUnionSum{
					Type: v1.EvaluationRuleUnionSumType("nosql-encryption-enabled"),
					NoSQLEncryptionEnabled: v1.NoSQLEncryptionEnabled{
						EvaluationRuleId: v1.NoSQLEncryptionEnabledEvaluationRuleId("nosql-encryption-enabled"),
						Parameter: v1.NewOptEvaluationRuleParametersZonedEvaluationTarget(
							v1.EvaluationRuleParametersZonedEvaluationTarget{
								ServicePrincipalId: v1.NewOptString(params.ServicePrincipalID),
								Zones:              params.Targets,
							},
						),
					},
				},
			},
		}
	default:
		panic("unsupported evaluation rule ID: " + params.ID)
	}
}
