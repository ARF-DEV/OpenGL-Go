#version 330 core
out vec4 FragColor;
  
// in vec3 ourColor;
in vec2 TexCords;
in vec3 Normal;
in vec3 FragPos;

// uniform sampler2D texture1;
// uniform sampler2D texture2;

uniform vec3 cameraPosition;
// uniform vec3 lightPos;
// uniform vec3 lightColor;

// uniform vec3 objectColor;

struct Material {
    sampler2D diffuse;
    sampler2D specular; 
    float shininess;
};

struct Light {
    float constant;
    float linear;
    float quadratic;

    vec3 position;
    vec3 direction;
    float cutoff;
    float outercutoff;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

uniform Material material;
uniform Light light;
void main()
{
    // FragColor = mix(texture(texture1, TexCoord), texture(texture2, TexCoord), 0.2);
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(light.position - FragPos);
    float diff = max(dot(norm, lightDir), 0);
    vec3 diffuse = diff * (light.diffuse );

    vec3 viewDir = normalize(cameraPosition - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);  
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
    vec3 specular = vec3(texture(material.specular, TexCords)) * spec * light.specular;  

    vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCords));

    float distance = length(light.position - FragPos);
    float attenuation = 1 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));
    // ambient *= attenuation;
    // diffuse *= attenuation;
    // specular *= attenuation;

    float theta = dot(lightDir, normalize(-light.direction));
    float epsinol = light.cutoff - light.outercutoff;
    float intesity = clamp((theta - light.outercutoff) / epsinol, 0, 1.0);
    vec3 result = vec3(0);
    if (theta > light.outercutoff) {
        diffuse *= intesity;
        specular *= intesity;
        result = ambient + diffuse + specular;
    } 
    else {
        result = ambient;
    }

    FragColor = vec4(result, 1) * texture(material.diffuse, TexCords);

}