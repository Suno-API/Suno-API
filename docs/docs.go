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
        "/suno/account": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "account"
                ],
                "summary": "Get Account config",
                "responses": {
                    "200": {
                        "description": "song task",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/ginplus.DataResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/main.Account"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/suno/fetch": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "suno"
                ],
                "summary": "Fetch task",
                "parameters": [
                    {
                        "description": "fetch task ids",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.FetchReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "song tasks",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/ginplus.DataResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/po.Task"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/suno/fetch/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "suno"
                ],
                "summary": "Fetch task by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "fetch single task by id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "song task",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/ginplus.DataResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/po.Task"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/suno/submit/lyrics": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "suno"
                ],
                "summary": "Submit Suno lyrics task",
                "parameters": [
                    {
                        "description": "sumbmit generate lyrics",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.SubmitGenLyricsReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "task_id",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/ginplus.DataResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/suno/submit/music": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "suno"
                ],
                "summary": "Submit Suno song task",
                "parameters": [
                    {
                        "description": "sumbmit generate song",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.SubmitGenSongReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "task_id",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/ginplus.DataResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ginplus.DataResult": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "main.Account": {
            "type": "object",
            "properties": {
                "certificate": {
                    "$ref": "#/definitions/main.SunoCert"
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "main.FetchReq": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "ids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "main.SubmitGenLyricsReq": {
            "type": "object",
            "properties": {
                "prompt": {
                    "type": "string"
                }
            }
        },
        "main.SubmitGenSongReq": {
            "type": "object",
            "properties": {
                "continue_at": {
                    "type": "number"
                },
                "continue_clip_id": {
                    "type": "string"
                },
                "gpt_description_prompt": {
                    "type": "string"
                },
                "make_instrumental": {
                    "type": "boolean"
                },
                "mv": {
                    "type": "string"
                },
                "prompt": {
                    "type": "string"
                },
                "tags": {
                    "type": "string"
                },
                "task_id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "main.SunoCert": {
            "type": "object",
            "properties": {
                "cookie": {
                    "type": "string"
                },
                "credits_left": {
                    "type": "integer"
                },
                "is_active": {
                    "type": "boolean"
                },
                "jwt": {
                    "type": "string"
                },
                "last_update": {
                    "description": "最后更新时间，小于5秒，可以直接使用",
                    "type": "integer"
                },
                "monthly_limit": {
                    "type": "integer"
                },
                "monthly_usage": {
                    "type": "integer"
                },
                "period": {
                    "type": "string"
                },
                "session_id": {
                    "type": "string"
                }
            }
        },
        "po.Task": {
            "type": "object",
            "properties": {
                "action": {
                    "description": "任务类型, song, lyrics, description-mode",
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "data": {},
                "fail_reason": {
                    "type": "string"
                },
                "finish_time": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "search_item": {
                    "description": "Progress   string     ` + "`" + `json:\"progress\" gorm:\"type:varchar(20);index\"` + "`" + `",
                    "type": "string"
                },
                "start_time": {
                    "type": "integer"
                },
                "status": {
                    "description": "任务状态, submitted, queueing, processing, success, failed",
                    "allOf": [
                        {
                            "$ref": "#/definitions/po.TaskStatus"
                        }
                    ]
                },
                "submit_time": {
                    "type": "integer"
                },
                "task_id": {
                    "description": "第三方id，不一定有",
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "po.TaskStatus": {
            "type": "string",
            "enum": [
                "NOT_START"
            ],
            "x-enum-varnames": [
                "TaskStatusNotStart"
            ]
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
