import requests
import json

# Define the base URL for the GitHub API
BASE_URL = "https://api.github.com"

# Function to get user information
def get_user_info(username):
    url = f"{BASE_URL}/users/{username}"
    response = requests.get(url)

    if response.status_code == 200:
        return response.json()
    else:
        print(f"Failed to fetch user info: {response.status_code}")
        return None

# Function to get user repositories
def get_user_repos(username):
    url = f"{BASE_URL}/users/{username}/repos"
    response = requests.get(url)

    if response.status_code == 200:
        return response.json()
    else:
        print(f"Failed to fetch user repos: {response.status_code}")
        return None

# Main function
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

if __name__ == "__main__":
    main()
