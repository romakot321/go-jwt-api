#!/usr/bin/python3
import requests

base_url = "http://127.0.0.1:8000/api"
account = {"username": "user", "password": "user", "guid": "62381ebf-7f8b-4ec4-96c5-0d8a0879ea21"}


def register(account: dict) -> dict | None:
    resp = requests.post(base_url + "/auth/register", json=account)
    if 'already exists' in resp.text:
        return
    assert resp.status_code == 200, resp.text
    created_account = resp.json()["user"]
    assert created_account.get("guid") == account["guid"], f"Invalid created account guid: {created_account}"
    return created_account


def login(account: dict) -> dict:
    resp = requests.post(base_url + "/auth/login", json=account)
    assert resp.status_code == 200, resp.text
    tokens = resp.json()["token"]
    assert tokens.get("access_token") is not None, resp.text
    assert tokens.get("refresh_token") is not None, resp.text
    return tokens 


def refresh(refresh_token: str) -> dict:
    resp = requests.post(base_url + "/auth/refresh", headers={"Authorization": "Bearer " + refresh_token})
    assert resp.status_code == 200, resp.text
    tokens = resp.json()["token"]
    assert tokens.get("access_token") is not None, resp.text
    assert tokens.get("refresh_token") in (None, ""), resp.text
    return tokens


def get_user_info(access_token: str) -> dict:
    resp = requests.get(base_url + "/user/me", headers={"Authorization": "Bearer " + access_token})
    assert resp.status_code == 200, resp.text
    return resp.json()


def v1_login(guid: str) -> dict:
    resp = requests.post(base_url + "/auth/v1/login", params={"guid": guid})
    assert resp.status_code == 200, resp.text
    tokens = resp.json()["token"]
    assert tokens.get("access_token") is not None, resp.text
    assert tokens.get("refresh_token") is not None, resp.text
    return tokens


def v1_refresh(access: str, refresh: str) -> dict:
    resp = requests.post(
        base_url + "/auth/v1/refresh",
        params={"accessToken": access, "refreshToken": refresh}
    )
    assert resp.status_code == 200, resp.text
    tokens = resp.json()["token"]
    assert tokens.get("access_token") is not None, resp.text
    assert tokens.get("refresh_token") in (None, ""), resp.text
    return tokens


def v1_refresh_unauthorized(access: str, refresh: str) -> None:
    resp = requests.post(
        base_url + "/auth/v1/refresh",
        params={"accessToken": access, "refreshToken": refresh}
    )
    assert resp.status_code == 400, "Expected 'Invalid tokens', return: " + resp.text


print("-----Run requests for auth test...")
print("Register:", isinstance(register(account), dict | None))
tokens = login(account)
print("Login:", isinstance(tokens, dict))
print("Get current user:", isinstance(get_user_info(tokens["access_token"]), dict))
new_tokens = refresh(tokens["refresh_token"])
print("Refresh:", isinstance(new_tokens, dict))
print("Get current user with refreshed token:", isinstance(get_user_info(new_tokens["access_token"]), dict))

print("\n-----Run requests for V1 endpoints...")
v1_tokens = v1_login(account["guid"])
print("Login V1:", isinstance(v1_tokens, dict))
print("Refresh V1:", isinstance(v1_refresh(v1_tokens["access_token"], v1_tokens["refresh_token"]), dict))
print("Refresh V1 with old refresh token:", v1_refresh_unauthorized(v1_tokens["access_token"], v1_tokens["refresh_token"]) is None)
