package docs

import (
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "This is a BankAccount server.",
        "title": "BANKACCOUNT API",
        "version": "1.0"
    },
    "host": "localhost:1323",
    "basePath": "/v1",
    "paths": {
        "/v1/users": {
            "get": {
                "description": "get All Users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Insert Into User",
						"name": "{}",
                        "in": "path",
                        "required": false,
                        "schema": {
                            "type": "int"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
							"type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
		},
		"/v1/user/{id}": {
		"get": {
			"description": "get All Users",
			"consumes": [
				"application/json"
			],
			"produces": [
				"application/json"
			],
			"parameters": [
				{
					"description": "Some ID",
					"name": "some_id",
					"in": "path",
					"required": true,
					"schema": {
						"type": "string"
					}
				},
				{
					"description": "Offset",
					"name": "offset",
					"in": "query",
					"required": true,
					"schema": {
						"type": "int"
					}
				},
				{
					"description": "Offset",
					"name": "limit",
					"in": "query",
					"required": true,
					"schema": {
						"type": "int"
					}
				}
			],
			"responses": {
				"200": {
					"description": "ok",
					"schema": {
						"type": "string"
					}
				},
				"400": {
					"description": "We need ID!!",
					"schema": {
						"type": "object",
						"$ref": "#/definitions/web.APIError"
					}
				},
				"404": {
					"description": "Can not find ID",
					"schema": {
						"type": "object",
						"$ref": "#/definitions/web.APIError"
					}
				}
			}
		}
	},
        "user/{id}/bankAccount": {
            "get": {
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "limit",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
		},
		"user/{id}/bankAccount{idBankAccount}": {
            "delete": {
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "limit",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
		},
		"user/{id}/bankAccount/{idBankAccount}/deposit": {
            "put": {
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "limit",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
		},
		"user/{id}/bankAccount/{idBankAccount}/withdraw": {
            "put": {
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "limit",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
		},
		"tranfers/from/:idFrom/to/:idTo": {
            "post": {
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Some ID",
                        "name": "some_id",
                        "in": "path",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    },
                    {
                        "description": "Offset",
                        "name": "limit",
                        "in": "query",
                        "required": true,
                        "schema": {
                            "type": "int"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "We need ID!!",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/web.APIError"
                        }
                    }
                }
            }
		},
    },
    "definitions": {
        "web.APIError": {
            "type": "object",
            "properties": {
                "ErrorCode": {
                    "type": "int"
                },
                "ErrorMessage": {
                    "type": "string"
                }
            }
        }
    }
}`

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}
func init() {
	swag.Register(swag.Name, &s{})
}
