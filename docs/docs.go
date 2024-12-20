// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Leonid Eliseev",
            "url": "https://t.me/Lenchiiiikkkk",
            "email": "leonid.2004eliseev@mail.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "Retrieves a list of songs with optional filters for group name, name, release date, text, and link.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get Songs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by group name",
                        "name": "group_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by song name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by release date (format: YYYY-MM-DD)",
                        "name": "release_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by text in the song",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by link",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Maximum number of items to retrieve (pagination)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items to skip (pagination)",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with a list of songs",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseGetSongs"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new song with the provided details.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Create Song",
                "parameters": [
                    {
                        "description": "Details of the song to be created",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RequestCreateSong"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Song created successfully",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseCreateSong"
                        }
                    },
                    "400": {
                        "description": "Bad request: validation error or song already exists",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict: song with the specified group and name already exists",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/songs/{uuid}": {
            "get": {
                "description": "Retrieves a song by its ID. Supports pagination for song text verses.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get Song",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Song ID (validated as UUID)",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Maximum number of verses to retrieve (pagination)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of verses to skip (pagination)",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with the song details",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseGetSong"
                        }
                    },
                    "400": {
                        "description": "Bad request: invalid ID, invalid limit/offset, or song not found",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not found: song with the specified ID does not exist",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a song by its ID.",
                "tags": [
                    "songs"
                ],
                "summary": "Delete Song",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Song ID (validated as UUID)",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Song successfully deleted"
                    },
                    "400": {
                        "description": "Bad request: invalid ID format",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not found: song with the specified ID does not exist",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Updates a song by its ID with new details.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Update Song",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Song ID (validated as UUID)",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Details for updating the song",
                        "name": "input",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/dto.RequestUpdateSong"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with updated song details",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseUpdateSong"
                        }
                    },
                    "400": {
                        "description": "Bad request: invalid id, input or song not changed",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not found: song with the specified ID does not exist",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict: song with the specified group and name already exists",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/handerr.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.RequestCreateSong": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "type": "string",
                    "example": "Imagine Dragons"
                },
                "song": {
                    "type": "string",
                    "example": "Thunder"
                }
            }
        },
        "dto.RequestUpdateSong": {
            "type": "object",
            "required": [
                "group",
                "name",
                "release_date",
                "text"
            ],
            "properties": {
                "group": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "dto.ResponseCreateSong": {
            "type": "object",
            "properties": {
                "song": {
                    "$ref": "#/definitions/dto.ResponseSong"
                }
            }
        },
        "dto.ResponseGetSong": {
            "type": "object",
            "properties": {
                "song": {
                    "$ref": "#/definitions/dto.ResponseSong"
                }
            }
        },
        "dto.ResponseGetSongs": {
            "type": "object",
            "properties": {
                "songs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.ResponseSong"
                    }
                }
            }
        },
        "dto.ResponseSong": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "dto.ResponseUpdateSong": {
            "type": "object",
            "properties": {
                "song": {
                    "$ref": "#/definitions/dto.ResponseSong"
                }
            }
        },
        "handerr.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Song Library API",
	Description:      "API for managing a library of songs, including creating, retrieving, updating, and deleting songs.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
