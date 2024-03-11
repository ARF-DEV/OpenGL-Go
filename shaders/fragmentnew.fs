#version 330 core

out vec4 FragColor;

in vec2 TexCords;
in vec3 Normal;
in vec3 FragPos;

uniform vec3 cameraPosition;

struct DirLight {
    vec3 direction;

    vec3 ambient;
    vec3 specular;
    vec3 diffuse;
};

uniform DirLight dirlight;

struct Material {
    sampler2D specular;
    sampler2D diffuse;

    float shininess;
};

uniform Material material;

vec3 CalcDirLighting(DirLight light, Material material, vec3 normal, vec3 viewDir) {
    // calculate diffuse light
    float diff = max(dot(light.direction, normal), 0);
    vec3 diffuse = diff * light.diffuse * vec3(texture(material.diffuse, TexCords));

    // calculate ambient light
    vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCords));

    // calculate specular light
    vec3 R = reflect(light.direction, normal);
    float spec = pow(max(dot(R, viewDir), 0), material.shininess);
    vec3 specular = spec * light.specular * vec3(texture(material.specular, TexCords));
    
    return (diffuse + ambient + specular);
}

struct PointLight {
    vec3 position;
    
    vec3 ambient;
    vec3 specular;
    vec3 diffuse;

    float constant;
    float linear;
    float quadratic;
};

#define NR_POINT_LIGHTS 4
uniform PointLight pointlights[NR_POINT_LIGHTS];

vec3 CalcPointLight(PointLight light, Material material, vec3 normal, vec3 viewDir, vec3 fragPos) {

    // calculate diffuse light
    vec3 lightDir = normalize(light.position - fragPos);
    float diff = max(dot(lightDir, normal), 0);
    vec3 diffuse = diff * light.diffuse * vec3(texture(material.diffuse, TexCords));

    // calculate ambient light
    vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCords));

    // calculate specular light
    vec3 R = reflect(lightDir, normal);
    float spec = pow(max(dot(R, viewDir), 0), material.shininess);
    vec3 specular = spec * light.specular * vec3(texture(material.specular, TexCords));

    float D = length(light.position - fragPos);
    float attenuation = 1 / (light.constant + light.linear * D + light.quadratic * (D * D));

    return (ambient + diffuse + specular) * attenuation;
}

struct SpotLight {
    vec3 direction;
    vec3 position;

    vec3 ambient;
    vec3 specular;
    vec3 diffuse;

    float constant;
    float linear;
    float quadratic;

    float outerCutoff;
    float innerCutoff;
};

uniform SpotLight spotlight;

vec3 CalcSpotLight(SpotLight light, Material material, vec3 normal, vec3 viewDir, vec3 fragPos) {

    // calculate diffuse light
    vec3 lightDir = normalize(light.position - fragPos);
    float diff = max(dot(lightDir, normal), 0);
    vec3 diffuse = diff * light.diffuse * vec3(texture(material.diffuse, TexCords));

    // calculate ambient light
    vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCords));

    // calculate specular light
    vec3 R = reflect(lightDir, normal);
    float spec = pow(max(dot(R, viewDir), 0), material.shininess);
    vec3 specular = spec * light.specular * vec3(texture(material.specular, TexCords));

    float D = length(light.position - fragPos);
    float attenuation = 1 / (light.constant + light.linear * D + light.quadratic * (D * D));

    float theta = dot(lightDir, normalize(-light.direction));
    float epsilon = light.innerCutoff - light.outerCutoff;
    float intensity = clamp((theta - light.outerCutoff) / epsilon, 0, 1); 
    if (theta > light.outerCutoff) {
        return (ambient + diffuse + specular) * attenuation * intensity;
    }
    return vec3(0);
}
void main() {
    vec3 norm = normalize(Normal);
    vec3 viewDir = normalize(cameraPosition - FragPos);

    // calculcate directional light
    vec3 result = CalcDirLighting(dirlight, material, norm, viewDir);
    
    // calculate point lights
    for(int i = 0; i < NR_POINT_LIGHTS; i++) {
        result += CalcPointLight(pointlights[i], material, norm, viewDir, FragPos);
    }

    // TOOD: calculcate spot light
  result += CalcSpotLight(spotlight, material, norm, viewDir, FragPos);

    FragColor = vec4(result, 1);

}
