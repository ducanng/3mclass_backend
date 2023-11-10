// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/v1/public/auth/login": {
            "post": {
                "description": "User login by phone or email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Public/Auth"
                ],
                "summary": "User login by phone or email",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "UserLoginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/authhandler.UserLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/authhandler.UserLoginResponse"
                        }
                    }
                }
            }
        },
        "/v1/public/auth/refresh_token": {
            "get": {
                "description": "Refresh Token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Public/Auth"
                ],
                "summary": "Refresh Token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/authhandler.RefreshTokenResponse"
                        }
                    }
                }
            }
        },
        "/v1/public/auth/registration": {
            "post": {
                "description": "User registration by phone, email, display name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Public/Auth"
                ],
                "summary": "User registration by phone, email, display name",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "UserRegistrationRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/authhandler.UserRegistrationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/authhandler.UserRegistrationResponse"
                        }
                    }
                }
            }
        },
        "/v1/public/auth/verify_otp": {
            "post": {
                "description": "Verify OTP",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Public/Auth"
                ],
                "summary": "Verify OTP",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "VerifyOTPRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/authhandler.VerifyOTPRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/authhandler.VerifyOTPResponse"
                        }
                    }
                }
            }
        },
        "/v1/public/u/user/logout": {
            "post": {
                "description": "Logout user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Public/Auth"
                ],
                "summary": "Logout user",
                "responses": {
                    "200": {
                        "description": "Logout ok!",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "authhandler.RefreshTokenResponse": {
            "type": "object"
        },
        "authhandler.UserLoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "authhandler.UserLoginResponse": {
            "type": "object"
        },
        "authhandler.UserRegistrationRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "phone_number"
            ],
            "properties": {
                "display_name": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string",
                    "maxLength": 12,
                    "minLength": 9
                }
            }
        },
        "authhandler.UserRegistrationResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "message": {
                    "description": "WaitingResendOTPSeconds uint64 ` + "`" + `json:\"waiting_resend_otp_seconds\"` + "`" + `",
                    "type": "string"
                },
                "next_action": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                },
                "session": {
                    "type": "string"
                }
            }
        },
        "authhandler.VerifyOTPRequest": {
            "type": "object"
        },
        "authhandler.VerifyOTPResponse": {
            "type": "object"
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