import requests
import json
from jinja2 import Environment, FileSystemLoader, select_autoescape

# Define the base URL for the GitHub API
BASE_URL = "https://api.github.com"

def get_user_info(username):
    url = f"{BASE_URL}/users/{username}"
    response = requests.get(url)

    if response.status_code == 200:
        return response.json()
    else:
        print(f"Failed to fetch user info: {response.status_code}")
        return None

def get_user_repos(username):
    url = f"{BASE_URL}/users/{username}/repos"
    response = requests.get(url)

    if response.status_code == 200:
        return response.json()
    else:
        print(f"Failed to fetch user repos: {response.status_code}")
        return None

def get_repo_languages(username, repo_name):
    url = f"{BASE_URL}/repos/{username}/{repo_name}/languages"
    response = requests.get(url)

    if response.status_code == 200:
        return response.json()
    else:
        print(f"Failed to fetch languages for repo {repo_name}: {response.status_code}")
        return None

def main():
    username = input("Enter GitHub username: ")

    # Get user information
    user_info = get_user_info(username)
    if user_info:
        print(json.dumps(user_info, indent=4))

    # Get user repositories
    user_repos = get_user_repos(username)
    if user_repos:
        print(f"\nRepositories of {username}:")
        for repo in user_repos:
            print(f"- {repo['name']}")

        # Get languages used in each repository
        language_usage = {}
        for repo in user_repos:
            repo_languages = get_repo_languages(username, repo['name'])
            if repo_languages:
                for language, amount in repo_languages.items():
                    if language in language_usage:
                        language_usage[language] += amount
                    else:
                        language_usage[language] = amount

        # Calculate total amount of code
        total_amount = sum(language_usage.values())

        # Prepare data for the template
        languages = [
            {
                "language": language,
                "icon": "",  # You can add logic to determine the icon if needed
                "stats": f"{amount} bytes ({(amount / total_amount) * 100:.2f}%)"
            }
            for language, amount in language_usage.items()
        ]

        # Load and render the template with autoescape enabled
        env = Environment(
            loader=FileSystemLoader('.'),
            autoescape=select_autoescape(['html', 'xml', 'j2', 'json'])
        )
        template = env.get_template('template.json.j2')
        output = template.render(languages=languages)

        # Save the rendered template to a file
        with open('config.json', 'w') as f:
            f.write(output)

        print("Output saved to output.json")

if __name__ == "__main__":
    main()
