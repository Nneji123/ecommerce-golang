// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "LeadzAura Support",
            "url": "http://leadzaura.com",
            "email": "contact@leadzaura.com"
        },
        "license": {
            "name": "MIT",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/generate-emails": {
            "post": {
                "description": "Generates email permutations based on the provided parameters.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Generate email permutations",
                "operationId": "handle-post-generate-emails",
                "parameters": [
                    {
                        "description": "Request body with parameters",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_emailapi_emailpermutator.GenerateEmailParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/leads": {
            "get": {
                "description": "Retrieves all leads from the database.",
                "summary": "Retrieve all leads",
                "operationId": "handle-get-all-leads",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Lead"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "OAuth2Application": [
                            "write"
                        ]
                    }
                ],
                "description": "Creates a new lead based on the provided parameters.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new lead",
                "operationId": "handle-post-create-lead",
                "parameters": [
                    {
                        "description": "Lead details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Lead"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/db.Lead"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/leads/{id}": {
            "get": {
                "description": "Retrieves a lead by its ID from the database.",
                "summary": "Retrieve a lead by ID",
                "operationId": "handle-get-lead",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Lead ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Lead"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates an existing lead with the provided information.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update an existing lead",
                "operationId": "handle-update-lead",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Lead ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated lead details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Lead"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Lead"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a lead with the specified ID from the database.",
                "summary": "Delete a lead by ID",
                "operationId": "handle-delete-lead",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Lead ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_leadsapi.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/send-email": {
            "post": {
                "description": "Parses the multipart form data into an EmailRequest struct and enqueues an email delivery task.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Send emails",
                "operationId": "handle-post-send-emails",
                "parameters": [
                    {
                        "type": "string",
                        "description": "From email address",
                        "name": "From",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "List of recipient email addresses",
                        "name": "To",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "List of CC email addresses",
                        "name": "CC",
                        "in": "formData"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "List of BCC email addresses",
                        "name": "BCC",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Email subject",
                        "name": "Subject",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Email body",
                        "name": "Body",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "file"
                        },
                        "collectionFormat": "csv",
                        "description": "List of file attachments",
                        "name": "Attachments",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "SMTP server address",
                        "name": "Server",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "SMTP server port",
                        "name": "Port",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "SMTP username",
                        "name": "Username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "SMTP password",
                        "name": "Password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_emailapi_emailsender.JSONResponse"
                        }
                    }
                }
            }
        },
        "/validate-email": {
            "get": {
                "security": [
                    {
                        "OAuth2Application": [
                            "callback"
                        ]
                    }
                ],
                "description": "Validates email addresses based on the provided parameters.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Validate emails",
                "operationId": "handle-get-validate-email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email address to validate",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "default": false,
                        "description": "Perform SMTP check",
                        "name": "smtp",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "default": false,
                        "description": "Perform SOCKS check",
                        "name": "socks",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_emailapi_emailvalidator.validateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal_emailapi_emailvalidator.validateResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal_emailapi_emailvalidator.validateResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.Lead": {
            "type": "object",
            "properties": {
                "emailAddress": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "leadListID": {
                    "description": "Foreign key for the lead list",
                    "type": "integer"
                },
                "linkedinURL": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "integer"
                },
                "scrapedDataFromLinkedin": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "internal_emailapi_emailpermutator.GenerateEmailParams": {
            "type": "object",
            "properties": {
                "domain1": {
                    "type": "string"
                },
                "domain2": {
                    "type": "string"
                },
                "domain3": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "middleName": {
                    "type": "string"
                },
                "nickName": {
                    "type": "string"
                }
            }
        },
        "internal_emailapi_emailsender.JSONResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "internal_emailapi_emailvalidator.validateResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "result": {}
            }
        },
        "internal_leadsapi.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "LeadzAura API",
	Description:      "Leadz Aura API is a service for leads generation and outreach..",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
