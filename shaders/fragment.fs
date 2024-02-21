#version 330 core
out vec4 FragColor;
  
// in vec3 ourColor;
// in vec2 TexCoord;
in vec3 Normal;
in vec3 FragPos;

// uniform sampler2D texture1;
// uniform sampler2D texture2;

uniform vec3 lightPos;
uniform vec3 lightColor;
uniform vec3 objectColor;
void main()
{
    // FragColor = mix(texture(texture1, TexCoord), texture(texture2, TexCoord), 0.2);
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0);
    vec3 diffuse = lightColor * diff;

    float ambientStrength = 0.1;
    vec3 ambient = lightColor * ambientStrength;
    vec3 result = (ambient + diffuse) * objectColor;
    FragColor = vec4(result, 1);

}