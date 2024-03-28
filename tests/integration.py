import requests

PROXIES = {
    "http": "socks5://127.0.0.1:1080",
    "https": "socks5://127.0.0.1:1080",
}


def test_ipv4():
    response = requests.get("http://www.google.com", proxies=PROXIES)
    print(response.content)


def exec_tests():
    tests = [test_ipv4]

    for testsuite in tests:
        testsuite()


if __name__ == "__main__":
    exec_tests()
