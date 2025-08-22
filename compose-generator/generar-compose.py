import sys

SERVER_YAML = """
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
"""

NETWORK_YAML = """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24

"""

CLIENT_YML = """
  client{id}:
    container_name: client{id}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={id}
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
    - ./client/config.yaml:/config.yaml
"""

def generate_script(file_destination, num_clients):

    with open(file_destination, 'w') as f:
        write_in_file(SERVER_YAML, f)
        write_dinamically(num_clients, f)
        write_in_file(NETWORK_YAML, f)

def write_in_file(yml: str, file):
    file.write(yml)

def write_dinamically(num_clients: int, file):
    for i in range(num_clients):
        yaml = CLIENT_YML.format(id=i + 1)
        write_in_file(yaml, file)

if __name__ == "__main__":
    try:
        if len(sys.argv) != 3:
            print("You should call this script as: ./generar-compose.py <output_file> <num_clients>")
            sys.exit(1)
        file_destination = sys.argv[1]
        client_nums = int(sys.argv[2])
        generate_script(file_destination, client_nums)
        print(f"Compose file '{file_destination}' generated with {client_nums} clients.")
        sys.exit(0)
    except ValueError as err: 
        print("You should provide a valid integer for the number of clients.")
        sys.exit(1)
    except Exception as e:
        print("An unexpected error occurred:", e)
        sys.exit(2)