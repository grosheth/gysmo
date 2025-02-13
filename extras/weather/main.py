import urllib.request
import urllib.parse
import json
from jinja2 import Template

# Define the URL and parameters for the Open-Meteo API
url = "https://api.open-meteo.com/v1/forecast"

# Change the coordinates to the desired location
params = {
    "latitude": 45.5088,
    "longitude": -73.5878,
    "current_weather": True,
    "timezone": "America/New_York"
}

# Encode the parameters and create the full URL
query_string = urllib.parse.urlencode(params)
full_url = f"{url}?{query_string}"

# Make the HTTP request to the Open-Meteo API
with urllib.request.urlopen(full_url) as response:
    data = response.read().decode()

# Parse the JSON response
weather_data = json.loads(data)

# Extract relevant information from the response
latitude = weather_data['latitude']
longitude = weather_data['longitude']
elevation = weather_data['elevation']
timezone = weather_data['timezone']
timezone_abbreviation = weather_data['timezone_abbreviation']
utc_offset_seconds = weather_data['utc_offset_seconds']

# Extract current weather data
current_weather = weather_data['current_weather']
current_temperature = current_weather['temperature']
current_windspeed = current_weather['windspeed']
current_weather_code = current_weather['weathercode']
current_time = current_weather['time']

# Map weather codes to descriptions (simplified example)
weather_descriptions = {
    0: "clear sky",
    1: "mainly clear",
    2: "partly-cloudy",
    3: "overcast",
    45: "fog",
    48: "depositing rime fog",
    51: "drizzle",
    53: "drizzle",
    55: "drizzle",
    56: "freezing drizzle",
    57: "freezing drizzle",
    61: "rainy",
    63: "rainy",
    65: "rainy",
    66: "freezing rain",
    67: "freezing rain",
    71: "snowy",
    73: "snowy",
    75: "snowy",
    77: "snow grains",
    80: "rain showers",
    81: "rain showers",
    82: "rain showers",
    85: "snow showers",
    86: "snow showers",
    95: "thunderstorm",
    96: "thunderstorm",
    99: "thunderstorm"
}

current_weather_description = weather_descriptions.get(current_weather_code, "unknown")

# Debugging: Print the values to ensure they are correct
print(f"Latitude: {latitude}")
print(f"Longitude: {longitude}")
print(f"Elevation: {elevation}")
print(f"Timezone: {timezone}")
print(f"Timezone Abbreviation: {timezone_abbreviation}")
print(f"UTC Offset Seconds: {utc_offset_seconds}")
print(f"Current Temperature: {current_temperature}")
print(f"Current Windspeed: {current_windspeed}")
print(f"Current Weather Description: {current_weather_description}")
print(f"Current Time: {current_time}")

# Read the JSON template file
with open('template.json.j2', 'r') as template_file:
    template_content = template_file.read()

# Create a Jinja2 template from the template content
template = Template(template_content)

# Render the template with actual values
output_content = template.render(
    temperature=current_temperature,
    windspeed=current_windspeed,
    time=current_time,
    weather_description=current_weather_description
)

# Debugging: Print the output content to ensure it is correct
print("Output Content:")
print(output_content)

# Parse the filled-in template content as JSON
output_data = json.loads(output_content)

# Write the output data to a new JSON file
with open('config.json', 'w') as output_file:
    json.dump(output_data, output_file, indent=4)
