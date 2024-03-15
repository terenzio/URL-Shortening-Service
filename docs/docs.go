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
            "name": "Terence Liu",
            "url": "https://github.com/terenzio/URL-Shortening-Service",
            "email": "terenzio@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/redirect/{shortcode}": {
            "get": {
                "description": "NOTE: Copy the full url including the short code to the browser to be redirected. Do not use the Swagger UI here as it does not support redirection.",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "REDIRECT"
                ],
                "summary": "Redirects the user to the original URL based on the input short code.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Short Code",
                        "name": "shortcode",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "307": {
                        "description": "Redirected to original url - example: http://localhost:9000/api/v1/redirect/2v5ompxD",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Parameter missing - enter the short code in the URL path",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "No original URL exists for the given short code",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/url/add": {
            "post": {
                "description": "NOTE 1: In the JSON body, the \"original_url\" should contain proper formatting with either http or https. Example: https://www.google.com.\nNOTE 2: In the JSON body, the \"expiry\" date is optional, with the default expiration set to 30 days from now. The expiry time can be customized like this example: 2024-04-02T00:00:00Z.\nNOTE 3: In the JSON body, the \"custom_short_code\" is also optional. A unique custom short code can be set for the shortened URL.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "URL"
                ],
                "summary": "Creates a shortened link for the given original URL.",
                "parameters": [
                    {
                        "description": "Original URL, Expiry Time (optional), Custom Short Code (optional)",
                        "name": "original_url",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.AddURLRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Shortened URL",
                        "schema": {
                            "$ref": "#/definitions/domain.AddSuccessResponse"
                        }
                    }
                }
            }
        },
        "/url/display": {
            "get": {
                "description": "Displays the list of all shortened URLs mapped to their original ones in JSON format.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "URL"
                ],
                "summary": "Displays the list of all shortened URLs mapped to their original ones in JSON format.",
                "responses": {
                    "200": {
                        "description": "URL Mappings",
                        "schema": {
                            "$ref": "#/definitions/domain.URLMapping"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.AddSuccessResponse": {
            "type": "object",
            "properties": {
                "expiry": {
                    "type": "string"
                },
                "original_url": {
                    "type": "string"
                },
                "shortened_url": {
                    "type": "string"
                }
            }
        },
        "domain.AddURLRequest": {
            "type": "object",
            "properties": {
                "custom_short_code": {
                    "type": "string"
                },
                "expiry": {
                    "type": "string"
                },
                "original_url": {
                    "type": "string"
                }
            }
        },
        "domain.URLMapping": {
            "type": "object",
            "properties": {
                "expiry": {
                    "type": "string"
                },
                "original_url": {
                    "type": "string"
                },
                "short_code": {
                    "type": "string"
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "URL Shortening Service",
	Description:      "This is a URL shortening service that allows users to shorten long URLs especially built for TSMC.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
