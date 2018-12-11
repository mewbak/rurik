#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// Output fragment color
out vec4 finalColor;

// NOTE: Add here your custom variables
uniform vec2 size = vec2(640, 480);

void main() 
{
	float x = 1.0/size.x;
	float y = 1.0/size.y;
    
	vec4 horizEdge = vec4(0.0);
	horizEdge -= texture2D(texture0, vec2(fragTexCoord.x - x, fragTexCoord.y - y))*1.0;
	horizEdge -= texture2D(texture0, vec2(fragTexCoord.x - x, fragTexCoord.y    ))*2.0;
	horizEdge -= texture2D(texture0, vec2(fragTexCoord.x - x, fragTexCoord.y + y))*1.0;
	horizEdge += texture2D(texture0, vec2(fragTexCoord.x + x, fragTexCoord.y - y))*1.0;
	horizEdge += texture2D(texture0, vec2(fragTexCoord.x + x, fragTexCoord.y    ))*2.0;
	horizEdge += texture2D(texture0, vec2(fragTexCoord.x + x, fragTexCoord.y + y))*1.0;
    
	vec4 vertEdge = vec4(0.0);
	vertEdge -= texture2D(texture0, vec2(fragTexCoord.x - x, fragTexCoord.y - y))*1.0;
	vertEdge -= texture2D(texture0, vec2(fragTexCoord.x    , fragTexCoord.y - y))*2.0;
	vertEdge -= texture2D(texture0, vec2(fragTexCoord.x + x, fragTexCoord.y - y))*1.0;
	vertEdge += texture2D(texture0, vec2(fragTexCoord.x - x, fragTexCoord.y + y))*1.0;
	vertEdge += texture2D(texture0, vec2(fragTexCoord.x    , fragTexCoord.y + y))*2.0;
	vertEdge += texture2D(texture0, vec2(fragTexCoord.x + x, fragTexCoord.y + y))*1.0;
    
	vec3 edge = sqrt((horizEdge.rgb*horizEdge.rgb) + (vertEdge.rgb*vertEdge.rgb));
	
	finalColor = vec4(edge, texture2D(texture0, fragTexCoord).a);
}