window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "openapi.yaml",
    dom_id: '#swagger-ui',
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    layout: "StandaloneLayout"
  });
};
