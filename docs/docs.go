// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/contests": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "创建比赛",
                "parameters": [
                    {
                        "description": "Name",
                        "name": "Name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Description",
                        "name": "Description",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Rounds",
                        "name": "Rounds",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/result.CreateContestRequestRound"
                            }
                        }
                    },
                    {
                        "description": "StartTime",
                        "name": "StartTime",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "EndTime",
                        "name": "EndTime",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.CreateContestRequest"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Project": {
            "type": "integer",
            "enum": [
                1,
                2,
                3,
                4,
                5,
                6,
                7,
                8,
                9,
                10,
                11,
                12,
                13,
                14,
                15,
                16,
                17,
                18,
                19
            ],
            "x-enum-varnames": [
                "Cube222",
                "Cube333",
                "Cube444",
                "Cube555",
                "Cube666",
                "Cube777",
                "CubeSk",
                "CubePy",
                "CubeSq1",
                "CubeMinx",
                "CubeClock",
                "Cube333OH",
                "Cube333FM",
                "Cube333BF",
                "Cube444BF",
                "Cube555BF",
                "Cube333MBF",
                "JuBaoHaoHao",
                "OtherCola"
            ]
        },
        "result.CreateContestRequest": {
            "type": "object",
            "properties": {
                "Description": {
                    "type": "string"
                },
                "EndTime": {
                    "type": "integer"
                },
                "Name": {
                    "type": "string"
                },
                "Rounds": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/result.CreateContestRequestRound"
                    }
                },
                "StartTime": {
                    "type": "integer"
                }
            }
        },
        "result.CreateContestRequestRound": {
            "type": "object",
            "properties": {
                "Final": {
                    "type": "boolean"
                },
                "Name": {
                    "type": "string"
                },
                "Number": {
                    "type": "integer"
                },
                "Project": {
                    "$ref": "#/definitions/model.Project"
                },
                "Upsets": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
